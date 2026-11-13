package main

import (
	"testing"
	"time"
)

func TestHumanize(t *testing.T) {
	type tt struct {
		duration time.Duration
		res      string
	}

	cases := []tt{
		{duration: 0, res: "00:00:00.0"},
		{duration: time.Millisecond * 5300, res: "00:00:05.3"},
		{duration: time.Millisecond * 59960, res: "00:00:59.9"},
		{duration: time.Second * 2000, res: "00:33:20.0"},
		{duration: time.Second * 3680, res: "01:01:20.0"},
		{duration: time.Millisecond * 3725990, res: "01:02:05.9"},
	}

	for _, c := range cases {
		if actual := humanizeDuration(c.duration); actual != c.res {
			t.Errorf("humanizeDuration(%v) = %s, want %s", c.duration, actual, c.res)
		}
	}
}
