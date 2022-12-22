//
//                    ////              SQLite Cloud
//        ////////////  ///
//      ///             ///  ///        Product     : SQLite Cloud Web Server
//     ///             ///  ///         Version     : 0.2.0
//     //             ///   ///  ///    Date        : 2022/12/08
//    ///             ///   ///  ///    Author      : Andrea Donetti
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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	// "github.com/gorilla/schema"
	// "github.com/gookit/validate"
	// "github.com/go-swagger/go-swagger"

	sqlitecloud "github.com/sqlitecloud/sdk"
)

const ApiKeyHeaderKey string = "X-SQLiteCloud-Api-Key"

type methodMask int

const (
	GET_BITMASK methodMask = 1 << iota
	POST_BITMASK
	PATCH_BITMASK
	DELETE_BITMASK
)

var (
	methodsMap map[string]methodMask = map[string]methodMask{
		http.MethodGet:    GET_BITMASK,
		http.MethodPost:   POST_BITMASK,
		http.MethodPatch:  PATCH_BITMASK,
		http.MethodDelete: DELETE_BITMASK,
	}
	localconn *sqlitecloud.SQCloud = nil
)

func (this *Server) serveApiRest(writer http.ResponseWriter, request *http.Request) {
	// get apikey.
	apikey, err := getApiKeyFromHeader(request)
	if err != nil {
		writeError(writer, http.StatusUnauthorized, err.Error(), "")
		return
	}
	// fmt.Printf("ApiRest apikey: %s", apikey)

	// get project, database, table, id
	vars := mux.Vars(request)
	projectID := sqlitecloud.SQCloudEnquoteString(vars["projectID"])
	databaseName := sqlitecloud.SQCloudEnquoteString(vars["databaseName"])
	tableName := sqlitecloud.SQCloudEnquoteString(vars["tableName"])
	idstr, idfound := vars["id"]

	// extract int value for the id
	id := 0
	if idfound {
		id, err = strconv.Atoi(idstr)
		if err != nil {
			writeError(writer, http.StatusUnprocessableEntity, err.Error(), "")
			return
		}
	}

	// check if HTTP verb is enabled for project/database/table on auth db
	// issues:
	// - the auth db is not linked with the dbref db, it is not updated when a database/table is created, renamed or deleted
	// - enabled or disabled by default?
	if !isApiRestEnabled(request.Method, projectID, databaseName, tableName) {
		// TODO: The origin server MUST generate an Allow header field in a 405 response containing a list of the target resource's currently supported methods.
		// writer.Header()["Access-Control-Allow-Methods"] = ""
		writeError(writer, http.StatusMethodNotAllowed, "Not Allowed", "")
		return
	}

	// prepare the SQL query
	query, queryargs, statusCode, err, allowedMethods := apiRestQuery(request, projectID, databaseName, tableName, id)
	if err != nil {
		writeError(writer, statusCode, err.Error(), allowedMethods)
		return
	}

	// prepare the full command: "SWITCH APIKEY ?; <SQL>"
	query = "SWITCH APIKEY ?; " + query
	queryargs = prependInterface(queryargs, apikey)

	// call pool's ExecuteSQLArray
	// fmt.Printf("query:%s, args:%v\n", query, queryargs)
	res, err, _, _, _ := apicm.ExecuteSQLArray(projectID, query, &queryargs)
	if err != nil {
		writeError(writer, http.StatusInternalServerError, err.Error(), "")
		return
	}

	// Convert sqlitecloud.Result to map with ResultToObj
	resobj, err := ResultToObj(res)
	if err != nil {
		writeError(writer, http.StatusInternalServerError, err.Error(), "")
		return
	}

	// Parse the result object
	var value interface{}
	switch resobj.(type) {
	case map[string]interface{}:
		// Extract only the "rows" part of the rowset
		resmap, _ := resobj.(map[string]interface{})
		value = resmap["rows"]

	default:
		value = resobj
	}

	// reply with a JSON using the map result, in case of Insert maybe fetch the ID or the full new object
	// how can I get the last inserted ID? with DATABASE GET ROWID?
	statusCode = http.StatusOK
	writer.WriteHeader(statusCode)
	response := map[string]interface{}{"status": statusCode, "message": "OK", "value": value}
	bresponse, err := json.Marshal(response)
	if err != nil {
		writeError(writer, http.StatusInternalServerError, err.Error(), "")
	}

	writer.Write(bresponse)
}

func prependInterface(slice []interface{}, value interface{}) []interface{} {
	slice = append(slice, 0)
	copy(slice[1:], slice)
	slice[0] = value
	return slice
}

func getApiKeyFromHeader(r *http.Request) (string, error) {
	k := r.Header.Get(ApiKeyHeaderKey)
	if len(k) == 0 {
		return "", fmt.Errorf("ApiKey header not found")
	}
	return k, nil
}

func writeError(writer http.ResponseWriter, statusCode int, message string, allowedMethods string) {
	if statusCode == http.StatusMethodNotAllowed {
		writer.Header().Set("Access-Control-Allow-Methods", allowedMethods)
		writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		writer.Header().Set("Pragma", "no-cache")
		writer.Header().Set("Expires", "0")
	}
	writer.WriteHeader(statusCode)
	writer.Write([]byte(fmt.Sprintf("{\"status\":%d,\"message\":\"%s\"}", statusCode, message)))
}

func isApiRestEnabled(method string, projectID string, databaseName string, tableName string) bool {
	methodBitmask, found := methodsMap[strings.ToUpper(method)]
	if !found {
		return false
	}

	query := fmt.Sprintf("SELECT methods_mask FROM RestApiSettings WHERE project_uuid=? AND database_name=? AND table_name=? AND methods_mask & %d = %d", methodBitmask, methodBitmask)
	queryargs := []interface{}{projectID, databaseName, tableName, methodBitmask, methodBitmask}
	res, err, _, _, _ := apicm.ExecuteSQLArray("auth", query, &queryargs)

	if err != nil || res.GetNumberOfRows() == 0 {
		return false
	}

	return true
}

// apiRestQuery returns the query that implements the requested CRUD operation.
// The first return value is the query. The second return value is the array of
// values that must be binded to the query parameters, if any. The third return
// value is the statusCode the must be returned to the client in case of error.
// The fourth return value is the error, if any. The fifth return value is the
// list of allowedMethods if the returned statusCode is http.StatusMethodNotAllowed
func apiRestQuery(request *http.Request, projectID string, databaseName string, tableName string, id int) (string, []interface{}, int, error, string) {
	query := ""
	queryargs := []interface{}{}

	switch request.Method {
	case http.MethodGet:
		if len(tableName) == 0 {
			query = "SWITCH DATABASE ?; LIST TABLES"
			queryargs = append(queryargs, databaseName)

		} else {
			query = fmt.Sprintf("SWITCH DATABASE ?; SELECT * FROM %s", tableName)
			queryargs = append(queryargs, databaseName)
			if id > 0 {
				query += " WHERE _rowid_ = ?"
				queryargs = append(queryargs, id)
			}
		}

	case http.MethodPost:
		if len(tableName) == 0 {
			return "", []interface{}{}, http.StatusMethodNotAllowed, fmt.Errorf("Not Allowed"), http.MethodGet
		} else {
			query = fmt.Sprintf("SWITCH DATABASE ?; INSERT INTO %s", tableName)
			queryargs = append(queryargs, databaseName)
			if id != 0 {
				return "", []interface{}{}, http.StatusBadRequest, fmt.Errorf("Invalid request"), http.MethodGet
			} else {
				bodybytes, err := ioutil.ReadAll(request.Body)
				if err != nil {
					return "", []interface{}{}, http.StatusBadRequest, fmt.Errorf("Error reading body: %v", err), ""
				}

				bodymap := map[string]string{}
				err = json.Unmarshal(bodybytes, &bodymap)
				if err != nil {
					return "", []interface{}{}, http.StatusBadRequest, fmt.Errorf("Invalid body"), ""
				}
				if len(bodymap) == 0 {
					return "", []interface{}{}, http.StatusBadRequest, fmt.Errorf("Empty body"), ""
				}

				fields := ""
				values := ""
				for key, value := range bodymap {
					separator := ","
					if len(values) == 0 {
						separator = ""
					}
					fields = fmt.Sprintf("%s%s%s", fields, separator, sqlitecloud.SQCloudEnquoteString(key))
					values = fmt.Sprintf("%s%s?", values, separator)
					queryargs = append(queryargs, value)
				}
				query = fmt.Sprintf("%s (%s) VALUES (%s)", query, fields, values)
			}

			// add the command to return the last inserted row
			query += "; DATABASE GET ROWID;"
		}

	case http.MethodPatch:
		if len(tableName) == 0 {
			return "", []interface{}{}, http.StatusMethodNotAllowed, fmt.Errorf("Not Allowed"), http.MethodGet
		} else {
			query = fmt.Sprintf("SWITCH DATABASE ?; UPDATE %s SET ", tableName)
			queryargs = append(queryargs, databaseName)

			bodybytes, err := ioutil.ReadAll(request.Body)
			if err != nil {
				return "", []interface{}{}, http.StatusBadRequest, fmt.Errorf("Error reading body: %v", err), ""
			}

			bodymap := map[string]string{}
			err = json.Unmarshal(bodybytes, &bodymap)
			if err != nil {
				return "", []interface{}{}, http.StatusBadRequest, fmt.Errorf("Invalid body"), ""
			}
			if len(bodymap) == 0 {
				return "", []interface{}{}, http.StatusBadRequest, fmt.Errorf("Empty body"), ""
			}

			nsetvalues := 0
			for key, value := range bodymap {
				separator := ","
				if nsetvalues == 0 {
					separator = ""
				}
				nsetvalues += 1

				query += fmt.Sprintf("%s%s=?", separator, sqlitecloud.SQCloudEnquoteString(key))
				queryargs = append(queryargs, value)
			}

			if id > 0 {
				query += " WHERE _rowid_=?;"
				queryargs = append(queryargs, id)
			} else {
				return "", []interface{}{}, http.StatusBadRequest, fmt.Errorf("Full table update without filters is not allowed"), ""
			}
		}

	case http.MethodDelete:
		if len(tableName) == 0 {
			return "", []interface{}{}, http.StatusMethodNotAllowed, fmt.Errorf("Not Allowed"), http.MethodGet
		} else {
			query = fmt.Sprintf("SWITCH DATABASE ?; DELETE FROM %s", tableName)
			queryargs = append(queryargs, databaseName)

			if id > 0 {
				query += " WHERE _rowid_=?;"
				queryargs = append(queryargs, id)
			} else {
				return "", []interface{}{}, http.StatusBadRequest, fmt.Errorf("Full table delete without filters is not allowed"), ""
			}
		}
	}

	return query, queryargs, http.StatusOK, nil, ""
}
