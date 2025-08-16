package clients

import (
	"net/http"
	"strings"
	"time"
)

type SendMessageRateLimitedRoundTripper struct {
	// Это ограничитель скорости запросов на отправку сообщений. Более элегантного способа для этой fsm библиотеки я не нашел
	transport http.RoundTripper
	ticker    *time.Ticker
}

func NewSendMessageRateLimitedRoundTripper(transport http.RoundTripper, rate int) *SendMessageRateLimitedRoundTripper {
	return &SendMessageRateLimitedRoundTripper{
		transport: transport,
		ticker:    time.NewTicker(time.Second / time.Duration(rate)),
	}
}

func (r *SendMessageRateLimitedRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if strings.Contains(req.URL.Path, "sendMessage") {
		<-r.ticker.C
	}
	return r.transport.RoundTrip(req)
}

func (r *SendMessageRateLimitedRoundTripper) Stop() {
	r.ticker.Stop()
}
