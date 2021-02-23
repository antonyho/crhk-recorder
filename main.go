package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/antonyho/crhk-recorder/pkg/dayofweek"

	"github.com/antonyho/crhk-recorder/pkg/stream/recorder"
)

func main() {
	var (
		channel   string
		startTime string
		endTime   string
		duration  time.Duration
		weekdays  string
		repeat    bool
	)
	flag.StringVar(&channel, "c", "881", "channel name in abbreviation")
	flag.StringVar(&startTime, "s", "", "start time with timezone abbreviation")
	flag.StringVar(&endTime, "e", "", "end time with timezone abbreviation")
	flag.DurationVar(&duration, "d", 0, "record duration [don't do this over 24 hours]")
	flag.StringVar(&weekdays, "w", "1,2,3,4,5", "day of week on scheduled recording [comma seperated] [Sunday=0]")
	flag.BoolVar(&repeat, "r", true, "repeat recording at scheduled time on next day")
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

	dowMask := dayofweek.New()
	if weekdays != "" {
		days := strings.Split(weekdays, ",")
		for _, day := range days {
			day = strings.TrimSpace(day)
			d, err := strconv.ParseUint(day, 10, 8)
			if err != nil {
				panic(fmt.Errorf("incorrect day of week parameter [%s]", weekdays))
			} else if d >= 0 && d <= 6 {
				dowMask.Enable(time.Weekday(d))
			} else {
				panic(fmt.Errorf("incorrect day of week parameter [%s]", weekdays))
			}
		}
	}

	if err := rcdr.Schedule(startTime, endTime, *dowMask, repeat); err != nil {
		panic(err)
	}
}
