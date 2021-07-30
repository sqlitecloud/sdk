package sqlitecloud

// #include <stdlib.h>
// #include "sqcloud.h"
import "C"

import "fmt"
//import "os"
//import "bufio"
import "strings"
import "errors"
import "time"
import "io"
//import "strconv"

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

type SQCloudResult struct {
  result *C.struct_SQCloudResult

  Rows            uint    // Must be set during ...
  Columns         uint    // Must be set

  ColumnWidth     []uint  // Must be set
  HeaderWidth     []uint  // Must be set
  MaxHeaderWidth  uint    // Must be set

  Type            uint    // Must be set during Initialisation
  ErrorCode       int
  ErrorMessage    string
}

// ResultSet Methods (100% GO)

func (this *SQCloudResult ) GetType() uint {
  return this.Type
}
func (this *SQCloudResult ) IsOK() bool {
  return this.Type == RESULT_OK
}
func (this *SQCloudResult ) GetNumberOfRows() uint {
  return this.Rows
}
func (this *SQCloudResult ) GetNumberOfColumns() uint {
  return this.Columns
}
func (this *SQCloudResult ) Dump() {
  this.DumpToScreen( 0 )
}

func (this *SQCloudResult ) ToJSON() string {
  return "todo" // Use Writer into Buffer
}

// Additional ResultSet Methods (100% GO)

func (this *SQCloudResult ) GetMaxColumnLength( Column uint ) uint {
  return this.ColumnWidth[ Column ]
}
func (this *SQCloudResult ) GetNameWidth( Column uint ) uint {
  return this.HeaderWidth[ Column ]
}
func (this *SQCloudResult ) GetMaxNameWidth() uint {
  return this.MaxHeaderWidth
}

// Additional Data Access Functions (100% GO)

func (this *SQCloudResult ) IsError() bool {
  return this.Type == RESULT_ERROR
}
func (this *SQCloudResult ) IsNull() bool {
  return this.Type == RESULT_NULL
}
func (this *SQCloudResult ) IsJson() bool {
  return this.Type == RESULT_JSON
}
func (this *SQCloudResult ) IsString() bool {
  return this.Type == RESULT_STRING
}
func (this *SQCloudResult ) IsInteger() bool {
  return this.Type == RESULT_INTEGER
}
func (this *SQCloudResult ) IsFloat() bool {
  return this.Type == RESULT_FLOAT
}
func (this *SQCloudResult ) IsRowSet() bool {
  return this.Type == RESULT_ROWSET
}
func (this *SQCloudResult ) IsTextual() bool {
  return this.IsJson() || this.IsString() || this.IsInteger() || this.IsFloat()
}

// Additional ResultSet Methods (100% GO)

func (this *SQCloudResult ) GetSQLDateTime( Row uint, Column uint ) time.Time {
  datetime, _ := time.Parse( "2006-01-02 15:04:05", this.CGetStringValue( Row, Column ) )
  return datetime
} 

// ResultSet Methods (C SDK)

func (this *SQCloudResult ) GetBuffer() string {
  return this.CGetResultBuffer()
}
func (this *SQCloudResult ) GetLength() uint {
  return this.CGetResultLen()
}
func (this *SQCloudResult ) GetMaxLength() uint32 {
  return this.CGetMaxLen()
}
func (this *SQCloudResult ) Free() {
  this.CFree()
}
func (this *SQCloudResult ) GetValueType( Row uint, Column uint ) int {
  return this.CGetValueType( Row, Column )
}
func (this *SQCloudResult ) GetColumnName( Column uint ) string {
  return this.CGetColumnName( Column )
}
func (this *SQCloudResult ) GetStringValue( Row uint, Column uint ) string {
  return this.CGetStringValue( Row, Column )
} 
func (this *SQCloudResult ) GetInt32Value( Row uint, Column uint ) int32 {
  return this.CGetInt32Value( Row, Column )
} 
func (this *SQCloudResult ) GetInt64Value( Row uint, Column uint ) int64 {
  return this.CGetInt64Value( Row, Column )
} 
func (this *SQCloudResult ) GetFloat32Value( Row uint, Column uint ) float32 {
  return this.CGetFloat32Value( Row, Column )
} 
func (this *SQCloudResult ) GetFloat64Value( Row uint, Column uint ) float64 {
  return this.CGetFloat64Value( Row, Column )
}
func (this *SQCloudResult ) DumpToScreen( MaxLineLength uint ) {
  this.CDump( MaxLineLength )
}

////// Row Methods (100% GO)

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
func (this *SQCloudResult ) GetFirstRow() *SQCloudResultRow {
  return this.GetRow( 0 )
}
func (this *SQCloudResult ) GetLastRow() *SQCloudResultRow {
  switch this.Rows {
  case 0:  return nil
  default: return this.GetRow( this.Rows - 1 )
  }
}

/// Dump Method (100% GO)

func trimStringToMaxLength( Buffer string, MaxLineLength uint ) string {
  switch MaxLineLength {
  case 0: return Buffer
  default:
    len := uint( len( Buffer ) )
    if len > MaxLineLength {
      return fmt.Sprintf( fmt.Sprintf( "%%%ds…", MaxLineLength - 1 ), Buffer )
    }
    return Buffer
  }
}

func renderCenteredString( Buffer string, Width int ) string {
  return fmt.Sprintf( "%[1]*s", -Width, fmt.Sprintf( "%[1]*s", ( Width + len( Buffer ) ) / 2, Buffer ) )
}

func (this *SQCloudResult) renderHorizontalTableLine( Left string, Fill string, Seperator string, Right string, MaxLineLength uint ) string {
  outBuffer := ""
  for _, columnWidth := range this.ColumnWidth {
    outBuffer = fmt.Sprintf( "%s%s%s", outBuffer, strings.Repeat( Fill, int( columnWidth ) ), Seperator )
  }
  outBuffer = fmt.Sprintf( "%s%s%s", Left, strings.TrimRight( outBuffer, Seperator ), Right )
  return trimStringToMaxLength( outBuffer, MaxLineLength )
}
func (this *SQCloudResult) renderTableColumnNames( Left string, Seperator string, Right string, MaxLineLength uint ) string {
  outBuffer := ""
  for forThisColumn, columnWidth := range this.ColumnWidth {
    outBuffer = fmt.Sprintf( "%s%s%s", outBuffer, renderCenteredString( this.GetColumnName( uint( forThisColumn ) ), int( columnWidth ) ), Seperator )
  }
  outBuffer = fmt.Sprintf( "%s%s%s", Left, strings.TrimRight( outBuffer, Seperator ), Right )
  return trimStringToMaxLength( outBuffer, MaxLineLength )
}
func (this *SQCloudResult) renderTableHeader( Format int, MaxLineLength uint ) string {
  switch( Format ) {
    case OUTFORMAT_JSON: return "["

    case OUTFORMAT_MARKDOWN:
      return this.renderTableColumnNames( "|", "|", "|", MaxLineLength )          + "\r\n" +
             this.renderHorizontalTableLine( "|", "-", "|", "|", MaxLineLength )  + "\r\n"

    case OUTFORMAT_TABLE:
      tableLine := this.renderHorizontalTableLine( "+", "-", "+", "+", MaxLineLength )
      return  tableLine                                                           + "\r\n" +
              this.renderTableColumnNames( "|", "|", "|", MaxLineLength )         + "\r\n" +
              tableLine                                                           + "\r\n"

    case OUTFORMAT_BOX:
      return this.renderHorizontalTableLine( "┌", "─", "┬", "┐", MaxLineLength )  + "\r\n" +
             this.renderTableColumnNames( "│", "│", "│", MaxLineLength )          + "\r\n" +
             this.renderHorizontalTableLine( "├", "─", "┼", "┤", MaxLineLength )  + "\r\n"
    default:
      return "" // No header
  }
  return ""
}
func (this *SQCloudResult) renderTableFooter( Format int, MaxLineLength uint ) string {
  switch( Format ) {
    case OUTFORMAT_JSON: return "]"

    case OUTFORMAT_TABLE:
      return this.renderHorizontalTableLine( "+", "-", "+", "+", MaxLineLength )

    case OUTFORMAT_BOX:
      return this.renderHorizontalTableLine( "└", "─", "┴", "┘", MaxLineLength )

    default:
      return "" // No footer
  }
}

func (this *SQCloudResult) DumpToWriter( Out io.Writer, Format int, Seperator string, MaxLineLength uint, SurpressOK bool ) ( int, error ) {
  switch this.Type {
  case RESULT_OK:
    if SurpressOK {
      return 0, nil
    } else {
      return io.WriteString( Out, "OK\r\n" )
    }
    
  case RESULT_NULL:
    return io.WriteString( Out, "NULL\r\n" )

  case RESULT_ERROR:
    return io.WriteString( Out, fmt.Sprintf( "ERROR: %s (%d)\r\n", this.ErrorMessage, this.ErrorCode ) )

  case RESULT_STRING, RESULT_INTEGER, RESULT_FLOAT, RESULT_JSON:
    return io.WriteString( Out, this.CGetResultBuffer() + "\r\n")

  case RESULT_ROWSET:
    var totalOutputLength int = 0

    // Render Table Header
    if len, err := io.WriteString( Out, this.renderTableHeader( Format, MaxLineLength ) ); err == nil {
      totalOutputLength += len

      // Render Table Body
      for row := this.GetFirstRow(); !row.IsEOF(); row = row.Next() {
        if len, err := row.DumpToWriter( Out, Format, Seperator, MaxLineLength ); err == nil {
          totalOutputLength += len
        } else {
          return len + totalOutputLength, err
        }
      }

      // Render Table Footer
      if len, err := io.WriteString( Out, this.renderTableFooter( Format, MaxLineLength ) ); err == nil {
        totalOutputLength += len
        return totalOutputLength, nil
      } else {
        return len + totalOutputLength, err
      }

    } else {
      return len + totalOutputLength, err
    }
    
  default:
    return 0, errors.New( "Unknown Output Format" )
  }
}