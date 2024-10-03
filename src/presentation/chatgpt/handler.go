package chatgpt

import (
	"github.com/labstack/echo/v4"

	"github.com/shinya-ac/server1Q1A/application/chatgpt"
	"github.com/shinya-ac/server1Q1A/pkg/logging"
	validator "github.com/shinya-ac/server1Q1A/pkg/validator"
	"github.com/shinya-ac/server1Q1A/presentation/settings"
)

type handler struct {
	ocrUseCase         *chatgpt.OcrUseCase
	generateQasUseCase *chatgpt.GenerateQasUseCase
}

func NewHandler(
	ocrUseCase *chatgpt.OcrUseCase,
	generateQasUseCase *chatgpt.GenerateQasUseCase,
) handler {
	return handler{
		ocrUseCase:         ocrUseCase,
		generateQasUseCase: generateQasUseCase,
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

func (h handler) GenerateQas(ctx echo.Context) error {
	var params generateQasParams

	err := ctx.Bind(&params)
	if err != nil {
		logging.Logger.Error("paramsの形式が不正", "error", err)
		settings.ReturnBadRequest(ctx, err)
		return err
	}

	qas, err := h.generateQasUseCase.Run(ctx.Request().Context(), params.Content)
	if err != nil {
		return settings.ReturnStatusInternalServerError(ctx, err)
	}

	response := mapQasToResponse(qas)
	return settings.ReturnStatusCreated(ctx, response)
}

func mapQasToResponse(qas []*chatgpt.Qas) generateQasResponse {
	questions := []string{}
	answers := []string{}

	for _, qa := range qas {
		questions = append(questions, qa.Question)
		answers = append(answers, qa.Answer)
	}

	return generateQasResponse{
		Questions: questions,
		Answers:   answers,
	}
}
