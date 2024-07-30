package midi

import (
	"fmt"
	"log"
)

// Input: Byte array starting with a track and ending with either end of
// track or end of file
func verifyTrackStart(tBin []byte) bool {
	mtrk := tBin[:4]
	return string(mtrk) == "MTrk"
}

func trackLen(tBin []byte) int {
	if !verifyTrackStart(tBin) {
		log.Fatal("Invalid start of track chunk")
	}
	return int(tBin[4])<<24 | int(tBin[5])<<16 | int(tBin[6])<<8 | int(tBin[7])
}

// Parse all track chunks and all of their track events and meta events.
// Input: full MIDI file contents
func ParseTracks(data []byte, nTracks int) []TrackChunk {
	var tracks []TrackChunk
	startPos := 14
	for n := 0; n < nTracks; n++ {
		log.Println("TRACK")
		tc := TrackChunk{}
		len := trackLen(data[startPos:]) // length
		tc.Length = len
		nextPos := startPos + 8 + len // next chunk starts after <MTrk>, <length>, and []<track_event>...

		// running status
		lastControlType := 0b0

		pos := startPos + 8
		for pos < nextPos {
			// variable length <v_time>
			te := TrackEvent{}
			vTime, p := variableLen(data, pos)
			pos = p
			te.VTime = vTime
			fmt.Println()
			log.Printf("v_time %v\n", vTime)
			// <midi_event>
			msgStatusNum := int(data[pos])

			if msgStatusNum>>7 == 0b1 { // control byte
				switch msgStatusNum {
				case 255: // meta message
					log.Println("META MESSAGE")
					pos++
					metaType := data[pos] // meta msg type
					pos++
					len, p := variableLen(data, pos) // length, data start pos
					pos = p
					switch metaType {
					case 3: // track name
						log.Println("TRACK NAME")
						trackName := string(data[pos : pos+len])
						log.Printf("name %v\n", trackName)
						tc.Name = trackName
					case 47: // undefined
						log.Println("UNDEFINED 47 (ignore)")
					case 88: // time signature
						log.Println("TIME SIGNATURE (ignore)")
					case 89: // key signature
						log.Println("KEY SIGNATURE (ignore)")
					// number of sharps/flats is 255, not within -7 to 7 range
					case 81: // tempo change
						log.Println("TEMPO CHANGE")
						value := int(data[pos])<<16 | int(data[pos+1])<<8 | int(data[pos+2])
						bpm := 60000000 / value
						log.Printf("bpm %v (value %v)\n", bpm, value)
						te.Tempo = bpm
						te.Type = "tempo_change"
						tc.Events = append(tc.Events, te)
					default:
						log.Printf("Unknown meta message type %v", metaType)
					}
					log.Printf("%+v\n", te)
					pos += len
				default:
					// TODO: switch to hex; eg: 0x9X
					highOrderBits := msgStatusNum >> 4 // 4 highest order bits
					lastControlType = highOrderBits    // allow running status
					switch highOrderBits {
					case 0b1011: // control change
						log.Println("CONTROL CHANGE (ignore)")
						pos += 3 // skip two bytes
					case 0b1000:
						log.Println("NOTE OFF")
						pos++
						noteNum := data[pos]
						pos++
						velocity := data[pos]

						te.Type = "note_off"
						te.Note = int(noteNum)
						te.Velocity = int(velocity)
						tc.Events = append(tc.Events, te)
						log.Printf("%+v", te)
						pos++
					case 0b1001:
						log.Println("NOTE ON")
						pos++
						noteNum := data[pos]
						pos++
						velocity := data[pos]

						te.Type = "note_on"
						te.Note = int(noteNum)
						te.Velocity = int(velocity)
						tc.Events = append(tc.Events, te)
						log.Printf("%+v", te)
						pos++
					case 0b1100:
						log.Println("PROGRAM CHANGE (ignore)")
						pos += 2 // skip 1 byte
					default:
						log.Fatalf("Unknown message status code %v", msgStatusNum)
					}
				}
			} else { // use running status
				switch lastControlType {
				case 0b1011:
					log.Println("CONTROL CHANGE (ignore)")
					pos += 3 - 1 // skip 1 byte
				case 0b1100:
					log.Println("PROGRAM CHANGE (ignore)")
					pos += 1 // skip byte
				case 0b1000:
					log.Println("NOTE OFF")
					// pos++
					noteNum := data[pos]
					pos++
					velocity := data[pos]

					te.Type = "note_off"
					te.Note = int(noteNum)
					te.Velocity = int(velocity)
					tc.Events = append(tc.Events, te)
					log.Printf("%+v", te)
					pos++
				case 0b1001:
					log.Println("NOTE ON")
					// pos++
					noteNum := data[pos]
					pos++
					velocity := data[pos]

					te.Type = "note_on"
					te.Note = int(noteNum)
					te.Velocity = int(velocity)
					tc.Events = append(tc.Events, te)
					log.Printf("%+v", te)
					pos++
				default:
					log.Fatalf("Unable to use running status on data byte %v\n", msgStatusNum)
				}
			}
		}
		startPos = nextPos
		log.Printf("%+v\n\n", tc)
	}
	return tracks
}
