package memory

type ISerializable interface {
	Serialize(ISerializationHeader, IByteBuffer) (int64, int32)
	Deserialize(ISerializationHeader, IByteBuffer) (int64, int32)
}

type ISerializationHeader interface {
	BeginSerializing(buffer IByteBuffer) (int64, int64, int32)
	EndSerializing(buffer IByteBuffer, headerPos int64, totalLength int64) (int64, int32)
	BeginDeserializing(buffer IByteBuffer, validate bool) (int64, int32)
	EndDeserializing(buffer IByteBuffer) int32
	HeaderLength() int64
	BodyLength() int64
	String() string
}
