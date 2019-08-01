package fbmes

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type processorMock struct {
}

func (m processorMock) Process(mes Messaging) error {
	logrus.Infof("messaging: %+v", mes)
	return nil
}

func TestVerificationHandler(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	SetDebugLogger(logrus.StandardLogger())

	r, err := http.NewRequest("", "?hub.verify_token=test&hub.challenge=CHALLENGE_ACCEPTED&hub.mode=subscribe", strings.NewReader(""))
	require.Nil(t, err)
	w := httptest.NewRecorder()

	h := VerificationHandler("test")
	h.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)

	body, err := w.Body.ReadString(0)
	require.Equal(t, io.EOF, err)
	assert.Equal(t, "CHALLENGE_ACCEPTED", body)

	r, err = http.NewRequest("", "?hub.verify_token=fail&hub.challenge=CHALLENGE_ACCEPTED&hub.mode=subscribe", strings.NewReader(""))
	require.Nil(t, err)
	w = httptest.NewRecorder()
	h.ServeHTTP(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	body, err = w.Body.ReadString('\n')
	require.Nil(t, err)
	assert.Equal(t, "wrong validation token\n", body)
}

func TestMessageHandler(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	SetDebugLogger(logrus.StandardLogger())

	r, err := http.NewRequest("", "", strings.NewReader(`
{  
   "object":"page",
   "entry":[  
      {  
         "id":"315216248953883",
         "time":1564645349030,
         "messaging":[  
            {  
               "sender":{  
                  "id":"1489738607768443"
               },
               "recipient":{  
                  "id":"315216248953883"
               },
               "timestamp":1564645348623,
               "message":{  
                  "mid":"DykAFwzfisNw2jkIbiJFCsL8UL5vmFWFBwzUvvrgVn055a9kGotUNSxDwMx4YETIxHKUn7HEe3cgCftttKr59Q",
                  "seq":0,
                  "text":"a"
               }
            }
         ]
      }
   ]
}`))
	require.Nil(t, err)
	w := httptest.NewRecorder()
	h := MessageHandler(processorMock{})

	h.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)

	r, err = http.NewRequest("", "", strings.NewReader(`{"object"`))
	require.Nil(t, err)
	w = httptest.NewRecorder()

	h.ServeHTTP(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	r, err = http.NewRequest("", "", strings.NewReader(`{"object": "fail"}`))
	require.Nil(t, err)
	w = httptest.NewRecorder()

	h.ServeHTTP(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)

}
