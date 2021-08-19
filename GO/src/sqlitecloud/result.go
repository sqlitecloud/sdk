//
//                    ////              SQLite Cloud
//        ////////////  ///             
//      ///             ///  ///        Product     : SQLite Cloud GO SDK
//     ///             ///  ///         Version     : 0.0.2
//     //             ///   ///  ///    Date        : 2021/08/18
//    ///             ///   ///  ///    Author      : Andreas Pfeil
//   ///             ///   ///  ///     
//   ///     //////////   ///  ///      Description : GO Methods related to the 
//   ////                ///  ///                     SQCloudResult class.
//     ////     //////////   ///        
//        ////            ////          
//          ////     /////              
//             ///                      Copyright   : 2021 by SQLite Cloud Inc.
//
// -----------------------------------------------------------------------TAB=2

package sqlitecloud

// #include <stdlib.h>
// #include "../../../C/sqcloud.h"
import "C"

import "fmt"
import "os"

import "bufio"
import "strings"
import "errors"
import "time"
import "io"
import "strconv"
import "encoding/json"
import "golang.org/x/term"

// SQCloudResType
const RESULT_OK           = C.RESULT_OK
const RESULT_ERROR        = C.RESULT_ERROR
const RESULT_STRING       = C.RESULT_STRING
const RESULT_INTEGER      = C.RESULT_INTEGER
const RESULT_FLOAT        = C.RESULT_FLOAT
const RESULT_ROWSET       = C.RESULT_ROWSET
const RESULT_NULL         = C.RESULT_NULL
const RESULT_JSON         = C.RESULT_JSON

// SQCloudValueType
const VALUE_INTEGER       = C.VALUE_INTEGER
const VALUE_FLOAT         = C.VALUE_FLOAT
const VALUE_TEXT          = C.VALUE_TEXT
const VALUE_BLOB          = C.VALUE_BLOB
const VALUE_NULL          = C.VALUE_NULL

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

type SQCloudResult struct {
  result *C.struct_SQCloudResult

  Rows            uint
  Columns         uint

  ColumnWidth     []uint
  HeaderWidth     []uint
  MaxHeaderWidth  uint

  Type            uint
  ErrorCode       int
  ErrorMessage    string
}

// ResultSet Methods (100% GO)

// GetType returns the type of this query result as an integer (see: RESULT_ constants).
func (this *SQCloudResult ) GetType() uint {
  return this.Type
}

// IsOK returns true if this query result if of type "RESULT_OK", false otherwise.
func (this *SQCloudResult ) IsOK() bool {
  return this.Type == RESULT_OK
}

// GetNumberOfRows returns the number of rows in this query result
func (this *SQCloudResult ) GetNumberOfRows() uint {
  return this.Rows
}

// GetNumberOfColumns returns the number of columns in this query result
func (this *SQCloudResult ) GetNumberOfColumns() uint {
  return this.Columns
}

// Dump outputs this query result to the screen.
// Warning: No line truncation is used. If you want to truncation the output to a certain width, use: SQCloudResult.DumpToScreen( width )
func (this *SQCloudResult ) Dump() {
  w := 0
  if width, _, err := term.GetSize( 0 ); err == nil { w = width }
  this.DumpToScreen( uint( w ) )
}

// ToJSON returns a JSON representation of this query result.
// BUG(andreas): The SQCloudResult.ToJSON method is not implemented yet.
func (this *SQCloudResult ) ToJSON() string {
  return "todo" // Use Writer into Buffer
}

// Additional ResultSet Methods (100% GO)

// GetMaxColumnLength returns the number of runes of the value in the specified column with the maximum length in this query result.
// The Column index is an unsigned int in the range of 0...GetNumberOfColumns() - 1.
func (this *SQCloudResult ) GetMaxColumnLength( Column uint ) uint {
  return this.ColumnWidth[ Column ]
}
// GetNameWidth returns the number of runes of the column name in the specified column.
// The Column index is an unsigned int in the range of 0...GetNumberOfColumns() - 1.
func (this *SQCloudResult ) GetNameWidth( Column uint ) uint {
  return this.HeaderWidth[ Column ]
}
// GetMaxNameWidth returns the number of runes of the longest column name.
func (this *SQCloudResult ) GetMaxNameWidth() uint {
  return this.MaxHeaderWidth
}

// Additional Data Access Functions (100% GO)

// IsError returns true if this query result is of type "RESULT_ERROR", false otherwise.
func (this *SQCloudResult ) IsError() bool {
  return this.Type == RESULT_ERROR
}

// IsNull returns true if this query result is of type "RESULT_NULL", false otherwise.
func (this *SQCloudResult ) IsNull() bool {
  return this.Type == RESULT_NULL
}

// IsJson returns true if this query result is of type "RESULT_JSON", false otherwise.
func (this *SQCloudResult ) IsJson() bool {
  return this.Type == RESULT_JSON
}

// IsString returns true if this query result is of type "RESULT_STRING", false otherwise.
func (this *SQCloudResult ) IsString() bool {
  return this.Type == RESULT_STRING
}

// IsInteger returns true if this query result is of type "RESULT_INTEGER", false otherwise.
func (this *SQCloudResult ) IsInteger() bool {
  return this.Type == RESULT_INTEGER
}

// IsFloat returns true if this query result is of type "RESULT_FLOAT", false otherwise.
func (this *SQCloudResult ) IsFloat() bool {
  return this.Type == RESULT_FLOAT
}

// IsRowSet returns true if this query result is of type "RESULT_ROWSET", false otherwise.
func (this *SQCloudResult ) IsRowSet() bool {
  return this.Type == RESULT_ROWSET
}

// IsTextual returns true if this query result is of type "RESULT_JSON", "RESULT_STRING", "RESULT_INTEGER" or "RESULT_FLOAT", false otherwise.
func (this *SQCloudResult ) IsTextual() bool {
  return this.IsJson() || this.IsString() || this.IsInteger() || this.IsFloat()
}

// Additional ResultSet Methods (100% GO)

// GetSQLDateTime parses this query result value in Row and Column as an SQL-DateTime and returns its value.
// The Row index is an unsigned int in the range of 0...GetNumberOfRows() - 1.
// The Column index is an unsigned int in the range of 0...GetNumberOfColumns() - 1.
func (this *SQCloudResult ) GetSQLDateTime( Row uint, Column uint ) time.Time {
  datetime, _ := time.Parse( "2006-01-02 15:04:05", this.CGetStringValue( Row, Column ) )
  return datetime
} 

// ResultSet Methods (C SDK)

// GetBuffer returns the buffer of this query result as string.
func (this *SQCloudResult ) GetBuffer() string {
  return this.CGetResultBuffer()
}

func (this *SQCloudResult ) GetBufferAsString() ( string, error ) {
  if this.IsString() {
    return this.GetBuffer(), nil
  }
  return "", errors.New( "Result is not a string" )
}

func (this *SQCloudResult ) GetBufferAsJSON() ( object interface{}, err error ) {
  if this.IsJson() {
    err = json.Unmarshal( []byte( this.GetBuffer() ), object)
    return
  }
  return "", errors.New( "Result is not a JSON object" )
}

func (this *SQCloudResult ) GetBufferAsInt32() ( int32, error ) {
  if this.IsInteger() {
    value64, err := strconv.ParseInt( this.GetBuffer(), 0, 32 )
    return int32( value64 ), err
  }
  return 0, errors.New( "Result is not an integer number" )
}
func (this *SQCloudResult ) GetBufferAsInt64() ( int64, error ) {
  if this.IsInteger() {
    return strconv.ParseInt( this.GetBuffer(), 0, 64 )
  }
  return 0, errors.New( "Result is not an integer number" )
}

func (this *SQCloudResult ) GetBufferAsFloat32() ( float32, error ) {
  if this.IsFloat() {
    value64, err := strconv.ParseFloat( this.GetBuffer(), 32 )
    return float32( value64 ), err
  }
  return 0, errors.New( "Result is not a float number" )
}
func (this *SQCloudResult ) GetBufferAsFloat64() ( float64, error ) {
  if this.IsFloat() {
    return strconv.ParseFloat( this.GetBuffer(), 64 )
  }
  return 0, errors.New( "Result is not a float number" )
}

// GetLength returns the length of the buffer of this query result.
func (this *SQCloudResult ) GetLength() uint {
  return this.CGetResultLen()
}

// GetLength returns the maximum length of the buffer of this query result.
// BUG(andreas): What is this GetMaxLength for?
func (this *SQCloudResult ) GetMaxLength() uint32 {
  return this.CGetMaxLen()
}

// Free frees all memory allocated by this query result.
func (this *SQCloudResult ) Free() {
  this.CFree()
  this.result         = nil
  this.Rows           = 0
  this.Columns        = 0
  this.ColumnWidth    = []uint{}
  this.HeaderWidth    = []uint{}
  this.MaxHeaderWidth = 0
  this.Type           = 0
  this.ErrorCode      = 0
  this.ErrorMessage   = ""
}

// GetValueType returns the type of the value in row Row and column Column of this query result.
// The Row index is an unsigned int in the range of 0...GetNumberOfRows() - 1.
// The Column index is an unsigned int in the range of 0...GetNumberOfColumns() - 1.
// Possible return types are: VALUE_INTEGER, VALUE_FLOAT, VALUE_TEXT, VALUE_BLOB, VALUE_NULL
func (this *SQCloudResult ) GetValueType( Row uint, Column uint ) int {
  return this.CGetValueType( Row, Column )
}

// GetColumnName returns the column name in column Column of this query result.
// The Column index is an unsigned int in the range of 0...GetNumberOfColumns() - 1.
func (this *SQCloudResult ) GetColumnName( Column uint ) string {
  return this.CGetColumnName( Column )
}

// GetStringValue returns the contents in row Row and column Column of this query result as string.
// The Row index is an unsigned int in the range of 0...GetNumberOfRows() - 1.
// The Column index is an unsigned int in the range of 0...GetNumberOfColumns() - 1.
func (this *SQCloudResult ) GetStringValue( Row uint, Column uint ) string {
  return this.CGetStringValue( Row, Column )
}

// GetInt32Value returns the contents in row Row and column Column of this query result as int32.
// The Row index is an unsigned int in the range of 0...GetNumberOfRows() - 1.
// The Column index is an unsigned int in the range of 0...GetNumberOfColumns() - 1.
func (this *SQCloudResult ) GetInt32Value( Row uint, Column uint ) int32 {
  return this.CGetInt32Value( Row, Column )
}

// GetInt64Value returns the contents in row Row and column Column of this query result as int64.
// The Row index is an unsigned int in the range of 0...GetNumberOfRows() - 1.
// The Column index is an unsigned int in the range of 0...GetNumberOfColumns() - 1.
func (this *SQCloudResult ) GetInt64Value( Row uint, Column uint ) int64 {
  return this.CGetInt64Value( Row, Column )
}

// GetFloat32Value returns the contents in row Row and column Column of this query result as float32.
// The Row index is an unsigned int in the range of 0...GetNumberOfRows() - 1.
// The Column index is an unsigned int in the range of 0...GetNumberOfColumns() - 1.
func (this *SQCloudResult ) GetFloat32Value( Row uint, Column uint ) float32 {
  return this.CGetFloat32Value( Row, Column )
}

// GetFloat64Value returns the contents in row Row and column Column of this query result as float64.
// The Row index is an unsigned int in the range of 0...GetNumberOfRows() - 1.
// The Column index is an unsigned int in the range of 0...GetNumberOfColumns() - 1.
func (this *SQCloudResult ) GetFloat64Value( Row uint, Column uint ) float64 {
  return this.CGetFloat64Value( Row, Column )
}

// DumpToScreen outputs this query result to the screen.
// The output is truncated at a maximum line width of MaxLineLength runes (compare: SQCloudResult.Dump())
func (this *SQCloudResult ) DumpToScreen( MaxLineLength uint ) {
  this.DumpToWriter( bufio.NewWriter( os.Stdout ), OUTFORMAT_BOX, false, "│", "NULL", "\r\n", MaxLineLength, false )
}

////// Row Methods (100% GO)

// GetRow returns a pointer to the row Row of this query result.
// The Row index is an unsigned int in the range of 0...GetNumberOfRows() - 1.
// If the index row can not be found, nil is returned instead.
func (this *SQCloudResult ) GetRow( Row uint ) ( *SQCloudResultRow ) {
  if Row >= this.Rows {
    return nil
  }
  row := SQCloudResultRow{
    result  : this,
    row     : Row,
    rows    : this.Rows,
    columns : this.Columns,
  }
  return &row
}

// GetFirstRow returns the first row of this query result.
// If this query result has no row's, nil is returned instead.
func (this *SQCloudResult ) GetFirstRow() *SQCloudResultRow {
  return this.GetRow( 0 )
}

// GetLastRow returns the first row of this query result.
// If this query result has no row's, nil is returned instead.
func (this *SQCloudResult ) GetLastRow() *SQCloudResultRow {
  switch this.Rows {
  case 0:  return nil
  default: return this.GetRow( this.Rows - 1 )
  }
}

/// Dump Method (100% GO)

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

func (this *SQCloudResult) renderHorizontalTableLine( Left string, Fill string, Separator string, Right string ) string {
  outBuffer := ""
  for _, columnWidth := range this.ColumnWidth {
    outBuffer = fmt.Sprintf( "%s%s%s", outBuffer, strings.Repeat( Fill, int( columnWidth + 2 ) ), Separator )
  }
  return fmt.Sprintf( "%s%s%s", Left, strings.TrimRight( outBuffer, Separator ), Right )
}
func (this *SQCloudResult) renderTableColumnNames( Left string, Separator string, Right string ) string {
  outBuffer := ""
  for forThisColumn, columnWidth := range this.ColumnWidth {
    outBuffer = fmt.Sprintf( "%s%s%s", outBuffer, renderCenteredString( this.GetColumnName( uint( forThisColumn ) ), int( columnWidth + 2 ) ), Separator )
  }
  return fmt.Sprintf( "%s%s%s", Left, strings.TrimRight( outBuffer, Separator ), Right )
}
func (this *SQCloudResult) renderTableHeader( Format int, Separator string, NewLine string, MaxLineLength uint ) string {
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
func (this *SQCloudResult) renderTableFooter( Format int, NewLine string, MaxLineLength uint ) string {
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
func (this *SQCloudResult) DumpToWriter( Out *bufio.Writer, Format int, NoHeader bool, Separator string, NullValue string, NewLine string, MaxLineLength uint, SuppressOK bool ) ( int, error ) {
  if sep, err := GetDefaultSeparatorForOutputFormat( Format ); err != nil {
    return 0, err
  } else if strings.ToUpper( strings.TrimSpace( Separator ) ) == "<AUTO>" {
    Separator = sep
  }

  if strings.TrimSpace( NullValue ) == "" { NullValue = "NULL" }

  // fmt.Printf( "Type = %d\r\n", this.Type )


  switch this.Type {
  case RESULT_OK:
    if SuppressOK {
      return 0, nil
    } else {
      return io.WriteString( Out, fmt.Sprintf( "OK%s", NewLine ) )
    }
    
  case RESULT_NULL:
    return io.WriteString( Out, fmt.Sprintf( "%s%s", NullValue, NewLine ) )

  case RESULT_ERROR:
    return io.WriteString( Out, fmt.Sprintf( "ERROR: %s (%d)%s", this.ErrorMessage, this.ErrorCode, NewLine ) )

  case RESULT_STRING, RESULT_INTEGER, RESULT_FLOAT, RESULT_JSON:
    return io.WriteString( Out, this.CGetResultBuffer() + NewLine )

  case RESULT_ROWSET: 
    var totalOutputLength int = 0

    if !NoHeader { // Render Table Header incl. new line
      if len, err := io.WriteString( Out, this.renderTableHeader( Format, Separator, NewLine, MaxLineLength ) ); err == nil {
        totalOutputLength += len
      } else {
        return len + totalOutputLength, err
      }
    }

    // Render Table Body incl. new line
    for row := this.GetFirstRow(); row != nil; row = row.Next() {
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