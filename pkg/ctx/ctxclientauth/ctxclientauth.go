package ctxclientauth

import (
	"context"

	"github.com/Hidayathamir/user-activity-tracking-go/internal/model"
)

type clientAuthContextKey struct{}

var clientAuthKey = clientAuthContextKey{}

func Set(ctx context.Context, clientAuth *model.ClientAuth) context.Context {
	return context.WithValue(ctx, clientAuthKey, clientAuth)
}

func Get(ctx context.Context) *model.ClientAuth {
	if val, ok := ctx.Value(clientAuthKey).(*model.ClientAuth); ok {
		return val
	}
	return nil
}
