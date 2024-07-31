package main

import (
	"log"

	"github.com/roy-rishi/midi-to-gcode/file"
	"github.com/roy-rishi/midi-to-gcode/gcode"
	"github.com/roy-rishi/midi-to-gcode/midi"
)

func main() {
	var data []byte = file.ReadBin("documents/starmachine-2k.mid")

	head := midi.ParseHeader(data)

	rawTracks := midi.ParseTracks(data, head.NumTracks)
	log.Printf("PARSED TRACKS %+v\n", rawTracks)
	gCodeRes := gcode.GenGCode(rawTracks, head.Division)
	log.Printf("GENERATED G-Code\n\n%v\n", gCodeRes)
}
