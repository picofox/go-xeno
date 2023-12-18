package transcomm

const (
	Initialized = uint8(0)
	Connecting  = uint8(1)
	Connected   = uint8(2)
	Closed      = uint8(3)
)

var sConnStateToString []string = []string{
	"Initialized",
	"Connecting",
	"Connected",
	"Closed",
}

func ConnStateCodeToString(c uint8) string {
	if c > Closed {
		return "NAState"
	}
	return sConnStateToString[c]

}
