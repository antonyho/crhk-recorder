package resolver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPlaylistLocatorPageURL(t *testing.T) {
	playlistLocatorURL, channelName, err := GetPlaylistLocatorPageURL("881")
	if assert.NoError(t, err) {
		t.Logf("Playlist Locator URL: %s", playlistLocatorURL)
		t.Logf("Channel Name: %s", channelName)
	}
}

func TestFind(t *testing.T) {
	channelName, playlist, _, err := Find("881")

	if assert.NoError(t, err) {
		t.Logf("Channel Name: %s", channelName)
		t.Logf("Playlist: %+v", playlist)
	}
}
