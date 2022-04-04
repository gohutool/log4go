package main

import "log4go"

func main() {
	// Load the configuration (isn't this easy?)
	log4go.LoadConfiguration("example.xml")

	// And now we're ready!
	log4go.Finest("This will only go to those of you really cool UDP kids!  If you change enabled=true.")
	log4go.Debug("Oh no!  %d + %d = %d!", 2, 2, 2+2)
	log4go.Info("About that time, eh chaps?")
}
