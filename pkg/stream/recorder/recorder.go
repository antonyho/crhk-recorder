package downloader

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/antonyho/crhk-recorder/pkg/stream/url"

	"github.com/antonyho/crhk-recorder/pkg/stream/resolver"
)

// Record the given channel
func Record(channel string, startFrom time.Time, until time.Time) {
	// TODO - Skip download if the same playlist has been downloaded
}

func Download(channel string, fileDest string, termination <-chan bool) error {
	for {
		select {
		case <- termination:
			return nil

		default:
			//f, err := os.Create(fileDest)
			//if err != nil {
			//	return err
			//}
			//defer f.Close()
			//bufFile := bufio.NewWriter(f)
			//defer bufFile.Flush()
			for {
				channelName, playlist, err := resolver.Find(channel)
				if err != nil {
					return err
				}
				for _, track := range playlist {
					resp, err := http.Get(url.StreamMediaURL(channelName, track.Path))
					if err != nil {
						return err
					}
					media, err := ioutil.ReadAll(resp.Body)
					resp.Body.Close()
					if err != nil {
						return err
					}
					if len(media) == 0 {
						return errors.New("empty media file")
					}

					if err := ioutil.WriteFile(fmt.Sprintf("%s/%s", fileDest, track.Path), media, 0644); err != nil {
						return err
					}
				}
			}
		}
	}

	return errors.New("unexpected termination")
}