package migrate

import (
	"context"
	"fmt"
	"time"

	"github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"
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

func connect(
	ctx context.Context,
	config Config,
) (*dgo.Dgraph, error) {

	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)),
	}

	connection, err := grpc.NewClient(addr, opts...)
	if err != nil {
		return nil, err
	}

	client := dgo.NewDgraphClient(
		api.NewDgraphClient(connection),
	)

	err = client.Login(ctx, config.User, config.Password)
	if err != nil {
		return nil, err
	}

	return client, nil
}
