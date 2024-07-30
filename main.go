package main

import (
	"github.com/roy-rishi/midi-to-gcode/file"
	"github.com/roy-rishi/midi-to-gcode/midi"
)

func main() {
	var data []byte = file.Read("documents/starmachine-2k.mid")

	head := midi.ParseHeader(data)
	midi.ParseTracks(data, head.NumTracks)
}
