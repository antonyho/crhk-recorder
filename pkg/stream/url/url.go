package url

import (
	"net/url"
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

	// PlaylistURLTemplate is the URL template for the radio stream playlist
	PlaylistURLTemplate = "https://live.881903.com/edge-aac/{{.Name}}/chunks.m3u8"

	// MediaURLTemplate is the URL for the stream media file
	StreamMediaURLTemplate = "https://live.881903.com/edge-aac/{{.ChannelName}}/{{.Filename}}"
)

// Station is the template structure for RadioStationPage URL
type Station struct {
	Name string
}

// StreamMedia is the template structure for StreamMediaURL
type StreamMedia struct {
	ChannelName, Filename string
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
func FetchPlaylistLocatorURL(radioStationPageHTML string) (locatorURL string, channelName string, found bool, err error) {
	matched := regexp.MustCompile(PlaylistLocatorURLJSON).FindStringSubmatch(radioStationPageHTML)
	if len(matched) < 2 {
		return
	}
	locatorURL, err = url.PathUnescape(matched[1])
	if err != nil {
		return
	}
	locURI, err := url.Parse(locatorURL)
	if err != nil {
		return
	}

	splitedPath := strings.Split(locURI.Path, "/")
	if len(splitedPath) < 4 {
		return
	}
	channelName = splitedPath[3]

	found = true
	return
}

// PlaylistURL builds the radio channel stream playlist URL
func PlaylistURL(channelName string) (*url.URL, error) {
	tpl, err := template.New("playlisturl").Parse(PlaylistURLTemplate)
	if err != nil {
		return nil, err
	}

	playlistURLStr := new(strings.Builder)
	err = tpl.Execute(playlistURLStr, Station{Name: channelName})
	if err != nil {
		return nil, err
	}
	playlistURL, err := url.Parse(playlistURLStr.String())
	if err != nil {
		return nil, err
	}

	return playlistURL, nil
}

// StreamMediaURL builds the stream media URL
func StreamMediaURL(channelName, filename string) string {
	tpl, err := template.New("mediaurl").Parse(StreamMediaURLTemplate)
	if err != nil {
		panic(err)
	}

	streamMediaURLStr := new(strings.Builder)
	err = tpl.Execute(streamMediaURLStr, StreamMedia{ChannelName: channelName, Filename: filename})
	if err != nil {
		panic(err)
	}

	return streamMediaURLStr.String()
}