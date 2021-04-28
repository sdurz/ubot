package ubot

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/sdurz/axon"
)

// ServerSource is ServerSource that receives updates by exposing an http endpoint.
// The endpoint is exposed at http://hostname:<port>/bot<apiToken>.
func ServerSource(bot *Bot, ctx context.Context, updatesChan chan axon.O) {
	serverHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			body      []byte
			rawUpdate interface{}
			update    axon.O
			err       error
			ok        bool
		)

		body, err = ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}
		log.Println(string(body))

		if err = json.Unmarshal(body, rawUpdate); err != nil {
			log.Printf("Error decoding body: %v", err)
			http.Error(w, "can't decode body", http.StatusBadRequest)
			return
		}

		if update, ok = rawUpdate.(map[string]interface{}); !ok {
			log.Printf("Error decoding body: %v", err)
			return
		}

		if err = bot.process(ctx, update); err != nil {
			log.Println("Update processing error: ", err)
		}
	})

	if bot.Configuration.WebhookUrl == "" {
		log.Fatal("empty webhook url")
	}

	mux := http.NewServeMux()
	mux.Handle("/bot"+bot.Configuration.APIToken, serverHandler)
	srv := &http.Server{
		Addr:    bot.Configuration.ServerPort,
		Handler: mux,
	}

	if ok, err := bot.SetWebhook(axon.O{"url": bot.Configuration.WebhookUrl}); !ok || err != nil {
		log.Fatal("Can't set webhook")
		return
	}
	go http.ListenAndServe(bot.Configuration.ServerPort, mux)
	<-ctx.Done()

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("server Shutdown Failed:%+s", err)
	}

	log.Println("Server stopped")
}
