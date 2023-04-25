package util

import (
	"context"

	"google.golang.org/grpc/metadata"
)

func GetSteamHeader(ctx context.Context, key string) string {
	var values []string
	var token string

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		values = md.Get(key)
	}

	if len(values) > 0 {
		token = values[0]
	}

	return token
}
