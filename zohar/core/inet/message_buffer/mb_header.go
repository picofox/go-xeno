package message_buffer

type MBHeader struct {
	_command int16
}

func (ego *MBHeader) Command() int16 {
	return ego._command
}

func (ego *MBHeader) SetCommand(cmd int16) {
	ego._command = cmd
}
