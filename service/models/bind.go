package models

import "time"

type Bind struct {
	Timeout time.Duration
	Path    string
	Method  string
	Topic   string
}
