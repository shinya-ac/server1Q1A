package error_test

import (
	"testing"

	errDomain "github.com/shinya-ac/server1Q1A/domain/error"
	"github.com/stretchr/testify/assert"
)

func TestNewError(t *testing.T) {
	t.Run("エラーメッセージが正しく設定されていることを確認", func(t *testing.T) {
		message := "エラーです。"
		err := errDomain.NewError(message)
		assert.Equal(t, message, err.Error())
	})
}

func TestNotFoundErr(t *testing.T) {
	t.Run("NotFoundErrのエラーメッセージが正しいことを確認", func(t *testing.T) {
		assert.Equal(t, "not found", errDomain.NotFoundErr.Error())
	})
}
