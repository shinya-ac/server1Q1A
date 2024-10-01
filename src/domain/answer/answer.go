package answer

import (
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	errDomain "github.com/shinya-ac/server1Q1A/domain/error"
	"github.com/shinya-ac/server1Q1A/pkg/logging"
)

type Answer struct {
	Id         string
	UserId     string
	QuestionId string
	FolderId   string
	Content    string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func NewAnswer(userId, questionId, folderId, content string) (*Answer, error) {
	if utf8.RuneCountInString(content) < titleLengthMin {
		err := errDomain.NewError("解答の値が不正です。")
		logging.Logger.Error("解答の値が不正", "error", err)
		return nil, err
	}
	id := uuid.New().String()
	return &Answer{
		Id:         id,
		UserId:     userId,
		QuestionId: questionId,
		FolderId:   folderId,
		Content:    content,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}, nil
}

const (
	titleLengthMin = 1
)
