package memory

import "xeno/zohar/core/datatype"

type BufferCache[T [32]byte | [64]byte | [128]byte | [256]byte | [512]byte | [datatype.SIZE_1K]byte |
	[datatype.SIZE_2K]byte | [datatype.SIZE_4K]byte | [datatype.SIZE_8K]byte | [datatype.SIZE_16K]byte | [datatype.SIZE_32K]byte |
	[datatype.SIZE_64K]byte | [datatype.SIZE_128K]byte | [datatype.SIZE_256K]byte | [datatype.SIZE_512K]byte |
	[datatype.SIZE_1M]byte] struct {
	_holder T
}
