package clients

import (
	"context"
	"net/http"
	"sync"
)

type Client struct {
	httpClient *http.Client
	ctx        context.Context
}

var client *Client

func NewClient(ctx context.Context) *Client {
	var once sync.Once
	once.Do(func() {
		client = &Client{
			httpClient: http.DefaultClient,
			ctx:        ctx,
		}
	})

	return client
}

func GetClient() *Client {
	return client
}
