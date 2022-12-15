//
//                    ////              SQLite Cloud
//        ////////////  ///
//      ///             ///  ///        Product     : SQLite Cloud GO SDK
//     ///             ///  ///         Version     : 1.0.1
//     //             ///   ///  ///    Date        : 2021/10/13
//    ///             ///   ///  ///    Author      : Andreas Pfeil
//   ///             ///   ///  ///
//   ///     //////////   ///  ///      Description : GO Methods related to the
//   ////                ///  ///                     Value class.
//     ////     //////////   ///
//        ////            ////
//          ////     /////
//             ///                      Copyright   : 2021 by SQLite Cloud Inc.
//
// -----------------------------------------------------------------------TAB=2

package sqlitecloud

import (
	"errors"
	"strconv"
	"time"
)

const (
	CMD_STRING       = '+'
	CMD_ZEROSTRING   = '!'
	CMD_ERROR        = '-'
	CMD_INT          = ':'
	CMD_FLOAT        = ','
	CMD_ROWSET       = '*'
	CMD_ROWSET_CHUNK = '/'
	CMD_JSON         = '#'
	CMD_RAWJSON      = '{'
	CMD_NULL         = '_'
	CMD_BLOB         = '$'
	CMD_COMPRESSED   = '%'
	CMD_PUBSUB       = '|'
	CMD_COMMAND      = '^'
	CMD_RECONNECT    = '@'
	CMD_ARRAY        = '='
)

const (
	NO_EXTCODE = 0
	NO_OFFCODE = -1
)

type Value struct {
	Type   byte // _ + # : , $ ^ @ -   /// Types that are not in this Buffer: ROWSET, PUBSUB
	Buffer []byte
}

func (this *Value) GetType() byte {
	switch this.Type {
	case CMD_ZEROSTRING:
		return CMD_STRING // Translate C-String to String
	case CMD_RAWJSON:
		return CMD_JSON // Translate RAW-JSON to JSON
	case CMD_ROWSET, CMD_ROWSET_CHUNK:
		return CMD_ROWSET // Translate to ROWSET
	case CMD_ARRAY:
		return CMD_ARRAY // Array
	case CMD_INT, CMD_FLOAT, CMD_STRING, CMD_JSON, CMD_BLOB,
		CMD_COMMAND, CMD_RECONNECT, CMD_ERROR, CMD_PUBSUB, CMD_NULL:
		if this.Buffer == nil {
			return CMD_NULL // unset buffer translates to NULL
		} else {
			return this.Type
		}
	default:
		return 0
	}
}
func (this *Value) IsSet() bool       { return this.GetType() != 0 }
func (this *Value) IsNULL() bool      { return this.GetType() == CMD_NULL }
func (this *Value) IsString() bool    { return this.GetType() == CMD_STRING }
func (this *Value) IsJSON() bool      { return this.GetType() == CMD_JSON }
func (this *Value) IsInteger() bool   { return this.GetType() == CMD_INT }
func (this *Value) IsFloat() bool     { return this.GetType() == CMD_FLOAT }
func (this *Value) IsBLOB() bool      { return this.GetType() == CMD_BLOB }
func (this *Value) IsPSUB() bool      { return this.GetType() == CMD_PUBSUB }
func (this *Value) IsCommand() bool   { return this.GetType() == CMD_COMMAND }
func (this *Value) IsReconnect() bool { return this.GetType() == CMD_RECONNECT }
func (this *Value) IsError() bool     { return this.GetType() == CMD_ERROR }
func (this *Value) IsRowSet() bool    { return this.GetType() == CMD_ROWSET }
func (this *Value) IsArray() bool     { return this.GetType() == CMD_ARRAY }

func (this *Value) IsText() bool {
	return this.IsString() || this.IsInteger() || this.IsFloat() || this.IsBLOB()
}

func (this *Value) GetLength() uint64 { return uint64(len(this.Buffer)) }
func (this *Value) GetBuffer() []byte { return this.Buffer } // Also good for BLOB

func (this *Value) GetString() string { return string(this.GetBuffer()) } // Also good for: JSON, BLOB, Command, Reconnect
func (this *Value) IsOK() bool        { return this.GetType() == '+' && this.GetString() == "OK" }

func (this *Value) GetInt32() (int32, error) {
	switch value, err := strconv.ParseInt(this.GetString(), 0, 32); {
	case err != nil:
		return 0, err
	default:
		return int32(value), nil
	}
}
func (this *Value) GetInt64() (int64, error) { return strconv.ParseInt(this.GetString(), 0, 64) }

func (this *Value) GetFloat32() (float32, error) {
	switch value, err := strconv.ParseFloat(this.GetString(), 32); {
	case err != nil:
		return 0, err
	default:
		return float32(value), nil
	}
}
func (this *Value) GetFloat64() (float64, error) { return strconv.ParseFloat(this.GetString(), 64) }

// GetError returns the ErrorCode, ExtErrorCode, ErrorOffset, ErrorMessage
// and the error object of the receiver
func (this *Value) GetError() (int, int, int, string, error) {
	ErrorCode := 0
	ExtErrorCode := NO_EXTCODE
	ErrorOffset := NO_OFFCODE
	nColons := 0
	for i, LEN, buffer := uint64(0), this.GetLength(), this.GetBuffer(); i < LEN; i++ {
		switch c := buffer[i]; c {
		case ':':
			nColons++
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			switch nColons {
			case 0:
				ErrorCode = ErrorCode*10 + int(c) - int('0')
			case 1:
				ExtErrorCode = ExtErrorCode*10 + int(c) - int('0')
			case 2:
				if ErrorOffset == NO_OFFCODE {
					ErrorOffset = 0
				}
				ErrorOffset = ErrorOffset*10 + int(c) - int('0')
			}

		default:
			return ErrorCode, ExtErrorCode, ErrorOffset, string(buffer[i+1:]), nil
		}
	}
	return -1, NO_EXTCODE, NO_OFFCODE, this.GetString(), errors.New("Invalid error format")
}

// Aux Methods

func (this *Value) GetSQLDateTime() (time.Time, error) {
	for _, format := range []string{
		"2006-01-02 15:04:05",
		"2006-01-02",
		"15:04:05",
	} {
		switch datetime, err := time.Parse(format, this.GetString()); {
		case err != nil:
			return time.Unix(0, 0), err
		default:
			return datetime, nil
		}
	}
	return time.Unix(0, 0), errors.New("Invalid Format")
}
