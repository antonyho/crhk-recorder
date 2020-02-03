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
)

// Station is the template structure for RadioStationPage URL
type Station struct {
	Name string
}

// RadioStationPageURL builds the radion station page URL
func RadioStationPageURL(stationName string) *url.URL {
	tpl, err := template.New("stationurl").Parse(RadioStationPage)
	if err != nil {
		panic(err)
	}

	stationURLStr := new(strings.Builder)
	err = tpl.Execute(stationURLStr, Station{Name: stationName})
	if err != nil {
		panic(err)
	}
	stationURL, err := url.Parse(stationURLStr.String())
	if err != nil {
		panic(err)
	}

	return stationURL
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
