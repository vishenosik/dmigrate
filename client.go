package migrate

import (
	"context"
	"fmt"
	"time"

	"github.com/dgraph-io/dgo/v240"
	"github.com/dgraph-io/dgo/v240/protos/api"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding/gzip"
)

type Config struct {
	User     string
	Password string
	Host     string
	Port     uint16
	Timeout  time.Duration
}

type cancelFunc = func()

func mustConnect(
	ctx context.Context,
	config Config,
) (*dgo.Dgraph, cancelFunc) {

	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)),
	}

	connection, err := grpc.NewClient(addr, opts...)
	if err != nil {
		panic(errors.Wrap(err, "grpc client connection"))
	}

	client := dgo.NewDgraphClient(
		api.NewDgraphClient(connection),
	)

	err = client.Login(ctx, config.User, config.Password)
	if err != nil {
		panic(errors.Wrap(err, "dgraph client login"))
	}

	return client, func() {
		err := connection.Close()
		if err != nil {
			panic(errors.Wrap(err, "grpc close connection"))
		}
	}
}
