package resolver

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ushis/m3u"

	crhk "github.com/antonyho/crhk-recorder/pkg/stream/url"
)

const (
	// PlaylistLocationHeaderName is the header name for playlist location
	PlaylistLocationHeaderName = "location"
)

const (
	// UserAgentCamouflage disguises our HTTP client as a common browser agent
	UserAgentCamouflage = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:71.0) Gecko/20100101 Firefox/71.0"
)

// Find channel M3U format playlist
func Find(channel string) (
	channelName string,
	livestreamServer string,
	cloudfrontCookie CloudfrontCookie,
	err error,
) {
	playlistCloudFrontURL, channelName, channelPageURL, err := GetCloudFrontResolverURL(channel)
	if err != nil {
		return
	}

	cloudfrontCookie, livestreamServer, err = GetPlaylistAuthentication(channelPageURL, playlistCloudFrontURL)
	if err != nil {
		return
	}

	return
}

// GetCloudFrontResolverURL finds the URL to visit in order to get CloudFront cookies
func GetCloudFrontResolverURL(channel string) (string, string, string, error) {
	channelPageURL := crhk.RadioChannelPageURL(channel)

	resp, err := http.Get(channelPageURL)
	if err != nil {
		return "", "", "", err
	}
	defer resp.Body.Close()

	channelPageBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", "", err
	}
	playlistCFURL, channelName, success, err := crhk.FetchPlaylistLocatorURL(string(channelPageBody))
	if err != nil {
		return "", "", "", err
	} else if !success {
		return "", "", "", errors.New("playlist URL not found")
	}

	return playlistCFURL, channelName, channelPageURL, nil
}

// GetPlaylistAuthentication gets the playlist access authentication cookies
// It accesses the playlistCloudFrontURL, then it will get 304 redirected to
// a new location which contains the CloudFront policy and key pair value in
// response headers.
func GetPlaylistAuthentication(
	refererURL, playlistCloudFrontURL string,
) (
	cloudfrontCookie CloudfrontCookie, livestreamServerHostname string, err error,
) {
	req, err := http.NewRequest(http.MethodGet, playlistCloudFrontURL, nil)
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", UserAgentCamouflage)
	req.Header.Set("Referer", refererURL)

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("Getting CloudFront cookies HTTP request failed. URL: %s \nResponse code: %d", playlistCloudFrontURL, resp.StatusCode)
		err = errors.New("http request getting CloudFront cookies failed")
		return
	}

	for _, c := range resp.Cookies() {
		switch c.Name {
		case CloudFrontCookieNamePolicy:
			cloudfrontCookie.Policy = c.Value

		case CloudFrontCookieNameKeyPairID:
			cloudfrontCookie.KeyPairID = c.Value

		case CloudFrontCookieNameSignature:
			cloudfrontCookie.Signature = c.Value
		}
	}

	livestreamServerHostname = resp.Request.URL.Host

	return
}

// GetPlaylist gets the playlist using the given authentication cookie values
func GetPlaylist(
	channelName string,
	streamServer string,
	cloudfrontCookie CloudfrontCookie,
) (m3u.Playlist, error) {
	playlistURL, err := crhk.PlaylistURL(channelName, streamServer)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, playlistURL.String(), nil)
	if err != nil {
		return nil, err
	}
	req.AddCookie(&http.Cookie{Name: CloudFrontCookieNamePolicy, Value: cloudfrontCookie.Policy})
	req.AddCookie(&http.Cookie{Name: CloudFrontCookieNameKeyPairID, Value: cloudfrontCookie.KeyPairID})
	req.AddCookie(&http.Cookie{Name: CloudFrontCookieNameSignature, Value: cloudfrontCookie.Signature})

	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("Playlist HTTP request failed. URL: %s \nResponse code: %d", playlistURL, resp.StatusCode)
		return nil, fmt.Errorf("playlist fetching failed")
	}

	return m3u.Parse(resp.Body)
}
