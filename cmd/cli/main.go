package main

import (
	"flag"
	"github.com/antonyho/crhk-recorder/pkg/stream/recorder"
)

func main() {
	var (
		channel = flag.String("c", "881", "channel name in abbreviation")
		startTime = flag.String("s", "23:00:00 +0800", "start time with timezone abbreviation")
		endTime = flag.String("e", "23:59:59 +0800", "end time with timezone abbreviation")
	)
	flag.Parse()

	rcdr := recorder.NewRecorder(*channel)
	if err := rcdr.Schedule(*startTime, *endTime); err != nil {
		panic(err)
	}
}
