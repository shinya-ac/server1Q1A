package qa

import (
	"context"

	answerDomain "github.com/shinya-ac/server1Q1A/domain/answer"
	questionDomain "github.com/shinya-ac/server1Q1A/domain/question"
	"github.com/shinya-ac/server1Q1A/pkg/auth"
	"github.com/shinya-ac/server1Q1A/pkg/logging"
)

type CreateQaPairUseCase struct {
	questionRepo questionDomain.QuestionRepository
	answerRepo   answerDomain.AnswerRepository
}

func NewCreateQaPairUseCase(qr questionDomain.QuestionRepository, ar answerDomain.AnswerRepository) *CreateQaPairUseCase {
	return &CreateQaPairUseCase{
		questionRepo: qr,
		answerRepo:   ar,
	}
}

type QaPairInputDto struct {
	Question string
	Answer   string
}

func (uc *CreateQaPairUseCase) Run(ctx context.Context, folderId string, qaPairs []QaPairInputDto) error {
	sub, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		logging.Logger.Error("userIdの取得に失敗", "error", err)
		return err
	}
	for _, qa := range qaPairs {
		question := questionDomain.NewQuestion(sub, folderId, qa.Question)
		err := uc.questionRepo.Create(ctx, question)
		if err != nil {
			logging.Logger.Error("question作成失敗", "error", err)
			return err
		}

		answer := answerDomain.NewAnswer(sub, question.Id, folderId, qa.Answer)
		err = uc.answerRepo.Create(ctx, answer)
		if err != nil {
			logging.Logger.Error("answer作成失敗", "error", err)
			return err
		}
	}
	return nil
}
