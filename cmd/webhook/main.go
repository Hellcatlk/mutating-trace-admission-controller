package main

import (
	"context"
	"crypto/tls"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"mutating-trace-admission-controller/pkg/config"
	"mutating-trace-admission-controller/pkg/server"

	"github.com/golang/glog"
)

func main() {
	// Read configuration location from command line arg
	var configPath string
	flag.StringVar(&configPath, "configPath", "", "Path that points to the YAML configuration for this webhook.")
	flag.Parse()

	// Parse configuration
	cfg, err := config.Parse(configPath)
	if err != nil {
		glog.Errorf("parse configuration failed: %v", err)
		return
	}

	// Load certificates
	pair, err := cfg.LoadX509KeyPair()
	if err != nil {
		glog.Errorf("load X509 key pair failed: %v", err)
		return
	}

	// Config webhook server
	whsvr := &server.WebhookServer{
		Server: &http.Server{
			Addr:      ":443",
			TLSConfig: &tls.Config{Certificates: []tls.Certificate{pair}},
		},
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/mutate", whsvr.Serve)
	whsvr.Server.Handler = mux

	// Begin webhook server
	go func() {
		err := whsvr.Server.ListenAndServeTLS("", "")
		if err != nil {
			glog.Errorf("listen and serve webhook server failed: %v", err)
			return
		}
	}()

	// Listening OS shutdown singal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	// Shutdown webhook server
	err = whsvr.Server.Shutdown(context.Background())
	if err != nil {
		glog.Error(err)
	}
}
