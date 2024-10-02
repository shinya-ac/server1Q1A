package qa

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/shinya-ac/server1Q1A/application/qa"
	"github.com/shinya-ac/server1Q1A/presentation/settings"
)

type Handler struct {
	createQaPairUseCase *qa.CreateQaPairUseCase
	getQaPairsUseCase   *qa.GetQaPairsUseCase
}

func NewHandler(cqapuc *qa.CreateQaPairUseCase, gqpus *qa.GetQaPairsUseCase) *Handler {
	return &Handler{
		createQaPairUseCase: cqapuc,
		getQaPairsUseCase:   gqpus,
	}
}

// CreateQaPairs godoc
// @Summary QAペアを作成する
// @Description 質問と回答のペアを作成する
// @Tags QA
// @Accept  json
// @Produce  json
// @Param folder_id path string true "Folder ID"
// @Param request body CreateQaPairRequest true "質問と回答のペアのリクエスト"
// @Success 201 {object} CreateQaResponse
// @Router /folders/{folder_id}/qa [post]
// @Security BearerAuth

func (h *Handler) CreateQaPairs(ctx echo.Context) error {
	var req CreateQaPairRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, "paramsの形式が不正")
	}

	folderId := ctx.Param("folder_id")

	qaPairs := []qa.QaPairInputDto{}
	for _, selectedQa := range req.SelectedQas {
		qaPairs = append(qaPairs, qa.QaPairInputDto{
			Question: selectedQa.Question,
			Answer:   selectedQa.Answer,
		})
	}

	if err := h.createQaPairUseCase.Run(ctx.Request().Context(), folderId, qaPairs); err != nil {
		return ctx.JSON(http.StatusInternalServerError, "QAペアの作成に失敗しました")
	}

	response := CreateQaResponse{
		Res: "一問一答の作成に成功しました",
	}
	return ctx.JSON(http.StatusCreated, response)
}

// GetQaPairs godoc
// @Summary QAペアの一覧を取得
// @Description 質問と回答のペアの一覧を取得する
// @Tags QA
// @Accept  json
// @Produce  json
// @Param folder_id path string true "Folder ID"
// @Success 200 {array} qa.QaPairOutputDto "質問と回答のペアのリスト"
// @Router /folders/{folder_id}/qa [get]
// @Security BearerAuth

func (h *Handler) GetQaPairs(ctx echo.Context) error {
	folderId := ctx.Param("folder_id")

	qaPairs, err := h.getQaPairsUseCase.Run(ctx.Request().Context(), folderId)
	if err != nil {
		return settings.ReturnStatusInternalServerError(ctx, err)
	}

	return settings.ReturnStatusOK(ctx, qaPairs)
}
