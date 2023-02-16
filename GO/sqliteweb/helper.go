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
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/smtp"
	"os"
	"strings"
	"text/template"

	sqlitecloud "github.com/sqlitecloud/sdk"
)

func Hash(data string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(data)))
}

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

func sendMailWithTemplate(from string, to string, body map[string]string, templateName string, language string) error {
	if language == "" {
		language = "en"
	}

	body["From"] = from
	body["To"] = to

	path := cfg.Section("mail").Key("mail.template.path").String()

	for _, templatePath := range []string{fmt.Sprintf("%s/%s/%s", path, language, templateName), fmt.Sprintf("%s/%s", path, templateName)} {
		if !PathExists(templatePath) {
			continue
		}

		if template, err := template.ParseFiles(templatePath); err == nil {
			var outBuffer bytes.Buffer
			if err = template.Execute(&outBuffer, body); err != nil {
				return fmt.Errorf("Error in template.ParseFiles: %s %s", templateName, err.Error())
			}

			if err := sendMail(from, []string{to}, outBuffer.Bytes()); err != nil {
				return fmt.Errorf("Error while sending mail: %s %s", templateName, err.Error())
			}

			return nil

		} else {
			return fmt.Errorf("Error in template.ParseFiles: %s %s", templateName, err.Error())
		}
	}

	return fmt.Errorf("Email template not found: %s", templateName)
}

func sendMail(from string, to []string, msg []byte) error {
	host, _, err := net.SplitHostPort(cfg.Section("mail").Key("mail.proxy.host").String())
	if err != nil {
		return err
	}

	auth := smtp.PlainAuth("", cfg.Section("mail").Key("mail.proxy.user").String(), cfg.Section("mail").Key("mail.proxy.password").String(), host)

	err = smtp.SendMail(cfg.Section("mail").Key("mail.proxy.host").String(), auth, from, to, msg)
	return err
}

func writeError(writer http.ResponseWriter, statusCode int, message string, allowedMethods *string) {
	if statusCode == http.StatusMethodNotAllowed && allowedMethods != nil {
		writer.Header().Set("Access-Control-Allow-Methods", *allowedMethods)
		writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		writer.Header().Set("Pragma", "no-cache")
		writer.Header().Set("Expires", "0")
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Content-Encoding", "utf-8")
	writer.WriteHeader(statusCode)
	writer.Write([]byte(fmt.Sprintf("{\"status\":%d,\"message\":\"%s\"}", statusCode, message)))
}
