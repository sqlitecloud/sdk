//
//                    ////              SQLite Cloud
//        ////////////  ///             
//      ///             ///  ///        Product     : SQLite Cloud GO SDK
//     ///             ///  ///         Version     : 0.0.1
//     //             ///   ///  ///    Date        : 2021/08/13
//    ///             ///   ///  ///    Author      : Andreas Pfeil
//   ///             ///   ///  ///     
//   ///     //////////   ///  ///      Description : GO methods related to the
//   ////                ///  ///                     SQCloudResultRow class.
//     ////     //////////   ///        
//        ////            ////          
//          ////     /////              
//             ///                      Copyright   : 2021 by SQLite Cloud Inc.
//
// -----------------------------------------------------------------------TAB=2

package sqlitecloud

import "fmt"
//import "os"
import "io"
//import "bufio"
import "strings"
//import "errors"
import "time"
//import "strconv"
import "html"

type SQCloudResultRow struct {
  result  *SQCloudResult
  row     uint // 0, 1, ... rows-1
  rows    uint 
  columns uint
}

// ToJSON returns a JSON representation of this query result row.
// BUG(andreas): The SQCloudResultRow.ToJSON method is not implemented yet.
func (this *SQCloudResultRow ) ToJSON() string {
  return "todo" // Use Writer into Buffer
}

// IsFirst returns true if this query result row is the first in the result set, false otherwise.
func (this *SQCloudResultRow ) IsFirst() bool {
  return this.row == 0
}

// IsLast returns true if this query result row is the last in the result set, false otherwise.
func (this *SQCloudResultRow ) IsLast() bool {
  return this.row == this.rows - 1
}

// IsEOF returns false if this query result row is in the result set, true otherwise.
func (this *SQCloudResultRow ) IsEOF() bool {
  return this.row >= this.rows
}

// Rewind resets the iterator and returns the first row in this query result. 
func (this *SQCloudResultRow ) Rewind() *SQCloudResultRow {
  this.row = 0
  return this
}

// Next fetches the next row in this query result and returns it, otherwise if there is no next row, nil is returned.
func (this *SQCloudResultRow ) Next() *SQCloudResultRow {
  if this.row < this.rows - 1 {
    this.row++
    return this
  } 
  return nil
}

// GetType returns the type of the value in column Column of this query result row.
// The Column index is an unsigned int in the range of 0...GetNumberOfColumns() - 1.
// Possible return types are: VALUE_INTEGER, VALUE_FLOAT, VALUE_TEXT, VALUE_BLOB, VALUE_NULL
func (this *SQCloudResultRow ) GetType( Column uint ) int {
  return this.result.GetValueType( this.row, Column )
}

// IsInteger returns true if this query result row column Column is of type "VALUE_INTEGER", false otherwise.
// The Column index is an unsigned int in the range of 0...GetNumberOfColumns() - 1.
func (this *SQCloudResultRow ) IsInteger( Column uint ) bool {
  return this.GetType( Column ) == RESULT_INTEGER
}

// IsFloat returns true if this query result row column Column is of type "VALUE_FLOAT", false otherwise.
// The Column index is an unsigned int in the range of 0...GetNumberOfColumns() - 1.
func (this *SQCloudResultRow ) IsFloat( Column uint ) bool {
  return this.GetType( Column ) == VALUE_FLOAT
}

// IsText returns true if this query result row column Column is of type "VALUE_TEXT", false otherwise.
// The Column index is an unsigned int in the range of 0...GetNumberOfColumns() - 1.
func (this *SQCloudResultRow ) IsText( Column uint ) bool {
  return this.GetType( Column ) == VALUE_TEXT
}

// IsBlob returns true if this query result row column Column is of type "VALUE_BLOB", false otherwise.
// The Column index is an unsigned int in the range of 0...GetNumberOfColumns() - 1.
func (this *SQCloudResultRow ) IsBlob( Column uint ) bool {
  return this.GetType( Column ) == VALUE_BLOB
}

// IsNull returns true if this query result row column Column is of type "VALUE_NULL", false otherwise.
// The Column index is an unsigned int in the range of 0...GetNumberOfColumns() - 1.
func (this *SQCloudResultRow ) IsNull( Column uint ) bool {
  return this.GetType( Column ) == RESULT_NULL
}

// IsTextual returns true if this query result row column Column is of type "VALUE_NULL" or "VALUE_BLOB", false otherwise.
// The Column index is an unsigned int in the range of 0...GetNumberOfColumns() - 1.
func (this *SQCloudResultRow ) IsTextual( Column uint ) bool {
  return this.IsText( Column ) || this.IsBlob( Column )
}

// GetMaxNameWidth returns the number of runes of the longest column name.
func (this *SQCloudResultRow ) GetMaxNameWidth() uint {
  return this.result.GetMaxNameWidth()
}

// GetNameWidth returns the number of runes of the column name in the specified column.
// The Column index is an unsigned int in the range of 0...GetNumberOfColumns() - 1.
func (this *SQCloudResultRow ) GetNameWidth( Column uint ) uint {
  return this.result.GetNameWidth( Column )
}

// GetName returns the column name in column Column of this query result.
// The Column index is an unsigned int in the range of 0...GetNumberOfColumns() - 1.
func (this *SQCloudResultRow ) GetName( Column uint ) string {
  return this.result.GetColumnName( Column )
}

// GetWidth returns the number of runes of the value in the specified column with the maximum length in this query result.
// The Column index is an unsigned int in the range of 0...GetNumberOfColumns() - 1.
// BUG(andreas): Rename GetWidth->GetmaxWidth
func (this *SQCloudResultRow ) GetWidth( Column uint ) uint {
  return this.result.GetMaxColumnLength( Column )
}

// GetStringValue returns the contents in column Column of this query result row as string.
// The Column index is an unsigned int in the range of 0...GetNumberOfColumns() - 1.
func (this *SQCloudResultRow ) GetStringValue( Column uint ) string {
  return this.result.GetStringValue( this.row, Column )
}

// GetInt32Value returns the contents in column Column of this query result row as int32.
// The Column index is an unsigned int in the range of 0...GetNumberOfColumns() - 1.
func (this *SQCloudResultRow ) GetInt32Value( Column uint ) int32 {
  return this.result.GetInt32Value( this.row, Column )
}

// GetInt64Value returns the contents in column Column of this query result row as int64.
// The Column index is an unsigned int in the range of 0...GetNumberOfColumns() - 1.
func (this *SQCloudResultRow ) GetInt64Value( Column uint ) int64 {
  return this.result.GetInt64Value( this.row, Column )
}

// GetFloat32Value returns the contents in column Column of this query result row as float32.
// The Column index is an unsigned int in the range of 0...GetNumberOfColumns() - 1.
func (this *SQCloudResultRow ) GetFloat32Value( Column uint ) float32 {
  return this.result.GetFloat32Value( this.row, Column )
}

// GetFloat64Value returns the contents in column Column of this query result row as float64.
// The Column index is an unsigned int in the range of 0...GetNumberOfColumns() - 1.
func (this *SQCloudResultRow ) GetFloat64Value( Column uint ) float64 {
  return this.result.GetFloat64Value( this.row, Column )
}

// GetSQLDateTime parses this query result value in column Column as an SQL-DateTime and returns its value.
// The Column index is an unsigned int in the range of 0...GetNumberOfColumns() - 1.
func (this *SQCloudResultRow ) GetSQLDateTime( Column uint ) time.Time {
  return this.result.GetSQLDateTime( this.row, Column )
} 


func (this *SQCloudResultRow) renderLine( Format int, Seperator string, MaxLineLength uint ) string {
  buffer := ""
  for forThisColumn := uint( 0 ); forThisColumn < this.columns; forThisColumn++ {
    switch Format {
    case OUTFORMAT_LIST, OUTFORMAT_TABS:                      buffer += fmt.Sprintf( "%s%s", this.GetStringValue( forThisColumn ), Seperator )
    case OUTFORMAT_MARKDOWN, OUTFORMAT_TABLE, OUTFORMAT_BOX:  buffer += fmt.Sprintf( fmt.Sprintf( " %%-%ds %s", this.GetWidth( forThisColumn ), Seperator ), this.GetStringValue( forThisColumn ) ) 
    case OUTFORMAT_CSV:                                       buffer += fmt.Sprintf( "%s,", SQCloudEnquoteString( this.GetStringValue( forThisColumn ) ) )
    case OUTFORMAT_LINE:                                      buffer += fmt.Sprintf( fmt.Sprintf( "%%%ds = %%s\r\n", this.result.MaxHeaderWidth ), this.GetName( forThisColumn ), this.GetStringValue( forThisColumn ) )
    case OUTFORMAT_HTML:                                      buffer += fmt.Sprintf( "  <TD>%s</TD>\r\n", html.EscapeString( this.GetStringValue( forThisColumn ) ) )
    case OUTFORMAT_XML:                                       buffer += fmt.Sprintf( "    <field name=\"%s\">%s</field>\r\n", this.GetName( forThisColumn ), html.EscapeString( this.GetStringValue( forThisColumn ) ) )
    }
  }
  return trimStringToMaxLength( strings.TrimRight( buffer, Seperator ), MaxLineLength )
}

// DumpToWriter renders this query result row into the buffer of an io.Writer.
// The output Format can be specified and must be one of the following values: OUTFORMAT_LIST, OUTFORMAT_CSV, OUTFORMAT_QUOTE, OUTFORMAT_TABS, OUTFORMAT_LINE, OUTFORMAT_JSON, OUTFORMAT_HTML, OUTFORMAT_MARKDOWN, OUTFORMAT_TABLE, OUTFORMAT_BOX
// The Separator argument specifies the column separating string (default: '|'). 
// All lines are truncated at MaxLineLeength. A MaxLineLangth of '0' means no truncation. 
func (this *SQCloudResultRow) DumpToWriter( Out io.Writer, Format int, Seperator string, MaxLineLength uint ) ( int, error ) {
  buffer := ""
  
  switch( Format ) {
  case OUTFORMAT_LIST:
    buffer = this.renderLine( Format, Seperator, MaxLineLength ) + "\r\n"

  case OUTFORMAT_CSV:
    buffer = this.renderLine( Format, ",", MaxLineLength ) + "\r\n"

  case OUTFORMAT_QUOTE: 
    for forThisColumn := uint( 0 ); forThisColumn < this.columns; forThisColumn++ {
      switch this.GetType( forThisColumn ) {
      case VALUE_TEXT, VALUE_BLOB:  buffer += fmt.Sprintf( "'%s',", strings.Replace( this.GetStringValue( forThisColumn ), "'", "\\'", -1 ) )
      default:                      buffer += fmt.Sprintf( "%s,", this.GetStringValue( forThisColumn ))
      }
    }
    buffer = trimStringToMaxLength( strings.TrimRight( buffer, "," ), MaxLineLength ) + "\r\n"

  case OUTFORMAT_TABS:
    buffer = this.renderLine( Format, "\t", MaxLineLength ) + "\r\n"
  
  case OUTFORMAT_LINE:
    buffer = this.renderLine( Format, "", MaxLineLength ) + "\r\n"
  
  case OUTFORMAT_JSON:
    for forThisColumn := uint( 0 ); forThisColumn < this.columns; forThisColumn++ {
      switch this.GetType( forThisColumn ) {
      case VALUE_TEXT, VALUE_BLOB:  buffer += fmt.Sprintf( "\"%s\":\"%s\",", strings.Replace( this.GetName( forThisColumn ), "\"", "\"\"", -1 ), strings.Replace( this.GetStringValue( forThisColumn ), "\"", "\"\"", -1 ) )
      default:                      buffer += fmt.Sprintf( "\"%s\":\"%s\",", strings.Replace( this.GetName( forThisColumn ), "\"", "\"\"", -1 ), this.GetStringValue( forThisColumn ) )
      }
    }
    buffer = fmt.Sprintf( "  {%s},\r\n", strings.TrimRight( buffer, "," ) )

  case OUTFORMAT_HTML:
    buffer = fmt.Sprintf( "<TR>\r\n%s</TR>\r\n", this.renderLine( Format, "", MaxLineLength ) )
  
  case OUTFORMAT_XML:
    buffer = fmt.Sprintf( "  <row>\r\n%s  </row>\r\n", this.renderLine( Format, "", MaxLineLength ) )

  case OUTFORMAT_MARKDOWN, OUTFORMAT_TABLE:
    buffer = "|" + this.renderLine( Format, "|", MaxLineLength ) + "|\r\n"
  
  case OUTFORMAT_BOX:
    buffer = "│" + this.renderLine( Format, "│", MaxLineLength ) + "│\r\n"

  }
  
  return io.WriteString( Out, buffer )
}