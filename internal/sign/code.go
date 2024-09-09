package sign

type Code struct {
	Code string `json:"CODE"`
}

var CodeCH chan string = make(chan string, 1)
