package memory

func NeoBuffer(bufType int32, capacity int64) IByteBuffer {
	if bufType == 0 {
		bf := NeoLinearBuffer(capacity)
		return bf
	} else if bufType == 1 {
		bf := NeoRingBuffer(capacity)
		return bf
	}
	return nil
}
