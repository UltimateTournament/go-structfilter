package main

import (
	"bytes"
	"strings"
	"time"
)

// Duration is time.Duration that uses its string representation in JSON
type Duration time.Duration

func (d *Duration) UnmarshalJSON(b []byte) error {
	if len(b) == 0 || bytes.Equal(b, []byte("null")) {
		return nil
	}
	s := strings.Trim(string(b), "\"")
	dur, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	*d = Duration(dur)
	return nil
}

func (d *Duration) MarshalJSON() ([]byte, error) {
	if d == nil {
		return []byte("null"), nil
	}
	durStr := time.Duration(*d).String()
	jsnStr := "\"" + durStr + "\""
	return []byte(jsnStr), nil
}
