package url

import (
	"net/url"
	"regexp"
	"strings"
	"text/template"
)

const (
	// RadioStationPage is the URL template for the radio channel
	// StreamName can be 881, 903, 864
	RadioStationPage = "https://www.881903.com/live/{{.Name}}"

	// PlaylistLocatorURLJSON is the pattern of JSON which contains playlist locator URL value
	PlaylistLocatorURLJSON = `"liveJsUrl":"(.*?)"`

	// PlaylistURLTemplate is the URL template for the radio stream playlist
	PlaylistURLTemplate = "https://{{.server.Hostname}}/edge-aac/{{.channel.Name}}/chunks.m3u8"

	// StreamMediaURLTemplate is the URL for the stream media file
	StreamMediaURLTemplate = "https://{{.StreamServer.Hostname}}/edge-aac/{{.ChannelName}}/{{.Filename}}"
)

// StreamServer is the template structure for URLs
type StreamServer struct {
	Hostname string
}

// Channel is the template structure for RadioChannelPage URL
type Channel struct {
	Name string
}

// StreamMedia is the template structure for StreamMediaURL
type StreamMedia struct {
	ChannelName, Filename string
	StreamServer          StreamServer
}

// RadioChannelPageURL builds the radio channel page URL
func RadioChannelPageURL(channel string) string {
	tpl, err := template.New("stationurl").Parse(RadioStationPage)
	if err != nil {
		panic(err)
	}

	channelURLStr := new(strings.Builder)
	err = tpl.Execute(channelURLStr, Channel{Name: channel})
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
func PlaylistURL(channelName, streamServer string) (*url.URL, error) {
	tpl, err := template.New("playlisturl").Parse(PlaylistURLTemplate)
	if err != nil {
		return nil, err
	}

	playlistURLStr := new(strings.Builder)
	values := map[string]interface{}{
		"channel": Channel{Name: channelName},
		"server":  StreamServer{Hostname: streamServer},
	}
	err = tpl.Execute(playlistURLStr, values)
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
func StreamMediaURL(channelName, streamServer, filename string) string {
	tpl, err := template.New("mediaurl").Parse(StreamMediaURLTemplate)
	if err != nil {
		panic(err)
	}

	streamMediaURLStr := new(strings.Builder)
	streamMedia := StreamMedia{
		ChannelName:  channelName,
		Filename:     filename,
		StreamServer: StreamServer{Hostname: streamServer},
	}
	err = tpl.Execute(streamMediaURLStr, streamMedia)
	if err != nil {
		panic(err)
	}

	return streamMediaURLStr.String()
}
