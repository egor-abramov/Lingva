package rest

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
)

type AnalyzeRequestDTO struct {
	Lang string `json:"Lang" validate:"required"`
	Code string `json:"Code" validate:"required"`
}

func (r *AnalyzeRequestDTO) UnmarshalJSON(data []byte) error {
	type Alias AnalyzeRequestDTO

	aux := (*Alias)(r)
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	r.Lang = strings.ToUpper(aux.Lang)
	return nil
}

// @Summary      Code analyzing
// @Tags         analyze
// @Accept       json
// @Produce      json
// @Param        request body AnalyzeRequestDTO true "Payload"
// @Success      200 {object} restLib.Response[gen.AnalyzeResponse]
// @Failure      400 {object} restLib.Response[any]
// @Router       /rest/analyze [post]
func NewCodeAnalyzeHandler(log *slog.Logger, analyzer gen.CodeAnalyzeServiceClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.NewAnalyzeHandler"

		reqDTO, ok := restLib.Decode[AnalyzeRequestDTO](log, w, r)
		if !ok {
			return
		}

		langInt, exists := gen.Language_value[reqDTO.Lang]
		if !exists {
			log.Error(fmt.Sprintf("%s: invalid language received: %s", op, reqDTO.Lang))
			render.JSON(w, r, restLib.Error("unsupported language"))
			return
		}
		grpcReq := &gen.AnalyzeRequest{
			Lang: gen.Language(langInt),
			Code: reqDTO.Code,
		}

		ctx := context.Background()
		resp, err := analyzer.Analyze(ctx, grpcReq)
		if err != nil {
			log.Error(fmt.Sprintf("%s: error analyzing code: %s", op, err.Error()))
			render.JSON(w, r, restLib.Error("error analyzing code"))
			return
		}
		render.JSON(w, r, restLib.OK(resp))
	}
}
