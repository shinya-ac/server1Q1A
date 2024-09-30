package qa

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/shinya-ac/server1Q1A/application/qa"
)

type Handler struct {
	createQaPairUseCase *qa.CreateQaPairUseCase
}

func NewHandler(cqapuc *qa.CreateQaPairUseCase) *Handler {
	return &Handler{
		createQaPairUseCase: cqapuc,
	}
}

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
