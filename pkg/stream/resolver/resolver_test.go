package resolver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCloudFrontResolverURL(t *testing.T) {
	playlistCFURL, chName, chPageURL, err := GetCloudFrontResolverURL("881")
	if assert.NoError(t, err) {
		t.Logf("Playlist CloudFront URL: %s", playlistCFURL)
		t.Logf("Channel Name: %s", chName)
		t.Logf("Channel Page URL: %s", chPageURL)
	}
}

func TestFind(t *testing.T) {
	channelName, livestreamServer, _, err := Find("881")

	if assert.NoError(t, err) {
		t.Logf("Channel Name: %s", channelName)
		t.Logf("Playlist: %+v", livestreamServer)
	}
}

func TestGetPlaylistAuthentication(t *testing.T) {
	playlistCFURL, _, channelPageURL, err := GetCloudFrontResolverURL("881")
	if err != nil {
		t.Errorf("test failed on the prerequisite step: %v", err)
		t.FailNow()
	}
	cfookies, streamServer, err := GetPlaylistAuthentication(channelPageURL, playlistCFURL)
	if assert.NoError(t, err) {
		t.Logf("Stream Server Hostname: %s", streamServer)
		t.Logf("CloudFront Cookies: %+v", cfookies)
	}
}

func TestGetPlaylist(t *testing.T) {
	playlistCFURL, chName, channelPageURL, err := GetCloudFrontResolverURL("881")
	if err != nil {
		t.Errorf("test failed on the prerequisite step - get CF URL: %v", err)
		t.FailNow()
	}
	cfookies, streamServer, err := GetPlaylistAuthentication(channelPageURL, playlistCFURL)
	if err != nil {
		t.Errorf("test failed on the prerequisite step - get CF cookies: %v", err)
		t.FailNow()
	}

	playlist, err := GetPlaylist(chName, streamServer, cfookies)
	if assert.NoError(t, err) {
		t.Logf("Playlist: %+v", playlist)
	}
}
