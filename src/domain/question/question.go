package question

import (
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	errDomain "github.com/shinya-ac/server1Q1A/domain/error"
	"github.com/shinya-ac/server1Q1A/pkg/logging"
)

type Question struct {
	Id        string
	UserId    string
	FolderId  string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewQuestion(userId, folderId, content string) (*Question, error) {
	if utf8.RuneCountInString(content) < titleLengthMin {
		err := errDomain.NewError("質問の値が不正です。")
		logging.Logger.Error("質問の値が不正", "error", err)
		return nil, err
	}
	id := uuid.New().String()
	return &Question{
		Id:        id,
		UserId:    userId,
		FolderId:  folderId,
		Content:   content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

const (
	titleLengthMin = 1
)

func (q *Question) ParseTimeFields(createdAt, updatedAt []byte) error {
	var err error
	q.CreatedAt, err = parseTime(createdAt)
	if err != nil {
		return fmt.Errorf("created_atのパースに失敗: %w", err)
	}

	q.UpdatedAt, err = parseTime(updatedAt)
	if err != nil {
		return fmt.Errorf("updated_atのパースに失敗: %w", err)
	}

	return nil
}

func parseTime(t []byte) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", string(t))
}
