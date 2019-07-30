package fbmes

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Webhook struct {
	verifyToken string
}

func NewWebhook(verifyToken string) *Webhook {
	w := &Webhook{
		verifyToken: verifyToken,
	}

	Debug("started webhook using token: %s", verifyToken)

	return w
}

func (wh *Webhook) Verification(w http.ResponseWriter, r *http.Request) {
	challenge := r.URL.Query().Get("hub.challenge")
	token := r.URL.Query().Get("hub.verify_token")
	mode := r.URL.Query().Get("hub.mode")

	if token != wh.verifyToken || mode != "subscribe" {
		http.Error(
			w,
			"wrong validation token",
			http.StatusNotFound,
		)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(challenge))
	Debug("Webhook verified")
}

func (*Webhook) MessageHandler() http.HandlerFunc {
	type request struct {
		Object string `json:"object"`
		Entry  []struct {
			Messaging []interface{} `json:"messaging"`
		} `json:"entry"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var body request
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if body.Object != "page" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		for _, entry := range body.Entry {
			fmt.Println(entry.Messaging[0])
		}
	}
}
