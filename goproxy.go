package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/elazarl/goproxy"
)

func main() {
	addr := flag.String("addr", ":8080", "proxy listen address")
	mitm := flag.Bool("m", false, "Enable MITM for HTTPS connections")
	verbose := flag.Bool("v", true, "should every proxy request be logged to stdout")
	flag.Parse()
	proxy := goproxy.NewProxyHttpServer()
	if *mitm {
		caCert, err := tls.LoadX509KeyPair("ca.cer", "ca.key")
		if err == nil {
			goproxy.MitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectMitm, TLSConfig: goproxy.TLSConfigFromCA(&caCert)}
		}
		proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
		proxy.CertStore = &certCache{}
	}
	proxy.Verbose = *verbose
	log.Fatal(http.ListenAndServe(*addr, proxy))
}

type certCache struct {
	cache  sync.Map
	mutext sync.Mutex
}

func (storage *certCache) Fetch(hostname string, gen func() (*tls.Certificate, error)) (*tls.Certificate, error) {
	if cachedCert, found := storage.cache.Load(hostname); found {
		cert := cachedCert.(*tls.Certificate)
		return cert, nil
	}

	storage.mutext.Lock()
	defer storage.mutext.Unlock()

	if cachedCert, found := storage.cache.Load(hostname); found {
		cert := cachedCert.(*tls.Certificate)
		return cert, nil
	}

	cert, err := gen()
	if err != nil {
		return nil, fmt.Errorf("Could sign certificate for %s: %v", hostname, err)
	}

	storage.cache.Store(hostname, cert)

	return cert, nil
}
