package dayofweek_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	dow "github.com/antonyho/crhk-recorder/pkg/dayofweek"
)

func TestEnable(t *testing.T) {
	m := dow.New()
	dow.Enable(3, m)
	dow.Enable(int(time.Saturday), m)
	assert.EqualValues(t, 0b01001000, *m)
}

func TestEnabled(t *testing.T) {
	m := dow.New()
	dow.Enable(int(time.Thursday), m)
	assert.True(t, dow.Enabled(time.Thursday, *m))
}

func TestMask_Enable(t *testing.T) {
	m := dow.New()
	m.Enable(time.Thursday)
	assert.EqualValues(t, 0b00010000, *m)
}

func TestMask_EnableAll(t *testing.T) {
	m := dow.New()
	m.EnableAll()
	assert.EqualValues(t, 0b01111111, *m)
}

func TestMask_Enabled(t *testing.T) {
	m := dow.New()
	m.Enable(time.Thursday)
	assert.True(t, dow.Enabled(time.Thursday, *m))
}

func TestMask_AllEnabled(t *testing.T) {
	m := dow.New()
	assert.True(t, m.AllEnabled())

	m.Enable(time.Monday)
	assert.False(t, m.AllEnabled())

	for d := time.Sunday; d <= time.Saturday; d++ {
		m.Enable(d)
	}
	assert.True(t, m.AllEnabled())
}
