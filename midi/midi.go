package midi

type Header struct {
	Length    int
	Format    int
	NumTracks int
	Division  int
}

type TrackChunk struct {
	Length int
	Name   string
	Events []TrackEvent
}

type TrackEvent struct {
	VTime    int
	Type     string
	Tempo    int
	Note     int
	Velocity int
}

// Get value of variable length int.
// Input: binary file, starting position
// Return: variable length value, position of next byte
func variableLen(data []byte, pos int) (int, int) {
	var vTime int
	for {
		if data[pos]&0b10000000 == 0b10000000 {
			// high-order bit is set, so there exists another byte
			vTime = vTime<<7 | int(data[pos]<<1)
			pos++
		} else {
			// last byte to include
			vTime = vTime<<7 | int(data[pos])
			pos++
			break
		}
	}
	return vTime, pos
}
