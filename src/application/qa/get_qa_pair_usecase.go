package qa

import (
	"context"
	"errors"

	answerDomain "github.com/shinya-ac/server1Q1A/domain/answer"
	questionDomain "github.com/shinya-ac/server1Q1A/domain/question"
	"github.com/shinya-ac/server1Q1A/pkg/auth"
	"github.com/shinya-ac/server1Q1A/pkg/logging"
)

type QaPair struct {
	Question string
	Answer   string
}

type GetQaPairsUseCase struct {
	questionRepo questionDomain.QuestionRepository
	answerRepo   answerDomain.AnswerRepository
}

func NewGetQaPairsUseCase(qr questionDomain.QuestionRepository, ar answerDomain.AnswerRepository) *GetQaPairsUseCase {
	return &GetQaPairsUseCase{
		questionRepo: qr,
		answerRepo:   ar,
	}
}

func (uc *GetQaPairsUseCase) Run(ctx context.Context, folderId string) ([]QaPair, error) {
	userId, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		logging.Logger.Error("ユーザーIDの取得に失敗", "error", err)
		return nil, err
	}

	questions, err := uc.questionRepo.GetQuestionsByFolderId(ctx, folderId)
	if err != nil {
		logging.Logger.Error("質問の取得に失敗", "error", err)
		return nil, err
	}

	var filteredQuestions []*questionDomain.Question
	for _, question := range questions {
		if question.UserId == userId {
			filteredQuestions = append(filteredQuestions, question)
		}
	}

	if len(filteredQuestions) == 0 {
		return nil, errors.New("このフォルダーにはこのユーザーによって作成された質問はありません")
	}

	var questionIds []string
	for _, question := range filteredQuestions {
		questionIds = append(questionIds, question.Id)
	}

	answers, err := uc.answerRepo.GetAnswersByQuestionIds(ctx, questionIds)
	if err != nil {
		logging.Logger.Error("回答の取得に失敗", "error", err)
		return nil, err
	}

	var qaPairs []QaPair
	for _, question := range filteredQuestions {
		for _, answer := range answers {
			if question.Id == answer.QuestionId && answer.UserId == userId {
				qaPairs = append(qaPairs, QaPair{
					Question: question.Content,
					Answer:   answer.Content,
				})
				break
			}
		}
	}

	return qaPairs, nil
}
