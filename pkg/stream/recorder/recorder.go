package recorder

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	dow "github.com/antonyho/crhk-recorder/pkg/dayofweek"
	"github.com/antonyho/crhk-recorder/pkg/stream/resolver"
	"github.com/antonyho/crhk-recorder/pkg/stream/url"
)

const (
	// ConsecutiveErrorTolerance is the number of failures on downloading a stream
	// which shall be tolerated
	ConsecutiveErrorTolerance = 15

	// OneDay time value
	OneDay = 24 * time.Hour

	// TwoSeconds time value
	TwoSeconds = 2 * time.Second
)

// Recorder CRHK radio channel broadcasted online
type Recorder struct {
	Channel                 string
	ChannelName             string // specifies with stream sound quality (e.g. 881HD)
	StreamServer            string
	cloudfrontSessionCookie *resolver.CloudfrontCookie
	downloaded              map[string]bool
}

// NewRecorder is a constructor for Recorder
func NewRecorder(channel string) *Recorder {
	return &Recorder{
		Channel:    channel,
		downloaded: make(map[string]bool),
	}
}

func (r *Recorder) cleanup() {
	r.downloaded = make(map[string]bool)
}

// Download the media from channel playlist
func (r *Recorder) Download(targetFile io.Writer) error {
	if r.ChannelName == "" ||
		r.StreamServer == "" ||
		r.cloudfrontSessionCookie == nil ||
		!r.cloudfrontSessionCookie.Assigned() {
		channelName, streamServer, cloudfrontCookie, err := resolver.Find(r.Channel)
		if err != nil {
			return err
		}
		r.ChannelName = channelName
		r.StreamServer = streamServer
		r.cloudfrontSessionCookie = &cloudfrontCookie
	}

	playlist, err := resolver.GetPlaylist(r.ChannelName, r.StreamServer, *r.cloudfrontSessionCookie)
	if err != nil {
		return err
	}

	var lastTrackDuration time.Duration
	playlistDownloadStartTime := time.Now()
	for _, track := range playlist {
		if downloaded, found := r.downloaded[track.Path]; found && downloaded {
			// Skip if the same track has been downloaded
			continue
		}

		// Add CloudFront headers to the request
		req, err := http.NewRequest(http.MethodGet, url.StreamMediaURL(r.ChannelName, r.StreamServer, track.Path), nil)
		if err != nil {
			return err
		}
		req.AddCookie(&http.Cookie{Name: resolver.CloudFrontCookieNamePolicy, Value: r.cloudfrontSessionCookie.Policy})
		req.AddCookie(&http.Cookie{Name: resolver.CloudFrontCookieNameKeyPairID, Value: r.cloudfrontSessionCookie.KeyPairID})
		req.AddCookie(&http.Cookie{Name: resolver.CloudFrontCookieNameSignature, Value: r.cloudfrontSessionCookie.Signature})

		c := &http.Client{}
		resp, err := c.Do(req)
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("media file: unsuccessful HTTP request. response code: %d", resp.StatusCode)
		}
		media, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		resp.Body.Close()
		contentSize := len(media)
		if contentSize == 0 {
			return errors.New("empty media file")
		}

		written, err := targetFile.Write(media)
		if err != nil {
			return err
		} else if written != contentSize {
			return fmt.Errorf("written byte size %d does not match with file size %d", written, contentSize)
		}
		r.downloaded[track.Path] = true
		lastTrackDuration = time.Duration(track.Time) * time.Second
	}

	if time.Since(playlistDownloadStartTime) < lastTrackDuration {
		time.Sleep(lastTrackDuration - TwoSeconds) // Wait 2 seconds less to be secure
	}

	return nil
}

// Record the given channel
func (r *Recorder) Record(startFrom, until time.Time) error {
	if startFrom.After(until) {
		panic("incorrect time sequence")
	}

	currExecDirPath, err := os.Getwd()
	if err != nil {
		return err
	}
	mediaFilename := fmt.Sprintf("%s-%s.aac", r.Channel, startFrom.Format("2006-01-02-150405"))
	fileDestPath := filepath.Join(currExecDirPath, mediaFilename)
	f, err := os.Create(fileDestPath)
	if err != nil {
		return err
	}
	defer f.Close()
	bufFile := bufio.NewWriter(f)
	defer bufFile.Flush()

	diffFromStartTime := time.Until(startFrom)
	diffFromEndTime := time.Until(until)

	termination := make(chan bool)
	go func() {
		<-time.After(diffFromEndTime)
		termination <- true
	}()

	failCount := 0

	<-time.After(diffFromStartTime)
	for {
		select {
		case <-termination:
			r.cleanup()
			return nil

		default:
			if err := r.Download(bufFile); err != nil {
				if failCount < ConsecutiveErrorTolerance {
					log.Printf("Download Error: %+v", err)
					time.Sleep(calculateRetryDelay(failCount))
					failCount++
					continue
				} else {
					return err
				}
			}
			if err := bufFile.Flush(); err != nil {
				return err
			}
			failCount = 0
		}
	}
}

// Schedule a time to start and end recording everyday
// wd is a flag mask to control which day of week should be recorded
// endless controls if the schedule would continue endlessly on next scheduled day
// startTime format: 13:23:45 +0100 (24H with timezone offset)
func (r *Recorder) Schedule(startTime, endTime string, wd dow.Bitmask, endless bool) error {
	var timeDelay time.Duration
	thisYear, thisMonth, thisDay := time.Now().Date()
	start, err := time.Parse("15:04:05 -0700", startTime)
	if err != nil {
		return err
	}
	start = start.AddDate(thisYear, int(thisMonth)-1, thisDay-1)
	end, err := time.Parse("15:04:05 -0700", endTime)
	if err != nil {
		return err
	}
	end = end.AddDate(thisYear, int(thisMonth)-1, thisDay-1)
	if end.Before(start) { // To cover an overnight recording
		end = end.Add(OneDay)
	}
	for !wd.AllEnabled() && !wd.Enabled(start.Weekday()) {
		start = start.Add(OneDay)
		end = end.Add(OneDay)
	}
	if start.Before(time.Now()) { // To cover start time already passed
		timeDelay = OneDay
		start = start.Add(timeDelay)
		end = end.Add(timeDelay)
	}

	for {
		log.Printf("The next recording schedule: %s - %s", start.Format("2006-01-02 15:04:05 -0700"), end.Format("2006-01-02 15:04:05 -0700"))
		if time.Until(start) > time.Minute {
			// Wait a bit if the start time to more than 1 minute apart
			<-time.After(time.Until(start.Add(-10 * time.Second)))
		}
		if err := r.Record(start, end); err != nil {
			return err
		}
		if endless {
			start = start.Add(OneDay)
			end = end.Add(OneDay)
			for !wd.AllEnabled() && !wd.Enabled(start.Weekday()) {
				start = start.Add(OneDay)
				end = end.Add(OneDay)
			}
		} else {
			break
		}
	}

	return nil
}

func calculateRetryDelay(count int) time.Duration {
	if count > 0 {
		return time.Duration(count) * time.Second
	}
	return 0
}
