package chatgpt

import (
	"github.com/labstack/echo/v4"

	"github.com/shinya-ac/server1Q1A/application/chatgpt"
	"github.com/shinya-ac/server1Q1A/pkg/logging"
	validator "github.com/shinya-ac/server1Q1A/pkg/validator"
	"github.com/shinya-ac/server1Q1A/presentation/settings"
)

type handler struct {
	ocrUseCase *chatgpt.OcrUseCase
}

func NewHandler(
	ocrUseCase *chatgpt.OcrUseCase,
) handler {
	return handler{
		ocrUseCase: ocrUseCase,
	}
}

func (h handler) Ocr(ctx echo.Context) error {
	logging.Logger.Info("Ocr 実行開始")

	var params struct {
		ImageURL string `json:"image_url" validate:"required,url"`
	}

	err := ctx.Bind(&params)
	if err != nil {
		logging.Logger.Error("リクエストパラメータの形式が不正", "error", err)
		settings.ReturnBadRequest(ctx, err)
		return err
	}

	validate := validator.GetValidator()
	err = validate.Struct(params)
	if err != nil {
		logging.Logger.Error("パラメータの内容が不正", "error", err)
		settings.ReturnStatusBadRequest(ctx, err)
		return err
	}

	response, err := h.ocrUseCase.HandleImageAnalysis(ctx.Request().Context(), params.ImageURL)
	if err != nil {
		logging.Logger.Error("ChatGPT APIの実行に失敗", "error", err)
		settings.ReturnError(ctx, err)
		return err
	}

	return settings.ReturnStatusOK(ctx, map[string]string{"response": response})
}
