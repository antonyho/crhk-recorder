package recorder_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/antonyho/crhk-recorder/pkg/stream/recorder"
)

const (
	channel  = "881"
	fileDest = "/tmp/881hd.aac"
)

func TestRecorder_Download(t *testing.T) {
	testFile, err := os.Create(fileDest)
	defer testFile.Close()
	if err != nil {
		t.Fatal(err)
	}

	rcdr := recorder.NewRecorder(channel)
	if err := rcdr.Download(testFile); err != nil {
		t.Fatal(err)
	}
	if err := testFile.Sync(); err != nil {
		t.Fatal(err)
	}
}

func TestRecorder_Record(t *testing.T) {
	rcdr := recorder.NewRecorder(channel)
	term := make(chan struct{})
	go func() {
		start := time.Now()
		for {
			select {
			case <-term:
				fmt.Println()
				return
			default:
				fmt.Printf("\rElapsed: %v", time.Since(start))
			}
		}
	}()
	if err := rcdr.Record(time.Now().Add(2*time.Second), time.Now().Add(10*time.Second)); err != nil {
		t.Error(err)
	}
	term <- struct{}{}
}

func TestRecorder_Schedule(t *testing.T) {
	terminate := make(chan bool)
	performTest := func() {
		rcdr := recorder.NewRecorder(channel)
		tf := "15:04:05 -0700"
		now := time.Now()
		startTime := now.Add(5 * time.Second).Format(tf)
		endTime := now.Add(30 * time.Second).Format(tf)
		t.Logf("Start time: %v | End time: %v", startTime, endTime)
		if err := rcdr.Schedule(startTime, endTime); err != nil {
			t.Error(err)
		}
		terminate <- true
	}
	timeout := time.After(40 * time.Second)
	go performTest()

	select {
	case <-timeout:
		t.Fatal("Test didn't finish in time")
	case <-terminate:
	}
}
