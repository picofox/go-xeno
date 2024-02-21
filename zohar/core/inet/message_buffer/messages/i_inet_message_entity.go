package messages

type INetDataEntity interface {
	O1L15O1T15Serialize(*O1L15O1T15SerializationHelper) int32
	O1L15O1T15Deserialize(*O1L15O1T15DeserializationHelper) int32
	String() string
}
