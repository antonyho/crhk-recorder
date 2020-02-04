package url

import (
	"regexp"
	"strings"
	"text/template"
)

const (
	// RadioStationPage is the URL template for the radio station
	// StreamName can be 881, 903, 864
	RadioStationPage = "https://www.881903.com/live/{{.Name}}"

	// PlaylistLocatorURLJSON is the pattern of JSON which contains playlist locator URL value
	PlaylistLocatorURLJSON = `"liveJsUrl":"(.*?)"`

	// Playlist is the URL template for the radio stream playlist
	Playlist = "https://live.881903.com/edge-aac/{{.Name}}/chunks.m3u8"
)

// Station is the template structure for RadioStationPage URL
type Station struct {
	Name string
}

// RadioChannelPageURL builds the radio channel page URL
func RadioChannelPageURL(channel string) string {
	tpl, err := template.New("stationurl").Parse(RadioStationPage)
	if err != nil {
		panic(err)
	}

	channelURLStr := new(strings.Builder)
	err = tpl.Execute(channelURLStr, Station{Name: channel})
	if err != nil {
		panic(err)
	}

	return channelURLStr.String()
}

// FetchPlaylistLocatorURL fetches the playlist locator URL from
// radio station page HTML code
func FetchPlaylistLocatorURL(radioStationPageHTML string) (string, bool) {
	matched := regexp.MustCompile(PlaylistLocatorURLJSON).FindStringSubmatch(radioStationPageHTML)
	if len(matched) < 2 {
		return "", false
	}
	return matched[1], true
}

// PlaylistURL builds the radio channel stream playlist URL
func PlaylistURL(channel string) string {
	tpl, err := template.New("playlisturl").Parse(Playlist)
	if err != nil {
		panic(err)
	}

	playlistURLStr := new(strings.Builder)
	err = tpl.Execute(playlistURLStr, Station{Name: channel})
	if err != nil {
		panic(err)
	}

	return playlistURLStr.String()
}
