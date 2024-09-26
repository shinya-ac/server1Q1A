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

func NewFolder(
	Title string,
	UserId string,
) (*Folder, error) {
	if utf8.RuneCountInString(Title) < titleLengthMin || utf8.RuneCountInString(Title) > titleLengthMax {
		err := errDomain.NewError("タイトルの値が不正です。")
		logging.Logger.Error("タイトルの値が不正", "error", err)
		return nil, err
	}
	id, err := uuid.NewRandom()
	if err != nil {
		logging.Logger.Error("UUIDの生成に失敗", "error", err)
		return nil, err
	}

	return &Folder{
		Id:     id.String(),
		Title:  Title,
		UserId: UserId,
	}, nil
}

func (f *Folder) GetId() string {
	return f.Id
}

const (
	titleLengthMin = 1
	titleLengthMax = 50
)
