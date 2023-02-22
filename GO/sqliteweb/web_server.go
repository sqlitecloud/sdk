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
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"

	"github.com/gobwas/glob"
)

var originChecker glob.Glob

func init() {
	initializeSQLiteWeb()
	originChecker = glob.MustCompile("{http://*.sqlitecloud.io:*,https://*.sqlitecloud.io:*,https://sqlitecloud.io,https://sqlitecloud.io:*,https://fabulous-arithmetic-c4a014.netlify.app}")
}

func initWeb() {
	if cfg.Section("web").Key("enabled").MustBool(false) {
		SQLiteWeb.router.HandleFunc("/web/v1/sendmail", SQLiteWeb.serveWebSendmail).Methods(http.MethodPost, http.MethodOptions)
		SQLiteWeb.router.HandleFunc("/web/v1/user", SQLiteWeb.servePostWebUser).Methods(http.MethodPost, http.MethodOptions)
	}
}

func (this *Server) serveWebSendmail(writer http.ResponseWriter, request *http.Request) {
	if !DEBUG_SQLITEWEB {
		// security check, only allow requests from https://sqlitecloud.io
		origin := request.Header.Get("Origin")
		if !originChecker.Match(origin) {
			SQLiteWeb.Logger.Errorf("[WebSendmail] Invalid origin in serveWebSendmail: %s", origin)
			writeError(writer, http.StatusBadRequest, "Bad Request", nil)
			return
		}
	}

	if request.Method == http.MethodOptions {
		writer.Header().Set("Access-Control-Allow-Methods", "POST")
		writer.WriteHeader(http.StatusOK)
		return
	}

	// vars := mux.Vars(request)
	body := make(map[string]string)

	if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
		SQLiteWeb.Logger.Errorf("[WebSendmail] Cannot decode mail body: %s", err.Error())
		writeError(writer, http.StatusBadRequest, "Bad Request", nil)
		return
	}

	if honeypot, hasHoneypot := body["Honeypot"]; hasHoneypot && honeypot != "" {
		SQLiteWeb.Logger.Errorf("[WebSendmail] Honeypot: %s", honeypot)
		writeError(writer, http.StatusBadRequest, "Bad Request", nil)
		return
	}

	from, _ := mail.ParseAddress(cfg.Section("mail").Key("mail.from").String())
	to, _ := mail.ParseAddress(cfg.Section("mail").Key("mail.contactus.to").String())
	subject := fmt.Sprintf("[Contact Us] Message From %s %s", body["FirstName"], body["LasttName"])
	templateName := body["Template"]
	language := body["Language"]

	if err := sendMailWithTemplate(from, to, subject, body, templateName, language); err != nil {
		SQLiteWeb.Logger.Errorf("[WebSendmail] %s", err.Error())
		writeError(writer, http.StatusInternalServerError, "Internal Error", nil)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Content-Encoding", "utf-8")
	writer.Write([]byte(fmt.Sprintf("{\"status\":%d,\"message\":\"%s\"}", http.StatusOK, "OK")))
}

func (this *Server) servePostWebUser(writer http.ResponseWriter, request *http.Request) {
	if !DEBUG_SQLITEWEB {
		// security check, only allow requests from https://sqlitecloud.io
		origin := request.Header.Get("Origin")
		if !originChecker.Match(origin) {
			SQLiteWeb.Logger.Errorf("[WebUser] Invalid origin in serveWebUser: %s", origin)
			writeError(writer, http.StatusBadRequest, "Bad Request", nil)
			return
		}
	}

	if request.Method == http.MethodOptions {
		writer.Header().Set("Access-Control-Allow-Methods", "POST")
		writer.WriteHeader(http.StatusOK)
		return
	}

	body := make(map[string]string)

	if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
		SQLiteWeb.Logger.Errorf("[WebUser]Cannot decode mail body: %s", err.Error())
		writeError(writer, http.StatusBadRequest, "Bad Request", nil)
		return
	}

	if honeypot, hasHoneypot := body["Honeypot"]; hasHoneypot && honeypot != "" {
		SQLiteWeb.Logger.Errorf("[WebUser] Honeypot: %s", honeypot)
		writeError(writer, http.StatusBadRequest, "Bad Request", nil)
		return
	}

	// validate mail address
	email, err := mail.ParseAddress(body["Email"])
	if err != nil {
		SQLiteWeb.Logger.Errorf("[WebUser] invalid mail: %s %s", email.Address, err.Error())
		writeError(writer, http.StatusBadRequest, "Bad Request", nil)
		return
	}

	// add the new user to the waitlist (WebUser table)
	sql := "INSERT INTO WebUser (first_name, last_name, company, email, referral, message) VALUES (?, ?, ?, ?, ?, ? );"
	if _, err, errcode, _, _ := dashboardcm.ExecuteSQLArray("auth", sql, &[]interface{}{body["FirstName"], body["LastName"], body["Company"], email.Address, body["Referral"], body["Message"]}); err != nil {
		if errcode == 19 {
			// (19) SQLITE_CONSTRAINT
			SQLiteWeb.Logger.Errorf("[WebUser] Email already exists: %s %s", email.Address, err.Error())
			writeError(writer, http.StatusConflict, "Email already exists", nil)
			return
		} else {
			SQLiteWeb.Logger.Errorf("[WebUser] Error while executing INSERT INTO WebUser: %s", err.Error())
			writeError(writer, http.StatusInternalServerError, "Internal error", nil)
			return
		}
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Content-Encoding", "utf-8")
	writer.Write([]byte(fmt.Sprintf("{\"status\":%d,\"message\":\"%s\"}", http.StatusOK, "OK")))
	return
}
