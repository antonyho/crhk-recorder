package recorder_test

import (
	"fmt"
	"os"
	"path"
	"testing"
	"time"

	"github.com/antonyho/crhk-recorder/pkg/dayofweek"
	"github.com/antonyho/crhk-recorder/pkg/stream/recorder"
)

const (
	channel  = "881"
	filename = "881hd.aac"
)

func TestRecorder_Download(t *testing.T) {
	tmpDirPath := t.TempDir()
	fileDest := path.Join(tmpDirPath, filename)
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
	tmpDirPath := t.TempDir()
	if err := os.MkdirAll(tmpDirPath, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(tmpDirPath); err != nil {
		t.Fatal(err)
	}

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

func TestRecorder_Schedule_once(t *testing.T) {
	tmpDirPath := t.TempDir()
	if err := os.MkdirAll(tmpDirPath, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(tmpDirPath); err != nil {
		t.Fatal(err)
	}

	rcdr := recorder.NewRecorder(channel)
	tf := "15:04:05 -0700"
	now := time.Now()
	startTime := now.Add(5 * time.Second).Format(tf)
	endTime := now.Add(30 * time.Second).Format(tf)
	dowMask := dayofweek.New()
	t.Logf("Start time: %v | End time: %v", startTime, endTime)
	if err := rcdr.Schedule(startTime, endTime, *dowMask, false); err != nil {
		t.Error(err)
	}
}

func TestRecorder_Schedule_endless(t *testing.T) {
	tmpDirPath := t.TempDir()
	if err := os.MkdirAll(tmpDirPath, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(tmpDirPath); err != nil {
		t.Fatal(err)
	}

	terminate := make(chan bool)
	performTest := func() {
		rcdr := recorder.NewRecorder(channel)
		tf := "15:04:05 -0700"
		now := time.Now()
		startTime := now.Add(5 * time.Second).Format(tf)
		endTime := now.Add(30 * time.Second).Format(tf)
		dowMask := dayofweek.New()
		t.Logf("Start time: %v | End time: %v", startTime, endTime)
		if err := rcdr.Schedule(startTime, endTime, *dowMask, true); err != nil {
			t.Error(err)
		}
		terminate <- true
	}
	timeout := time.After(40 * time.Second)
	go performTest()

	select {
	case <-timeout:
		t.Log("Breaking the endless schedule")
	case <-terminate:
		t.Log("Endless schedule shall not terminate")
		t.FailNow()
	}
}

func TestRecorder_Schedule_DayOfWeek(t *testing.T) {
	tmpDirPath := t.TempDir()
	if err := os.MkdirAll(tmpDirPath, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(tmpDirPath); err != nil {
		t.Fatal(err)
	}

	terminate := make(chan bool)
	performTest := func() {
		rcdr := recorder.NewRecorder(channel)
		tf := "15:04:05 -0700"
		now := time.Date(2021, time.January, 24, 15, 04, 05, 0, time.UTC) // Sunday
		startTime := now.Add(5 * time.Second).Format(tf)
		endTime := now.Add(30 * time.Second).Format(tf)
		dowMask := dayofweek.New()
		dowMask.Enable(time.Tuesday)
		t.Logf("Start time: %v | End time: %v", startTime, endTime)
		// The next recording schedule: 2021-01-26 15:04:10 +0000 - 2021-01-26 15:04:35 +0000
		if err := rcdr.Schedule(startTime, endTime, *dowMask, true); err != nil {
			t.Error(err)
		}
		terminate <- true
	}
	timeout := time.After(40 * time.Second)
	go performTest()

	select {
	case <-timeout:
		t.Log("Breaking the endless schedule")
	case <-terminate:
		t.Log("Endless schedule shall not terminate")
		t.FailNow()
	}
}
