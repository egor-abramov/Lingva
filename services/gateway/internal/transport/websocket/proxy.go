package websocket

import (
	"lingva/api/gen"
)

type WSRequest struct {
	Setup      *WSInitialSetup `json:"setup,omitempty"`
	StdinChunk *string         `json:"stdin,omitempty"`
}

type WSInitialSetup struct {
	Lang gen.Language `json:"lang"`
	Code string       `json:"code"`
}

type WSResponse struct {
	Stdout     string `json:"stdout"`
	Stderr     string `json:"stderr"`
	IsFinished bool   `json:"is_finished"`
}

func WSReq2GRPC(wsReq *WSRequest) *gen.CodeRunRequest {
	var gRpcReq gen.CodeRunRequest

	if wsReq.Setup != nil {
		gRpcReq.Payload = &gen.CodeRunRequest_Setup{
			Setup: &gen.InitialSetup{
				Lang: wsReq.Setup.Lang,
				Code: wsReq.Setup.Code,
			},
		}
	} else if wsReq.StdinChunk != nil {
		gRpcReq.Payload = &gen.CodeRunRequest_StdinChunk{
			StdinChunk: *wsReq.StdinChunk,
		}
	}
	return &gRpcReq
}

func GRPCResp2WS(resp *gen.ExecutionResponse) *WSResponse {
	return &WSResponse{
		Stdout:     resp.Stdout,
		Stderr:     resp.Stderr,
		IsFinished: resp.IsFinished,
	}
}
