package core

import "fmt"

const (
	EC_OK int32 = iota
	EC_NOOP
	EC_CREATE_DIR_FAILED
	EC_DIR_ALREADY_EXIST
	EC_DELETE_DIR_FAILED
	EC_ENSURE_DIR_FAILED
	EC_SET_CWD_FAILED
	EC_TO_ABS_PATH_FAILED
	EC_TYPE_CONVERT_FAILED
	EC_TYPE_MISMATCH
	EC_NULL_VALUE
	EC_HANDLER_NOT_FOUND
	EC_INDEX_OOB
	EC_LOCK_FILE_FAILED
	EC_UNLOCK_FILE_FAILED
	EC_TRY_LOCK_FILE_FAILED
	EC_FILE_OPEN_FAILED
	EC_FILE_CLOSE_FAILED
	EC_FILE_READ_FAILED
	EC_FILE_WRITE_FAILED
	EC_FILE_SYNC_FAILED
	EC_FILE_SEEK_FAILED
	EC_FILE_STAT_FAILED
	EC_EOF
	EC_JSON_MARSHAL_FAILED
	EC_JSON_UNMARSHAL_FAILED
	EC_ELEMENT_NOT_FOUND
	EC_CONNECT_DB_FAILED
	EC_CREATE_DB_FAILED
	EC_PING_DB_FAILED
	EC_DB_RETRIVE_DATA_FAILED
	EC_DB_GET_COL_INFO_FAILED
	EC_DB_SET_AUTOCOMMIT_FAILED
	EC_DB_BEGIN_TRANS_FAILED
	EC_DB_COMMIT_TRANS_FAILED
	EC_DB_ROLLBACL_TRANS_FAILED
	EC_DB_DELETE_FAIELD
	EC_DB_PREPARE_FAILED
	EC_DB_INSERT_FAILED
	EC_STRING_TOO_LONG
	EC_RESPACE_FAILED
	EC_INCOMPLETE_DATA
	EC_ELEMENT_EXIST
	EC_INVALID_STATE
	EC_REACH_LIMIT
	EC_TRY_AGAIN
	EC_EPOLL_WAIT_ERROR
	EC_ACCEPT_ERROR
	EC_TCP_SEND_FAILED
	EC_TCP_CONNECT_ERROR
	EC_SET_NONBLOCK_ERROR
	EC_LISTEN_ERROR
	EC_GET_LOW_FD_ERROR
	EC_TCO_RECV_ERROR
	EC_MESSAGE_HANDLING_ERROR
	EC_CREATE_SOCKET_ERROR
	EC_SET_SOCKOPT_ERROR
	EC_TCP_DISCONNECTED
	EC_ALREADY_DONE
	EC_DESERIALIZE_FIELD_FAIELD
	EC_SERIALIZE_FIELD_FAIELD
	EC_MIN_VALUE_FIND_ERROR
	EC_READ_HEADER_ERROR
	EC_WRITE_HEADER_ERROR
	EC_ERROR_COUNT
)

var g_error_str = [EC_ERROR_COUNT]string{
	"OK",
	"NOOP",
	"CREATE_DIR_FAILED",
	"EC_DIR_ALREADY_EXIST",
	"DELETE_DIR_FAILED",
	"ENSURE_DIR_FAILED",
	"SET_CWD_FAILED",
	"TO_ABS_PATH_FAILED",
	"EC_TYPE_CONVERT_FAILED",
	"EC_TYPE_MISMATCH",
	"EC_NULL_VALUE",
	"EC_HANDLER_NOT_FOUND",
	"EC_INDEX_OOB",
	"EC_LOCK_FILE_FAILED",
	"EC_UNLOCK_FILE_FAILED",
	"EC_TRY_LOCK_FILE_FAILED",
	"EC_FILE_OPEN_FAILED",
	"EC_FILE_CLOSE_FAILED",
	"EC_FILE_READ_FAILED",
	"EC_FILE_WRITE_FAILED",
	"EC_FILE_SYNC_FAILED",
	"EC_FILE_SEEK_FAILED",
	"EC_FILE_STAT_FAILED",
	"EC_EOF",
	"EC_JSON_MARSHAL_FAILED",
	"EC_JSON_UNMARSHAL_FAILED",
	"EC_ELEMENT_NOT_FOUND",
	"EC_CONNECT_DB_FAILED",
	"EC_CREATE_DB_FAILED",
	"EC_PING_DB_FAILED",
	"EC_DB_RETRIVE_DATA_FAILED",
	"EC_DB_GET_COL_INFO_FAILED",
	"EC_DB_SET_AUTOCOMMIT_FAILED",
	"EC_DB_BEGIN_TRANS_FAILED",
	"EC_DB_COMMIT_TRANS_FAILED",
	"EC_DB_ROLLBACK_TRANS_FAILED",
	"EC_DB_DELETE_FAIELD",
	"EC_DB_PREPARE_FAILED",
	"EC_DB_INSERT_FAILED",
	"EC_STRING_TOO_LONG",
	"EC_RESPACE_FAILED",
	"EC_INCOMPLETE_DATA",
	"EC_ELEMENT_EXIST",
	"EC_INVALID_STATE",
	"EC_REACH_LIMIT",
	"EC_TRY_AGAIN",
	"EC_EPOLL_WAIT_ERROR",
	"EC_ACCEPT_ERROR",
	"EC_TCP_SEND_FAILED",
	"EC_TCP_CONNECT_ERROR",
	"EC_SET_NONBLOCK_ERROR",
	"EC_LISTEN_ERROR",
	"EC_GET_LOW_FD_ERROR",
	"EC_TCO_RECV_ERROR",
	"EC_MESSAGE_HANDLING_ERROR",
	"EC_CREATE_SOCKET_ERROR",
	"EC_SET_SOCKOPT_ERROR",
	"EC_TCP_DISCONNECTED",
	"EC_ALREADY_DONE",
	"EC_DESERIALIZE_FIELD_FAIELD",
	"EC_SERIALIZE_FIELD_FAIELD",
	"EC_MIN_VALUE_FIND_ERROR",
	"EC_READ_HEADER_ERROR",
	"EC_WRITE_HEADER_ERROR",
}

func MkErr(et int32, mark int32) int32 {
	rc := ((et&0x7FFF)|(1<<15))<<16 | (mark & 0xFFFF)
	return rc
}

func MkSuccess(mark int32) int32 {
	return mark & 0xFFFF
}

func ExErr(ec int32) (int32, int32) {
	ec &= 0x7FFFFFFF
	et := (ec >> 16) & 0x7FFF
	em := ec & 0xFFFF
	return et, em
}

func IsErrType(ec int32, et int32) bool {
	t, _ := ExErr(ec)
	return t == et
}

func Err(ec int32) bool {
	et, _ := ExErr(ec)
	if et != 0 {
		return true
	}
	return false
}

func ErrStr(ec int32) string {
	et, em := ExErr(ec)
	return fmt.Sprintf("%s(%d).%d", g_error_str[et], et, em)
}
