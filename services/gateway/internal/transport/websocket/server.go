package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"lingva/api/gen"
	restLib "lingva/pkg/rest"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	"github.com/gorilla/websocket"
)

type WSRequestDTO struct {
	Setup      *WSSetupDTO `json:"setup,omitempty"`
	StdinChunk *string     `json:"stdin,omitempty"`
}

type WSSetupDTO struct {
	Lang string `json:"Lang" validate:"required"`
	Code string `json:"Code" validate:"required"`
}

func (s *WSSetupDTO) UnmarshalJSON(data []byte) error {
	type Alias WSSetupDTO

	aux := (*Alias)(s)
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	s.Lang = strings.ToUpper(aux.Lang)
	return nil
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// @Summary      Run code
// @Tags         run
// @Description  Connect to websocket
// @Accept       json
// @Produce      json
// @Success      200 {object} WSResponse
// @Failure      400 {object} WSResponse
// @Router       /ws/run [get]
func NewCodeRunHandler(log *slog.Logger, client gen.CodeRunServiceClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "transport.websocket.NewCodeRunHandler"
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Error(fmt.Sprintf("%s: failed to upgrade websocket connection: %s", op, err))
			return
		}
		defer func() { _ = conn.Close() }()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		stream, err := client.ExecuteCode(ctx)
		if err != nil {
			log.Error(fmt.Sprintf("%s: failed to execute code: %s", op, err.Error()))
			return
		}
		done := make(chan struct{})

		// ws -> grpc
		go func() {
			defer cancel()

			for {
				var wsReqDTO WSRequestDTO
				if err := conn.ReadJSON(&wsReqDTO); err != nil {
					if ctx.Err() != nil {
						break
					}
					if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
						log.Error(fmt.Sprintf("%s: failed to read websocket request: %s", op, err.Error()))
					}
					break
				}

				var wsReq WSRequest
				if wsReqDTO.Setup != nil {
					langInt, exists := gen.Language_value[wsReqDTO.Setup.Lang]
					if !exists {
						log.Error(fmt.Sprintf("%s: invalid language received: %s", op, wsReqDTO.Setup.Lang))
						render.JSON(w, r, restLib.Error("unsupported language"))
						return
					}
					wsReq.Setup = &WSInitialSetup{
						Lang: gen.Language(langInt),
						Code: wsReqDTO.Setup.Code,
					}
				}
				wsReq.StdinChunk = wsReqDTO.StdinChunk

				grpcReq := WSReq2GRPC(&wsReq)
				if err := stream.Send(grpcReq); err != nil {
					log.Error(fmt.Sprintf("%s: failed to send websocket request: %s", op, err.Error()))
					break
				}
			}
		}()

		// grpc -> ws
		go func() {
			defer close(done)

			for {
				resp, err := stream.Recv()
				if err != nil {
					log.Error(fmt.Sprintf("%s: failed to read websocket response: %s", op, err.Error()))
					break
				}
				wsResp := GRPCResp2WS(resp)
				if err := conn.WriteJSON(wsResp); err != nil {
					log.Error(fmt.Sprintf("%s: failed to write websocket response: %s", op, err.Error()))
					break
				}
				if resp.GetIsFinished() {
					break
				}
			}
		}()

		<-done
		_ = stream.CloseSend()
	}
}
