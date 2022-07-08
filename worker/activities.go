package worker

import (
	"context"
)

func SayHello(ctx context.Context, name string) (string, error) {
	return "Hello " + name, nil
}
