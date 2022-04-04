package main

import (
	"log4go"
	"testing"
	"time"
)

func TestXMLConfiguration(t *testing.T) {
	// Load the configuration (isn't this easy?)
	log4go.LoadConfiguration("example.xml")
	
	// And now we're ready!
	log4go.Finest("This will only go to those of you really cool UDP kids!  If you change enabled=true.")
	log4go.Debug("Oh no!  %d + %d = %d!", 2, 2, 2+2)
	log4go.Info("1. About that time, eh chaps?")
	log4go.Info("2. About that time, eh chaps?")
	log4go.Info("3. About that time, eh chaps?")

	time.Sleep(1 * time.Second)

	log4go.Info("1. About that time, eh chaps?")
	log4go.Info("2. About that time, eh chaps?")
	log4go.Info("3. About that time, eh chaps?")

	time.Sleep(10 * time.Second)

}
