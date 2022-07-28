package main

import (
	"bytes"
	"encoding/json"
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
	go serve(log.WithField("app", "Server"), addr)
	consume(log.WithField("app", "Consumer"), addr)
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

func consume(log logrus.FieldLogger, host string) {
	type pload struct {
		Value string
	}
	post := func(p any) (status int, resp []byte) {
		pb, err := json.Marshal(p)
		if err != nil {
			log.Errorf("Unexpected error: %v", err)
			return 0, nil
		}
		r, err := http.Post("http://"+host+"/demo", "application/json", bytes.NewBuffer(pb))
		if err != nil {
			log.Errorf("Unexpected error: %v", err)
			return 0, nil
		}
		defer r.Body.Close()
		b, err := io.ReadAll(r.Body)
		if err != nil {
			log.Errorf("Unexpected error: %v", err)
			return 0, nil
		}
		return r.StatusCode, b
	}

	time.Sleep(100 * time.Millisecond)

	// Post for an error masked response.
	s, b := post(pload{Value: "mask"})
	log.Infof("Response: status %d, body %s", s, string(b))

	// Post for an error not masked response.
	s, b = post(pload{Value: "dont mask"})
	log.Infof("Response: status %d, body %s", s, string(b))
}
