//
//                    ////              SQLite Cloud
//        ////////////  ///
//      ///             ///  ///        Product     : SQLite Cloud Web Server
//     ///             ///  ///         Version     : 0.2.0
//     //             ///   ///  ///    Date        : 2023/02/13
//    ///             ///   ///  ///    Author      : Andreas Pfeil
//   ///             ///   ///  ///
//   ///     //////////   ///  ///      Description : Endpoints used by sqlitecloud.io
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
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"

	"github.com/gobwas/glob"
)

var originChecker glob.Glob

func init() {
	initializeSQLiteWeb()
	originChecker = glob.MustCompile("{https://*.sqlitecloud.io,https://*.sqlitecloud.io:*,https://sqlitecloud.io,https://sqlitecloud.io:*,https://fabulous-arithmetic-c4a014.netlify.app}")
}

func initWeb() {
	if cfg.Section("web").Key("enabled").MustBool(false) {
		SQLiteWeb.router.HandleFunc("/web/v1/sendmail", SQLiteWeb.serveWebSendmail).Methods("POST")
	}
}

func (this *Server) serveWebSendmail(writer http.ResponseWriter, request *http.Request) {
	// security check, only allow requests from https://sqlitecloud.io
	origin := request.Header.Get("Origin")
	if !originChecker.Match(origin) {
		SQLiteWeb.Logger.Errorf("Invalid origin in serveWebSendmail: %s", origin)
		writeError(writer, http.StatusBadRequest, "Bad Request", nil)
		return
	}

	// vars := mux.Vars(request)
	body := make(map[string]string)

	if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
		SQLiteWeb.Logger.Errorf("Cannot decode mail body: %s", err.Error())
		writeError(writer, http.StatusBadRequest, "Bad Request", nil)
		return
	}

	templateName := body["Template"]
	language := body["Language"]

	if language == "" {
		language = "en"
	}

	from := cfg.Section("mail").Key("mail.from").String()
	to := cfg.Section("mail").Key("mail.contactus.to").String()
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
				SQLiteWeb.Logger.Errorf("Error in template.ParseFiles: %s %s", templateName, err.Error())
				writeError(writer, http.StatusInternalServerError, "Internal Error (1)", nil)
				return
			}

			if err := mail(from, []string{to}, outBuffer.Bytes()); err != nil {
				SQLiteWeb.Logger.Errorf("Error while sending mail: %s %s", templateName, err.Error())
				writeError(writer, http.StatusInternalServerError, "Internal Error (2)", nil)
				return
			}

			writer.Header().Set("Content-Type", "application/json")
			writer.Header().Set("Content-Encoding", "utf-8")
			writer.Write([]byte(fmt.Sprintf("{\"status\":%d,\"message\":\"%s\"}", http.StatusOK, "OK")))
			return

		} else {
			SQLiteWeb.Logger.Errorf("Error in template.ParseFiles: %s %s", templateName, err.Error())
			writeError(writer, http.StatusInternalServerError, "Internal Error (3)", nil)
			return
		}
	}

	SQLiteWeb.Logger.Errorf("Email template not found: %s", templateName)
	writeError(writer, http.StatusInternalServerError, "Internal Error (4)", nil)
}
