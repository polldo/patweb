package main

import (
	"io"
	"net/http"
	"os"
	"time"

	"github.com/polldo/patweb/api"

	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	log.SetOutput(os.Stdout)

	log.Info("start demo")
	defer log.Info("demo complete")

	addr := "localhost:33888"
	go serve(log, addr)
	consume(log, addr)
}

func consume(log logrus.FieldLogger, host string) {
	time.Sleep(100 * time.Millisecond)
	r, err := http.Get("http://" + host + "/health")
	if err != nil {
		log.Error(err)
		return
	}
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		return
	}
	log.Infof("Response: %s", string(b))

}

func serve(log logrus.FieldLogger, addr string) {

	// Construct the mux for the API calls.
	mux := api.APIMux(api.APIConfig{
		Log: log,
	})

	// Construct a server to service the requests against the mux.
	api := http.Server{
		Addr:    addr,
		Handler: mux,
	}
	log.Error(api.ListenAndServe())
}
