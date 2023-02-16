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
	"github.com/sethvargo/go-password/password"
)

var originChecker glob.Glob

func init() {
	initializeSQLiteWeb()
	originChecker = glob.MustCompile("{http://*.sqlitecloud.io:*,https://*.sqlitecloud.io:*,https://sqlitecloud.io,https://sqlitecloud.io:*,https://fabulous-arithmetic-c4a014.netlify.app}")
}

func initWeb() {
	if cfg.Section("web").Key("enabled").MustBool(false) {
		SQLiteWeb.router.HandleFunc("/web/v1/sendmail", SQLiteWeb.serveWebSendmail).Methods(http.MethodPost, http.MethodOptions)
		SQLiteWeb.router.HandleFunc("/web/v1/user", SQLiteWeb.serveWebUser).Methods(http.MethodPost, http.MethodOptions)
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

	from := cfg.Section("mail").Key("mail.from").String()
	to := cfg.Section("mail").Key("mail.contactus.to").String()
	templateName := body["Template"]
	language := body["Language"]

	if err := sendMailWithTemplate(from, to, body, templateName, language); err != nil {
		SQLiteWeb.Logger.Errorf("[WebSendmail] %s", err.Error())
		writeError(writer, http.StatusInternalServerError, "Internal Error", nil)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Content-Encoding", "utf-8")
	writer.Write([]byte(fmt.Sprintf("{\"status\":%d,\"message\":\"%s\"}", http.StatusOK, "OK")))

	/*
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
					SQLiteWeb.Logger.Errorf("[WebSendmail] Error in template.ParseFiles: %s %s", templateName, err.Error())
					writeError(writer, http.StatusInternalServerError, "Internal Error (1)", nil)
					return
				}

				if err := mail(from, []string{to}, outBuffer.Bytes()); err != nil {
					SQLiteWeb.Logger.Errorf("[WebSendmail] Error while sending mail: %s %s", templateName, err.Error())
					writeError(writer, http.StatusInternalServerError, "Internal Error (2)", nil)
					return
				}

				writer.Header().Set("Content-Type", "application/json")
				writer.Header().Set("Content-Encoding", "utf-8")
				writer.Write([]byte(fmt.Sprintf("{\"status\":%d,\"message\":\"%s\"}", http.StatusOK, "OK")))
				return

			} else {
				SQLiteWeb.Logger.Errorf("[WebSendmail] Error in template.ParseFiles: %s %s", templateName, err.Error())
				writeError(writer, http.StatusInternalServerError, "Internal Error (3)", nil)
				return
			}
		}

		SQLiteWeb.Logger.Errorf("[WebSendmail] Email template not found: %s", templateName)
		writeError(writer, http.StatusInternalServerError, "Internal Error (4)", nil)
	*/
}

func (this *Server) serveWebUser(writer http.ResponseWriter, request *http.Request) {
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

	// vars := mux.Vars(request)
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
	email := body["Email"]
	if _, err := mail.ParseAddress(email); err != nil {
		SQLiteWeb.Logger.Errorf("[WebUser] invalid mail: %s %s", email, err.Error())
		writeError(writer, http.StatusBadRequest, "Bad Request", nil)
		return
	}

	// add the new user to the waitlist (WebUser table)
	sql := "INSERT INTO WebUser (first_name, last_name, company, email, referral, message) VALUES (?, ?, ?, ?, ?, ? );"
	if _, err, errcode, _, _ := dashboardcm.ExecuteSQLArray("auth", sql, &[]interface{}{body["FirstName"], body["LastName"], body["Company"], email, body["Referral"], body["Message"]}); err != nil {
		// (19) SQLITE_CONSTRAINT
		if errcode == 19 {
			SQLiteWeb.Logger.Errorf("[WebUser] Email already exists: %s %s", email, err.Error())
			writeError(writer, http.StatusConflict, "Email already exists", nil)
			return
		} else {
			SQLiteWeb.Logger.Errorf("[WebUser] Error while executing INSERT INTO WebUser: %s", err.Error())
			writeError(writer, http.StatusInternalServerError, "Internal error", nil)
			return
		}
	}

	// https://github.com/sethvargo/go-password ?
	password, err := password.Generate(10, 5, 3, false, false)
	if err != nil {
		SQLiteWeb.Logger.Errorf("[WebUser] Error while generating new password: %s", err.Error())
		writeError(writer, http.StatusInternalServerError, "Internal Server Error", nil)
		return
	}
	// passwordhash := // if the password is hashed, I have to implement a new password recovery procedure, I cannot send the hashed password

	// add a new dashboard user (User table)
	sql = "INSERT INTO Company (name) VALUES (?) RETURNING id"
	res, err, _, _, _ := dashboardcm.ExecuteSQLArray("auth", sql, &[]interface{}{body["Company"]})
	if err != nil {
		SQLiteWeb.Logger.Errorf("[WebUser] Error while executing INSERT INTO Company: %s", err.Error())
		writeError(writer, http.StatusInternalServerError, "Internal Server Error", nil)
		return
	}

	sql = "INSERT INTO User (company_id, first_name, last_name, email, password) VALUES (?, ?, ?, ?, ?);"
	if _, err, _, _, _ := dashboardcm.ExecuteSQLArray("auth", sql, &[]interface{}{res.GetInt32Value_(0, 0), body["FirstName"], body["LastName"], email, password, body["Message"]}); err != nil {
		SQLiteWeb.Logger.Errorf("[WebUser] Error while executing INSERT INTO WebUser: %s", err.Error())
		writeError(writer, http.StatusInternalServerError, "Internal Server Error", nil)
		return
	}

	// send a welcome mail with the dashboard password
	from := cfg.Section("mail").Key("mail.from").String()
	templateName := "welcome_web.eml"
	language := "en"
	body["Password"] = password
	if err := sendMailWithTemplate(from, email, body, templateName, language); err != nil {
		SQLiteWeb.Logger.Errorf("[WebUser] Error while sending welcome mail to %s: %s", email, err.Error())
		writeError(writer, http.StatusInternalServerError, "Internal Server Error", nil)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Content-Encoding", "utf-8")
	writer.Write([]byte(fmt.Sprintf("{\"status\":%d,\"message\":\"%s\"}", http.StatusOK, "OK")))
	return
}
