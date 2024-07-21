package midi

type header struct {
	length   int
	format   int
	nTracks  int
	division int
}

type track struct {
	length int
	events []event
}

type event struct {
	vTime int
}
