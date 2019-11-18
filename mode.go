package main

import "fmt"

type RunningMode string

const (
	ModeHTTP RunningMode = "http"
	ModeTCP  RunningMode = "tcp"
)

//UnmarshalFlag Unmarshal flag
func (n *RunningMode) UnmarshalFlag(value string) error {
	switch value {
	case string(ModeHTTP), string(ModeTCP):
		*n = RunningMode(value)
	default:
		return fmt.Errorf("mode can be only '%s' or '%s'", ModeHTTP, ModeTCP)
	}
	return nil
}

//MarshalFlag Marshals flag
func (n RunningMode) MarshalFlag() (string, error) {
	return string(n), nil
}
