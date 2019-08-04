package fbmes

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
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
			debug(`fbmes: VerificationHandler: invalid token passed, expected: "%s", received: "%s"`, verifyToken, token)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(challenge))
		debug("fbmes: VerificationHandler: webhook verified")
	}
}

type User struct {
	ID string `json:"id"`
}

type Timestamp time.Time

func (t *Timestamp) MarshalJSON() ([]byte, error) {
	ts := time.Time(*t).Unix()
	return []byte(string(ts)), nil
}

func (t *Timestamp) UnmarshalJSON(b []byte) error {
	ts, err := strconv.Atoi(string(b))
	if err != nil {
		return err
	}
	*t = Timestamp(time.Unix(int64(ts), 0))
	return nil
}

type Message struct {
	MID        string `json:"mid,omitempty"`
	Text       string `json:"text,omitempty"`
	QuickReply *struct {
		Payload string `json:"payload,omitempty"`
	} `json:"quick_reply,omitempty"`
	Attachments *[]Attachment `json:"attachments,omitempty"`
	Attachment  *Attachment   `json:"attachment,omitempty"`
}

type Attachment struct {
	Type    string  `json:"type,omitempty"`
	Payload Payload `json:"payload,omitempty"`
}

type Payload struct {
	URL string `json:"url,omitempty"`
}

type Delivery struct {
	Seq       int64    `json:"seq"`
	Mids      []string `json:"mids"`
	Watermark int64    `json:"watermark"`
}

type Messaging struct {
	Sender    User      `json:"sender"`
	Recipient User      `json:"recipient"`
	Timestamp Timestamp `json:"timestamp"`
	Message   *Message  `json:"message,omitempty"`
	Delivery  *Delivery `json:"delivery,omitempty"`
}

type MessagingProcessor interface {
	Process(m Messaging) error
}

func MessageHandler(m MessagingProcessor) http.HandlerFunc {
	type request struct {
		Object string `json:"object"`
		Entry  []struct {
			Messaging []Messaging `json:"messaging"`
		} `json:"entry"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		var body request
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			debug(`fbmes: MessageHandler: error while decoding request body: "%s"`, err.Error())
			return
		}
		if body.Object != "page" {
			w.WriteHeader(http.StatusBadRequest)
			debug(`fbmes: MessageHandler: received invalid body.object: "%s"`, body.Object)
			return
		}

		debug("fbmes: MessageHandler: successfully received message: %+v", body)

		for _, entry := range body.Entry {
			for _, messaging := range entry.Messaging {
				err := m.Process(messaging)
				if err != nil {
					logrus.Errorf("message processor: %s", err.Error())
					continue
				}
				debug("MessageHandler: processed messaging: %+v", messaging)
			}
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("message received"))
	}
}
