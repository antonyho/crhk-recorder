package downloader

import (
	"github.com/antonyho/crhk-recorder/pkg/stream/url"
	"net/http"

	"github.com/antonyho/crhk-recorder/pkg/stream/resolver"
)

// Record the given channel
func Record(channel string) {
	playlist, err := resolver.Find(channel)
	if err != nil {
		panic(err)
	}
	for _, track := range playlist {
		resp, err := http.Get(url.StreamMediaURL("881hd", track.Path))
	}
}
