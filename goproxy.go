package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/elazarl/goproxy"
)

func main() {
	addr := flag.String("addr", ":8080", "proxy listen address")
	mitm := flag.Bool("m", false, "Enable MITM for HTTPS connections")
	verbose := flag.Bool("v", true, "should every proxy request be logged to stdout")
	flag.Parse()
	proxy := goproxy.NewProxyHttpServer()
	if *mitm {
		proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	}
	proxy.Verbose = *verbose
	log.Fatal(http.ListenAndServe(*addr, proxy))
}
