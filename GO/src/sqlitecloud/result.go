//
//                    ////              SQLite Cloud
//        ////////////  ///
//      ///             ///  ///        Product     : SQLite Cloud GO SDK
//     ///             ///  ///         Version     : 1.1.0
//     //             ///   ///  ///    Date        : 2021/10/01
//    ///             ///   ///  ///    Author      : Andreas Pfeil
//   ///             ///   ///  ///
//   ///     //////////   ///  ///      Description : GO Methods related to the
//   ////                ///  ///                     Result class.
//     ////     //////////   ///
//        ////            ////
//          ////     /////
//             ///                      Copyright   : 2021 by SQLite Cloud Inc.
//
// -----------------------------------------------------------------------TAB=2

package sqlitecloud

import "fmt"
import "os"

import "bufio"
import "bytes"
import "strings"
import "errors"
import "time"
import "io"
import "encoding/json"
import "golang.org/x/term"

var rowsetChunkEndPatterns = []([]byte){ []byte( "5 0 0 0" ), []byte( "6 0 0 0 " ), []byte( "8 0 0 0 0 " ) }

const OUTFORMAT_LIST      = 0
const OUTFORMAT_CSV       = 1
const OUTFORMAT_QUOTE     = 2
const OUTFORMAT_TABS      = 3
const OUTFORMAT_LINE      = 4
const OUTFORMAT_JSON      = 5
const OUTFORMAT_HTML      = 6
const OUTFORMAT_MARKDOWN  = 7
const OUTFORMAT_TABLE     = 8
const OUTFORMAT_BOX       = 9
const OUTFORMAT_XML       = 10

// The Result is either a Literal or a RowSet
type Result struct {
  uncompressedChuckSizeSum  uint64

  value                     Value

  rows                      []ResultRow
  ColumnNames               []string
  ColumnWidth               []uint64
  MaxHeaderWidth            uint64
}

func (this *Result ) Rows() []ResultRow { return this.rows }

// ResultSet Methods (100% GO)

// GetType returns the type of this query result as an integer (see: RESULT_ constants).
func (this *Result ) GetType() byte { return this.value.GetType() }

// IsOK returns true if this query result if of type "RESULT_OK", false otherwise.
func (this *Result ) IsOK()     bool { return this.value.IsOK() }

// GetNumberOfRows returns the number of rows in this query result
func (this *Result ) GetNumberOfRows() uint64 {
  switch {
  case !this.IsRowSet():  return 0
  default:                return uint64( len( this.rows ) )
  }
}

// GetNumberOfColumns returns the number of columns in this query result
func (this *Result ) GetNumberOfColumns() uint64 {
  switch {
  case !this.IsRowSet():  return 0
  default:                return uint64( len( this.ColumnWidth ) )
  }
}

// Dump outputs this query result to the screen.
// Warning: No line truncation is used. If you want to truncation the output to a certain width, use: Result.DumpToScreen( width )
func (this *Result ) Dump() {
  w := 0
  if width, _, err := term.GetSize( 0 ); err == nil { w = width }
  this.DumpToScreen( uint( w ) )
}

// ToJSON returns a JSON representation of this query result.
// BUG(andreas): The Result.ToJSON method is not implemented yet.
func (this *Result ) ToJSON() string {
  return "todo" // Use Writer into Buffer
}

// Additional ResultSet Methods (100% GO)

// GetMaxColumnLength returns the number of runes of the value in the specified column with the maximum length in this query result.
// The Column index is an unsigned int in the range of 0...GetNumberOfColumns() - 1.
func (this *Result ) GetMaxColumnWidth( Column uint64 ) ( uint64, error ) {
  switch {
  case !this.IsRowSet():                    return 0, errors.New( "Not a RowSet" )
  case Column >= this.GetNumberOfColumns(): return 0, errors.New( "Column Index out of bounds" )
  default:                                  return this.ColumnWidth[ Column ], nil
  }
}
// GetNameWidth returns the number of runes of the column name in the specified column.
// The Column index is an unsigned int in the range of 0...GetNumberOfColumns() - 1.
func (this *Result ) GetNameLength( Column uint64 ) ( uint64, error ) {
  switch {
  case !this.IsRowSet():                    return 0, errors.New( "Not a RowSet" )
  case Column >= this.GetNumberOfColumns(): return 0, errors.New( "Column Index out of bounds" )
  default:                                  return uint64( len( this.ColumnNames[ Column ] ) ), nil
  }
}
// GetMaxNameWidth returns the number of runes of the longest column name.
func (this *Result ) GetMaxNameWidth() uint64 { return this.MaxHeaderWidth }

// Additional Data Access Functions (100% GO)

// IsError returns true if this query result is of type "RESULT_ERROR", false otherwise.
func (this *Result ) IsError()        bool { return this.value.IsError() }

// IsNull returns true if this query result is of type "RESULT_NULL", false otherwise.
func (this *Result ) IsNULL()           bool { return this.value.IsNULL() }

// IsJson returns true if this query result is of type "RESULT_JSON", false otherwise.
func (this *Result ) IsJSON()           bool { return this.value.IsJSON() }

// IsString returns true if this query result is of type "RESULT_STRING", false otherwise.
func (this *Result ) IsString()         bool { return this.value.IsString() }

// IsInteger returns true if this query result is of type "RESULT_INTEGER", false otherwise.
func (this *Result ) IsInteger()      bool { return this.value.IsInteger() }

// IsFloat returns true if this query result is of type "RESULT_FLOAT", false otherwise.
func (this *Result ) IsFloat()        bool { return this.value.IsFloat() }

// IsPSUB returns true if this query result is of type "RESULT_XXXXXXX", false otherwise.
func (this *Result ) IsPSUB()           bool { return this.value.IsPSUB() }

// IsCommand returns true if this query result is of type "RESULT_XXXXXXX", false otherwise.
func (this *Result ) IsCommand()      bool { return this.value.IsCommand() }

// IsReconnect returns true if this query result is of type "RESULT_XXXXXXX", false otherwise.
func (this *Result ) IsReconnect()    bool { return this.value.IsReconnect() }

func (this *Result ) IsBLOB()               bool { return this.value.IsBLOB() }

// IsText returns true if this query result is of type "RESULT_JSON", "RESULT_STRING", "RESULT_INTEGER" or "RESULT_FLOAT", false otherwise.
func (this *Result ) IsText()       bool { return this.value.IsText() }

// IsRowSet returns true if this query result is of type "RESULT_ROWSET", false otherwise.
func (this *Result ) IsRowSet()        bool {
  switch {
  case !this.value.IsRowSet():    return false
  case this.rows == nil:          return false
  case this.ColumnNames == nil:   return false
  case this.ColumnWidth == nil:   return false
  case this.MaxHeaderWidth == 0:  return false
  default:                        return true
  }
}
func (this *Result ) IsLiteral()        bool { return !this.IsRowSet() }

// ResultSet Buffer/Scalar Methods

// GetUncompressedChuckSizeSum returns the
func (this *Result ) GetUncompressedChuckSizeSum() uint64 { return this.uncompressedChuckSizeSum }


// GetBufferLength returns the length of the buffer of this query result.
func (this *Result ) GetBufferLength() ( uint64, error ) {
  switch {
  case this.IsRowSet(): return 0, errors.New( "Not a scalar value" )
  default:              return this.value.GetLength(), nil
  }
}

// GetBuffer returns the buffer of this query result as string.
func (this *Result ) GetBuffer() []byte { return this.value.GetBuffer() }

func (this *Result ) GetString() ( string, error ) {
  switch {
  case this.IsRowSet(): return "", errors.New( "Not a literal" )
  default:              return this.value.GetString(), nil
  }
}
func (this *Result ) GetString_() string { 
  value, _ := this.GetString()
  return value
}

func (this *Result ) GetJSON() ( object interface{}, err error ) {
  switch {
  case !this.IsJSON(): return nil, errors.New( "Not a JSON object" )
  default:
    err = json.Unmarshal( this.value.GetBuffer(), object )
    return
  }
}
func (this *Result ) GetJSON_() ( object interface{} ) { 
  value, _ := this.GetJSON()
  return value
}

func (this *Result ) GetInt32() ( int32, error ) {
  switch {
  case !this.IsInteger(): return 0, errors.New( "Not an integer value" )
  default:                return this.value.GetInt32()
  }
}
func (this *Result ) GetInt32_() int32 { 
  value, _ := this.GetInt32()
  return value
}

func (this *Result ) GetInt64() ( int64, error ) {
  switch {
  case !this.IsInteger(): return 0, errors.New( "Not an integer value" )
  default:                return this.value.GetInt64()
  }
}
func (this *Result ) GetInt64_() int64 { 
  value, _ := this.GetInt64()
  return value
}

func (this *Result ) GetFloat32() ( float32, error ) {
  switch {
  case !this.IsFloat():   return 0, errors.New( "Not a float value" )
  default:                return this.value.GetFloat32()
  }
}
func (this *Result ) GetFloat32_() float32 { 
  value, _ := this.GetFloat32()
  return value
}

func (this *Result ) GetFloat64() ( float64, error ) {
  switch {
  case !this.IsFloat():   return 0, errors.New( "Not a float value" )
  default:                return this.value.GetFloat64()
  }
}
func (this *Result ) GetFloat64_() float64 { 
  value, _ := this.GetFloat64()
  return value
}

func (this *Result ) GetError() ( int, string, error ) {
  switch {
  case !this.IsError():   return 0, "", errors.New( "Not an error" )
  default:                return this.value.GetError()
  }
}
func (this *Result ) GetError_() ( int, string ) { 
  code, message, _ := this.GetError()
  return code, message
}

func (this *Result ) GetErrorAsString() string {
  switch code, message, err := this.GetError(); {
  case err != nil:  return fmt.Sprintf( "INTERNAL ERROR: %s", err.Error() )
  default:          return fmt.Sprintf( "ERROR: %s (%d)", message, code )
  }
}

// Free frees all memory allocated by this query result.
func (this *Result ) Free() {
  this.value          = Value{ Type: 0, Buffer: nil } // GC
  this.rows           = []ResultRow{}                 // GC
  this.ColumnNames    = []string{}                    // GC
  this.ColumnWidth    = []uint64{}                    // GC
  this.MaxHeaderWidth = 0
}


// GetName returns the column name in column Column of this query result.
// The Column index is an unsigned int in the range of 0...GetNumberOfColumns() - 1.
func (this *Result ) GetName( Column uint64 ) ( string, error ) {
  switch {
  case Column >= this.GetNumberOfColumns(): return "", errors.New( "Column Index out of bounds" )
  default:                                  return this.ColumnNames[ Column ], nil
  }
}
func (this *Result ) GetName_( Column uint64 ) string { 
  value, _ := this.GetName( Column )
  return value
}

// DumpToScreen outputs this query result to the screen.
// The output is truncated at a maximum line width of MaxLineLength runes (compare: Result.Dump())
func (this *Result ) DumpToScreen( MaxLineLength uint ) {
  this.DumpToWriter( bufio.NewWriter( os.Stdout ), OUTFORMAT_BOX, false, "│", "NULL", "\r\n", MaxLineLength, false )
}

////// Row Methods (100% GO)

// GetRow returns a pointer to the row Row of this query result.
// The Row index is an unsigned int in the range of 0...GetNumberOfRows() - 1.
// If the index row can not be found, nil is returned instead.
func (this *Result ) GetRow( Row uint64 ) ( *ResultRow, error ) {
  switch {
  case !this.IsRowSet():              return nil, errors.New( "Not a Rowset" )
  case Row >= this.GetNumberOfRows(): return nil, errors.New( "Row index is out of bounds" )
  default:                            return &this.rows[ Row ], nil
  }
}

// GetFirstRow returns the first row of this query result.
// If this query result has no row's, nil is returned instead.
func (this *Result ) GetFirstRow() ( *ResultRow, error ) { return this.GetRow( 0 ) }

// GetLastRow returns the first row of this query result.
// If this query result has no row's, nil is returned instead.
func (this *Result ) GetLastRow() ( *ResultRow, error ) { return this.GetRow( this.GetNumberOfRows() - 1 ) }





// Additional Row Methods <- sollte es eigentlich nicht geben!!!!!

func (this *Result ) GetValue( Row uint64, Column uint64 ) ( *Value, error ) {
  switch row, err := this.GetRow( Row ); {
  case err != nil:                      return nil, err
  default:                              return row.GetValue( Column )
  }
}

// GetValueType returns the type of the value in row Row and column Column of this query result.
// The Row index is an unsigned int in the range of 0...GetNumberOfRows() - 1.
// The Column index is an unsigned int in the range of 0...GetNumberOfColumns() - 1.
// Possible return types are: VALUE_INTEGER, VALUE_FLOAT, VALUE_TEXT, VALUE_BLOB, VALUE_NULL
func (this *Result ) GetValueType( Row uint64, Column uint64 ) ( byte, error ) {
  switch value, err := this.GetValue( Row, Column ); {
  case err != nil:    return '_', err
  default:            return value.GetType(), nil
  }
}
func (this *Result ) GetValueType_( Row uint64, Column uint64 ) byte { 
  value, _ := this.GetValueType( Row, Column )
  return value
}

// GetStringValue returns the contents in row Row and column Column of this query result as string.
// The Row index is an unsigned int in the range of 0...GetNumberOfRows() - 1.
// The Column index is an unsigned int in the range of 0...GetNumberOfColumns() - 1.
func (this *Result ) GetStringValue( Row uint64, Column uint64 ) ( string, error ) {
  switch value, err := this.GetValue( Row, Column ); {
  case err != nil:    return "", err
  default:            return value.GetString(), nil
  }
}
func (this *Result ) GetStringValue_( Row uint64, Column uint64 ) string { 
  value, _ := this.GetStringValue( Row, Column )
  return value
}

// GetInt32Value returns the contents in row Row and column Column of this query result as int32.
// The Row index is an unsigned int in the range of 0...GetNumberOfRows() - 1.
// The Column index is an unsigned int in the range of 0...GetNumberOfColumns() - 1.
func (this *Result ) GetInt32Value( Row uint64, Column uint64 ) ( int32, error ) {
  switch value, err := this.GetValue( Row, Column ); {
  case err != nil:    return 0, err
  default:            return value.GetInt32()
  }
}
func (this *Result ) GetInt32Value_( Row uint64, Column uint64 ) int32 { 
  value, _ := this.GetInt32Value( Row, Column )
  return value
}

// GetInt64Value returns the contents in row Row and column Column of this query result as int64.
// The Row index is an unsigned int in the range of 0...GetNumberOfRows() - 1.
// The Column index is an unsigned int in the range of 0...GetNumberOfColumns() - 1.
func (this *Result ) GetInt64Value( Row uint64, Column uint64 ) ( int64, error ) {
  switch value, err := this.GetValue( Row, Column ); {
    case err != nil:    return 0, err
    default:            return value.GetInt64()
    }
}
func (this *Result ) GetInt64Value_( Row uint64, Column uint64 ) int64 { 
  value, _ := this.GetInt64Value( Row, Column )
  return value
}

// GetFloat32Value returns the contents in row Row and column Column of this query result as float32.
// The Row index is an unsigned int in the range of 0...GetNumberOfRows() - 1.
// The Column index is an unsigned int in the range of 0...GetNumberOfColumns() - 1.
func (this *Result ) GetFloat32Value( Row uint64, Column uint64 ) ( float32, error ) {
  switch value, err := this.GetValue( Row, Column ); {
  case err != nil:    return 0, err
  default:            return value.GetFloat32()
  }
}
func (this *Result ) GetFloat32Value_( Row uint64, Column uint64 ) float32 { 
  value, _ := this.GetFloat32Value( Row, Column )
  return value
}

// GetFloat64Value returns the contents in row Row and column Column of this query result as float64.
// The Row index is an unsigned int in the range of 0...GetNumberOfRows() - 1.
// The Column index is an unsigned int in the range of 0...GetNumberOfColumns() - 1.
func (this *Result ) GetFloat64Value( Row uint64, Column uint64 ) ( float64, error ) {
  switch value, err := this.GetValue( Row, Column ); {
  case err != nil:    return 0, err
  default:            return value.GetFloat64()
  }
}
func (this *Result ) GetFloat64Value_( Row uint64, Column uint64 ) float64 { 
  value, _ := this.GetFloat64Value( Row, Column )
  return value
}

// GetSQLDateTime parses this query result value in Row and Column as an SQL-DateTime and returns its value.
// The Row index is an unsigned int in the range of 0...GetNumberOfRows() - 1.
// The Column index is an unsigned int in the range of 0...GetNumberOfColumns() - 1.
func (this *Result ) GetSQLDateTime( Row uint64, Column uint64 ) ( time.Time, error ) {
  switch value, err := this.GetValue( Row, Column ); {
  case err != nil:    return time.Unix( 0, 0 ), err
  default:            return value.GetSQLDateTime()
  }
}
func (this *Result ) GetSQLDateTime_( Row uint64, Column uint64 ) time.Time { 
  value, _ := this.GetSQLDateTime( Row, Column )
  return value
}


////////////////////////////

func trimStringToMaxLength( Buffer string, MaxLineLength uint ) string {
  switch {
  case MaxLineLength == 0:                            return Buffer
  case MaxLineLength >= uint( len([]rune(Buffer) ) ): return Buffer
  default:                                            return fmt.Sprintf( fmt.Sprintf( "%%.%ds…", MaxLineLength - 1 ), Buffer )
  }
}
func renderCenteredString( Buffer string, Width int ) string {
  return fmt.Sprintf( "%[1]*s", -Width, fmt.Sprintf( "%[1]*s", ( Width + len( Buffer ) ) / 2, Buffer ) )
}

func (this *Result) renderHorizontalTableLine( Left string, Fill string, Separator string, Right string ) string {
  outBuffer := ""
  for _, columnWidth := range this.ColumnWidth {
    outBuffer = fmt.Sprintf( "%s%s%s", outBuffer, strings.Repeat( Fill, int( columnWidth + 2 ) ), Separator )
  }
  return fmt.Sprintf( "%s%s%s", Left, strings.TrimRight( outBuffer, Separator ), Right )
}
func (this *Result) renderTableColumnNames( Left string, Separator string, Right string ) string {
  outBuffer := ""
  for forThisColumn, columnWidth := range this.ColumnWidth {
    columnName, _ := this.GetName( uint64( forThisColumn ) )
    outBuffer      = fmt.Sprintf( "%s%s%s", outBuffer, renderCenteredString( columnName, int( columnWidth + 2 ) ), Separator )
  }
  return fmt.Sprintf( "%s%s%s", Left, strings.TrimRight( outBuffer, Separator ), Right )
}
func (this *Result) renderTableHeader( Format int, Separator string, NewLine string, MaxLineLength uint ) string {
  switch( Format ) {
    case OUTFORMAT_JSON: return fmt.Sprintf( "[%s", NewLine )

    case OUTFORMAT_MARKDOWN:
      return trimStringToMaxLength( this.renderTableColumnNames( Separator, Separator, Separator ), MaxLineLength )         + NewLine +
             trimStringToMaxLength( this.renderHorizontalTableLine( Separator, "-", Separator, Separator ), MaxLineLength ) + NewLine

    case OUTFORMAT_TABLE:
      tableLine := trimStringToMaxLength( this.renderHorizontalTableLine( "+", "-", "+", "+" ), MaxLineLength )             + NewLine
      return  tableLine                                                                                                     +
              trimStringToMaxLength( this.renderTableColumnNames( Separator, Separator, Separator ), MaxLineLength )        + NewLine +
              tableLine

    case OUTFORMAT_BOX:
      return trimStringToMaxLength( this.renderHorizontalTableLine( "┌", "─", "┬", "┐" ), MaxLineLength )                   + NewLine +
             trimStringToMaxLength( this.renderTableColumnNames( Separator, Separator, Separator ), MaxLineLength )         + NewLine +
             trimStringToMaxLength( this.renderHorizontalTableLine( "├", "─", "┼", "┤" ), MaxLineLength )                   + NewLine
    case OUTFORMAT_XML:
      return trimStringToMaxLength( "<?xml version=\"1.0\" encoding=\"UTF-8\" standalone=\"no\"?>", MaxLineLength )         + NewLine +
             trimStringToMaxLength( fmt.Sprintf( "<resultset statement=\"%s\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\">", Separator ), MaxLineLength ) + NewLine

    default:
      return "" // No header
  }
}
func (this *Result) renderTableFooter( Format int, NewLine string, MaxLineLength uint ) string {
  switch( Format ) {
    case OUTFORMAT_JSON:  return trimStringToMaxLength( "]", MaxLineLength )                                                  + NewLine
    case OUTFORMAT_TABLE: return trimStringToMaxLength( this.renderHorizontalTableLine( "+", "-", "+", "+" ), MaxLineLength ) + NewLine
    case OUTFORMAT_BOX:   return trimStringToMaxLength( this.renderHorizontalTableLine( "└", "─", "┴", "┘" ), MaxLineLength ) + NewLine
    case OUTFORMAT_XML:   return trimStringToMaxLength( "</resultset>", MaxLineLength )                                       + NewLine
    default:              return "" // No footer
  }
}

// DumpToWriter renders this query result into the buffer of an io.Writer.
// The output Format can be specified and must be one of the following values: OUTFORMAT_LIST, OUTFORMAT_CSV, OUTFORMAT_QUOTE, OUTFORMAT_TABS, OUTFORMAT_LINE, OUTFORMAT_JSON, OUTFORMAT_HTML, OUTFORMAT_MARKDOWN, OUTFORMAT_TABLE, OUTFORMAT_BOX
// The Separator argument specifies the column separating string (default: '|').
// All lines are truncated at MaxLineLeength. A MaxLineLangth of '0' means no truncation.
// If this query result is of type RESULT_OK and SuppressOK is set to false, an "OK" string is written to the buffer, otherwise nothing is written to the buffer.
func (this *Result) DumpToWriter( Out *bufio.Writer, Format int, NoHeader bool, Separator string, NullValue string, NewLine string, MaxLineLength uint, SuppressOK bool ) ( int, error ) {
  if sep, err := GetDefaultSeparatorForOutputFormat( Format ); err != nil {
    return 0, err
  } else if strings.ToUpper( strings.TrimSpace( Separator ) ) == "<AUTO>" {
    Separator = sep
  }

  if strings.TrimSpace( NullValue ) == "" { NullValue = "NULL" }

  // fmt.Printf( "Type = %d\r\n", this.Type )

  switch {
  case this.IsOK():
    if SuppressOK {
      return 0, nil
    } else {
      return io.WriteString( Out, fmt.Sprintf( "OK%s", NewLine ) )
    }

  case this.IsNULL():
    return io.WriteString( Out, fmt.Sprintf( "%s%s", NullValue, NewLine ) )

  case this.IsError():
    return io.WriteString( Out, fmt.Sprintf( "%s%s", this.GetErrorAsString(), NewLine ) )

  case this.IsString(), this.IsInteger(), this.IsFloat(), this.IsJSON():
    return io.WriteString( Out, string( this.GetBuffer() ) + NewLine )
    return 0, nil

  case this.IsRowSet():
    var totalOutputLength int = 0

    if !NoHeader { // Render Table Header incl. new line
      if len, err := io.WriteString( Out, this.renderTableHeader( Format, Separator, NewLine, MaxLineLength ) ); err == nil {
        totalOutputLength += len
      } else {
        return len + totalOutputLength, err
      }
    }

    // Render Table Body incl. new line
    for row, err := this.GetFirstRow(); err == nil && row != nil; row = row.Next() {
      if len, err := row.DumpToWriter( Out, Format, Separator, NullValue, NewLine, MaxLineLength ); err == nil {
        totalOutputLength += len
      } else {
        return len + totalOutputLength, err
      }
    }

    if !NoHeader { // Render Table Footer
      if len, err := io.WriteString( Out, this.renderTableFooter( Format, NewLine, MaxLineLength ) ); err == nil {
        totalOutputLength += len
      } else {
        return len + totalOutputLength, err
      }
    }

    Out.Flush()
    return totalOutputLength, nil

  default:
    return 0, errors.New( "Unknown Output Format" )
  }
}

func GetOutputFormatFromString( Format string ) ( int, error ) {
  switch strings.ToUpper( strings.TrimSpace( Format ) ) {
  case "LIST":      return OUTFORMAT_LIST,      nil
  case "CSV":       return OUTFORMAT_CSV,       nil
  case "QUOTE":     return OUTFORMAT_QUOTE,     nil
  case "TABS":      return OUTFORMAT_TABS,      nil
  case "LINE":      return OUTFORMAT_LINE,      nil
  case "JSON":      return OUTFORMAT_JSON,      nil
  case "HTML":      return OUTFORMAT_HTML,      nil
  case "MARKDOWN":  return OUTFORMAT_MARKDOWN,  nil
  case "TABLE":     return OUTFORMAT_TABLE,     nil
  case "BOX":       return OUTFORMAT_BOX,       nil
  case "XML":       return OUTFORMAT_XML,       nil
  case "":          return -1,                  errors.New( "Missing output format" )
  default:          return -1,                  errors.New( "Unknown output format" )
  }
}

func GetDefaultSeparatorForOutputFormat( Format int ) ( string, error ) {
  switch Format {
  case OUTFORMAT_LIST, OUTFORMAT_MARKDOWN, OUTFORMAT_TABLE: return "|",   nil
  case OUTFORMAT_CSV, OUTFORMAT_QUOTE, OUTFORMAT_JSON:      return ",",   nil
  case OUTFORMAT_TABS:                                      return "\t",  nil
  case OUTFORMAT_LINE:                                      return "=",   nil
  case OUTFORMAT_HTML, OUTFORMAT_XML:                       return "",    nil
  case OUTFORMAT_BOX:                                       return "│",   nil
  default:                                                  return "",    errors.New( "Unknown output format" )
  }
}

//////
// is called from connection.Select
func( this *SQCloud ) readResult() ( *Result, error ) {
  ErrorResult := Result{
    value:                    Value{ Type: '-', Buffer: []byte( "100000 Unknown internal error" ) }, // This is an unset Value
    rows:                     nil,

    ColumnNames:              nil,
    ColumnWidth:              nil,
    MaxHeaderWidth:           0,

    uncompressedChuckSizeSum: 0,
  }
  result := ErrorResult

	
  var rowIndex      uint64 = 0

  for { // loop through all chunks

    if chunk, err := this.readNextRawChunk(); err != nil {
      ErrorResult.uncompressedChuckSizeSum = chunk.LEN
      ErrorResult.value.Buffer             = []byte( fmt.Sprintf( "100001 Internal Error: SQCloud.readNextRawChunk (%s)", err.Error() ) )
      return &ErrorResult, err

    } else {

      if err := chunk.Uncompress(); err != nil {
        ErrorResult.uncompressedChuckSizeSum = chunk.LEN
        ErrorResult.value.Buffer             = []byte( fmt.Sprintf( "100002 Internal Error: Chunk.Uncompress (%s)", err.Error() ) )
        return &ErrorResult, err

      } else {
        result.uncompressedChuckSizeSum += chunk.LEN

        switch Type := chunk.GetType(); Type {
        case '%':                     return nil, errors.New( "Nested compression" )

          // Values
        case '_':                     fallthrough // NULL
        case ':', ',':                fallthrough // INT, FLOAT
        case '+', '!', '$', '-', '#': fallthrough // String, C-String, BLOB, Error-String, JSON-String
        case '|':                     fallthrough // PSUB
        case '^':                     fallthrough // Command
        case '@':                                 // Reconnect
          result.value.Type = Type
          switch bytesRead, err := result.value.readBufferAt( chunk, 1 ); {
          case err != nil:            return nil, err
          case bytesRead == 0:        return nil, errors.New( "No Data" )
          case Type == '|':
            println( "Do the PSUB magic, open a second connection to the server and enter: " + result.value.GetString() )
                                      fallthrough
          default:                    return &result, nil
          }

          // RowSet
        case '/', '*':

          var offset        uint64 = 1 // skip the first type byte
          var bytesRead     uint64 = 0
          var LEN           uint64 = 0
          var IDX           uint64 = 1
          var NROWS         uint64 = 0
          var NCOLS         uint64 = 0

					// Detect end of Rowset Chunk directly without parsing...
					if Type == '/' {
						for _, pattern := range rowsetChunkEndPatterns {
							if chunk.RAW[ offset ] == pattern[ 0 ] && bytes.Equal( chunk.RAW[ offset : offset + uint64( len( pattern ) ) ], pattern ) { return &result, nil }
					} }
					
          if   LEN, bytesRead, err = chunk.readUInt64At( offset ); err != nil { return nil, err }
          offset += bytesRead

          if Type == '/' {
            if IDX, bytesRead, err = chunk.readUInt64At( offset ); err != nil { return nil, err }
            offset += bytesRead
          }

          if NROWS, bytesRead, err = chunk.readUInt64At( offset ); err != nil { return nil, err } // 0..rows-1
          offset += bytesRead

          if NCOLS, bytesRead, err = chunk.readUInt64At( offset ); err != nil { return nil, err } // 0..columns-1
          offset += bytesRead

          LEN = LEN + offset // check for overreading...

          if Type == '/' && NROWS == 0 && NCOLS == 0 { return &result, nil }

          if IDX == 1 {
            result.rows           = []ResultRow{}
            result.ColumnNames    = make( []string,     int( NCOLS ) )
            result.ColumnWidth    = make( []uint64,     int( NCOLS ) )
            result.MaxHeaderWidth = 0

            for column := uint64( 0 ); column < NCOLS; column++ { // Read in the column names, use the result.value as scratch variable
              switch val, bytesRead, err := chunk.readValueAt( offset ); {
              case err != nil:      return nil, err
              case !val.IsString(): return nil, errors.New( "Invalid Column name" )
              default:
                result.ColumnNames[ column ] = val.GetString()
                result.ColumnWidth[ column ] = val.GetLength()
                if result.MaxHeaderWidth < result.ColumnWidth[ column ] { result.MaxHeaderWidth = result.ColumnWidth[ column ] }
                offset += bytesRead
              }
            }
          }

          // read all the rows from this chunk
          rows := make( []ResultRow, int( NROWS ) )
          for row := uint64( 0 ); row < NROWS; row++ {

            rows[ row ].result  = &result
            rows[ row ].index   = rowIndex
            rows[ row ].columns = make( []Value, int( NCOLS ) )

            rowIndex++

            for column := uint64( 0 ); column < NCOLS; column++ {
              switch rows[ row ].columns[ column ], bytesRead, err = chunk.readValueAt( offset ); {
              case err != nil: return nil, err
              default:
                columnLength := rows[ row ].columns[ column ].GetLength()
                if result.ColumnWidth[ column ] < columnLength { result.ColumnWidth[ column ] = columnLength }
                if result.MaxHeaderWidth        < columnLength { result.MaxHeaderWidth        = columnLength }
                offset += bytesRead
              }
            }
          }

          result.rows = append( result.rows, rows... )

          result.value.Type     = '*'
          result.value.Buffer   = nil

          if Type == '*' { return &result, nil } // return if it is a rowset
          this.sendString( "OK" )                // ask the server for the next chunk and loop (Thank's Andrea)

        case '{':
          result.value.Type = '#' // translate JSON Type to uniform '#'
          result.value.Buffer = chunk.GetData()
          return &result, nil

        default:
          return nil, errors.New( "Unknown response type" )
        }
      }
    }
  }
}