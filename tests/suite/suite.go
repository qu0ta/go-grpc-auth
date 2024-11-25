package suite

import (
	"context"
	"github.com/qu0ta/go-grpc-auth/internal/config"
	authv1 "github.com/qu0ta/pet-proto/gen/go/auth"
	"testing"
)

type Suite struct {
	*testing.T
	Cfg        *config.Config
	AuthClient authv1.AuthClient
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

}
