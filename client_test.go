package migrate

import (
	"context"
	"log"

	"github.com/dgraph-io/dgo/v240"
	"github.com/dgraph-io/dgo/v240/protos/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding/gzip"
)

const (
	addr     = "localhost:9180"
	user     = "groot"
	password = "password"
)

type cancelFunc func()

func getTestingDgraphClient() (*dgo.Dgraph, cancelFunc) {
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

	return client, func() { connection.Close() }
}
