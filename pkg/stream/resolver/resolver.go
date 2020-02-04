package resolver

import (
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/antonyho/crhk-recorder/pkg/stream/url"
	"github.com/ushis/m3u"
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
)

// Find channel M3U format playlist
func Find(channel string) (m3u.Playlist, error) {
	channelPageURL := url.RadioChannelPageURL(channel)

	resp, err := http.Get(channelPageURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	channelPageBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	playlistLocatorURL, success := url.FetchPlaylistLocatorURL(string(channelPageBody))
	if !success {
		return nil, errors.New("playlist URL not found")
	}

	playlistAuthURL, err := GetPlaylistAuthenticationURL(playlistLocatorURL)
	if err != nil {
		return nil, err
	}

	policy, keypair, sig, err := GetPlaylistAuthentication(playlistAuthURL)
	if err != nil {
		return nil, err
	}

	return GetPlaylist(channel, policy, keypair, sig)
}

// GetPlaylistAuthenticationURL finds the authentication URL for playlist
func GetPlaylistAuthenticationURL(locatorURL string) (string, error) {
	resp, err := http.Get(locatorURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	authenticatorLocation := resp.Header.Get(PlaylistLocationHeaderName)
	if authenticatorLocation == "" {
		return "", errors.New("playlist authenticator location URL header not found")
	}

	return authenticatorLocation, nil
}

// GetPlaylistAuthentication gets the playlist access authentication cookies
func GetPlaylistAuthentication(authURL string) (policy, keypair, sig string, err error) {
	resp, err := http.Get(authURL)
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
func GetPlaylist(channel, policy, keypair, sig string) (m3u.Playlist, error) {
	playlistURL := url.PlaylistURL(channel)

	req, err := http.NewRequest(http.MethodGet, playlistURL, nil)
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
