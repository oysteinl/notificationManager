package notifier

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

type MockClient struct{}

func (m *MockClient) Get(string) (*http.Response, error) {
	resp := &http.Response{Body: ioutil.NopCloser(bytes.NewBuffer([]byte("Test response"))), StatusCode: 200}
	return resp, nil
}
func (m *MockClient) Do(_ *http.Request) (*http.Response, error) {
	resp := &http.Response{Body: ioutil.NopCloser(bytes.NewBuffer([]byte("Test response"))), StatusCode: 200}
	return resp, nil
}

func TestTelegram_TextMessage(t *testing.T) {

	t.Run("Test empty message", func(t *testing.T) {
		n := TelegramNotification{Method: SendMessage, MessageToSend: "", Config: TelegramConfig{Token: "testtoken", ChatId: "testid"}}
		err := n.Send(&MockClient{})
		if err == nil {
			t.Fatal("did not get an error but wanted one")
		}
		if err != EmptyMessageErr {
			t.Errorf("got %q, wanted %q", err, EmptyMessageErr)
		}
	})

	t.Run("Test missing config", func(t *testing.T) {
		n := TelegramNotification{MessageToSend: "Test", Config: TelegramConfig{}}
		err := n.Send(&MockClient{})
		if err == nil {
			t.Fatal("did not get an error but wanted one")
		}
		if err != EmptyConfigErr {
			t.Errorf("got %q, but wanted %q", err, EmptyConfigErr)
		}
	})

	t.Run("Test successful sending", func(t *testing.T) {
		n := TelegramNotification{Method: SendMessage, MessageToSend: "Test", Config: TelegramConfig{ChatId: "id...", Token: "..token"}}
		err := n.Send(&MockClient{})
		if err != nil {
			t.Error(err)
		}
	})
}

func TestTelegram_PhotoMessage(t *testing.T) {

	t.Run("Test missing photo", func(t *testing.T) {
		n := TelegramNotification{Method: SendPhoto, MessageToSend: "Test", Config: TelegramConfig{ChatId: "id...", Token: "..token"}}
		err := n.Send(&MockClient{})
		if err == nil {
			t.Fatal("did not get an error but wanted one")
		}
		if err != NoPhotoErr {
			t.Errorf("got %q, wanted %q", err, NoPhotoErr)
		}
	})

	t.Run("Test successful Photo sending", func(t *testing.T) {
		n := TelegramNotification{Method: SendPhoto, Image: []byte("Here is a string...."), MessageToSend: "Test", Config: TelegramConfig{ChatId: "id...", Token: "..token"}}
		err := n.Send(&MockClient{})
		if err != nil {
			t.Error(err)
		}
	})
}
