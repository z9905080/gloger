package gloger

import "testing"

func TestNewDebug(t *testing.T) {
	glog := NewLogger()
	glog.SetCurrentLevel(DEBUG)
	glog.Debug("AAA")
}

func TestDebug(t *testing.T) {
	SetCurrentLevel(DEBUG)
	Debug("AAA")
}