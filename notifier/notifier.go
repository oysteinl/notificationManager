package notifier

import "net/http"

type HttpClient interface {
	Get(string) (*http.Response, error)
	Do(req *http.Request) (*http.Response, error)
}
