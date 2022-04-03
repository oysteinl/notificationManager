package notifier

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

var EmptyMessageErr = errors.New("empty message")
var NoPhotoErr = errors.New("photo missing")
var EmptyConfigErr = errors.New("config missing")
var InvalidMethodErr = errors.New("invalid method")

type TelegramNotification struct {
	MessageToSend string
	Image         []byte
	Config        TelegramConfig
	Method        int8
}

type TelegramConfig struct {
	Token  string
	ChatId string
}

const (
	SendMessage = iota + 1
	SendPhoto
)

func (t TelegramNotification) Send(client HttpClient) error {
	if t.Config.ChatId == "" || t.Config.Token == "" {
		return EmptyConfigErr
	}
	var err error

	switch t.Method {
	case SendMessage:
		err = t.sendMessage(client)
	case SendPhoto:
		err = t.sendImage(client)
	default:
		return InvalidMethodErr
	}

	return err
}

func (t TelegramNotification) sendImage(client HttpClient) error {

	if t.Image == nil {
		return NoPhotoErr
	}
	// Create buffer
	buf := new(bytes.Buffer) // caveat IMO dont use this for large files, \
	// create a tmpfile and assemble your multipart from there (not tested)
	w := multipart.NewWriter(buf)

	// Create file field
	fw, err := w.CreateFormFile("photo", "snap.jpeg")
	if err != nil {
		return err
	}
	_, err = fw.Write(t.Image)
	if err != nil {
		return err
	}
	// Important if you do not close the multipart writer you will not have a
	// terminating boundary
	w.Close()
	repoUrl := fmt.Sprintf("https://api.telegram.org/bot%s/sendPhoto?chat_id=%s&caption=%s",
		t.Config.Token,
		t.Config.ChatId,
		t.MessageToSend)
	req, err := http.NewRequest("POST", repoUrl, buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		b, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		return errors.New(fmt.Sprintf("Telegram notify did not work. HTTP: %d, %s", res.StatusCode, string(b)))
	}
	return nil
}

func (t TelegramNotification) sendMessage(client HttpClient) error {
	if t.MessageToSend == "" {
		return EmptyMessageErr
	}
	resp, err := client.Get(fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s",
		t.Config.Token,
		t.Config.ChatId,
		t.MessageToSend))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return errors.New(fmt.Sprintf("Telegram notify did not work. HTTP: %d", resp.StatusCode))
	}
	return nil
}
