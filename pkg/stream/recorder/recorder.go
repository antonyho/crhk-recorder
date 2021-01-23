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

	"github.com/antonyho/crhk-recorder/pkg/stream/url"

	"github.com/antonyho/crhk-recorder/pkg/stream/resolver"
)

// Recorder CRHK radio channel broadcasted online
type Recorder struct {
	Channel                 string
	cloudfrontSessionCookie *resolver.CloudfrontCookie
	downloaded              map[string]bool
	targetFile              io.Writer
}

// NewRecorder is a constructor for Recorder
func NewRecorder(channel string) *Recorder {
	return &Recorder{
		Channel:    channel,
		downloaded: make(map[string]bool),
	}
}

// Download the media from channel playlist
func (r Recorder) Download(targetFile io.Writer) error {
	channelName, playlist, cloudfrontCookie, err := resolver.Find(r.Channel)
	if err != nil {
		return err
	}
	if r.cloudfrontSessionCookie == nil {
		r.cloudfrontSessionCookie = &cloudfrontCookie
	}
	for _, track := range playlist {
		if downloaded, found := r.downloaded[track.Path]; found && downloaded {
			// Skip if the same track has been downloaded
			continue
		}

		// Add CloudFront headers to the request
		req, err := http.NewRequest(http.MethodGet, url.StreamMediaURL(channelName, track.Path), nil)
		if err != nil {
			return err
		}
		req.AddCookie(&http.Cookie{Name: resolver.CloudFrontCookieNamePolicy, Value: cloudfrontCookie.Policy})
		req.AddCookie(&http.Cookie{Name: resolver.CloudFrontCookieNameKeyPairID, Value: cloudfrontCookie.KeyPairID})
		req.AddCookie(&http.Cookie{Name: resolver.CloudFrontCookieNameSignature, Value: cloudfrontCookie.Signature})

		c := &http.Client{}
		resp, err := c.Do(req)
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("unsuccessful HTTP request. response code: %d", resp.StatusCode)
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
	}

	return nil
}

// Record the given channel
func (r Recorder) Record(startFrom, until time.Time) error {
	if startFrom.After(until) {
		panic("incorrect time sequence")
	}

	currExecDirPath, err := os.Getwd()
	if err != nil {
		return err
	}
	mediaFilename := fmt.Sprintf("%s-%s.aac", r.Channel, startFrom.Format("2006-01-02"))
	fileDestPath := filepath.Join(currExecDirPath, mediaFilename)
	log.Println(fileDestPath)
	f, err := os.Create(fileDestPath)
	if err != nil {
		return err
	}
	defer f.Close()
	bufFile := bufio.NewWriter(f)
	defer bufFile.Flush()
	r.targetFile = bufFile

	diffFromStartTime := time.Until(startFrom)
	diffFromEndTime := time.Until(until)

	termination := make(chan bool)
	go func() {
		<-time.After(diffFromEndTime)
		termination <- true
	}()

	<-time.After(diffFromStartTime)
	for {
		select {
		case <-termination:
			return nil

		default:
			if err := r.Download(r.targetFile); err != nil {
				return err
			}
			if err := bufFile.Flush(); err != nil {
				return err
			}
		}
	}
}

// Schedule a time to start and end recording everyday
// startTime format: 13:23:45 +0100 (24H with timezone offset)
// startTime format: 18:54:21 +0100 (24H with timezone offset)
func (r Recorder) Schedule(startTime, endTime string) error {
	var timeDelay time.Duration
	thisYear, thisMonth, thisDay := time.Now().Date()
	start, err := time.Parse("15:04:05 -0700", startTime)
	if err != nil {
		return err
	}
	start = start.AddDate(thisYear, int(thisMonth)-1, thisDay-1)
	if start.Before(time.Now()) {
		timeDelay = 24 * time.Hour
		start = start.Add(timeDelay)
	}
	end, err := time.Parse("15:04:05 -0700", endTime)
	if err != nil {
		return err
	}
	end = end.AddDate(thisYear, int(thisMonth)-1, thisDay-1)
	if end.Before(start) {
		end = end.Add(24 * time.Hour)
	}
	if timeDelay > 0 {
		end = end.Add(timeDelay)
	}

	for { // When it starts, it never ends...
		log.Printf("The next recording schedule: %s - %s", start.Format("2006-01-02 15:04:05 -0700"), end.Format("2006-01-02 15:04:05 -0700"))
		if time.Until(start) > time.Minute {
			// Wait a bit if the start time to more than 1 minute apart
			<-time.After(time.Until(start.Add(-10 * time.Second)))
		}
		if err := r.Record(start, end); err != nil {
			return err
		}
		start = start.Add(24 * time.Hour)
		end = end.Add(24 * time.Hour)
	}
}
