package sqlitecloud

import "fmt"
//import "os"
import "io"
//import "bufio"
import "strings"
//import "errors"
import "time"
//import "strconv"

type SQCloudResultRow struct {
	result 	*SQCloudResult
	row 		uint
	rows 		uint
	columns uint
}



func (this *SQCloudResultRow ) IsFirst() bool {
	return this.row == 0
}
func (this *SQCloudResultRow ) IsLast() bool {
	return this.row == this.rows
}

func (this *SQCloudResultRow ) IsEOF() bool {
	return this.row >= this.rows
}

func (this *SQCloudResultRow ) Rewind() *SQCloudResultRow {
	this.row = 0
	return this
}

func (this *SQCloudResultRow ) Next() *SQCloudResultRow {
	if this.row < this.rows {
		this.row++
		return this
	} 
	return nil
}





func (this *SQCloudResultRow ) GetType( Column uint ) int {
	return this.result.GetValueType( this.row, Column )
} 
func (this *SQCloudResultRow ) GetNameWidth( Column uint ) uint {
	return this.result.GetNameWidth( Column )
}
func (this *SQCloudResultRow ) GetName( Column uint ) string {
	return this.result.GetColumnName( Column )
}
func (this *SQCloudResultRow ) GetMaxNameWidth() uint {
	return this.result.GetMaxNameWidth()
}



func (this *SQCloudResultRow ) GetWidth( Column uint ) uint {
	return this.result.GetMaxColumnLength( Column )
}
func (this *SQCloudResultRow ) GetStringValue( Column uint ) string {
	return this.result.GetStringValue( this.row, Column )
} 
func (this *SQCloudResultRow ) GetInt32Value( Column uint ) int32 {
	return this.result.GetInt32Value( this.row, Column )
} 
func (this *SQCloudResultRow ) GetInt64Value( Column uint ) int64 {
	return this.result.GetInt64Value( this.row, Column )
} 
func (this *SQCloudResultRow ) GetFloat32Value( Column uint ) float32 {
	return this.result.GetFloat32Value( this.row, Column )
} 
func (this *SQCloudResultRow ) GetFloat64Value( Column uint ) float64 {
	return this.result.GetFloat64Value( this.row, Column )
}
func (this *SQCloudResultRow ) GetSQLDateTime( Column uint ) time.Time {
	return this.result.GetSQLDateTime( this.row, Column )
} 




func (this *SQCloudResultRow) renderLine( Format int, Seperator string, MaxLineLength uint ) string {
	buffer := ""
	for c := uint( 0 ); c < this.columns; c++ {
		switch Format {
		case OUTFORMAT_LIST, OUTFORMAT_TABS: 											buffer += fmt.Sprintf( "%s%s", this.GetStringValue( c ), Seperator )
		case OUTFORMAT_MARKDOWN, OUTFORMAT_TABLE, OUTFORMAT_BOX: 	buffer +=	fmt.Sprintf( fmt.Sprintf( "%%%ds|", this.GetWidth( c ) ), this.GetStringValue( c ) ) 
		case OUTFORMAT_CSV: 																			buffer += fmt.Sprintf( "%s,", SQCloudEnquoteString( this.GetStringValue( c ) ) )
		case OUTFORMAT_LINE:																			buffer += fmt.Sprintf( fmt.Sprintf( "%%%ds = %s\r\n", this.result.MaxHeaderWidth ), this.GetName( c ), this.GetStringValue( c ) )
		case OUTFORMAT_HTML: 																			buffer += fmt.Sprintf( "  <TD>%s</TD>\r\n", this.GetStringValue( c ) )
		}
	}
	return trimStringLength( strings.TrimRight( buffer, Seperator ), MaxLineLength )
}


func (this *SQCloudResultRow) DumpToWriter( Out io.Writer, Format int, Seperator string, MaxLineLength uint ) ( int, error ) {
	buffer := ""
	
	switch( Format ) {
	case OUTFORMAT_LIST:
		buffer = this.renderLine( Format, Seperator, MaxLineLength )

	case OUTFORMAT_CSV:
		buffer = this.renderLine( Format, ",", MaxLineLength )
		//for c := uint( 0 ); c < this.columns; c++ {
		//	buffer += fmt.Sprintf( "%s,", SQCloudEnquoteString( this.GetStringValue( c ) ) )
		//}
		//buffer = trimStringLength( strings.TrimRight( buffer, "," ), MaxLineLength )
	
	case OUTFORMAT_QUOTE:	
		for forThisColumn := uint( 0 ); forThisColumn < this.columns; forThisColumn++ {
			switch this.GetType( forThisColumn ) {
			case VALUE_TEXT, VALUE_BLOB: 	buffer += fmt.Sprintf( "'%s',", strings.Replace( this.GetStringValue( forThisColumn ), "'", "\\'", -1 ) )
			default: 											buffer += fmt.Sprintf( "%s,", this.GetStringValue( forThisColumn ))
			}
		}
		buffer = trimStringLength( strings.TrimRight( buffer, "," ), MaxLineLength )

	case OUTFORMAT_TABS:
		buffer = this.renderLine( Format, "\t", MaxLineLength )
	
	case OUTFORMAT_LINE:
		buffer = this.renderLine( Format, "", MaxLineLength )
		// for c := uint( 0 ); c < this.columns; c++ {
		// 	buffer += fmt.Sprintf( fmt.Sprintf( "%%%ds = %s\r\n", this.result.MaxHeaderWidth ), this.GetName( c ), this.GetStringValue( c ) )
		// }
	
	case OUTFORMAT_JSON:
		for forThisColumn := uint( 0 ); forThisColumn < this.columns; forThisColumn++ {
			switch this.GetType( forThisColumn ) {
			case VALUE_TEXT, VALUE_BLOB: 	buffer += fmt.Sprintf( "\"%s\":\"%s\",", strings.Replace( this.GetName( forThisColumn ), "\"", "\"\"", -1 ), strings.Replace( this.GetStringValue( forThisColumn ), "\"", "\"\"", -1 ) )
			default: 											buffer += fmt.Sprintf( "\"%s\":\"%s\",", strings.Replace( this.GetName( forThisColumn ), "\"", "\"\"", -1 ), this.GetStringValue( forThisColumn ) )
			}
		}
		buffer = fmt.Sprintf( "{%s},", strings.TrimRight( buffer, "," ) )

	case OUTFORMAT_HTML:
		// for c := uint( 0 ); c < this.columns; c++ {
		// 	buffer += fmt.Sprintf( "  <TD>%s</TD>\r\n", this.GetStringValue( c ) )
		// }
		buffer = fmt.Sprintf( "<TR>\r\n%s</TR>", this.renderLine( Format, "", MaxLineLength ) )
	
	case OUTFORMAT_MARKDOWN, OUTFORMAT_TABLE:
		buffer = "|" + this.renderLine( Format, "|", MaxLineLength ) + "|"
		// buffer = "|"
		// for c := uint( 0 ); c < this.columns; c++ {
		// 	buffer +=	fmt.Sprintf( fmt.Sprintf( "%%%ds|", this.GetWidth( c ) ), this.GetStringValue( c ) ) 
		// }
	
	case OUTFORMAT_BOX:
		buffer = "│" + this.renderLine( Format, "│", MaxLineLength ) + "│"
		// buffer = "│"
		// for c := uint( 0 ); c < this.columns; c++ {
		// 	buffer +=	fmt.Sprintf( fmt.Sprintf( "%%%ds│", this.GetWidth( c ) ), this.GetStringValue( c ) ) 
		// }
	}
	
	return io.WriteString( Out, buffer )
}
