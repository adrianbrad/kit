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

func TestVerificationHandler(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	DEBUG = true
	Logger = logrus.StandardLogger()

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

	assert.Equal(t, http.StatusNotFound, w.Code)
	body, err = w.Body.ReadString('\n')
	require.Nil(t, err)
	assert.Equal(t, "wrong validation token\n", body)
}

func TestMessageHandler(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	DEBUG = true
	Logger = logrus.StandardLogger()

	r, err := http.NewRequest("", "", strings.NewReader(`{"object": "page", "entry": [{"messaging": [{"message": "TEST_MESSAGE"}]}]}`))
	require.Nil(t, err)
	w := httptest.NewRecorder()
	h := MessageHandler()

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
	assert.Equal(t, http.StatusNotFound, w.Code)

}
