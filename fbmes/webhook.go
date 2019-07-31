package fbmes

import (
	"encoding/json"
	"net/http"
)

func VerificationHandler(verifyToken string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body.Close()

		challenge := r.URL.Query().Get("hub.challenge")
		token := r.URL.Query().Get("hub.verify_token")
		mode := r.URL.Query().Get("hub.mode")

		if token != verifyToken || mode != "subscribe" {
			http.Error(
				w,
				"wrong validation token",
				http.StatusBadRequest,
			)
			debug(`VerificationHandler: invalid token passed, expected: "%s", received: "%s"`, verifyToken, token)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(challenge))
		debug("VerificationHandler: webhook verified")
	}
}

type Messaging struct {
}

type MessagingProcessor interface {
}

func MessageHandler() http.HandlerFunc {
	type request struct {
		Object string `json:"object"`
		Entry  []struct {
			Messaging []interface{} `json:"messaging"`
		} `json:"entry"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var body request
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			debug(`MessageHandler: error while decoding request body: "%s"`, err.Error())
			return
		}
		if body.Object != "page" {
			w.WriteHeader(http.StatusBadRequest)
			debug(`MessageHandler: received invalid body.object: "%s"`, body.Object)
			return
		}

		debug("MessageHandler: successfully received message: %+v", body)

		for _, entry := range body.Entry {
			debug("MessageHandler: processed entry: %+v", entry)
		}
	}
}
