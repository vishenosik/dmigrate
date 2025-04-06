package migrate

import (
	"context"
	"log"
	"testing"

	"github.com/dgraph-io/dgo/v240"
	"github.com/dgraph-io/dgo/v240/protos/api"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding/gzip"
)

const (
	addr     string = "localhost:9180"
	host     string = "localhost"
	port     uint16 = 9180
	user     string = "groot"
	password string = "password"
)

type clientSuite struct {
	client *dgo.Dgraph
	ctx    context.Context
	cancel cancelFunc
}

func newClientSuite(t *testing.T) *clientSuite {
	client, cancelCli := getTestingClient(t)
	ctx, cancelCtx := context.WithCancel(context.Background())

	return &clientSuite{
		client: client,
		ctx:    ctx,
		cancel: func() {
			cancelCtx()
			cancelCli()
		},
	}
}

func getTestingClient(t *testing.T) (*dgo.Dgraph, cancelFunc) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)),
	}

	connection, err := grpc.NewClient(addr, opts...)
	if err != nil {
		log.Fatal(err)
	}

	client := dgo.NewDgraphClient(
		api.NewDgraphClient(connection),
	)

	err = client.Login(context.Background(), user, password)
	if err != nil {
		log.Fatal(err)
	}

	return client, func() {
		err := client.Alter(context.Background(), &api.Operation{DropAll: true})
		require.NoError(t, err)
		err = connection.Close()
		require.NoError(t, err)
	}
}
