package clients

import (
	"lingva/api/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewAnalyzer(analyzerAddr string) (gen.CodeAnalyzeServiceClient, error) {
	analyzerConn, err := grpc.NewClient(analyzerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	analyzerClient := gen.NewCodeAnalyzeServiceClient(analyzerConn)

	return analyzerClient, nil
}
