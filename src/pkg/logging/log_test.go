package logging_test

import (
	"bytes"
	"io"
	"log/slog"
	"os"
	"sync"
	"testing"

	"github.com/shinya-ac/server1Q1A/pkg/logging"
	"github.com/stretchr/testify/assert"
)

func TestInitLogger(t *testing.T) {
	originalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	var once sync.Once
	logging.Logger = nil
	logging.InitLogger()
	once.Do(func() {
		logging.Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	})

	assert.NotNil(t, logging.Logger)

	oldLogger := logging.Logger
	logging.InitLogger()
	assert.Equal(t, oldLogger, logging.Logger)

	logging.Logger.Info("Test message")

	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	os.Stdout = originalStdout

	assert.Contains(t, buf.String(), "Test message")
}
