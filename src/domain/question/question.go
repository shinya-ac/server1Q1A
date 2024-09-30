package question

import (
	"time"

	"github.com/google/uuid"
)

type Question struct {
	Id        string
	UserId    string
	FolderId  string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewQuestion(userId, folderId, content string) *Question {
	id := uuid.New().String()
	return &Question{
		Id:        id,
		UserId:    userId,
		FolderId:  folderId,
		Content:   content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
