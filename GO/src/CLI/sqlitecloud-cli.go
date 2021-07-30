package main

import "fmt"
import "os"
import "strings"
import "sqlitecloud"

import "github.com/docopt/docopt-go"      // https://github.com/docopt/docopt-go
import "github.com/deadsy/go-cli"         // https://github.com/deadsy/go-cli

var history_file = ".sqlitecloud_history.txt"
var version      = "version 0.0.2, (c) 2021 by SQLite Cloud Inc."
var usage        = `SQLite Cloud Command Line Interface.

Usage:
  sqlitecloud-cli [--host=<hostname>] [--port=<port>] --user=<username> --password=<password> [--database=<database>] [OPTIONS] [-|<FILE>...]
  sqlitecloud-cli URL [OPTIONS] [-|<FILE>...]
  sqlitecloud-cli -?|--help|--version

Arguments:
  URL               Connect to SQLite Cloud database server with an URL as connection string
  FILE...           After connection to the SQLite Cloud database server, execute SQL Command(s) from Script-FILE(s)
  -                 After connection to the SQLite Cloud database server, execute SQL Command(s) from stdin

Examples:
  sqlitecloud-cli --host ***REMOVED*** -c
  sqlitecloud-cli "sqlitecloud://user:pass@host.com:8860/dbname?timeout=10&uuid=12342&compress=lz4&sslmode=disabled"
  sqlitecloud-cli --help

General Options:
  --cmd COMMAND              Run "COMMAND" before reading stdin
  -l, --list                 List available databases, then exit
  -d, --dbname DBNAME        Open database DBNAME
  -b, --bail                 Stop after hitting an error
  -?, --help                 Show this screen
  --version                  Display version information   
  
Output Format Options:  
  -o, --output FILE          Send query results to FILE
  --echo                     Print commands before execution
  --quiet                    Run quietly (no messages, only query output)
  --noheader                 Turn headers off
  --nullvalue TEXT           Set text string for NULL values [default '']
  --newline SEP              Set output row separator [default: '\n']
  --separator SEP            Set output column separator[ default: '|']
  --mode [LIST|CSV|QUOTE|TABS|LINE|JSON|HTML|MARKDOWN|TABLE|BOX]        
                             Specify the Output mode [default 'LIST']

Connection Options:
  -h, --host HOSTNAME        Connect to SQLite Cloud database server host name [default: localhost]
  -p, --port PORT            Use specified port to connect to SQLIte Cloud database server [default: 8860]
  -n, --user USERNAME        Use USERNAME for authentication 
  -w, --password PASSWORD    Use PASSWORD for authentication   
  -u, --uuid          
  
  -t, --timeout SECS         Set Timeout for network operations to SECS seconds [default: 10]
  -e, --sslmode              Use SSL Encryption to connect
  -c, --compress (NONE|LZ4)  Use line compression [default: NONE]
` 

func main() {
  arguments, err := docopt.ParseArgs( usage, nil, "sqlitecloud-cli " + version )
  if err == nil {

    var config struct {
      Host     string   `docopt:"--host"`
      Port     int      `docopt:"--port"`
      Compress bool     `docopt:"-c"`
      UseStdIn bool     `docopt:"-"`
      Files    []string `docopt:"<FILE>"`
    }
    err = arguments.Bind( &config )
    if err == nil {
      var db *sqlitecloud.SQCloud
      db, err = sqlitecloud.Connect( config.Host )
      if err == nil {
        defer db.Close()

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