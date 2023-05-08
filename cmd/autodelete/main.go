package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
	rdebug "runtime/debug"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	autodelete "github.com/spike01/AutoDelete"
)

var flagShardID = flag.Int("shard", -1, "shard ID of this bot")
var flagNoHttp = flag.Bool("nohttp", false, "skip http handler")
var flagMetricsPort = flag.Int("metrics", 6130, "port for metrics listener; shard ID is added")
var flagMetricsListen = flag.String("metricslisten", "127.0.0.4", "addr to listen on for metrics handler")

func main() {
	var conf autodelete.Config

	// Hardcoded because whatever, I make the rules here
	conf.ClientId = os.Getenv("CLIENT_ID")
	conf.ClientSecret = os.Getenv("CLIENT_SECRET")
	conf.BotToken = os.Getenv("BOT_TOKEN")
	conf.AdminUser = os.Getenv("ADMIN_USER")
	conf.Http = HTTP{
		Listen: "localhost:2202",
		Public: "https://home.riking.org",
	}
	conf.Backlog_limit = 1000
	conf.ErrorLogCh = ""
	conf.StatusMessage = "in the garbage"

	flag.Parse()

	if conf.BotToken == "" {
		fmt.Println("bot token must be specified")
	}
	if conf.Shards > 0 && *flagShardID == -1 {
		fmt.Println("This AutoDelete instance is configured to be sharded; please specify --shard=n")
		return
	}
	if *flagShardID > conf.Shards {
		fmt.Println("error: shard nbr is > shard count")
		return
	}

	b := autodelete.New(conf)

	err = b.ConnectDiscord(*flagShardID, conf.Shards)
	if err != nil {
		fmt.Println("connect error:", err)
		return
	}

	var privHttp http.ServeMux
	var pubHttp http.ServeMux

	go func() {
		for {
			time.Sleep(time.Hour * 1)
			rdebug.FreeOSMemory()
		}
	}()
	go func() {
		privHttp.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
		privHttp.Handle("/metrics", promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{}))
		metricSvr := &http.Server{
			Handler: &privHttp,
			Addr:    fmt.Sprintf("%s:%d", *flagMetricsListen, *flagMetricsPort+*flagShardID),
		}

		err := metricSvr.ListenAndServe()
		fmt.Println("exiting metric server", err)
	}()

	if !*flagNoHttp {
		fmt.Printf("url: %s%s\n", conf.HTTP.Public, "/discord_auto_delete/oauth/start")
		pubHttp.HandleFunc("/discord_auto_delete/oauth/start", b.HTTPOAuthStart)
		pubHttp.HandleFunc("/discord_auto_delete/oauth/callback", b.HTTPOAuthCallback)
		pubSrv := &http.Server{
			Handler: &pubHttp,
			Addr:    conf.HTTP.Listen,
		}
		err = pubSrv.ListenAndServe()
		fmt.Println("exiting main()", err)
	} else {
		select {}
	}
}
