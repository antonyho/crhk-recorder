package main

import (
	"flag"
	"time"

	"github.com/antonyho/crhk-recorder/pkg/stream/recorder"
)

func main() {
	var (
		channel   string
		startTime string
		endTime   string
		duration  time.Duration
	)
	flag.StringVar(&channel, "c", "881", "channel name in abbreviation")
	flag.StringVar(&startTime, "s", "", "start time with timezone abbreviation")
	flag.StringVar(&endTime, "e", "", "end time with timezone abbreviation")
	flag.DurationVar(&duration, "d", 0, "record duration (don't do this over 24 hours)")
	flag.Parse()

	if duration == 0 {
		if startTime == "" && endTime == "" {
			panic("record time value must be provided")
		}
	}

	rcdr := recorder.NewRecorder(channel)
	if duration > time.Duration(0) && endTime == "" {
		// Change Schedule() accepts parameter types and create endTime
		// With duration parameter, it will override the endTime
		start, err := time.Parse("15:04:05 -0700", startTime)
		if err != nil {
			panic(err)
		}
		endTime = start.Add(duration).Format("15:04:05 -0700")
	}

	if err := rcdr.Schedule(startTime, endTime); err != nil {
		panic(err)
	}
}
