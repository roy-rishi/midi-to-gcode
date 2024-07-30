package midi

import (
	"log"
)

func hVerify(hBin []byte) bool {
	mthd := string(hBin[:4])
	return mthd == "MThd"
}

func hLen(hBin []byte) int {
	return int(hBin[4])<<24 | int(hBin[5])<<16 | int(hBin[6])<<8 | int(hBin[7])
}

func hFormat(hBin []byte) int {
	return int(hBin[8])<<8 | int(hBin[9])
}

func hNumTracks(hBin []byte) int {
	return int(hBin[10])<<8 | int(hBin[11])
}

func hDivision(hBin []byte) int {
	return int(hBin[12])<<8 | int(hBin[13])
}

// Parse the MIDI header for file metadata.
// Input: full MIDI file contents
func ParseHeader(hBin []byte) Header {
	h := Header{}
	// MThd
	if !hVerify(hBin) {
		log.Fatal("File not recognized as MIDI")
	}
	// length
	len := hLen(hBin)
	if len != 6 {
		log.Fatalf("Header must be 6 bytes long; it is %d", len)
	}
	h.Length = len
	// format
	format := hFormat(hBin)
	if format != 1 {
		log.Fatalf("Format must be 1; it is %v", format)
	}
	h.Format = format
	h.NumTracks = hNumTracks(hBin) // n
	h.Division = hDivision(hBin)   // division

	log.Printf("HEADER %+v\n\n", h)
	return h
}
