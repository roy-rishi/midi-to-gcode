package main

import (
	"log"

	"github.com/roy-rishi/midi-to-gcode/file"

	"github.com/roy-rishi/midi-to-gcode/converter"
)

func main() {
	var midi []byte = file.ReadBin("../inputs/starmachine-2k.mid")

	headerValid, err := converter.ValidateHeader(midi)
	if !headerValid {
		log.Fatal(err)
	}
	numTracks, timeDivision := converter.ParseHeader(midi)
	converter.ParseNoteEvents(midi, numTracks, timeDivision)
}
