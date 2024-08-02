package converter

import (
	"encoding/binary"
	"log"
)

// Parse the MIDI header for file metadata.
// Params: MIDI file
// Return: number of tracks, time division
func ParseHeader(midi []byte) (int, int) {
	numTracks := int(midi[10])<<8 | int(midi[11])
	timeDivision := int(midi[12])<<8 | int(midi[13])
	return numTracks, timeDivision
}

// Validate the header (start) of MIDI file for compatibility
// Params: MIDI file
// Return: if MIDI is valid; an error string if MIDI is not valid
func ValidateHeader(midi []byte) (bool, string) {
	first4Chars := string(midi[:4])
	if first4Chars != "MThd" {
		return false, "File not recognized as MIDI"
	}

	headerLen := int(midi[4])<<24 | int(midi[5])<<16 | int(midi[6])<<8 | int(midi[7])
	if headerLen != 6 {
		return false, "Header must be 6 bytes long; it is " + string(rune(headerLen))
	}

	format := int(midi[8])<<8 | int(midi[9])
	if format != 1 {
		return false, "Unsupported MIDI type of " + string(rune(format)) + "; MIDI 1 format required"
	}
	return true, ""
}

// Generate a list of track events.
func ParseNoteEvents(midi []byte, numTracks int, timeDivison int) {
	headerLen := 14
	pos := headerLen
	// loop over each track
	for t := 0; t < numTracks; t++ {
		log.Printf("TRACK %v", t+1)
		var tempo int
		var lastEventStartByte byte

		// verify starting pos. of track (MTrk)
		mtrkLen := 4
		if string(midi[pos:pos+mtrkLen]) != "MTrk" {
			log.Fatal("Reached invalid starting position of track")
		}
		pos += mtrkLen

		// read track length in bytes
		trackLen := int(midi[pos])<<24 | int(midi[pos+1])<<16 | int(midi[pos+2])<<8 | int(midi[pos+3])
		pos += 4

		// loop over track event bytes
		endOfTrack := pos + trackLen
		for pos < endOfTrack {
			log.Printf("TRACK EVENT pos=%v", pos)
			// v time
			vTime, bytesRead := binary.Uvarint(midi[pos:]) // TODO: verify
			log.Printf("v_time %v", vTime)
			pos += bytesRead

			var eventStartByte byte
			if midi[pos]>>7 == 0b1 {
				eventStartByte = midi[pos]
				pos++
			} else {
				// if on data byte, enable running status using last command type
				eventStartByte = lastEventStartByte
			}
			lastEventStartByte = eventStartByte
			if eventStartByte == 255 { // meta event
				log.Println("meta event")
				metaType := midi[pos]
				pos++
				dataLen, bytesRead := binary.Uvarint(midi[pos:])
				pos += bytesRead
				switch metaType {
				case 0x03: // track name
					log.Println(string(midi[pos : pos+int(dataLen)]))
					pos += int(dataLen)
				case 0x51: // tempo
					tempo = int(midi[pos])<<16 | int(midi[pos+1])<<8 | int(midi[pos+2])
					log.Printf("tempo %v\n", tempo)
					pos += 3
				default:
					pos += int(dataLen)
				}
			} else {
				command := eventStartByte >> 4
				// channel := eventType & 0b00001111
				switch command {
				case 0x8: // note off
					log.Println("NOTE OFF")
					keyCode := midi[pos]
					pos++
					velocity := midi[pos]
					pos++
					log.Printf("key_code %v vel. %v\n", keyCode, velocity)
				case 0x9: // note on OR note off
					log.Println("NOTE ON")
					keyCode := midi[pos]
					pos++
					velocity := midi[pos]
					pos++
					log.Printf("key_code %v vel. %v\n", keyCode, velocity)
				case 0xA: // aftertouch (ignore)
					log.Println("AFTERTOUCH")
					pos += 2
				case 0xB: // continuous controller (ignore)
					log.Println("CONTINUOUS CONTROLLER")
					pos += 2
				case 0xC: // patch change (ignore)
					log.Println("PATCH CHANGE")
					pos += 1
				case 0xD: // channel pressure (ignore)
					log.Println("CHANNEL PRESSURE")
					pos += 1
				case 0xE: // pitch bend (ignore)
					log.Println("PITCH BEND")
					pos += 2
				case 0xF: // non-musical command (not implemented)
					log.Fatal("NON-MUSICAL COMMAND not implemented")
				}
			}
			log.Println()
		}
	}
}
