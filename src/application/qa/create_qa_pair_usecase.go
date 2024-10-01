package qa

import (
	"context"

	"github.com/shinya-ac/server1Q1A/application/transaction"
	answerDomain "github.com/shinya-ac/server1Q1A/domain/answer"
	questionDomain "github.com/shinya-ac/server1Q1A/domain/question"
	"github.com/shinya-ac/server1Q1A/pkg/auth"
	"github.com/shinya-ac/server1Q1A/pkg/logging"
)

type CreateQaPairUseCase struct {
	questionRepo questionDomain.QuestionRepository
	answerRepo   answerDomain.AnswerRepository
	txManager    transaction.TransactionManager
}

func NewCreateQaPairUseCase(qr questionDomain.QuestionRepository, ar answerDomain.AnswerRepository, txManager transaction.TransactionManager) *CreateQaPairUseCase {
	return &CreateQaPairUseCase{
		questionRepo: qr,
		answerRepo:   ar,
		txManager:    txManager,
	}
}

type QaPairInputDto struct {
	Question string
	Answer   string
}

func (uc *CreateQaPairUseCase) Run(ctx context.Context, folderId string, qaPairs []QaPairInputDto) error {
	return uc.txManager.RunInTransaction(ctx, func(ctx context.Context) error {
		sub, err := auth.GetUserIDFromContext(ctx)
		if err != nil {
			logging.Logger.Error("userIdの取得に失敗", "error", err)
			return err
		}
		var questions []*questionDomain.Question
		var answers []*answerDomain.Answer
		for _, qa := range qaPairs {
			question, err := questionDomain.NewQuestion(sub, folderId, qa.Question)
			if err != nil {
				logging.Logger.Error("質問の作成に失敗しました", "error", err)
				return err
			}
			questions = append(questions, question)

			answer, err := answerDomain.NewAnswer(sub, question.Id, folderId, qa.Answer)
			if err != nil {
				logging.Logger.Error("解答の作成に失敗しました", "error", err)
				return err
			}
			answers = append(answers, answer)
		}

		err = uc.questionRepo.BulkCreate(ctx, questions)
		if err != nil {
			logging.Logger.Error("questionsのバルク作成失敗", "error", err)
			return err
		}

		err = uc.answerRepo.BulkCreate(ctx, answers)
		if err != nil {
			logging.Logger.Error("answersのバルク作成失敗", "error", err)
			return err
		}

		return nil
	})
}
