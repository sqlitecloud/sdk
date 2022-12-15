//
//                    ////              SQLite Cloud
//        ////////////  ///
//      ///             ///  ///        Product     : SQLite Cloud Web Server
//     ///             ///  ///         Version     : 0.1.2
//     //             ///   ///  ///    Date        : 2022/02/0
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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/smtp"
	"strings"
	"text/template" // html/template
	"time"

	"github.com/Shopify/go-lua"
	sqlitecloud "github.com/sqlitecloud/sdk"
	//"github.com/gorilla/mux"
)

//import "bytes"
//import "time"
//import "context"
//import "errors"

//import "strconv"

//import "reflect"

//import "github.com/kardianos/service"

// import "github.com/gorilla/websocket"

// var db *sqlitecloud.SQCloud
// var out = bufio.NewWriter( os.Stdout )

const internalLuaFunctions = `
function filter( tree, map )
  newTree = {}
  for rowIndex = 1, #tree do
    newRow = {}
    for from, to in pairs( map ) do newRow[ to ] = tree[ rowIndex ][ from ] end
    newTree[ #newTree + 1 ] = newRow
  end
  return newTree
end
`

// Helper Functions

// Interface map to LUA Object (Table)
func MI2LUATable(L *lua.State, x map[string]interface{}) {
	// SQLiteWeb.Logger.Infof("MI2LUATable map: %v", x)
	L.NewTable()
	for k, m := range x {
		L.PushString(k)
		switch m.(type) {
		case bool:
			L.PushBoolean(m.(bool))
		case float64:
			L.PushNumber(m.(float64))
		case string:
			L.PushString(m.(string))
		case map[string]interface{}:
			MI2LUATable(L, m.(map[string]interface{}))
		case []interface{}:
			MI2LUAArray(L, m.([]interface{}))
		default:
			L.PushNil()
		}
		L.SetTable(-3)
	}
}

func MI2LUAArray(L *lua.State, x []interface{}) {
	// SQLiteWeb.Logger.Infof("MI2LUATable array: %v", x)
	L.NewTable()
	for k, m := range x {
		// SQLiteWeb.Logger.Infof("MI2LUAArray k: %v", k)
		L.PushInteger(k + 1)
		switch m.(type) {
		case bool:
			L.PushBoolean(m.(bool))
		case float64:
			L.PushNumber(m.(float64))
		case string:
			L.PushString(m.(string))
		case map[string]interface{}:
			MI2LUATable(L, m.(map[string]interface{}))
		case []interface{}:
			MI2LUAArray(L, m.([]interface{}))
		default:
			L.PushNil()
		}
		L.SetTable(-3)
	}
}

// LUA Array to Interface array
func LUAArray2IA(L *lua.State) []interface{} {
	array := make([]interface{}, 0)

	for i := 1; ; i++ {
		L.PushInteger(i) // get element at index i
		L.Table(-2)      // load it to the top of the stack
		switch L.TypeOf(-1) {
		case lua.TypeNil: // null is not allowed in an array and therefore: this markes the end...
			L.Pop(1)
			return array
		case lua.TypeBoolean:
			array = append(array, L.ToBoolean(-1))
		case lua.TypeNumber:
			array = append(array, lua.CheckNumber(L, -1))
		case lua.TypeString:
			array = append(array, lua.CheckString(L, -1))
		case lua.TypeTable:
			array = append(array, LUATable2MI(L))
		}
		L.Pop(1)
	}
	return array
}

// LUA Object to Interface map
func LUATable2MI(L *lua.State) map[string]interface{} {
	tree := make(map[string]interface{})

	switch L.TypeOf(1) {
	case lua.TypeTable:
		L.PushNil() // Add nil entry on stack (need 2 free slots).
		for L.Next(-2) {
			key := lua.CheckString(L, -2)
			switch L.TypeOf(-1) {
			case lua.TypeNil:
				tree[key] = nil
			case lua.TypeBoolean:
				tree[key] = L.ToBoolean(-1)
			case lua.TypeNumber:
				tree[key] = lua.CheckNumber(L, -1)
			case lua.TypeString:
				tree[key] = lua.CheckString(L, -1)
			case lua.TypeTable:
				L.PushNil() // Or Array??? -> Do a look ahead...
				// sanity check, avoid infinite loop in case of empty table
				if L.Next(-2) {
					// not empty table
					isArray := lua.CheckString(L, -2) == "1" // if nested key == 1 (this is prohipited and therefore an indicator for an array)
					L.Pop(2)

					switch isArray {
					case true:
						tree[key] = LUAArray2IA(L)
					default:
						tree[key] = LUATable2MI(L)
					}
				}
			}
			L.Pop(1)
		}
	}
	return tree
}

// json Functions

func lua_jsonEncode(L *lua.State) int {
	if L.TypeOf(1) == lua.TypeTable {
		jsonString, _ := json.Marshal(LUATable2MI(L))
		L.PushString(string(jsonString))
		return 1
	}
	L.PushString("")
	return 1
}

func lua_jsonDecode(L *lua.State) int {
	if L.TypeOf(1) == lua.TypeString {
		var x interface{} // map[string]interface{}

		if err := json.Unmarshal([]byte(lua.CheckString(L, 1)), &x); err == nil {
			// type-specific logic
			switch val := x.(type) {
			case map[string]interface{}:
				MI2LUATable(L, val)
				// SQLiteWeb.Logger.Infof("lua_jsonDecode map: %v", val)
			case []interface{}:
				MI2LUAArray(L, val)
				// SQLiteWeb.Logger.Infof("lua_jsonDecode arr: %v", val)
			}
			return 1
		} else {
			SQLiteWeb.Logger.Errorf("Error in lua_jsonDecode: %s", err)
		}
	}

	L.PushNil()
	return 1
}

// .ini File functions
func lua_getINIString(L *lua.State) int {
	switch {
	case L.TypeOf(1) != lua.TypeString:
		fallthrough // section
	case L.TypeOf(2) != lua.TypeString:
		fallthrough // key
	case L.TypeOf(3) != lua.TypeString:
		L.PushNil() // defaultValue
	default:
		L.PushString(GetINIString(lua.CheckString(L, 1), lua.CheckString(L, 2), lua.CheckString(L, 3)))

	}
	return 1
}
func lua_getINIBoolean(L *lua.State) int {
	X := L.TypeOf(3)
	print(X)
	switch {
	case L.TypeOf(1) != lua.TypeString:
		fallthrough // section
	case L.TypeOf(2) != lua.TypeString:
		fallthrough // key
	case L.TypeOf(3) != lua.TypeBoolean:
		L.PushNil() // default
	default:
		switch strings.ToLower(GetINIString(lua.CheckString(L, 1), lua.CheckString(L, 2), fmt.Sprintf("%t", L.ToBoolean(3)))) {
		case "1":
			fallthrough
		case "on":
			fallthrough
		case "true":
			fallthrough
		case "enable":
			fallthrough
		case "enabled":
			L.PushBoolean(true)
		default:
			L.PushBoolean(false)
		}
	}
	return 1
}

func lua_getINIArray(L *lua.State) int {
	switch {
	case L.TypeOf(1) != lua.TypeString:
		fallthrough // section
	case L.TypeOf(2) != lua.TypeString:
		fallthrough // key
	case L.TypeOf(3) != lua.TypeString:
		L.PushNil() // defaultValue
	default:
		serverList := []string{}
		for _, server := range strings.Split(L.PushString(GetINIString(lua.CheckString(L, 1), lua.CheckString(L, 2), lua.CheckString(L, 3))), ",") {
			server = strings.TrimSpace(server)
			if server != "" {
				serverList = append(serverList, server)
			}
		}
		if len(serverList) > 0 {
			L.NewTable()
			for i, server := range serverList {
				L.PushInteger(i + 1)
				L.PushString(server)
				L.SetTable(-3)
			}
		} else {
			L.PushNil()
		}
	}
	return 1
}

func lua_listINIProjects(L *lua.State) int {
	L.NewTable()
	i := 1
	for _, s := range cfg.SectionStrings() {
		switch {
		case len(s) != 36:
			continue // TODO: Check if section name matches regexp of uuid
		default:
			L.PushInteger(i)
			L.PushString(s)
			L.SetTable(-3)
			i++
		}
	}
	return 1
}

func lua_parseConnectionString(L *lua.State) int {
	switch {
	case L.TypeOf(1) != lua.TypeString:
		L.PushNil() // defaultValue
	default:
		if config, err := sqlitecloud.ParseConnectionString(lua.CheckString(L, 1)); err == nil {
			L.NewTable()
			L.PushString("Host")
			L.PushString(config.Host)
			L.SetTable(-3)
			L.PushString("Port")
			L.PushInteger(config.Port)
			L.SetTable(-3)
			L.PushString("Username")
			L.PushString(config.Username)
			L.SetTable(-3)
			L.PushString("Password")
			L.PushString(config.Password)
			L.SetTable(-3)
			L.PushString("Database")
			L.PushString(config.Database)
			L.SetTable(-3)
			L.PushString("Timeout")
			L.PushInteger(int(config.Timeout.Seconds()))
			L.SetTable(-3)
			L.PushString("Compress")
			L.PushString(config.CompressMode)
			L.SetTable(-3)
			L.PushString("Pem")
			L.PushString(config.Pem)
			L.SetTable(-3)
			L.PushString("ApiKey")
			L.PushString(config.ApiKey)
			L.SetTable(-3)
			L.PushString("NoBlob")
			L.PushBoolean(config.NoBlob)
			L.SetTable(-3)
			L.PushString("MaxData")
			L.PushInteger(config.MaxData)
			L.SetTable(-3)
			L.PushString("MaxRows")
			L.PushInteger(config.MaxRows)
			L.SetTable(-3)
			L.PushString("MaxRowset")
			L.PushInteger(config.MaxRowset)
			L.SetTable(-3)
		} else {
			L.PushNil()
		}
	}
	return 1
}

// SQLiteCloud helper

func lua_enquoteSQL(L *lua.State) int {
	switch L.TypeOf(1) {
	case lua.TypeNil:
		L.PushString("null")
	case lua.TypeBoolean:
		L.PushString(fmt.Sprintf("%t", L.ToBoolean(1)))
	case lua.TypeNumber:
		L.PushString(fmt.Sprintf("%f", lua.CheckNumber(L, 1)))
	case lua.TypeString:
		data := sqlitecloud.SQCloudEnquoteString(lua.CheckString(L, 1))
		if strings.HasPrefix(data, "\"") && strings.HasSuffix(data, "\"") {
			data = data[1 : len(data)-1]
		}
		L.PushString(data)
	default:
		L.PushNil()
	}
	return 1
}

func lua_reloadNodes(L *lua.State) int {
	if L.TypeOf(1) == lua.TypeString {
		uuid := lua.CheckString(L, 1)
		_, _ = cm.getNextServer(uuid, true)
		return 1
	}

	return 0
}

func lua_executeSQL(L *lua.State) int {
	if L.TypeOf(1) == lua.TypeString && L.TypeOf(2) == lua.TypeString {
		uuid := lua.CheckString(L, 1)
		query := lua.CheckString(L, 2)

		// get the optional arguments for the query
		// if the third argument is an array then convert it to a go slice
		var pargs *[]interface{} = nil
		if L.TypeOf(3) == lua.TypeTable {
			L.PushValue(3)
			L.PushNil()
			// sanity check, avoid infinite loop in case of empty table
			if L.Next(-2) {
				// not empty table
				isArray := lua.CheckString(L, -2) == "1" // if nested key == 1 (this is prohipited and therefore an indicator for an array)
				L.Pop(2)

				if isArray {
					args := LUAArray2IA(L)
					pargs = &args
				}
			}
		}
		// SQLiteWeb.Logger.Debugf("lua_executeSQL '%s', args: %v\n", query, pargs)

		res, err, errCode, extErrCode := cm.ExecuteSQLArray(uuid, query, pargs)

		if res != nil {
			defer res.Free()

			if err == nil {
				L.NewTable()

				errorNumber, extErrorCode, errorOffset, errorMessage, _ := res.GetError()
				L.PushString("ErrorNumber")
				L.PushInteger(errorNumber)
				L.SetTable(-3)
				L.PushString("ExtendedErrorNumber")
				L.PushInteger(extErrorCode)
				L.SetTable(-3)
				L.PushString("ErrorOffset")
				L.PushInteger(errorOffset)
				L.SetTable(-3)
				L.PushString("ErrorMessage")
				L.PushString(errorMessage)
				L.SetTable(-3)

				L.PushString("Value")
				if errorNumber == 0 && res.GetNumberOfRows() == 0 && res.GetNumberOfColumns() == 0 {
					L.PushString(res.GetString_())
				} else {
					L.PushNil()
				}
				L.SetTable(-3)

				L.PushString("NumberOfRows")
				L.PushInteger(int(res.GetNumberOfRows()))
				L.SetTable(-3)

				L.PushString("NumberOfColumns")
				L.PushInteger(int(res.GetNumberOfColumns()))
				L.SetTable(-3)

				null := uint64(0)

				L.PushString("Rows")
				if numRows := res.GetNumberOfRows(); numRows > 0 {
					L.NewTable() // row

					for r, R := null, numRows; r < R; r++ {
						L.PushInteger(int(r) + 1)

						L.NewTable() // columns
						for c, C := null, res.GetNumberOfColumns(); c < C; c++ {
							// L.PushInteger( int( c ) + 1 )
							L.PushString(res.GetName_(c))
							switch res.GetValueType_(r, c) {
							case '_':
								L.PushNil()
							case ':':
								L.PushInteger(int(res.GetInt32Value_(r, c)))
							case ',':
								L.PushNumber(res.GetFloat64Value_(r, c))
							default:
								L.PushString(res.GetStringValue_(r, c))
							}
							L.SetTable(-3)
						}
						L.SetTable(-3)
					}
				} else {
					L.PushNil()
				}
				L.SetTable(-3)

				L.PushString("Columns")
				if numCols := res.GetNumberOfColumns(); numCols > 0 {
					L.NewTable()
					for c, C := null, numCols; c < C; c++ {
						L.PushInteger(int(c) + 1)
						colname, _ := res.GetName(c)
						L.PushString(colname)
						L.SetTable(-3)
					}
				} else {
					L.PushNil()
				}
				L.SetTable(-3)

				//res.DumpToWriter( out, sqlitecloud.OUTFORMAT_LIST, false, "|", "NULL", "\r\n", 0, false )
				return 1
			}
		}

		if err != nil {
			L.NewTable()
			L.PushString("ErrorNumber")
			L.PushInteger(errCode)
			L.SetTable(-3)

			L.PushString("ExtendedErrorNumber")
			L.PushInteger(extErrCode)
			L.SetTable(-3)

			L.PushString("ErrorMessage")
			L.PushString(err.Error())
			L.SetTable(-3)

			L.PushString("Value")
			L.PushInteger(0)
			L.SetTable(-3)

			L.PushString("NumberOfRows")
			L.PushInteger(0)
			L.SetTable(-3)

			L.PushString("NumberOfColumns")
			L.PushInteger(0)
			L.SetTable(-3)

			L.PushString("Rows")
			L.NewTable() // row
			L.SetTable(-3)

			SQLiteWeb.Logger.Errorf("Error in ExecuteSQL: %s (%d:%d)", err, errCode, extErrCode)

			return 1
		}
	}

	return 0
}

// Email & Template helper

// mailTo( mailTo, Subject, template_data ), z.B. mailTo( "andreas@byte.watch", "welcome.eml", { To = "andreas.pfeil@..." } )
func mail(L *lua.State) int {
	host, _, err := net.SplitHostPort(cfg.Section("lua").Key("mail.proxy.host").String())

	switch {
	case err != nil:
		goto fail
	case L.TypeOf(1) != lua.TypeString:
		goto fail
	case L.TypeOf(2) != lua.TypeString:
		goto fail
	case L.TypeOf(3) != lua.TypeTable:
		goto fail

	default:
		auth := smtp.PlainAuth("", cfg.Section("lua").Key("mail.proxy.user").String(), cfg.Section("lua").Key("mail.proxy.password").String(), host)

		path := cfg.Section("lua").Key("mail.template.path").String()
		tempName := lua.CheckString(L, 1)
		language := lua.CheckString(L, 2)

		data := make(map[string]string)
		L.PushNil() // Add nil entry on stack (need 2 free slots).
		for L.Next(-2) {
			if L.TypeOf(-2) == lua.TypeString && L.TypeOf(-1) == lua.TypeString {
				data[lua.CheckString(L, -2)] = lua.CheckString(L, -1)
			}
			L.Pop(1)
		}

		now := time.Now()
		data["From"] = cfg.Section("lua").Key("mail.from").String()
		data["Time"] = now.Format("15:04:05")
		data["Date"] = now.Format("2006-01-02")
		data["Year"] = now.Format("2006")
		data["Month"] = now.Format("01")
		data["Day"] = now.Format("02")

		for _, path := range []string{fmt.Sprintf("%s/%s/%s", path, language, tempName), fmt.Sprintf("%s/%s", path, tempName)} {
			if !PathExists(path) {
				continue
			}

			if temp, err := template.ParseFiles(path); err == nil {

				var outBuffer bytes.Buffer
				err = temp.Execute(&outBuffer, data)
				fmt.Printf("%v (%s)\r\n", err, string(outBuffer.Bytes()))

				if err = smtp.SendMail(cfg.Section("lua").Key("mail.proxy.host").String(), auth, data["From"], []string{data["To"]}, outBuffer.Bytes()); err == nil {
					L.PushBoolean(true)
					return 1
				}
				// fmt.Printf( "%v\r\n", err )
			}
		}
	}

fail:
	L.PushBoolean(false)
	return 1
}

func (this *Server) executeLua(basePath string, endpoint string, userID int64, writer http.ResponseWriter, request *http.Request) {
	path := basePath

	args := []string{""}
	globals := make(map[string]string)

	now := time.Now()

	globals["userid"] = fmt.Sprintf("%d", userID)
	globals["method"] = strings.TrimSpace(strings.ToUpper(request.Method))
	globals["host"] = request.Host
	globals["client"] = request.RemoteAddr
	globals["uri"] = request.RequestURI
	globals["body"] = ""
	globals["now"] = now.Format("2006-01-02 15:04:05")
	globals["now_1h"] = now.Add(-1 * time.Hour).Format("2006-01-02 15:04:05")

	if body, err := ioutil.ReadAll(request.Body); err == nil {
		globals["body"] = string(body)
	}

NextPart:
	for _, part := range strings.Split(endpoint, "/") {
		if strings.TrimSpace(part) != "" {

			if dir, err := ioutil.ReadDir(path); err == nil {
				for _, fileinfo := range dir {
					fileName := fileinfo.Name()
					if fileName == part {
						args = append(args, part)
						args[0] = fmt.Sprintf("%s/%s", args[0], fileName)
						path = fmt.Sprintf("%s/%s", path, fileName)
						continue NextPart
					}
				}

				for _, fileinfo := range dir {
					fileName := fileinfo.Name()
					if strings.HasPrefix(fileName, "{") && strings.HasSuffix(fileName, "}") {
						globals[fileName[1:len(fileName)-1]] = part
						args = append(args, part)
						args[0] = fmt.Sprintf("%s/%s", args[0], fileName)
						path = fmt.Sprintf("%s/%s", path, fileName)
						continue NextPart
					}
				}
				sendError(writer, "Endpoint not found.", http.StatusNotFound)
				return
			}
		}
		//fmt.Printf( "i=%d, path=%s, part=%s\r\n", i, path, part )
		//fmt.Printf( "%s/{%d}\r\n", part, len( args ) )
	}

	// LUA only!!!
	script := fmt.Sprintf("%s.lua", globals["method"])
	if len(args) == 0 || len(args[0]) == 0 {
		sendError(writer, "Endpoint not found.", http.StatusNotFound)
		return
	}
	args = append(args, script)
	path = fmt.Sprintf("%s/%s", path, script)
	args[0] = fmt.Sprintf("%s/%s", args[0][1:], script)

	//fmt.Printf( "PATH=%s\r\n", path )

	if PathExists(path) {
		l := lua.NewState()
		lua.OpenLibraries(l)

		// register internal sql specific functons
		l.Register("parseConnectionString", lua_parseConnectionString)
		l.Register("executeSQL", lua_executeSQL)
		l.Register("reloadNodes", lua_reloadNodes)
		l.Register("enquoteSQL", lua_enquoteSQL)

		// register internal json related functions
		l.Register("jsonEncode", lua_jsonEncode)
		l.Register("jsonDecode", lua_jsonDecode)

		// register internal .ini file functions
		l.Register("getINIString", lua_getINIString)
		l.Register("listINIProjects", lua_listINIProjects)
		l.Register("getINIArray", lua_getINIArray)
		l.Register("getINIBoolean", lua_getINIBoolean)

		// register internal mail related functions
		l.Register("mail", mail)

		// register context related functions
		l.Register("SetStatus", func(L *lua.State) int {
			switch {
			case L.TypeOf(1) != lua.TypeNumber:
				break
			default:
				writer.WriteHeader(lua.CheckInteger(L, 1))
			}
			return 0
		})
		l.Register("SetHeader", func(L *lua.State) int {
			switch {
			case L.TypeOf(1) != lua.TypeString:
				break
			case L.TypeOf(2) != lua.TypeString:
				break
			default:
				writer.Header().Set(lua.CheckString(L, 1), lua.CheckString(L, 2))
			}
			return 0
		})
		l.Register("Write", func(L *lua.State) int {
			switch {
			case L.TypeOf(1) != lua.TypeString:
				break
			default:
				writer.Write([]byte(lua.CheckString(L, 1)))
			}
			return 0
		})

		// initialize the lua path
		l.NewTable()
		l.PushInteger(0)
		l.PushString(path) //l.PushString( fmt.Sprintf( "%s.lua", path ) )
		l.SetTable(-3)

		// create and populate the "query" array
		l.NewTable()
		for k, v := range request.URL.Query() {
			l.PushString(k)
			if len(v) > 1 {
				l.NewTable()
				for i, vv := range v {
					l.PushInteger(int(i) + 1)
					l.PushString(vv)
					l.SetTable(-3)
				}
				l.SetTable(-3)
			} else {
				l.PushString(v[0])
				l.SetTable(-3)
			}
			// fmt.Printf( "%s = %s\r\n", k, v )
		}
		l.SetGlobal("query")

		// create and populate the "globals" array
		for varName, varValue := range globals {
			l.PushString(varValue)
			l.SetGlobal(varName)
		}

		// create and populate the "args" array
		for i, arg := range args {
			l.PushInteger(i + 1)
			l.PushString(arg)
			l.SetTable(-3)
		}
		l.SetGlobal("args")

		// create and populate the "header" array
		l.NewTable()
		for k, m := range request.Header {
			l.PushString(k)
			l.NewTable() // row
			for i, v := range m {
				l.PushInteger(int(i) + 1)
				l.PushString(v)
				l.SetTable(-3)
			}
			l.SetTable(-3)
		}
		l.SetGlobal("header")

		// fmt.Printf( "will execute lua script: '%s'\r\n", path )
		// fmt.Printf( "%v\r\n", args )

		// set the lua libary path
		if PathExists(cfg.Section("lua").Key("package.path").String()) {
			lua.DoString(l, fmt.Sprintf(`package.path = "%s/?.lua"`, cfg.Section("lua").Key("package.path").String()))
		}

		// load internal lua functions into context
		lua.DoString(l, internalLuaFunctions)

		// execute the lua file
		err := lua.DoFile(l, path)
		if err != nil {
			SQLiteWeb.Logger.Errorf("panic in %s: %s", path, err)
			panic(err)
		}
	} else {
		sendError(writer, "Endpoint not found.", http.StatusNotFound)
	}
}
