package context

import (
	"context"
	"net/http"
)

type Key string

const (
	HTTPRequestContextKey     = Key("httpRequest")
	HTTPRequestBodyContextKey = Key("httpRequestBody")
)

func WithHTTPRequest(ctx context.Context, r *http.Request) context.Context {
	return context.WithValue(ctx, HTTPRequestContextKey, r)
}

func WithHTTPRequestBody(ctx context.Context, body []byte) context.Context {
	return context.WithValue(ctx, HTTPRequestBodyContextKey, body)
}
