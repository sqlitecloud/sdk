//
//                    ////              SQLite Cloud
//        ////////////  ///
//      ///             ///  ///        Product     : SQLite Cloud Web Server
//     ///             ///  ///         Version     : 0.0.1
//     //             ///   ///  ///    Date        : 2021/11/17
//    ///             ///   ///  ///    Author      : Andreas Pfeil
//   ///             ///   ///  ///
//   ///     //////////   ///  ///      Description :
//   ////                ///  ///
//     ////     //////////   ///
//        ////            ////
//          ////     /////
//             ///                      Copyright   : 2021 by SQLite Cloud Inc.
//
// -----------------------------------------------------------------------TAB=2

package main

import (
	"errors"
	"os"
	"strings"

	sqlitecloud "github.com/sqlitecloud/sdk"
)

func PathExists(path string) bool {
	path = strings.TrimSpace(path)
	if path == "" {
		return false
	}

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// resultToObj is a helper method to convert the sqlitecloud.Result to a Map.
// The resulting map can be added to the response object before getting the final JSON response.
func ResultToObj(result *sqlitecloud.Result) (interface{}, error) {
	switch {
	case result.IsOK():
		return "OK", nil

	case result.IsNULL():
		return nil, nil

	case result.IsError():
		_, _, _, errMsg, _ := result.GetError()
		return nil, errors.New(errMsg)

	case result.IsString(), result.IsJSON():
		return result.GetString()

	case result.IsInteger():
		return result.GetInt32()

	case result.IsFloat():
		return result.GetFloat64()

	case result.IsArray():
		fallthrough

	case result.IsRowSet():
		var value = make(map[string]interface{}, 2)

		if numCols := result.GetNumberOfColumns(); numCols > 0 {
			var cols = make([]string, 0, numCols)
			for c := uint64(0); c < numCols; c++ {
				cols = append(cols, result.GetName_(c))
			}
			value["columns"] = cols
		}

		if numRows := result.GetNumberOfRows(); numRows > 0 {
			var rows = make([]map[string]interface{}, 0, numRows)

			for r := uint64(0); r < numRows; r++ {
				var row = make(map[string]interface{})

				for c := uint64(0); c < result.GetNumberOfColumns(); c++ {
					var v interface{}
					// L.PushString(result.GetName(c))
					switch result.GetValueType_(r, c) {
					case ':':
						v = result.GetInt32Value_(r, c)
					case ',':
						v = result.GetFloat64Value_(r, c)
					default:
						v = result.GetStringValue_(r, c)
					}
					row[result.GetName_(c)] = v
				}
				rows = append(rows, row)
			}
			value["rows"] = rows
		}
		return value, nil
	default:
		return 0, errors.New("Unknown Output Format")
	}
}
