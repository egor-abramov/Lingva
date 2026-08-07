package clients

import (
	"lingva/api/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewRunner(runnerAddr string) (gen.CodeRunServiceClient, error) {
	runnerConn, err := grpc.NewClient(runnerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	runnerClient := gen.NewCodeRunServiceClient(runnerConn)
	return runnerClient, nil
}
