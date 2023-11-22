package memory

func NumberOfOneInInt32(a int32) int8 {
	count := int8(0)
	for a != 0 {
		a = a & (a - 1)
		count++
	}
	return count
}
