//
//                    ////              SQLite Cloud
//        ////////////  ///
//      ///             ///  ///        Product     : SQLite Cloud GO SDK
//     ///             ///  ///         Version     : 1.0.0
//     //             ///   ///  ///    Date        : 2021/08/31
//    ///             ///   ///  ///    Author      : Andreas Pfeil
//   ///             ///   ///  ///
//   ///     //////////   ///  ///      Description : GO Methods related to the
//   ////                ///  ///                     ResultRow class.
//     ////     //////////   ///
//        ////            ////
//          ////     /////
//             ///                      Copyright   : 2021 by SQLite Cloud Inc.
//
// -----------------------------------------------------------------------TAB=2

package sqlitecloud

import (
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"strings"
	"time"
)

type ResultRow struct {
	result  *Result
	index   uint64  `json:"Index"` // 0, 1, ... rows-1
	columns []Value `json:"ColumnValues"`
}

// ToJSON returns a JSON representation of this query result row.
// BUG(andreas): The ResultRow.ToJSON method is not implemented yet.
func (this *ResultRow) ToJSON() ([]byte, error) { return json.Marshal(this) }

// IsFirst returns true if this query result row is the first in the result set, false otherwise.
func (this *ResultRow) IsFirst() bool {
	switch {
	case this.result.GetNumberOfRows() < 1:
		return false
	default:
		return this.index == 0
	}
}

// IsLast returns true if this query result row is the last in the result set, false otherwise.
func (this *ResultRow) IsLast() bool {
	switch {
	case this.result.GetNumberOfRows() < 1:
		return false
	default:
		return this.index == this.result.GetNumberOfRows()-1
	}
}

// IsEOF returns false if this query result row is in the result set, true otherwise.
func (this *ResultRow) IsEOF() bool {
	switch {
	case this.result.GetNumberOfRows() < 1:
		return true
	default:
		return this.index >= this.result.GetNumberOfRows()
	}
}

// Rewind resets the iterator and returns the first row in this query result.
func (this *ResultRow) Rewind() *ResultRow {
	switch row, err := this.result.GetFirstRow(); {
	case err != nil:
		return nil
	default:
		return row
	}
}

// Next fetches the next row in this query result and returns it, otherwise if there is no next row, nil is returned.
func (this *ResultRow) Next() *ResultRow {
	switch row, err := this.result.GetRow(this.index + 1); {
	case err != nil:
		return nil
	default:
		return row
	}
}

func (this *ResultRow) GetNumberOfColumns() uint64 { return uint64(len(this.columns)) }

func (this *ResultRow) GetValue(Column uint64) (*Value, error) {
	switch {
	case Column >= this.GetNumberOfColumns():
		return nil, errors.New("Column index out of bounds")
	default:
		return &this.columns[Column], nil
	}
}

// GetMaxNameLength returns the number of runes of the longest column name.
func (this *ResultRow) GetMaxNameLength() uint64 { return this.result.GetMaxNameWidth() }

// GetNameLength returns the number of runes of the column name in the specified column.
// The Column index is an unsigned int in the range of 0...GetNumberOfColumns() - 1.
func (this *ResultRow) GetNameLength(Column uint64) (uint64, error) {
	return this.result.GetNameLength(Column)
}

// GetName returns the column name in column Column of this query result.
// The Column index is an unsigned int in the range of 0...GetNumberOfColumns() - 1.
func (this *ResultRow) GetName(Column uint64) (string, error) { return this.result.GetName(Column) }

// GetMaxWidth returns the number of runes of the value in the specified column with the maximum length in this query result.
// The Column index is an unsigned int in the range of 0...GetNumberOfColumns() - 1.
// BUG(andreas): Rename GetWidth->GetmaxWidth
func (this *ResultRow) GetMaxWidth(Column uint64) (uint64, error) {
	return this.result.GetMaxColumnWidth(Column)
}

// Die folgenden Methoden sollten überflüssig sein...

// GetType returns the type of the value in column Column of this query result row.
// The Column index is an unsigned int in the range of 0...GetNumberOfColumns() - 1.
// Possible return types are: VALUE_INTEGER, VALUE_FLOAT, VALUE_TEXT, VALUE_BLOB, VALUE_NULL
func (this *ResultRow) GetType(Column uint64) (byte, error) {
	switch value, err := this.GetValue(Column); {
	case err != nil:
		return CMD_NULL, err
	case value == nil:
		return CMD_NULL, errors.New("Column index out of bounds")
	default:
		return value.GetType(), nil
	}
}

// IsInteger returns true if this query result row column Column is of type "VALUE_INTEGER", false otherwise.
// The Column index is an unsigned int in the range of 0...GetNumberOfColumns() - 1.
func (this *ResultRow) IsInteger(Column uint64) bool {
	switch value, err := this.GetValue(Column); {
	case err != nil:
		return false
	case value == nil:
		return false
	default:
		return value.IsInteger()
	}
}

// IsFloat returns true if this query result row column Column is of type "VALUE_FLOAT", false otherwise.
// The Column index is an unsigned int in the range of 0...GetNumberOfColumns() - 1.
func (this *ResultRow) IsFloat(Column uint64) bool {
	switch value, err := this.GetValue(Column); {
	case err != nil:
		return false
	case value == nil:
		return false
	default:
		return value.IsFloat()
	}
}

// IsString returns true if this query result row column Column is of type "VALUE_TEXT", false otherwise.
// The Column index is an unsigned int in the range of 0...GetNumberOfColumns() - 1.
func (this *ResultRow) IsString(Column uint64) bool {
	switch value, err := this.GetValue(Column); {
	case err != nil:
		return false
	case value == nil:
		return false
	default:
		return value.IsString()
	}
}

// IsBLOB returns true if this query result row column Column is of type "VALUE_BLOB", false otherwise.
// The Column index is an unsigned int in the range of 0...GetNumberOfColumns() - 1.
func (this *ResultRow) IsBLOB(Column uint64) bool {
	switch value, err := this.GetValue(Column); {
	case err != nil:
		return false
	case value == nil:
		return false
	default:
		return value.IsBLOB()
	}
}

// IsNULL returns true if this query result row column Column is of type "VALUE_NULL", false otherwise.
// The Column index is an unsigned int in the range of 0...GetNumberOfColumns() - 1.
func (this *ResultRow) IsNULL(Column uint64) bool {
	switch value, err := this.GetValue(Column); {
	case err != nil:
		return false
	case value == nil:
		return false
	default:
		return value.IsNULL()
	}
}

// IsText returns true if this query result row column Column is of type "VALUE_TEXT" or "VALUE_BLOB", false otherwise.
// The Column index is an unsigned int in the range of 0...GetNumberOfColumns() - 1.
func (this *ResultRow) IsText(Column uint64) bool {
	switch value, err := this.GetValue(Column); {
	case err != nil:
		return false
	case value == nil:
		return false
	default:
		return value.IsText()
	}
}

// GetStringValue returns the contents in column Column of this query result row as string.
// The Column index is an unsigned int in the range of 0...GetNumberOfColumns() - 1.
func (this *ResultRow) GetString(Column uint64) (string, error) {
	return this.result.GetStringValue(this.index, Column)
}

// GetInt32Value returns the contents in column Column of this query result row as int32.
// The Column index is an unsigned int in the range of 0...GetNumberOfColumns() - 1.
func (this *ResultRow) GetInt32(Column uint64) (int32, error) {
	return this.result.GetInt32Value(this.index, Column)
}

// GetInt64Value returns the contents in column Column of this query result row as int64.
// The Column index is an unsigned int in the range of 0...GetNumberOfColumns() - 1.
func (this *ResultRow) GetInt64(Column uint64) (int64, error) {
	return this.result.GetInt64Value(this.index, Column)
}

// GetFloat32Value returns the contents in column Column of this query result row as float32.
// The Column index is an unsigned int in the range of 0...GetNumberOfColumns() - 1.
func (this *ResultRow) GetFloat32(Column uint64) (float32, error) {
	return this.result.GetFloat32Value(this.index, Column)
}

// GetFloat64Value returns the contents in column Column of this query result row as float64.
// The Column index is an unsigned int in the range of 0...GetNumberOfColumns() - 1.
func (this *ResultRow) GetFloat64(Column uint64) (float64, error) {
	return this.result.GetFloat64Value(this.index, Column)
}

// GetSQLDateTime parses this query result value in column Column as an SQL-DateTime and returns its value.
// The Column index is an unsigned int in the range of 0...GetNumberOfColumns() - 1.
func (this *ResultRow) GetSQLDateTime(Column uint64) (time.Time, error) {
	return this.result.GetSQLDateTime(this.index, Column)
}

////////

func (this *ResultRow) renderValue(Column uint64, Quotation string, NullValue string) string {
	//fmt.Printf( "renderValue, col = %d\r\n", Column )

	if val, err := this.GetValue(Column); err != nil {
		return ""
	} else {
		switch {
		case val.IsNULL():
			return NullValue
		case val.IsText():
			return val.GetString()
		default:
			return fmt.Sprintf("%s%s%s", Quotation, val.GetString(), Quotation)
		}
	}
}

func (this *ResultRow) renderLine(Format int, Separator string, NullValue string, NewLine string, MaxLineLength uint) string {
	buffer := ""
	for forThisColumn := uint64(0); forThisColumn < this.result.GetNumberOfColumns(); forThisColumn++ {
		maxWidth, _ := this.GetMaxWidth(forThisColumn)
		columnName, _ := this.GetName(forThisColumn)

		switch Format {
		case OUTFORMAT_LIST, OUTFORMAT_TABS:
			buffer += fmt.Sprintf("%s%s", this.renderValue(forThisColumn, "", NullValue), Separator)
		case OUTFORMAT_MARKDOWN, OUTFORMAT_TABLE, OUTFORMAT_BOX:
			buffer += fmt.Sprintf(fmt.Sprintf(" %%-%ds %s", maxWidth, Separator), this.renderValue(forThisColumn, "", NullValue))
		case OUTFORMAT_CSV:
			buffer += fmt.Sprintf("%s,", SQCloudEnquoteString(this.renderValue(forThisColumn, "", NullValue)))
		case OUTFORMAT_LINE:
			buffer += trimStringToMaxLength(fmt.Sprintf(fmt.Sprintf("%%%ds = %%s", this.result.MaxHeaderWidth), columnName, this.renderValue(forThisColumn, "", NullValue)), MaxLineLength) + NewLine
		case OUTFORMAT_HTML:
			buffer += trimStringToMaxLength(fmt.Sprintf("  <TD>%s</TD>", html.EscapeString(this.renderValue(forThisColumn, "", NullValue))), MaxLineLength) + NewLine
		case OUTFORMAT_XML:
			buffer += trimStringToMaxLength(fmt.Sprintf("    <field name=\"%s\">%s</field>", columnName, html.EscapeString(this.renderValue(forThisColumn, "", NullValue))), MaxLineLength) + NewLine
		}
	}
	switch Format {
	case OUTFORMAT_LINE, OUTFORMAT_HTML, OUTFORMAT_XML:
		return buffer // Multiline output was truncated already
	default:
		return trimStringToMaxLength(strings.TrimRight(buffer, Separator), MaxLineLength)
	}
}

// DumpToWriter renders this query result row into the buffer of an io.Writer.
// The output Format can be specified and must be one of the following values: OUTFORMAT_LIST, OUTFORMAT_CSV, OUTFORMAT_QUOTE, OUTFORMAT_TABS, OUTFORMAT_LINE, OUTFORMAT_JSON, OUTFORMAT_HTML, OUTFORMAT_MARKDOWN, OUTFORMAT_TABLE, OUTFORMAT_BOX
// The Separator argument specifies the column separating string (default: '|').
// All lines are truncated at MaxLineLeength. A MaxLineLangth of '0' means no truncation.
func (this *ResultRow) DumpToWriter(Out io.Writer, Format int, Separator string, NullValue string, NewLine string, MaxLineLength uint) (int, error) {
	buffer := ""

	switch Format {
	case OUTFORMAT_LIST, OUTFORMAT_CSV, OUTFORMAT_TABS, OUTFORMAT_LINE:
		buffer = this.renderLine(Format, Separator, NullValue, NewLine, MaxLineLength) + NewLine

	case OUTFORMAT_QUOTE:
		for forThisColumn := uint64(0); forThisColumn < this.GetNumberOfColumns(); forThisColumn++ {
			val, _ := this.GetString(forThisColumn)
			switch Type, err := this.GetType(forThisColumn); {
			case err != nil:
			case Type == CMD_STRING:
				fallthrough
			case Type == CMD_ZEROSTRING:
				fallthrough
			case Type == CMD_BLOB:
				buffer += fmt.Sprintf("'%s'%s", strings.Replace(val, "'", "\\'", -1), Separator)
			default:
				buffer += fmt.Sprintf("%s%s", val, Separator)
			}
		}
		buffer = trimStringToMaxLength(strings.TrimRight(buffer, Separator), MaxLineLength) + NewLine

	case OUTFORMAT_JSON:
		sep := Separator
		for forThisColumn := uint64(0); forThisColumn < this.GetNumberOfColumns(); forThisColumn++ {
			if forThisColumn == this.GetNumberOfColumns()-1 {
				sep = ""
			}
			columnName, _ := this.GetName(forThisColumn)
			switch Type, err := this.GetType(forThisColumn); {
			case err != nil:
			case Type == CMD_STRING:
				fallthrough
			case Type == CMD_ZEROSTRING:
				fallthrough
			case Type == CMD_BLOB:
				buffer += fmt.Sprintf("\"%s\":\"%s\"%s", strings.Replace(columnName, "\"", "\\\"", -1), strings.Replace(this.renderValue(forThisColumn, "", NullValue), "\"", "\\\"", -1), sep)
			default:
				buffer += fmt.Sprintf("\"%s\":%s%s", strings.Replace(columnName, "\"", "\\\"", -1), this.renderValue(forThisColumn, "", NullValue), sep)
			}
		}
		buffer = trimStringToMaxLength(fmt.Sprintf("  {%s}", buffer), MaxLineLength) + NewLine

	case OUTFORMAT_HTML:
		buffer = trimStringToMaxLength("<TR>", MaxLineLength) + NewLine +
			this.renderLine(Format, Separator, NullValue, NewLine, MaxLineLength) +
			trimStringToMaxLength("</TR>", MaxLineLength) + NewLine

	case OUTFORMAT_XML:
		buffer = trimStringToMaxLength("  <row>", MaxLineLength) + NewLine +
			this.renderLine(Format, Separator, NullValue, NewLine, MaxLineLength) +
			trimStringToMaxLength("  </row>", MaxLineLength) + NewLine

	case OUTFORMAT_MARKDOWN, OUTFORMAT_TABLE, OUTFORMAT_BOX:
		buffer = trimStringToMaxLength(fmt.Sprintf("%s%s%s", Separator, this.renderLine(Format, Separator, NullValue, NewLine, MaxLineLength), Separator), MaxLineLength) + NewLine
	}

	return io.WriteString(Out, buffer)
}
