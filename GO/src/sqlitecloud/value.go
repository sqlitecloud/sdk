//
//                    ////              SQLite Cloud
//        ////////////  ///
//      ///             ///  ///        Product     : SQLite Cloud GO SDK
//     ///             ///  ///         Version     : 1.0.0
//     //             ///   ///  ///    Date        : 2021/08/31
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

type Value struct {
  Type     byte // _ + # : , $ ^ @ -   /// Types that are not in this Buffer: ROWSET, PUBSUB
  Buffer []byte
}

func (this *Value) GetType() byte     { 
  switch this.Type {
  case '!':                             return '+'  // Translate C-String to String
  case '{':                             return '#'  // Translate RAW-JSON to JSON 
  case '*', '/':                        return '*'  // Translate to ROWSET
  case ':', ',', '+', '#', '$', 
       '^', '@', '-', '|', '_':   
    if this.Buffer == nil {             return '_'  // unset buffer translates to NULL
    } else {                            return this.Type
    }
  default:                              return 0
  }
}
func (this *Value) IsSet()       bool { return this.GetType() != 0   }
func (this *Value) IsNULL()      bool { return this.GetType() == '_' }
func (this *Value) IsString()    bool { return this.GetType() == '+' }
func (this *Value) IsJSON()      bool { return this.GetType() == '#' }
func (this *Value) IsInteger()   bool { return this.GetType() == ':' }
func (this *Value) IsFloat()     bool { return this.GetType() == ',' }
func (this *Value) IsBLOB()      bool { return this.GetType() == '$' }
func (this *Value) IsPSUB()      bool { return this.GetType() == '|' }
func (this *Value) IsCommand()   bool { return this.GetType() == '^' }
func (this *Value) IsReconnect() bool { return this.GetType() == '@' }
func (this *Value) IsError()     bool { return this.GetType() == '-' }
func (this *Value) IsRowSet()    bool { return this.GetType() == '*' }

func (this *Value) IsText()      bool { return this.IsString() || this.IsInteger() || this.IsFloat() || this.IsBLOB() }

func (this *Value) GetLength() uint64 { return uint64( len( this.Buffer ) ) }
func (this *Value) GetBuffer() []byte { return this.Buffer } // Also good for BLOB

func (this *Value) GetString() string { return string(this.GetBuffer()) } // Also good for: JSON, BLOB, Command, Reconnect
func (this *Value) IsOK()        bool { return this.GetType() == '+' && this.GetString() == "OK" }

func (this *Value) GetInt32() (int32, error)  {
  switch value, err := strconv.ParseInt(this.GetString(), 0, 32); {
  case err != nil:                      return 0, err
  default:                              return int32(value), nil
  }
}
func (this *Value) GetInt64() (int64, error) { return strconv.ParseInt( this.GetString(), 0, 64) }

func (this *Value) GetFloat32() (float32, error) {
  switch value, err := strconv.ParseFloat(this.GetString(), 32); {
  case err != nil:                      return 0, err
  default:                              return float32(value), nil
  }
}
func (this *Value) GetFloat64() (float64, error) { return strconv.ParseFloat( this.GetString(), 64) }

func (this *Value) GetError() (ErrorCode int, ErrorMessage string, err error) {
  ErrorCode = 0
  for i, LEN, buffer := uint64(0), this.GetLength(), this.GetBuffer(); i < LEN; i++ {
    switch c := buffer[i]; c {
    case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':  ErrorCode = ErrorCode*10 + int(c) - int('0')
    default:                                                return ErrorCode, string(buffer[i+1:]), nil
    }
  }
  return -1, this.GetString(), errors.New("Invalid error format")
}

// Aux Methods

func (this *Value) GetSQLDateTime() (time.Time, error) {
  for _, format := range []string{
    "2006-01-02 15:04:05",
    "2006-01-02",
    "15:04:05",
  } {
    switch datetime, err := time.Parse(format, this.GetString()); {
    case err != nil:                      return time.Unix(0, 0), err
    default:                              return datetime, nil
    }
  }
  return time.Unix(0, 0), errors.New("Invalid Format")
}