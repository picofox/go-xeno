package unique

type IIdGenerator interface {
	Next() int64
}
