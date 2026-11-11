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
		{duration: time.Second * 2000, res: "00:33:20.0"},
		{duration: time.Second * 3680, res: "01:01:20.0"},
		{duration: time.Second * 3680, res: "01:01:20.0"},
	}

	for _, c := range cases {
		if actual := humanizeDuration(c.duration); actual != c.res {
			t.Errorf("test failed. expected %s, got %s", actual, c.res)
		}
	}
}
