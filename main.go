package main

import (
	"log"

	"github.com/roy-rishi/midi-to-gcode/file"
	"github.com/roy-rishi/midi-to-gcode/midi"
)

func main() {
	var data []byte = file.Read("documents/starmachine-2k.mid")
	// parse header
	head := midi.ParseHeader(data)
	log.Printf("Header: %+v", head)
}
