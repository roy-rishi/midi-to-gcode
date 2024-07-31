package gcode

import (
	"fmt"
	"log"
	"math"

	"github.com/roy-rishi/midi-to-gcode/file"
	"github.com/roy-rishi/midi-to-gcode/midi"
)

// return: pitch of note (Hz)
func decodeMidiNoteNum(midiNum int) float64 {
	// f(n) = 440 * 2 ^ [(n - 69) / 12]
	return 440 * math.Pow(2, float64(midiNum-69)/12)
}

func vTimeToSeconds(vTime int, division int, tempo float64) float64 {
	ticksPerQuarter := division
	microsPerQuarter := tempo
	microsPerTick := microsPerQuarter / float64(ticksPerQuarter)
	secsPerTick := microsPerTick / 1000000
	secs := float64(vTime) * secsPerTick
	return secs
}

func GenGCode(trackChunks []midi.TrackChunk, division int) string {
	res := ""

	var univStartGCode string = string(file.ReadBin("config/start.gcode"))
	res += univStartGCode

	// TODO: don't hard-code machine-specific config path
	var machineStartGCode string = string(file.ReadBin("config/machine-profiles/prusa-mk3(s+).gcode"))
	res += machineStartGCode

	// TODO: don't hard-code which track to play
	trackToPlay := 0
	i := 0
	// TODO: don't hard-code tempo
	var tempo float64
	// var lastPitch float64
	// var lastTime float64
	var curTime float64 = 0
	for i < len(trackChunks[trackToPlay].Events) {
		te := trackChunks[trackToPlay].Events[i]
		log.Printf("TRACK EVENT %+v\n", te)
		curTime += vTimeToSeconds(te.VTime, division, tempo)

		if te.Type == "note_off" || (te.Type == "note_on" && te.Velocity == 0) {
			// note off
			// pitch := lastPitch
			// durationSecs := (curTime - lastTime)
		} else {
			switch te.Type {
			case "note_on":
				// note on
				midiNum := te.Note
				pitch := decodeMidiNoteNum(midiNum)
				log.Printf("note midi %v = %v Hz\n", te.Note, pitch)
				log.Printf("at seconds %v\n", curTime)

				// lastPitch = pitch
				// lastTime = curTime
			case "tempo_change":
				tempo = float64(te.Tempo)
			}
		}
		i++
		fmt.Println()
	}

	return res
}
