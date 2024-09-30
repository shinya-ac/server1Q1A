package answer

import (
	"time"

	"github.com/google/uuid"
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

func NewAnswer(userId, questionId, folderId, content string) *Answer {
	id := uuid.New().String()
	return &Answer{
		Id:         id,
		UserId:     userId,
		QuestionId: questionId,
		FolderId:   folderId,
		Content:    content,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}
