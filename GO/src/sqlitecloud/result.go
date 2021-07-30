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

const OUTFORMAT_LIST 	    = 0
const OUTFORMAT_CSV 	    = 1
const OUTFORMAT_QUOTE	    = 2
const OUTFORMAT_TABS	    = 3
const OUTFORMAT_LINE	    = 4
const OUTFORMAT_JSON	    = 5
const OUTFORMAT_HTML	    = 6
const OUTFORMAT_MARKDOWN	= 7
const OUTFORMAT_TABLE	    = 8
const OUTFORMAT_BOX	      = 9

type SQCloudResult struct {
	result *C.struct_SQCloudResult

	Rows 						uint		// Must be set during ...
	Columns      		uint		// Must be set

	ColumnWidth  		[]uint 	// Must be set
	HeaderWidth  		[]uint 	// Must be set
	MaxHeaderWidth 	uint   	// Must be set

	Type 						uint 		// Must be set during Initialisation
	ErrorCode 			int
	ErrorMessage 		string
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





////// Row Methods


func (this *SQCloudResult ) GetRow( Row uint ) ( *SQCloudResultRow ) {
	if this.Rows < 1 {
		return nil
	}
	if Row >= this.Rows {
		return nil
	}
	if Row < 0 {
		return nil
	}

	row := SQCloudResultRow{
		result	: this,
		row			: 0,
		rows		: this.Rows,
		columns	: this.Columns,
	}
	return &row
}
func (this *SQCloudResult ) GetFirstRow() *SQCloudResultRow {
	return this.GetRow( 0 )
}
func (this *SQCloudResult ) GetLastRow() *SQCloudResultRow {
	if this.Rows > 0 {
		return this.GetRow( this.Rows - 1 )
	}
	return nil
}

func centerString( Buffer string, Width int ) string {
	return fmt.Sprintf( "%[1]*s", -Width, fmt.Sprintf( "%[1]*s", ( Width + len( Buffer ) ) / 2, Buffer ) )
}
func trimStringLength( Buffer string, MaxLineLength uint ) string {
	len := uint( len( Buffer ) )
	if MaxLineLength == 0 {
		MaxLineLength = len
	}
	if MaxLineLength < len {
		len = MaxLineLength
	}
	return Buffer[ 0 : len ] // return fmt.Sprintf( fmt.Sprintf( "%%%ds", MaxLineLength ), line )
}
func createTableLine( Left string, Fill string, Seperator string, Right string, ColumnWidth []uint, MaxLineLength uint ) string {
	outBuffer := ""
	for _, len := range ColumnWidth {
		outBuffer += fmt.Sprintf( "%s%s%s", outBuffer, strings.Repeat( Fill, int( len ) ), Seperator )
	}
	outBuffer = fmt.Sprintf( "%s%s%s", Left, strings.TrimRight( outBuffer, Seperator ), Right )
	return trimStringLength( outBuffer, MaxLineLength )
}

func (this *SQCloudResult) DumpHeaderToWriter( Out io.Writer, Format int, MaxLineLength uint ) ( int, error ) {
	header := ""
	switch( Format ) {
		case OUTFORMAT_LIST, OUTFORMAT_CSV,	OUTFORMAT_QUOTE, OUTFORMAT_TABS, OUTFORMAT_LINE, OUTFORMAT_HTML: return 0, nil // No header
		case OUTFORMAT_JSON: return io.WriteString( Out, "[" )

		case OUTFORMAT_MARKDOWN:
			header  = ""
			header += createTableLine( "|", "-", "|", "|", this.ColumnWidth, MaxLineLength )

		case OUTFORMAT_TABLE:
			header  = createTableLine( "+", "-", "+", "+", this.ColumnWidth, MaxLineLength )
			header += ""
			header += createTableLine( "+", "-", "+", "+", this.ColumnWidth, MaxLineLength )

		case OUTFORMAT_BOX:
			header  = createTableLine( "┌", "─", "┬", "┐", this.ColumnWidth, MaxLineLength )
			header += ""
			header += createTableLine( "├", "─", "┼", "┤", this.ColumnWidth, MaxLineLength )
	}
	return io.WriteString( Out, header )
}

func (this *SQCloudResult) DumpFooterToWriter( Out io.Writer, Format int, MaxLineLength uint ) ( int, error ) {
	footer := ""
	switch( Format ) {
		case OUTFORMAT_LIST, OUTFORMAT_CSV,	OUTFORMAT_QUOTE, OUTFORMAT_TABS, OUTFORMAT_LINE, OUTFORMAT_HTML, OUTFORMAT_MARKDOWN: return 0, nil // No footer
		case OUTFORMAT_JSON: return io.WriteString( Out, "]" )	
		case OUTFORMAT_TABLE:
			footer = createTableLine( "+", "-", "+", "+", this.ColumnWidth, MaxLineLength )
		case OUTFORMAT_BOX:
			footer = createTableLine( "└", "─", "┴", "┘", this.ColumnWidth, MaxLineLength )
	}
	return io.WriteString( Out, footer )
}

func (this *SQCloudResult) DumpToWriter( Out io.Writer, Format int, Seperator string, MaxLineLength uint, SurpressOK bool ) ( int, error ) {
	switch this.Type {
	case RESULT_OK:
		return io.WriteString( Out, "OK" )

	case RESULT_NULL:
		return io.WriteString( Out, "NULL" )

	case RESULT_ERROR:
		return io.WriteString( Out, fmt.Sprintf( "ERROR: %s (%d)", this.ErrorMessage, this.ErrorCode ) )

	case RESULT_STRING, RESULT_INTEGER, RESULT_FLOAT, RESULT_JSON:
		return io.WriteString( Out, this.CGetResultBuffer() )

	case RESULT_ROWSET:
		var totalOutputLength int = 0

		// Render Table Header
		if len, err := this.DumpHeaderToWriter( Out, Format, MaxLineLength ); err == nil  {
			totalOutputLength += len
		} else {
			return totalOutputLength + len, err
		}

		// Render Table Body
		for row := this.GetFirstRow(); !row.IsEOF(); row = row.Next() {
			if len, err := row.DumpToWriter( Out, Format, Seperator, MaxLineLength ); err == nil {
				totalOutputLength += len
			} else {
				return totalOutputLength + len, err
			}
		}

		
		// Render Table Footer
		if len, err := this.DumpFooterToWriter( Out, Format, MaxLineLength ); err == nil  {
			totalOutputLength += len
		} else {
			return totalOutputLength + len, err
		}

		return totalOutputLength, nil

	default:
		return 0, errors.New( "Unknown Output Format" )
	}
}