package downloader

import (
	"errors"
	"github.com/antonyho/crhk-recorder/pkg/stream/url"
	"io/ioutil"
	"net/http"

	"github.com/antonyho/crhk-recorder/pkg/stream/resolver"
)

// Record the given channel
func Record(channel string) {
	channelName, playlist, err := resolver.Find(channel)
	if err != nil {
		panic(err)
	}
	for _, track := range playlist {
		resp, err := http.Get(url.StreamMediaURL(channelName, track.Path))
		if err != nil {
			panic(err)
		}
		media, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			panic(err)
		}

		if len(media) == 0 {
			panic(errors.New("Empty media file"))
		}
	}
}
