package folder

import (
	"unicode/utf8"

	"github.com/google/uuid"

	errDomain "github.com/shinya-ac/server1Q1A/domain/error"
	"github.com/shinya-ac/server1Q1A/pkg/logging"
)

type Folder struct {
	Id     string
	Title  string
	UserId string
}

func newFolder(id, title, userId string) (*Folder, error) {
	if utf8.RuneCountInString(title) < titleLengthMin || utf8.RuneCountInString(title) > titleLengthMax {
		err := errDomain.NewError("タイトルの値が不正です。")
		logging.Logger.Error("タイトルの値が不正", "error", err)
		return nil, err
	}

	return &Folder{
		Id:     id,
		Title:  title,
		UserId: userId,
	}, nil
}

func NewFolder(userId, title string) (*Folder, error) {
	id := uuid.New().String()
	return newFolder(id, title, userId)
}

func ReconstructFolder(id, title, userId string) (*Folder, error) {
	return newFolder(id, title, userId)
}

func (f *Folder) GetId() string {
	return f.Id
}

const (
	titleLengthMin = 1
	titleLengthMax = 50
)
