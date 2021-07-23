package main

import "fmt"
import "os"
import "strings"
import "sqlitecloud"

import "github.com/docopt/docopt-go" 			// https://github.com/docopt/docopt-go
import "github.com/deadsy/go-cli" 				// https://github.com/deadsy/go-cli

var history_file = ".sqlitecloud_history.txt"
var version      = "version 0.0.1, (c) 2021 by SQLite Cloud Inc."
var usage        = `SQLite Cloud Command Line Interface.

Usage:
  sqlitecloud-cli --host=<hostname> [--port=<port>] [-c] [-|<FILE>...]
  sqlitecloud-cli -h|--help|--version

Arguments:
  FILE              Optional SQL Command script(s) to execute.
                    If no FILE is given, sqliteclout-cli will try to read commands from stdin.

Examples:
  sqlitecloud-cli --help
  sqlitecloud-cli --host ***REMOVED*** -c

Options:
  --host HOSTNAME   Connect to SQLite Cloud database server host name [default: localhost]
  --port PORT       Use specified port to connect to SQLIte Cloud database server [default: 8860]
  -c                Use compression [defsult: false]
  
  -h, --help        Show this screen.
  --version         Version Information.`

func main() {
  arguments, err := docopt.ParseArgs( usage, nil, "sqlitecloud-cli " + version )
  if err == nil {

		var config struct {
			Host     string   `docopt:"--host"`
			Port     int      `docopt:"--port"`
			Compress bool     `docopt:"-c"`
			UseStdIn bool		  `docopt:"-"`
			Files    []string `docopt:"<FILE>"`
		}
		err = arguments.Bind( &config )
		if err == nil {
			var db *sqlitecloud.SQCloud
			db, err = sqlitecloud.Connect( config.Host, config.Port )
			if err == nil {
				defer db.Disconnect()

				fmt.Printf( "Connection to %s OK...\r\n", config.Host);

				db.Compress( config.Compress )
				db.ExecuteFiles( config.Files )

				editor := cli.NewLineNoise()
				editor.SetMultiline( true )
				editor.HistoryLoad( history_file )

				var err error
				for command := ""; 
					!strings.HasPrefix( strings.ToUpper( command ), "EXIT") && 
					!strings.HasPrefix( strings.ToUpper( command ), "QUIT") && 
					err != cli.ErrQuit; 
					command, err = editor.Read( ">> ", "" ) {

					if err == nil {
						db.Execute( command )
					} else {
						fmt.Printf( "ERROR: %s\r\n", err )
					}	
				}
			
				editor.HistorySave( history_file )
				return
			}
			bailout( fmt.Sprintf( "Could not connect to database host (%s).", err.Error() ) )
		}
  }
	bailout( "Could not process the command line arguments." )
}

func bailout( Message string ) {
	fmt.Printf( "ERROR: %s\r\n", Message )
	os.Exit( 1 )
}