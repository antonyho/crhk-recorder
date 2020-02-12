package downloader

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDownload(t *testing.T) {
	term := make(chan bool)
	err := Download("881", "/tmp/881hd", term)

	for {
		select {
		case <-time.After(10 * time.Second):
			term <- true
		}
	}

	assert.NoError(t, err)
}