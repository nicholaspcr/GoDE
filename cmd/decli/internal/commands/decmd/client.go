package decmd

import (
	"context"
	"fmt"

	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func getClientAndContext(
	ctx context.Context,
) (
	context.Context,
	api.DifferentialEvolutionServiceClient,
	*grpc.ClientConn,
	error,
) {
	authToken, err := db.GetAuthToken()
	if err != nil {
		return nil, nil, nil, err
	}
	ctx = metadata.NewOutgoingContext(ctx, metadata.MD{
		"authorization": []string{fmt.Sprintf("Bearer %s", authToken)},
	})

	conn, err := grpc.NewClient(
		cfg.Server.GRPCAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, nil, nil, err
	}

	client := api.NewDifferentialEvolutionServiceClient(conn)
	return ctx, client, conn, nil
}
