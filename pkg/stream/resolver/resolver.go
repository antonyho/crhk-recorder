package resolver

import (
	"errors"
	"github.com/ushis/m3u"
	"io/ioutil"
	"net/http"

	"github.com/antonyho/crhk-recorder/pkg/stream/url"
)

const (
	// PlaylistLocationHeaderName is the header name for playlist location
	PlaylistLocationHeaderName = "location"
)

const (
	// CloudFrontCookieNamePolicy is the cookie name for CloudFront policy
	CloudFrontCookieNamePolicy = "CloudFront-Policy"

	// CloudFrontCookieNameKeyPairID is the cookie name for CloudFront key pair ID
	CloudFrontCookieNameKeyPairID = "CloudFront-Key-Pair-Id"

	// CloudFrontCookieNameSignature is the cookie name for CloudFront signature
	CloudFrontCookieNameSignature = "CloudFront-Signature"

	// UserAgentCamouflage disguises our HTTP client as a common browser agent
	UserAgentCamouflage = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:71.0) Gecko/20100101 Firefox/71.0"
)

// Find channel M3U format playlist
func Find(channel string) (channelName string, playlist m3u.Playlist, err error) {
	playlistLocatorURL, channelName, err := GetPlaylistLocatorPageURL(channel)
	if err != nil {
		return
	}

	policy, keypair, sig, err := GetPlaylistAuthentication(playlistLocatorURL)
	if err != nil {
		return
	}

	playlist, err = GetPlaylist(channelName, policy, keypair, sig)

	return
}

func GetPlaylistLocatorPageURL(channel string) (string, string, error) {
	channelPageURL := url.RadioChannelPageURL(channel)

	resp, err := http.Get(channelPageURL)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	channelPageBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	playlistLocatorURL, channelName, success, err := url.FetchPlaylistLocatorURL(string(channelPageBody))
	if err != nil {
		return "", "", err
	} else if !success {
		return "", "", errors.New("playlist URL not found")
	}

	return playlistLocatorURL, channelName, nil
}

// GetPlaylistAuthentication gets the playlist access authentication cookies
func GetPlaylistAuthentication(locatorURL string) (policy, keypair, sig string, err error) {
	req, err := http.NewRequest(http.MethodGet, locatorURL, nil)
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", UserAgentCamouflage)
	req.Header.Set("Referer", "https://www.881903.com/live/881")

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	for _, c := range resp.Cookies() {
		switch c.Name {
		case CloudFrontCookieNamePolicy:
			policy = c.Value

		case CloudFrontCookieNameKeyPairID:
			keypair = c.Value

		case CloudFrontCookieNameSignature:
			sig = c.Value
		}
	}

	return
}

// GetPlaylist gets the playlist using the given authentication cookie values
func GetPlaylist(channelName, policy, keypair, sig string) (m3u.Playlist, error) {
	playlistURL, err := url.PlaylistURL(channelName)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, playlistURL.String(), nil)
	if err != nil {
		return nil, err
	}
	req.AddCookie(&http.Cookie{Name: CloudFrontCookieNamePolicy, Value: policy})
	req.AddCookie(&http.Cookie{Name: CloudFrontCookieNameKeyPairID, Value: keypair})
	req.AddCookie(&http.Cookie{Name: CloudFrontCookieNameSignature, Value: sig})

	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return m3u.Parse(resp.Body)
}
