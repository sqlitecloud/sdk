//
//                    ////              SQLite Cloud
//        ////////////  ///
//      ///             ///  ///        Product     : SQLite Cloud CLI Application
//     ///             ///  ///         Version     : 1.1.1
//     //             ///   ///  ///    Date        : 2021/10/08
//    ///             ///   ///  ///    Author      : Andreas Pfeil
//   ///             ///   ///  ///
//   ///     //////////   ///  ///      Description : Features: Connection Strings,
//   ////                ///  ///                     Batch processing, Many output
//     ////     //////////   ///                      formats, Line truncation for Terminals,
//        ////            ////                        History, Static & Dynamic Autocomplete
//          ////     /////
//             ///                      Copyright   : 2021 by SQLite Cloud Inc.
//
// -----------------------------------------------------------------------TAB=2

package main

import (
	sqlitecloud "github.com/sqlitecloud/sdk"

	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/docopt/docopt-go"
	"github.com/peterh/liner"
	"golang.org/x/term"
)

var app_name = "sqlc"
var long_name = "SQLite Cloud Command Line Application"
var version = "version 1.1.1"
var copyright = "(c) 2021 by SQLite Cloud Inc."
var history_file = fmt.Sprintf("~/.%s_history.txt", app_name)

var banner = fmt.Sprintf(`   _____
  /    /     %s, %s
 / ___/ /    %s
 \  ___/ /
  \_ ___/    Enter ".help" for usage hints.`, long_name, version, copyright)

var usage = long_name + ` Command Line Interface.

Usage:
  sqlc [URL] [options] [<FILE>...]
  sqlc -?|--help|--version

Arguments:
  URL                      "sqlitecloud://user:pass@host.com:port/dbname?timeout=10&compress=NO"
  FILE...                  Execute SQL commands from FILE(s) after connecting to the SQLite Cloud database

Examples:
  sqlc "sqlitecloud://user:pass@host.com:8860/dbname?timeout=10&compress=lz4&tls=intern"
  sqlc --host dev1.sqlitecloud.io -u user --password=pass -d dbname -c LZ4 --tls=no
  sqlc --version
  sqlc -?

General Options:
  --cmd COMMAND            Run "COMMAND" before executing FILE... or reading from stdin
  -l, --list               List available databases, then exit
  -d, --dbname NAME        Use database NAME
  -b, --bail               Stop after hitting an error
  -?, --help               Show this screen
  --version                Display version information

Output Format Options:
  -o, --output FILE        Switch to BATCH mode, execute SQL Commands and send output to FILE, then exit.
                           In BATCH mode, the default output format is switched to QUOTE.
  
  --echo                   Disables --quiet, print command(s) before execution
  --quiet                  Disables --echo, run command(s) quietly (no messages, only query output)
  --noheader               Turn headers off
  --nullvalue TEXT         Set text string for NULL values [default: "NULL"]
  --newline SEP            Set output row separator [default: "\r\n"]
  --separator SEP          Set output column separator [default::"|"]
  --format (LIST|CSV|QUOTE|TABS|LINE|JSON|HTML|XML|MARKDOWN|TABLE|BOX)
                           Specify the Output mode [default::BOX]

Connection Options:
  -h, --host HOSTNAME      Connect to SQLite Cloud database server host name [default::localhost]
  -p, --port PORT          Use specified port to connect to SQLIte Cloud database server [default::8860]
  -u, --user USERNAME      Use USERNAME for authentication
  -w, --password PASSWORD  Use PASSWORD for authentication
  -t, --timeout SECS       Set Timeout for network operations to SECS seconds [default::10]
  -c, --compress (NO|LZ4)  Use line compression [default::NO]
  --tls [YES|NO|INTERN|FILE] Encrypt the database connection using the host's root CA set (YES), a custom CA with a PEM from FILE (FILE), the internal SQLiteCloud CA (INTERN), or disable the encryption (NO) [default::YES]
`

var help = `
.help                Show this message
.bail [on|off]       Stop after hitting an error [default: off]
.echo [on|off]       Print command(s) before execution [default: off]
.quiet [on|off]      Run command(s) quietly (no messages, only query output) [default: on]
.noheader [on|off]   Turn table headers off or on [default: off]
.nullvalue TEXT      Set TEXT string for NULL values [default: "NULL"]
.newline TEXT        Set output row separator [default: "\r\n"]
.separator TEXT      Set output column separator [default: "<auto>"]
.format [LIST|CSV|QUOTE|TABS|LINE|JSON|HTML|XML|MARKDOWN|TABLE|BOX]
                     Specify the Output mode [default: BOX]
.width [-1|0|<num>]  Sets the maximum allowed query result length per line to the
                     terminal width(-1), unlimited (0) or any other width(<num>) [default: -1]
.timeout             Set Timeout for network operations to SECS seconds [default: 10]
.compress            Use line compression [default: NO]
.exit, .quit         Exit this program

If no parameter is specified, then the default value is used as the parameter value.
Boolean settings are toggled if no parameter is specified.
`

type Parameter struct {
	URL string `docopt:"URL"`

	OutFile      string `docopt:"--output"`
	Command      string `docopt:"--cmd"`
	List         bool   `docopt:"--list"`
	Bail         bool   `docopt:"--bail"`
	Echo         bool   `docopt:"--echo"`
	Quiet        bool   `docopt:"--quiet"`
	NoHeader     bool   `docopt:"--noheader"`
	NullText     string `docopt:"--nullvalue"`
	NewLine      string `docopt:"--newline"`
	Separator    string `docopt:"--separator"`
	Format       string `docopt:"--format"`
	OutPutFormat int    `docopt:"--outputformat"`

	Host      string `docopt:"--host"`
	Port      int    `docopt:"--port"`
	User      string `docopt:"--user"`
	Password  string `docopt:"--password"`
	Database  string `docopt:"--dbname"`
	ApiKey    string `docopt:"--apikey"`
	NoBlob    bool   `docopt:"--noblob"`
	MaxData   int    `docopt:"--maxdata"`
	MaxRows   int    `docopt:"--maxrows"`
	MaxRowset int    `docopt:"--maxrowset"`

	Timeout  int      `docopt:"--timeout"`
	Compress string   `docopt:"--compress"`
	Tls      string   `docopt:"--tls"`
	UseStdIn bool     `docopt:"-"`
	Files    []string `docopt:"<FILE>"`
}

var tokens = []string{".echo ", ".help ", ".bail ", ".quiet ", ".noheader ", ".nullvalue ", ".newline ", ".separator ", ".format ",
	".width ", ".quit ",
	"SELECT ", "AS ", "FROM ", "JOIN ", "ON ", "USING ", "WHERE ", "LIKE ", "OR ", "AND ", "GROUP ", "ORDER ", "BY ", "ASC ", "DESC ", "LIMIT ", "TO ",
	"INSERT ", "UPDATE ", "DROP ", "IF ", "NOT ", "EXISTS ", "FAIL ", "IGNORE ", "TABLE ", "VALUES ", "SET ", "INTO ",
	"CREATE ", "ALTER ", "NULL ", "INTEGER ", "TEXT ", "PRIMARY ", "UNIQUE ", "DEFAULT ",
	"ABORT ", "ACTION ", "AFTER ", "ALL ", "ALWAYS ", "ANALYZE ", "ADD ", "ATTACH ",
	"AUTOINCREMENT ", "BEFORE ", "BEGIN ", "BETWEEN ", "CASCADE ", "CASE ", "CAST ", "CHECK ", "COLLATE ", "COLUMN ",
	"COMMIT ", "CONFLICT ", "CONSTRAINT ", "CROSS ", "CURRENT ", "CURRENT_DATE ", "CURRENT_TIME ",
	"CURRENT_TIMESTAMP ", "DATABASE ", "DEFERRABLE ", "DEFERRED ", "DELETE ", "DETACH ", "DISTINCT ",
	"DO ", "EACH ", "ELSE ", "END ", "ESCAPE ", "EXCEPT ", "EXCLUDE ", "EXCLUSIVE ", "EXPLAIN ",
	"FILTER ", "FIRST ", "FOLLOWING ", "FOR ", "FOREIGN ", "FULL ", "GENERATED ", "GLOB ", "GROUPS ",
	"HAVING ", "IMMEDIATE ", "IN ", "INDEX ", "INDEXED ", "INITIALLY ", "INNER ", "INSTEAD ",
	"INTERSECT ", "IS ", "ISNULL ", "KEY ", "LAST ", "LEFT ", "MATCH ", "MATERIALIZED ",
	"NATURAL ", "NO ", "NOTHING ", "NOTNULL ", "NULLS ", "OF ", "OFFSET ",
	"OTHERS ", "OUTER ", "OVER ", "PARTITION ", "PLAN ", "PRAGMA ", "PRECEDING ", "QUERY ", "RAISE ", "RANGE ",
	"RECURSIVE ", "REFERENCES ", "REGEXP ", "REINDEX ", "RELEASE ", "RENAME ", "REPLACE ", "RESTRICT ", "RETURNING ",
	"RIGHT ", "ROLLBACK ", "ROW ", "ROWS ", "SAVEPOINT ", "TEMP ", "TEMPORARY ", "THEN ",
	"TIES ", "TRANSACTION ", "TRIGGER ", "UNBOUNDED ", "UNION ", "VACUUM ",
	"VIEW ", "VIRTUAL ", "WHEN ", "WINDOW ", "WITH ", "WITHOUT ",
	"ABS( ", "CHANGES( ", "CHAR( ", "COALESCE( ", "GLOB( ", "HEX( ", "IFNULL( ", "IIF( ", "INSTR( ", "LAST_INSERT_ROWID( ",
	"LENGTH( ", "LIKE( ", "LIKELIHOOD( ", "LIKELY( ", "LOAD_EXTENSION( ", "LOWER( ", "LTRIM( ", "MAX( ", "MIN( ",
	"NULLIF( ", "PRINTF( ", "QUOTE( ", "RANDOM() ", "RANDOMBLOB( ", "REPLACE( ", "ROUND( ", "RTRIM( ", "SIGN( ",
	"SOUNDEX( ", "SQLITE_COMPILEOPTION_GET( ", "SQLITE_COMPILEOPTION_USED( ", "SQLITE_OFFSET( ", "SQLITE_SOURCE_ID() ",
	"SQLITE_VERSION() ", "SUBSTR( ", "SUBSTRING( ", "TOTAL_CHANGES() ", "TRIM( ", "TYPEOF( ", "UNICODE( ", "UNLIKELY( ",
	"UPPER( ", "ZEROBLOB( ", "DATE( ", "TIME( ", "DATETIME( ", "JULIANDAY( ", "STRFTIME( ", "DAYS ", "HOURS ", "MINUTES ",
	"NOW", "SECONDS ", "MONTHS ", "YEARS ", "START OF MONTH ", "START OF YEAR ", "START OF DAY ", "WEEKDAY ", "UNIXEPOCH ",
	"LOCALTIME ", "UTC ", "AVG( ", "COUNT( * ) ", "COUNT( ", "GROUP_CONCAT( ", "SUM( ", "TOTAL( ", "ACOS( ", "ACOSH( ",
	"ASIN( ", "ASINH( ", "ATAN( ", "ATAN2( ", "ATANH( ", "CEIL( ", "CEILING( ", "COS( ", "COSH( ", "DEGREES( ", "EXP( ",
	"FLOOR( ", "LN( ", "LOG( ", "LOG10( ", "LOG2( ", "MOD( ", "PI() ", "POW( ", "POWER( ", "RADIANS( ", "SIN( ", "SINH( ",
	"SQRT( ", "TAN( ", "TANH( ", "TRUNC( ", "JSON( ", "JSON_ARRAY( ", "JSON_ARRAY_LENGTH( ", "JSON_EXTRACT( ",
	"JSON_INSERT( ", "JSON_OBJECT( ", "JSON_PATCH( ", "JSON_REMOVE( ", "JSON_REPLACE( ", "JSON_SET( ", "JSON_TYPE( ",
	"JSON_VALID( ", "JSON_QUOTE( ", "JSON_GROUP_ARRAY( ", "JSON_GROUP_OBJECT( ", "JSON_EACH( ", "JSON_TREE( ",
	"LIST TABLES", "TABLES", "LIST DATABASES", "DATABASES", "LIST COMMANDS", "COMMANDS", "LIST INFO", "INFO",
	"AUTH USER ", "USER", "PASS", "CLOSE CONNECTION ", "CONNECTION", "CREATE DATABASE ", "DATABASE",
	"DISABLE PLUGIN ", "PLUGIN", "DROP DATABASE ", "DROP KEY ", "ENABLE PLUGIN ",
	"GET DATABASE", "GET DATABASE ID", "ID", "GET KEY ", "LIST CONNECTIONS", "CONNECTIONS", "LIST DATABASE CONNECTIONS ",
	"LIST DATABASE CONNECTIONS ID ", "LIST NODES", "LIST PLUGINS", "LIST CLIENT KEYS",
	"LIST DATABASE KEYS", "LISTEN ", "NOTIFY ", "PING", "REMOVE NODE ", "SET KEY ", "UNLISTEN ", "UNUSE DATABASE", "USE DATABASE",
	"on", "off", "enable", "disable", "true", "false",
}
var dynamic_tokens = []string{}

func replaceControlChars(in string) string {
	for from, to := range map[string]string{"\\0": string(0), "\\a": "\a", "\\b": "\b", "\\t": "\t", "\\n": "\n", "\\v": "\v", "\\f": "\f", "\\r": "\r"} {
		in = strings.ReplaceAll(in, from, to)
	}
	return in
}
func getFirstNoneEmptyString(args []string) string {
	for _, v := range args {
		v = strings.TrimSpace(v)
		if v != "" {
			return v
		}
	}
	return ""
}
func parseParameters() (Parameter, error) {
	parameter := Parameter{}
	var outputformat int

	// Parse Command Line Parameter (Attention: "::" -> ":<CTRL+Space>")
	if p, err := docopt.ParseArgs(strings.ReplaceAll(usage, "::", ": "), nil, fmt.Sprintf("%s %s, %s", app_name, version, copyright)); err == nil {

		// Preprocessing...
		if format, _ := p.String("--format"); format == "" { // If --format was not specified, use default values...
			if list, _ := p.Bool("--list"); list {
				p["--format"] = "LIST"
			} // use --format=LIST when in list mode
			if output, _ := p.String("--output"); output != "" {
				p["--format"] = "QUOTE"
			} // use --format=QUOTE when in batch mode
		}
		if format, _ := p.String("--format"); format == "" {
			p["--format"] = "BOX"
		} // use --format=BOX when --format is stil not specified

		format, err := p.String("--format")
		if err != nil {
			return Parameter{}, err
		}

		outputformat, err = sqlitecloud.GetOutputFormatFromString(format)
		if err != nil {
			return Parameter{}, err
		}

		p["--outputformat"] = outputformat

		// If the connection string is set, parse and apply the connection string...
		if url, isSet := p["URL"]; isSet && url != "<nil>" {
			if conf, err := sqlitecloud.ParseConnectionString(reflect.ValueOf(url).String()); err == nil {
				p["--host"] = getFirstNoneEmptyString([]string{dropError(p.String("--host")), conf.Host})
				p["--user"] = getFirstNoneEmptyString([]string{dropError(p.String("--user")), conf.Username})
				p["--password"] = getFirstNoneEmptyString([]string{dropError(p.String("--password")), conf.Password})
				p["--dbname"] = getFirstNoneEmptyString([]string{dropError(p.String("--dbname")), conf.Database})
				p["--host"] = getFirstNoneEmptyString([]string{dropError(p.String("--host")), conf.Host})
				p["--compress"] = getFirstNoneEmptyString([]string{dropError(p.String("--compress")), conf.CompressMode})
				if conf.Port > 0 {
					p["--port"] = getFirstNoneEmptyString([]string{dropError(p.String("--port")), fmt.Sprintf("%d", conf.Port)})
				}
				if conf.Timeout > 0 {
					p["--timeout"] = getFirstNoneEmptyString([]string{dropError(p.String("--timeout")), fmt.Sprintf("%d", conf.Timeout)})
				}
				p["--tls"] = getFirstNoneEmptyString([]string{dropError(p.String("--tls")), conf.Pem})
				p["--apikey"] = getFirstNoneEmptyString([]string{dropError(p.String("--apikey")), conf.ApiKey})
				if conf.NoBlob {
					p["--noblob"] = getFirstNoneEmptyString([]string{dropError(p.String("--noblob")), strconv.FormatBool(conf.NoBlob)})
				}
				if conf.MaxData > 0 {
					p["--maxdata"] = getFirstNoneEmptyString([]string{dropError(p.String("--maxdata")), fmt.Sprintf("%d", conf.MaxData)})
				}
				if conf.MaxRows > 0 {
					p["--maxrows"] = getFirstNoneEmptyString([]string{dropError(p.String("--maxrows")), fmt.Sprintf("%d", conf.MaxRows)})
				}
				if conf.MaxRowset > 0 {
					p["--maxrowset"] = getFirstNoneEmptyString([]string{dropError(p.String("--maxrowset")), fmt.Sprintf("%d", conf.MaxRowset)})
				}
			}
		} else {
			return Parameter{}, err
		}

		// Set default Values
		p["--host"] = getFirstNoneEmptyString([]string{dropError(p.String("--host")), "localhost"})
		p["--port"] = getFirstNoneEmptyString([]string{dropError(p.String("--port")), "8860"})
		p["--timeout"] = getFirstNoneEmptyString([]string{dropError(p.String("--timeout")), "10"})
		p["--compress"] = getFirstNoneEmptyString([]string{dropError(p.String("--compress")), "NO"})
		p["--tls"] = getFirstNoneEmptyString([]string{dropError(p.String("--tls")), "YES"})
		p["--separator"] = getFirstNoneEmptyString([]string{dropError(p.String("--separator")), dropError(sqlitecloud.GetDefaultSeparatorForOutputFormat(outputformat)), "|"})

		// Fix invalid(=unset) parameters, quotation & control-chars
		for k, v := range p {
			switch reflect.ValueOf(v).Kind() {
			case reflect.Invalid:
				p[k] = ""
			case reflect.String:
				p[k] = replaceControlChars(strings.Trim(reflect.ValueOf(v).String(), "'\""))
			default:
			}
		}

		// for k, v := range p { fmt.Printf( "%s='%v'\r\n", k, v ) }

		// Copy map data into Object
		if err := p.Bind(&parameter); err != nil {
			return Parameter{}, err
		}

		// Postprocessing...
		if parameter.OutFile != "" {
			parameter.Echo = false
			parameter.Quiet = true
		}
		if q, _ := p.Bool("--echo"); q {
			parameter.Echo = true
		}
		if parameter.Echo {
			parameter.Quiet = false
		}
		if q, _ := p.Bool("--quiet"); q {
			parameter.Quiet = true
			parameter.Echo = false
		}

	} else {
		return Parameter{}, err
	}
	return parameter, nil
}

func autocomplete(line string, pos int) (head string, suggestions []string, tail string) {
	start := line[0:pos]
	tail = strings.TrimPrefix(line[pos:], " ")
	split := strings.LastIndex(start, " ")
	head = ""
	line = start
	if split > 0 {
		head = start[0 : split+1]
		line = start[split+1:]
	}
	// fmt.Printf( "start=>%s<, tail=>%s<, head=>%s<, line=>%s<\r\n", start, tail, head, line )

	for _, token := range dynamic_tokens {
		if strings.HasPrefix(strings.ToLower(token), strings.ToLower(line)) {
			suggestions = append(suggestions, token)
		}
	}
	for _, token := range tokens {
		if strings.HasPrefix(strings.ToLower(token), strings.ToLower(line)) {
			suggestions = append(suggestions, token)
		}
	}
	return
}

func main() {
	width := -1
	out := bufio.NewWriter(os.Stdout)

	if parameter, err := parseParameters(); err != nil {
		fatal(out, fmt.Sprintf("ERROR: Could not parse arguments (%s).", err.Error()), &parameter)
	} else {
		if parameter.OutFile != "" {
			_ = os.Remove(parameter.OutFile)
			if file, err := os.OpenFile(parameter.OutFile, os.O_CREATE|os.O_RDWR, 0664); err != nil {
				fatal(out, fmt.Sprintf("ERROR: Could not create '%s' for writing", parameter.OutFile), &parameter)
			} else {
				out = bufio.NewWriter(file)
				width = 0
				defer func() {
					if err := file.Close(); err != nil {
						fatal(out, fmt.Sprintf("ERROR: Could not close file '%s': %s", parameter.OutFile, err), &parameter)
					}
				}()
			}
		}

		// print( out, fmt.Sprintf( "%s %s, %s", long_name, version, copyright ), &parameter )

		config := sqlitecloud.SQCloudConfig{
			Host:         parameter.Host,
			Port:         parameter.Port,
			Username:     parameter.User,
			Password:     parameter.Password,
			Database:     parameter.Database,
			Timeout:      time.Duration(parameter.Timeout) * time.Second,
			CompressMode: parameter.Compress,
			ApiKey:       parameter.ApiKey,
			NoBlob:       parameter.NoBlob,
			MaxData:      parameter.MaxData,
			MaxRows:      parameter.MaxRows,
			MaxRowset:    parameter.MaxRowset,
		}

		config.Secure, config.Pem = sqlitecloud.ParseTlsString(parameter.Tls)
		var db *sqlitecloud.SQCloud = sqlitecloud.New(config)

		if err := db.Connect(); err != nil {
			fatal(out, fmt.Sprintf("ERROR: %s.", err.Error()), &parameter)
		} else {
			defer db.Close()

			db.Callback = func(conn *sqlitecloud.SQCloud, json string) {
				print(out, json, &parameter)
				out.Flush()
			}

			//print( out, fmt.Sprintf( "Connected to %s.", parameter.Host ), &parameter )
			//print( out, strings.Split( help, "\n" )[ 0 ], &parameter )

			print(out, strings.ReplaceAll(banner, "<HOST>", parameter.Host), &parameter)
			print(out, "", &parameter)

			if parameter.List {
				Execute(db, out, "LIST DATABASES", width, &parameter)
				os.Exit(0)
			}

			// Execute single Command
			if parameter.Command != "" {
				Execute(db, out, parameter.Command, width, &parameter)
			}

			// Batch Mode starts here ///////////////////////////////////////////////////////////////

			// Execute Files
			if len(parameter.Files) > 0 {
				ExecuteFiles(db, out, parameter.Files, width, &parameter)
			}
			// Execute Stdin
			if parameter.UseStdIn {
				if err := ExecuteBuffer(db, out, os.Stdin, width, &parameter); err != nil {
					bail(out, fmt.Sprintf("Could not execute (%s)", err.Error()), &parameter)
				}
			}
			// End Batch Mode
			if parameter.OutFile != "" {
				os.Exit(0)
			}

			// Interactive Mode starts here /////////////////////////////////////////////////////////

			editor := liner.NewLiner()
			defer editor.Close()

			editor.SetCtrlCAborts(true)
			editor.SetMultiLineMode(true)
			editor.SetWordCompleter(autocomplete)

			if strings.HasPrefix(history_file, "~/") { // fix Home directory for POSIX and Windows
				if dir, err := os.UserHomeDir(); err == nil {
					history_file = fmt.Sprintf("%s/%s", dir, strings.TrimPrefix(history_file, "~/"))
				}
			}
			if f, err := os.Open(history_file); err == nil {
				editor.ReadHistory(f)
				f.Close()
			}

			prompt := "sqlc > "
			prompt = "\\H:\\p/\\d\\u > "

		Loop:
			for {
				out.Flush()
				db.Database, _ = db.GetDatabase()

				renderdPrompt := prompt
				renderdPrompt = strings.ReplaceAll(renderdPrompt, "\\H", parameter.Host)
				renderdPrompt = strings.ReplaceAll(renderdPrompt, "\\p", fmt.Sprintf("%d", parameter.Port))
				renderdPrompt = strings.ReplaceAll(renderdPrompt, "\\u", parameter.User)
				renderdPrompt = strings.ReplaceAll(renderdPrompt, "\\d", db.Database)

				renderdPrompt = strings.ReplaceAll(renderdPrompt, "\\T", time.Now().Format("2006-01-02"))
				renderdPrompt = strings.ReplaceAll(renderdPrompt, "\\t", time.Now().Format("15:04:05"))
				renderdPrompt = strings.ReplaceAll(renderdPrompt, "\\w", fmt.Sprintf("/%s", dropError(os.Getwd())))

				go func() { dynamic_tokens = db.GetAutocompleteTokens() }() // Update the dynamic tokens in the background...
				command, err := editor.Prompt(renderdPrompt)

				switch err {
				case io.EOF:
					break
				case nil:
					switch tokens := strings.Split(strings.ToLower(strings.TrimSpace(command)), " "); tokens[0] {
					case "":
						continue Loop
					case ".exit", "exit":
						break Loop
					case ".help":
						print(out, help, &parameter)
					case ".bail":
						parameter.Bail = getNextTokenValueAsBool(out, parameter.Bail, tokens, &parameter)
					case ".echo":
						parameter.Echo = getNextTokenValueAsBool(out, parameter.Echo, tokens, &parameter)
					case ".quiet":
						parameter.Quiet = getNextTokenValueAsBool(out, parameter.Quiet, tokens, &parameter)
					case ".noheader":
						parameter.NoHeader = getNextTokenValueAsBool(out, parameter.NoHeader, tokens, &parameter)
					case ".width":
						width = getNextTokenValueAsInteger(out, width, -1, tokens, &parameter)
					case ".nullvalue":
						parameter.NullText = getNextTokenValueAsString(out, parameter.NullText, "NULL", "", tokens, &parameter)
					case ".newline":
						parameter.NewLine = getNextTokenValueAsString(out, parameter.NewLine, "\r\n", "", tokens, &parameter)
					case ".prompt":
						prompt = getNextTokenValueAsString(out, prompt, "sqlc >", "", tokens, &parameter)
					case ".separator":
						parameter.Separator = getNextTokenValueAsString(out, parameter.Separator, "<auto>", "", tokens, &parameter)
					case ".timeout":
						parameter.Timeout = getNextTokenValueAsInteger(out, parameter.Timeout, 10, tokens, &parameter)
					case ".compress":
						parameter.Compress = getNextTokenValueAsString(out, parameter.Compress, "NO", "|no|lz4|", tokens, &parameter)
						db.Compress(parameter.Compress)

					case ".format":
						newFormat := getNextTokenValueAsString(out, parameter.Format, "BOX", "|list|csv|quote|tabs|line|json|html|xml|markdown|table|box|", tokens, &parameter)
						if newFormat != parameter.Format {
							parameter.Format = newFormat
							parameter.Separator = "<auto>"
							parameter.OutPutFormat, _ = sqlitecloud.GetOutputFormatFromString(newFormat)
						}

					default:
						Execute(db, out, command, width, &parameter)
						switch tokens[0] {
						case ".quit", "quit":
							break Loop
						default:
							editor.AppendHistory(command)
						}
					}
				default:
					bail(out, err.Error(), &parameter)
				}
			}

			if f, err := os.Create(history_file); err == nil {
				_, _ = editor.WriteHistory(f)
				_ = f.Close()
			}
		}
	}
}

func Execute(db *sqlitecloud.SQCloud, out *bufio.Writer, cmd string, width int, Settings *Parameter) {
	Seperator := Settings.Separator
	if cmd == "" {
		return
	}
	if Settings.OutPutFormat == sqlitecloud.OUTFORMAT_XML {
		Seperator = cmd
	}
	if width < 0 {
		if w, _, err := term.GetSize(0); err == nil {
			width = w
		}
	}
	if Settings.Echo {
		print(out, cmd, Settings)
	}

	start := time.Now()
	if res, err := db.Select(cmd); res != nil {
		defer res.Free()

		if err == nil {
			if !res.IsRowSet() || (res.GetNumberOfRows() > 0 && res.GetNumberOfColumns() > 0) {
				res.DumpToWriter(out, Settings.OutPutFormat, Settings.NoHeader, Seperator, Settings.NullText, Settings.NewLine, uint(width), Settings.Quiet)
			}
			if res.IsRowSet() {
				print(out, fmt.Sprintf("Rows: %d - Cols: %d: %s Time: %s", res.GetNumberOfRows(), res.GetNumberOfColumns(), renderByteCount(int64(res.GetUncompressedChuckSizeSum())), time.Since(start)), Settings)
			}
		}
	} else {
		bail(out, err.Error(), Settings)
		// bail(out, err.Error()+fmt.Sprintf(" (%d:%d:%d)", db.GetErrorCode(), db.GetExtErrorCode(), db.GetErrorOffset()), Settings)
	}
	print(out, "", Settings) // Empty line
}

func ExecuteBuffer(db *sqlitecloud.SQCloud, out *bufio.Writer, in *os.File, width int, Settings *Parameter) error {
	if scanner := bufio.NewScanner(in); scanner != nil {
		for scanner.Scan() {
			line := scanner.Text()
			if strings.ToUpper(line) == ".GOTO_PROMPT" {
				return nil // break out of sql script
			}
			if strings.TrimSpace(line) != "" {
				Execute(db, out, line, width, Settings)
			}
		}
		return scanner.Err() // nil or some error
	}
	return errors.New("Could not instanciate the line scanner")
}
func ExecuteFile(db *sqlitecloud.SQCloud, out *bufio.Writer, FilePath string, width int, Settings *Parameter) error {
	file, err := os.Open(FilePath)
	if err == nil {
		defer func() {
			_ = file.Close()
		}()
		if err := ExecuteBuffer(db, out, file, width, Settings); err != nil {
			return err
		}
	}
	return err
}
func ExecuteFiles(db *sqlitecloud.SQCloud, out *bufio.Writer, FilePathes []string, width int, Settings *Parameter) error {
	for _, file := range FilePathes {
		err := ExecuteFile(db, out, file, width, Settings)
		if err != nil {
			return err
		}
	}
	return nil
}

func print(out *bufio.Writer, Message string, Settings *Parameter) {
	if !Settings.Quiet {
		_, _ = io.WriteString(out, Message)
		_, _ = io.WriteString(out, Settings.NewLine)
	}
}
func bail(out *bufio.Writer, Message string, Settings *Parameter) {
	print(out, Message, Settings)
	if Settings.Bail {
		os.Exit(1)
	}
}
func fatal(out *bufio.Writer, Message string, Settings *Parameter) {
	_, _ = io.WriteString(out, Message)
	_, _ = io.WriteString(out, Settings.NewLine)
	_ = out.Flush()
	os.Exit(1)
}

func renderByteCount(count int64) string {
	const unit = 1000
	if count < unit {
		return fmt.Sprintf("%d Bytes", count)
	}
	div, exp := int64(unit), 0
	for n := count / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cBytes", float64(count)/float64(div), "kMGTPE"[exp])
}

func dropError(val string, err error) string { return val }
func getNextTokenValueAsBool(out *bufio.Writer, oldValue bool, tokens []string, Settings *Parameter) bool {
	switch len(tokens) {
	case 1:
		return !oldValue
	case 2:
		switch tokens[1] {
		case "on", "true", "1", "enable":
			return true
		case "off", "false", "0", "disable":
			return false
		default:
			bail(out, fmt.Sprintf("SYNTHAX ERROR: Invalid parameter \"%s\".", tokens[1]), Settings)
		}
	default:
		bail(out, "SYNTHAX ERROR: Wrong number of parameters.", Settings)
	}
	return oldValue
}
func getNextTokenValueAsInteger(out *bufio.Writer, oldValue int, defaultValue int, tokens []string, Settings *Parameter) int {
	switch len(tokens) {
	case 1:
		return defaultValue
	case 2:
		if i, err := strconv.Atoi(tokens[1]); err == nil {
			return i
		}
		bail(out, fmt.Sprintf("SYNTHAX ERROR: \"%s\" is not a valid number.", tokens[1]), Settings)
	default:
		bail(out, "SYNTHAX ERROR: Wrong number of parameters.", Settings)
	}
	return oldValue
}
func getNextTokenValueAsString(out *bufio.Writer, oldValue string, defaultValue string, allowedValues string, tokens []string, Settings *Parameter) string {
	switch len(tokens) {
	case 1:
		return defaultValue
	case 2:
		if allowedValues == "" {
			return replaceControlChars(tokens[1])
		}
		if strings.Contains(allowedValues, fmt.Sprintf("|%s|", tokens[1])) {
			return tokens[1]
		}
		bail(out, fmt.Sprintf("SYNTHAX ERROR: Invalid parameter \"%s\".", tokens[1]), Settings)
	default:
		bail(out, "SYNTHAX ERROR: Wrong number of parameters.", Settings)
	}
	return oldValue
}
