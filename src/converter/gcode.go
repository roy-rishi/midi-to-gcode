package converter

import "math"

// Decode MIDI note number (key code) as frequency (Hz). Eg: 69 -> a4 -> 440 Hz
// Return: pitch of note (Hz)
func noteNumToFreq(noteNum int) float64 {
	// f(n) = 440 * 2 ^ [(n - 69) / 12]
	return 440 * math.Pow(2, float64(noteNum-69)/12)
}
