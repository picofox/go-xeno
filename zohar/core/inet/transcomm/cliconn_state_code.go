package transcomm

const (
	Uninitialized = uint8(0)
	Initialized   = uint8(1)
	Connecting    = uint8(2)
	Connected     = uint8(3)
	Closed        = uint8(4)
	Finalizing    = uint8(5)
)

var sCliConnStateToString []string = []string{
	"Uninitialized",
	"Initialized",
	"Connecting",
	"Connected",
	"Closed",
	"Finalizing",
}

func CliConnStateCodeToString(c uint8) string {
	if c > Finalizing {
		return "NAState"
	}
	return sCliConnStateToString[c]

}
