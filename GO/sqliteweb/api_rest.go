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
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/swaggest/openapi-go/openapi3"
	"golang.org/x/exp/maps"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	// "github.com/gorilla/schema"
	// "github.com/gookit/validate"
	// "github.com/go-swagger/go-swagger"

	sqlitecloud "github.com/sqlitecloud/sdk"
)

const ApiKeyHeaderKey string = "X-SQLiteCloud-Api-Key"
const rootTableName string = ""

type tableColumnMetadata struct {
	Name       string
	DataType   string
	ColSeq     string
	NotNull    bool
	PrimaryKey bool
	AutoInc    bool
}

type methodMask int

const (
	GET_BITMASK methodMask = 1 << iota
	POST_BITMASK
	PATCH_BITMASK
	DELETE_BITMASK
)

var (
	methodBitmaskMap map[string]methodMask = map[string]methodMask{
		http.MethodGet:    GET_BITMASK,
		http.MethodPost:   POST_BITMASK,
		http.MethodPatch:  PATCH_BITMASK,
		http.MethodDelete: DELETE_BITMASK,
	}

	localconn *sqlitecloud.SQCloud = nil
)

func (this *Server) serveApiRest(writer http.ResponseWriter, request *http.Request) {
	this.Auth.cors(writer, request)

	start := time.Now()
	apikey := ""
	useridjwt := int64(-1)
	status := http.StatusOK
	var err error = nil

	defer func() {
		t := time.Since(start)
		errstring := ""
		if err != nil {
			errstring = err.Error()
		}
		user := ""
		if apikey != "" {
			user = hiddenApiKey(apikey)
		} else if useridjwt != -1 {
			user = fmt.Sprintf("%d", useridjwt)
		}
		SQLiteWeb.Logger.Infof("REST API: \"%s %s\" addr:%s user:%s exec_time:%s status:%d err:%s", request.Method, request.URL, request.RemoteAddr, user, t, status, errstring)
	}()

	// get project, database, table, id
	vars := mux.Vars(request)
	projectID := sqlitecloud.SQCloudEnquoteString(vars["projectID"])
	databaseName := sqlitecloud.SQCloudEnquoteString(vars["databaseName"])
	tableName := sqlitecloud.SQCloudEnquoteString(vars["tableName"])
	idstr, idfound := vars["id"]

	// OPTIONS method is used by preflight request (without auth)
	if request.Method == http.MethodOptions {
		methods := ""
		methods, err = optionsAllowedMethods(writer, request, projectID, databaseName, tableName)
		if err != nil {
			status = http.StatusInternalServerError
			writeError(writer, status, err.Error(), "")
			return
		}

		writer.Header().Set("Access-Control-Allow-Methods", methods)
		writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		writer.Header().Set("Pragma", "no-cache")
		writer.Header().Set("Expires", "0")
		writer.WriteHeader(status)
		return
	}

	// get auth (apikey or userid from jwt token)
	apikey, useridjwt, err = apiRestAuth(request)
	if err != nil {
		status = http.StatusUnauthorized
		writeError(writer, status, err.Error(), "")
		return
	}

	if useridjwt != -1 {
		projectID, _, err = verifyProjectID(useridjwt, projectID, apicm)
		if err != nil {
			status = http.StatusUnauthorized
			writeError(writer, status, err.Error(), "")
			return
		}
	}

	// extract int value for the id
	id := 0
	if idfound {
		id, err = strconv.Atoi(idstr)
		if err != nil {
			status = http.StatusUnprocessableEntity
			writeError(writer, status, err.Error(), "")
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
		status = http.StatusMethodNotAllowed
		err = errors.New("Not Allowed")
		writeError(writer, status, err.Error(), "")
		return
	}

	// prepare the SQL query
	query, queryargs, statusCode, err, allowedMethods := apiRestQuery(request, projectID, databaseName, tableName, id)
	if err != nil {
		status = statusCode
		writeError(writer, statusCode, err.Error(), allowedMethods)
		return
	}

	// prepare the full command: "SWITCH APIKEY ?; <SQL>"
	if apikey != "" {
		query = "SWITCH APIKEY ?; " + query
		queryargs = prependInterface(queryargs, apikey)
	}

	// call pool's ExecuteSQLArray
	// fmt.Printf("query:%s, args:%v\n", query, queryargs)
	res, err, _, _, _ := apicm.ExecuteSQLArray(projectID, query, &queryargs)
	if err != nil {
		status = http.StatusInternalServerError
		writeError(writer, status, err.Error(), "")
		return
	}

	jsonbytes, err := responseValueFromResult(request, projectID, databaseName, tableName, res)
	if err != nil {
		status = http.StatusInternalServerError
		writeError(writer, status, err.Error(), "")
		return
	}

	// reply with a JSON representation response value, in case of Insert maybe fetch the ID or the full new object
	// how can I get the last inserted ID? with DATABASE GET ROWID?
	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Content-Encoding", "utf-8")
	writer.WriteHeader(http.StatusOK)
	if err != nil {
		status = http.StatusInternalServerError
		writeError(writer, status, err.Error(), "")
	}

	writer.Write(jsonbytes)
}

func apiRestAuth(request *http.Request) (apikey string, useridjwt int64, err error) {
	apikey = ""
	useridjwt = -1
	err = nil

	apikey, err = getApiKeyFromHeader(request)
	if err != nil {
		if SQLiteWeb != nil {
			var jwterr error
			useridjwt, jwterr = SQLiteWeb.Auth.GetUserID(SQLiteWeb.Auth.getTokenFromAuthorization, request)
			if jwterr != nil {
				err = fmt.Errorf("%s, %s", err.Error(), jwterr.Error())
			} else {
				err = nil
			}
		}
	}

	return
}

func hiddenApiKey(k string) string {
	if len(k) <= 3 {
		return ""
	}
	return fmt.Sprintf("•••%s", k[len(k)-3:])
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
	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Content-Encoding", "utf-8")
	writer.WriteHeader(statusCode)
	writer.Write([]byte(fmt.Sprintf("{\"status\":%d,\"message\":\"%s\"}", statusCode, message)))
}

func optionsAllowedMethods(writer http.ResponseWriter, request *http.Request, projectID string, databaseName string, tableName string) (string, error) {
	query := fmt.Sprintf("SELECT methods_mask FROM RestApiSettings WHERE project_uuid=? AND database_name=? AND table_name=?")
	queryargs := []interface{}{projectID, databaseName, tableName}
	res, err, _, _, _ := apicm.ExecuteSQLArray("auth", query, &queryargs)
	if err != nil {
		return "", err
	}

	methodsmask := methodMask(0)
	if res.GetNumberOfRows() > 0 {
		methodsmask = methodMask(res.GetInt32Value_(0, 0))
	}

	methods := make([]string, 0)
	for verb, bitmask := range methodBitmaskMap {
		if methodsmask&bitmask == bitmask {
			methods = append(methods, verb)
		}
	}

	return strings.Join(methods, ", "), nil
}

func isApiRestEnabled(method string, projectID string, databaseName string, tableName string) bool {
	methodBitmask, found := methodBitmaskMap[strings.ToUpper(method)]
	if !found {
		return false
	}

	query := fmt.Sprintf("SELECT methods_mask FROM RestApiSettings WHERE project_uuid=? AND database_name=? AND table_name=? AND methods_mask & %d = %d", methodBitmask, methodBitmask)
	queryargs := []interface{}{projectID, databaseName, tableName}
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
		if tableName == rootTableName {
			query = "SWITCH DATABASE ?; LIST METADATA"
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
		if tableName == rootTableName {
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
			query += fmt.Sprintf("; SELECT * FROM %s ORDER BY _rowid_ DESC LIMIT 1", tableName)
		}

	case http.MethodPatch:
		if tableName == rootTableName {
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

			// add the command to return the modified row
			query += fmt.Sprintf("; SELECT * FROM %s WHERE _rowid_=?", tableName)
			queryargs = append(queryargs, id)
		}

	case http.MethodDelete:
		if tableName == rootTableName {
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

func responseValueFromResult(request *http.Request, projectID string, databaseName string, tableName string, result *sqlitecloud.Result) ([]byte, error) {
	if tableName == rootTableName {
		// returns a full OpenAPI description on the root path
		return openapiDocumentation(request, projectID, databaseName, result)
	}

	// Convert sqlitecloud.Result to map with ResultToObj
	resobj, err := ResultToObj(result)
	if err != nil {
		return nil, err
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

	jsonbytes, err := json.Marshal(value)
	return jsonbytes, err
}

// OpenAPI functions

type idPathRequest struct {
	ID string `path:"rowid"`
}

type defaultTableBody struct {
	Column string `json:"{column}"`
}

type errorResponse struct {
	Status  uint   `json:"status"`
	Message string `json:"message"`
}

type operation struct {
	path string
	op   openapi3.Operation
}

// openapiSetExposedMethods creates a map with all the possible ope.
// The returned string uses Go escape sequences (\t, \n, \xFF, \u0100)
// for control characters and non-printable characters as defined by IsPrint.
func openapiSetExposedMethods(reflector *openapi3.Reflector) (methodPaths map[string][]operation) {
	summaryGetAllRows := "Get all rows"
	summaryGetOneRow := "Get one row"
	summaryAddOneRow := "Add one row"
	summaryUpdateOneRow := "Update one row"
	summaryDeleteOneRow := "Delete one row"

	getAllOp := openapi3.Operation{Summary: &summaryGetAllRows}
	reflector.SetJSONResponse(&getAllOp, new([]defaultTableBody), http.StatusOK)
	reflector.SetJSONResponse(&getAllOp, new(errorResponse), http.StatusUnauthorized)
	reflector.SetJSONResponse(&getAllOp, new(errorResponse), http.StatusMethodNotAllowed)
	reflector.SetJSONResponse(&getAllOp, new(errorResponse), http.StatusInternalServerError)

	getOneOp := openapi3.Operation{Summary: &summaryGetOneRow}
	reflector.SetRequest(&getOneOp, new(idPathRequest), http.MethodGet)
	reflector.SetJSONResponse(&getOneOp, new([]defaultTableBody), http.StatusOK)
	reflector.SetJSONResponse(&getOneOp, new(errorResponse), http.StatusUnauthorized)
	reflector.SetJSONResponse(&getOneOp, new(errorResponse), http.StatusMethodNotAllowed)
	reflector.SetJSONResponse(&getOneOp, new(errorResponse), http.StatusInternalServerError)
	reflector.SetJSONResponse(&getOneOp, new(errorResponse), http.StatusUnprocessableEntity)

	postOp := openapi3.Operation{Summary: &summaryAddOneRow}
	reflector.SetRequest(&postOp, new(defaultTableBody), http.MethodPost)
	reflector.SetJSONResponse(&postOp, new([]defaultTableBody), http.StatusOK)
	reflector.SetJSONResponse(&postOp, new(errorResponse), http.StatusUnauthorized)
	reflector.SetJSONResponse(&postOp, new(errorResponse), http.StatusMethodNotAllowed)
	reflector.SetJSONResponse(&postOp, new(errorResponse), http.StatusInternalServerError)
	reflector.SetJSONResponse(&postOp, new(errorResponse), http.StatusBadRequest)

	patchOp := openapi3.Operation{Summary: &summaryUpdateOneRow}
	reflector.SetRequest(&patchOp, new(idPathRequest), http.MethodPatch)
	reflector.SetRequest(&patchOp, new(defaultTableBody), http.MethodPatch)
	reflector.SetJSONResponse(&patchOp, new([]defaultTableBody), http.StatusOK)
	reflector.SetJSONResponse(&patchOp, new(errorResponse), http.StatusUnauthorized)
	reflector.SetJSONResponse(&patchOp, new(errorResponse), http.StatusMethodNotAllowed)
	reflector.SetJSONResponse(&patchOp, new(errorResponse), http.StatusInternalServerError)
	reflector.SetJSONResponse(&patchOp, new(errorResponse), http.StatusUnprocessableEntity)
	reflector.SetJSONResponse(&patchOp, new(errorResponse), http.StatusBadRequest)

	deleteOp := openapi3.Operation{Summary: &summaryDeleteOneRow}
	reflector.SetRequest(&deleteOp, new(idPathRequest), http.MethodDelete)
	reflector.SetJSONResponse(&deleteOp, new(errorResponse), http.StatusUnauthorized)
	reflector.SetJSONResponse(&deleteOp, new(errorResponse), http.StatusMethodNotAllowed)
	reflector.SetJSONResponse(&deleteOp, new(errorResponse), http.StatusInternalServerError)
	reflector.SetJSONResponse(&deleteOp, new(errorResponse), http.StatusUnprocessableEntity)
	reflector.SetJSONResponse(&deleteOp, new(errorResponse), http.StatusBadRequest)

	// methodPaths cannot be a global var initializated once because I can't use a single
	// reflector object to create this map and a different reflector for each request,
	// otherwise the components schemas section would not be generated
	methodPaths = map[string][]operation{
		http.MethodGet: {
			{path: "", op: getAllOp},
			{path: "/{rowid}", op: getOneOp},
		},
		http.MethodPost: {
			{path: "", op: postOp},
		},
		http.MethodPatch: {
			{path: "/{rowid}", op: patchOp},
		},
		http.MethodDelete: {
			{path: "/{rowid}", op: deleteOp},
		},
	}

	return methodPaths
}

func dataTypeToType(dataType string) (t openapi3.SchemaType, err error) {
	t = openapi3.SchemaType("")
	uDataType := strings.ToUpper(dataType)
	switch uDataType {
	case "INT", "INTEGER", "TINYINT", "SMALLINT", "MEDIUMINT", "BIGINT", "UNSIGNED BIG INT", "INT2", "INT8":
		t = openapi3.SchemaTypeInteger
	case "BLOB":
		t = openapi3.SchemaTypeString
	case "REAL", "DOUBLE", "DOUBLE PRECISION", "FLOAT":
		t = openapi3.SchemaTypeNumber
	case "NUMERIC", "BOOLEAN", "DATE", "DATETIME":
		t = openapi3.SchemaTypeNumber
	default:
		switch {
		case strings.HasPrefix(uDataType, "CHARACTER"), strings.HasPrefix(uDataType, "VARYING CHARACTER"), strings.HasPrefix(uDataType, "NCHAR"), strings.HasPrefix(uDataType, "NATIVE CHARACTER"), strings.HasPrefix(uDataType, "NVARCHAR"), strings.HasPrefix(uDataType, "TEXT"), strings.HasPrefix(uDataType, "CLOB"):
			t = openapi3.SchemaTypeString
		case strings.HasPrefix(uDataType, "DECIMAL"):
			t = openapi3.SchemaTypeNumber
		default:
			t = openapi3.SchemaTypeObject
		}
	}

	return
}

func openapiServerURL(request *http.Request) string {
	url := *request.URL
	if url.Host == "" {
		url.Host = request.Host
	}
	if url.Scheme == "" {
		url.Scheme = "https"
	}
	return url.String()
}

func openapiDocumentation(request *http.Request, projectID string, databaseName string, metadataResult *sqlitecloud.Result) ([]byte, error) {
	// configure the openapi3 reflector
	reflector := openapi3.Reflector{}
	reflector.Spec = &openapi3.Spec{Openapi: "3.0.3"}
	reflector.Spec.Info.
		WithTitle(fmt.Sprintf("REST API for SQLiteCloud db %s", databaseName)).
		WithVersion("1.0.0")
		// .WithDescription("")
	reflector.Spec.WithServers(openapi3.Server{URL: openapiServerURL(request)})

	// Add security requirement
	securityApiKeyName := "api_key"
	securityJWTName := "bearer_token"

	// Declare security scheme.
	reflector.SpecEns().ComponentsEns().SecuritySchemesEns().WithMapOfSecuritySchemeOrRefValuesItem(
		securityApiKeyName,
		openapi3.SecuritySchemeOrRef{
			SecurityScheme: &openapi3.SecurityScheme{
				APIKeySecurityScheme: (&openapi3.APIKeySecurityScheme{}).
					WithName(ApiKeyHeaderKey).
					WithIn("header").
					WithDescription("API KEY Access"),
			},
		},
	).WithMapOfSecuritySchemeOrRefValuesItem(
		securityJWTName,
		openapi3.SecuritySchemeOrRef{
			SecurityScheme: &openapi3.SecurityScheme{
				HTTPSecurityScheme: (&openapi3.HTTPSecurityScheme{}).
					WithScheme("bearer").
					WithBearerFormat("JWT").
					WithDescription("JWT Access"),
			},
		},
	)
	reflector.Spec.WithSecurity(map[string][]string{securityApiKeyName: {}, securityJWTName: {}})

	// prepare an operation for each exposed method/table
	methodPaths := openapiSetExposedMethods(&reflector)
	tags := map[string]openapi3.Tag{}

	// get methods_mask for every table
	query := "SELECT table_name, methods_mask FROM RestApiSettings WHERE project_uuid=? AND database_name=?"
	queryargs := []interface{}{projectID, databaseName}
	resmethods, err, _, _, _ := apicm.ExecuteSQLArray("auth", query, &queryargs)
	if err != nil {
		return nil, err // errors.New("Cannot get REST API settings")
	}
	tablemethodsmasks := map[string]int32{}
	for r := uint64(0); r < resmethods.GetNumberOfRows(); r++ {
		tablename := resmethods.GetStringValue_(r, 0)
		methodsmask := resmethods.GetInt32Value_(r, 1)
		if tablename != "" {
			tablemethodsmasks[tablename] = methodsmask
		}
	}

	switch {
	case metadataResult.IsRowSet():
	default:
		return nil, errors.New("Unknown response format")
	}

	metadata := make(map[string][]tableColumnMetadata)
	for r := uint64(0); r < metadataResult.GetNumberOfRows(); r++ {
		tablename := metadataResult.GetStringValue_(r, 6)
		columnsMetadata := metadata[tablename]
		if columnsMetadata == nil {
			columnsMetadata = make([]tableColumnMetadata, 0, 5)
		}

		colMetadata := tableColumnMetadata{
			Name:       metadataResult.GetStringValue_(r, 0),
			DataType:   metadataResult.GetStringValue_(r, 1),
			ColSeq:     metadataResult.GetStringValue_(r, 2),
			NotNull:    metadataResult.GetInt32Value_(r, 3) == 1,
			PrimaryKey: metadataResult.GetInt32Value_(r, 4) == 1,
			AutoInc:    metadataResult.GetInt32Value_(r, 5) == 1,
		}

		columnsMetadata = append(columnsMetadata, colMetadata)
		metadata[tablename] = columnsMetadata
	}

	titlecaser := cases.Title(language.AmericanEnglish, cases.NoLower)
	defaultTableBodySchemaName := titlecaser.String(reflect.TypeOf(defaultTableBody{}).Name())
	defaultTableBodySchemaOrRef := reflector.Spec.Components.Schemas.MapOfSchemaOrRefValues[defaultTableBodySchemaName]

	for tablename, columnsMetadata := range metadata {
		mask := tablemethodsmasks[tablename]

		// customize the request and responses for each enabled operation with data from columnsMetadata
		for verb, bitmask := range methodBitmaskMap {
			if mask&int32(bitmask) == int32(bitmask) {
				// the method is enabled
				schemaName := fmt.Sprintf("%sRow", titlecaser.String(tablename))
				schemaNameOrRef := reflector.Spec.Components.Schemas.MapOfSchemaOrRefValues[schemaName]
				if defaultTableBodySchemaOrRef.Schema != nil && schemaNameOrRef.Schema == nil {
					// create a new schema for the operation's request
					// copy the schema from the default one (and remove the placeholder "{column}")
					newSchema := *(defaultTableBodySchemaOrRef.Schema)
					reflector.Spec.Components.Schemas.MapOfSchemaOrRefValues[schemaName] = openapi3.SchemaOrRef{Schema: &newSchema}
					delete(newSchema.Properties, "{column}")

					// add a field for each tables' column
					for _, colMetadata := range columnsMetadata {
						// fmt.Printf("colMetadata: %v\n", colMetadata)
						t, _ := dataTypeToType(colMetadata.DataType)
						t1 := colMetadata.DataType
						nullable := !colMetadata.NotNull
						schema := openapi3.Schema{
							Type:        &t,
							Nullable:    &nullable,
							Description: &t1,
						}
						newSchema.Properties[colMetadata.Name] = openapi3.SchemaOrRef{Schema: &schema}
					}
				}

				// replace the default schema name for each operation with the custom schema name
				for _, operation := range methodPaths[verb] {
					// customize the operation
					operation.op.Tags = []string{tablename}
					if _, found := tags[tablename]; !found {
						tdesc := fmt.Sprintf("Table %s", tablename)
						tags[tablename] = openapi3.Tag{Name: tablename, Description: &tdesc}
					}
					if operation.op.RequestBody != nil &&
						operation.op.RequestBody.RequestBody != nil {
						s, found := operation.op.RequestBody.RequestBody.Content["application/json"]
						if found && s.Schema != nil && s.Schema.SchemaReference != nil {
							s.Schema.SchemaReference.WithRef(filepath.Join("#/components/schemas", schemaName))
						}
					}
					if m := operation.op.Responses.MapOfResponseOrRefValues; m != nil {
						res, found := m["200"]
						if found && res.Response != nil {
							mt, found := res.Response.Content["application/json"]
							if found && mt.Schema != nil && mt.Schema.Schema != nil && mt.Schema.Schema.Items != nil && mt.Schema.Schema.Items.SchemaReference != nil {
								mt.Schema.Schema.Items.SchemaReference.WithRef(filepath.Join("#/components/schemas", schemaName))
							}
						}
					}

					// add the customized operation to the reflector's Spec
					path := filepath.Join("/", tablename, operation.path)
					err := reflector.Spec.AddOperation(verb, path, operation.op)
					if err != nil {
						return nil, err
					}
				}

			}
		}
	}

	// Add computed tags to Spec, one for each exposed table
	reflector.Spec.WithTags(maps.Values(tags)...)

	// remove the defaultTableBody schema
	delete(reflector.Spec.Components.Schemas.MapOfSchemaOrRefValues, defaultTableBodySchemaName)
	defaultTableBodySchemaOrRef.Schema = nil

	jsonbytes, err := reflector.Spec.MarshalJSON()
	return jsonbytes, nil
}
