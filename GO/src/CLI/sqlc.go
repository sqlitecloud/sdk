package main

import "fmt"
import "os"
import "io"
import "bufio"
import "errors"
import "strings"
import "reflect"
import "sqlitecloud"
import "encoding/json"

import "github.com/docopt/docopt-go"      // https://github.com/docopt/docopt-go
import "github.com/deadsy/go-cli"         // https://github.com/deadsy/go-cli

var app_name     = "sqlc"
var history_file = "." + app_name +".history"
var version      = "version 0.1.0, (c) 2021 by SQLite Cloud Inc."
var usage        = `SQLite Cloud Command Line Interface.

Usage:
  sqlc [options] [URL] [-|<FILE>...]
  sqlc -?|--help|--version

Arguments:
  URL                      Connection String in the form of: 
                           "sqlitecloud://user:pass@host.com:port/dbname?timeout=10&compress=NO"
  FILE...                  Execute SQL Commands from FILE(s) after connecting to the SQLite Cloud database
  -                        Execute SQL Commands from STDIN after connecting to the SQLite Cloud database

Examples:
  sqlc "sqlitecloud://user:pass@host.com:8860/dbname?timeout=10&compress=lz4"
  sqlc --host ***REMOVED*** -u user --password=pass -d dbname -c LZ4
  sqlc --version
  sqlc --help

General Options:
  --cmd COMMAND            Run "COMMAND" before reading from stdin
  -l, --list               List available databases, then exit
  -d, --dbname NAME        Use database NAME
  -b, --bail               Stop after hitting an error
  -?, --help               Show this screen
  --version                Display version information   
  
Output Format Options:  
  -o, --output FILE        Execute SQL Commands quietly from (--cmd,FILE,stdin) in batch mode 
                           and send results to FILE, then exit.
  --echo                   Print commands before execution
  --quiet                  Run quietly (no messages, only query output)
  --noheader               Turn headers off
  --nullvalue TEXT         Set text string for NULL values [default: "NULL"]
  --newline SEP            Set output row separator [default: "\r\n"]
  --separator SEP          Set output column separator [default: "|"]
  --mode (LIST|CSV|QUOTE|TABS|LINE|JSON|HTML|XML|MARKDOWN|TABLE|BOX)
                           Specify the Output mode [default: LIST]

Connection Options:
  -h, --host HOSTNAME      Connect to SQLite Cloud database server host name [default: localhost]
  -p, --port PORT          Use specified port to connect to SQLIte Cloud database server [default: 8860]
  -u, --user USERNAME      Use USERNAME for authentication 
  -w, --password PASSWORD  Use PASSWORD for authentication
  -t, --timeout SECS       Set Timeout for network operations to SECS seconds [default: 10]
  -c, --compress (NO|LZ4)  Use line compression [default: NO]
`

type Parameter struct {
  URL         string      `docopt:"URL"`

  OutFile     string      `docopt:"--output"`
  Command     string      `docopt:"--cmd"`
  List        bool        `docopt:"--list"`
  Bail        bool        `docopt:"--bail"`
  Echo        bool        `docopt:"--echo"`
  Quiet       bool        `docopt:"--quiet"`
  NoHeader    bool        `docopt:"--noheader"`
  NullText    string      `docopt:"--nullvalue"`
  NewLine     string      `docopt:"--newline"`
  Separator   string      `docopt:"--separator"`
  Mode        string      `docopt:"--mode"`

  Host        string      `docopt:"--host"`
  Port        int         `docopt:"--port"`
  User        string      `docopt:"--user"`
  Password    string      `docopt:"--password"`
  Database    string      `docopt:"--dbname"`

  Timeout     int         `docopt:"--timeout"`
  Compress    string      `docopt:"--compress"`
  UseStdIn    bool        `docopt:"-"`
  Files       []string    `docopt:"<FILE>"`
}
func (this *Parameter) ToJSON() string {
  jParameter, _ := json.Marshal( this )
  return string( jParameter )
}

func replaceControlChars( in string ) string {
  out := strings.ReplaceAll( in,  "\\0", string( 0 ) )
  out  = strings.ReplaceAll( out, "\\a", "\a" )
  out  = strings.ReplaceAll( out, "\\b", "\b" )
  out  = strings.ReplaceAll( out, "\\t", "\t" )
  out  = strings.ReplaceAll( out, "\\n", "\n" )
  out  = strings.ReplaceAll( out, "\\v", "\v" )
  out  = strings.ReplaceAll( out, "\\f", "\f" )
  out  = strings.ReplaceAll( out, "\\r", "\r" )
  return out
}
func parseParameters() ( Parameter, error ) {
  parameter:= Parameter{}
  // Parse Command Line Parameter
  if p, err := docopt.ParseArgs( usage, nil, fmt.Sprintf( "%s %s", app_name, version ) ); err == nil {
    
    // If set, parse and apply the Connection string....
    if url, isSet := p[ "URL" ]; isSet && url != "<nil>" {
      if Host, Port, Username, Password, Database, Timeout, Compress, err := sqlitecloud.ParseConnectionString( reflect.ValueOf( url ).String() ); err == nil {
        if Host != "" {
          p[ "--host" ] = Host
        }
        if Port != -1 {
          p[ "--port" ] = Port
        }
        if Username != "" {
          p[ "--user" ] = Username
        }
        if Password != "" {
          p[ "--password" ] = Password 
        }
        if Database != "" {
          p[ "--dbname" ] = Database 
        }
        if Timeout != -1 {
          p[ "--timeout" ] = Timeout 
        }
        if Compress != "" {
          p[ "--compress" ] = Compress
        }
      }
    } else {
      return Parameter{}, err
    }

    // Fix invalid(=unset) parameters, Quotation & Control-Chars
    for k, v := range p {
      switch reflect.ValueOf( v ).Kind() {
      case reflect.Invalid: p[ k ] = ""
      case reflect.String:  p[ k ] = replaceControlChars( strings.Trim( reflect.ValueOf( v ).String(), "'\"" ) )
      default:
      }
    }

    for k, v := range p {
      fmt.Printf( "%s='%v'\r\n", k, v )
    }

    // Copy map data into Object
    if err := p.Bind( &parameter ); err != nil {
      return Parameter{}, err
    }

  } else {
    return Parameter{}, err
  }

  // Postprocessing....
  if parameter.OutFile != "" { // batch mode
    parameter.Quiet = true
  } else {                     // interactive mode

  }


  if parameter.Quiet {
    parameter.Echo = false
  }

  return parameter, nil
}

func main() {
  out := bufio.NewWriter( os.Stdout )

  if parameter, err := parseParameters(); err != nil {
    bail( out, fmt.Sprintf( "Could not parse arguments (%s)", err.Error() ), &parameter )
    os.Exit( 1 )
  } else {
    if parameter.OutFile != "" {
      println( "outfile=" + parameter.OutFile )
      if file, err := os.OpenFile( parameter.OutFile, os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644 ); err != nil {
        bail( out, fmt.Sprintf( "Could not open '%s' for writing", parameter.OutFile ), &parameter )
      } else {
        out = bufio.NewWriter( file )
        defer file.Close()
      }
    }
    print( out, "trying to connect...\r\n", &parameter )

    var db *sqlitecloud.SQCloud = sqlitecloud.New()
    if err := db.Connect( parameter.Host, parameter.Port, parameter.User, parameter.Password, parameter.Database, parameter.Timeout, parameter.Compress, 0 ); err != nil {
      bail( out, fmt.Sprintf( "Could not connect (%s)", err.Error() ), &parameter )
      os.Exit( 1 )
    } else {
      defer db.Close()

      print( out, fmt.Sprintf( "Connection to %s OK...\r\n", parameter.Host ), &parameter )
      
      if parameter.List {
        Execute( db, out, "LIST DATABASES", &parameter )
        os.Exit( 0 )
      } 

      // Batch Mode starts here

      // Execute single Command   
      if parameter.Command != "" {
        Execute( db, out, parameter.Command, &parameter )
      }
      // Execute Files  
      if len( parameter.Files ) > 0 {
        ExecuteFiles( db, out, parameter.Files, &parameter )
      }
      // Execute Stdin  
      if parameter.UseStdIn {
        if err := ExecuteBuffer( db, out, os.Stdin, &parameter ); err != nil {
          bail( out, fmt.Sprintf( "Could not execute (%s)", err.Error() ), &parameter )
        }
      }
      // End Batch Mode
      if parameter.OutFile != "" {
        os.Exit( 0 )
      }

      // Interactive Mode starts here

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
          Execute( db, out, command, &parameter )
        } else {
          bail( out, err.Error(), &parameter )
        }
      }
      editor.HistorySave( history_file )
    }
  }
}



func Execute( db *sqlitecloud.SQCloud, out *bufio.Writer, cmd string, Settings *Parameter ) {
  if Settings.Echo {
    print( out, cmd, Settings )
  }

  println( "Execute:" + cmd ) // Use Dump...

  if !Settings.Quiet {
    print( out, "OK", Settings ) 
  }
}

func ExecuteBuffer( db *sqlitecloud.SQCloud, out *bufio.Writer, in *os.File, Settings *Parameter ) error {
  if scanner := bufio.NewScanner( in ); scanner != nil {
    for scanner.Scan() {
      line := scanner.Text()
      if strings.ToUpper( line ) == ".PROMPT" {
        return nil // break out of sql script
      }
      fmt.Println( ">> %s\r\n", line )
      Execute( db, out, line, Settings )
    }
    return scanner.Err() // nil or some error
  }
  return errors.New( "Could not instanciate the line scanner" )
}

func ExecuteFile( db *sqlitecloud.SQCloud, out *bufio.Writer, FilePath string, Settings *Parameter ) error {
  file, err := os.Open( FilePath )
  if err == nil {
    defer file.Close()
    if err := ExecuteBuffer( db, out, file, Settings ); err != nil {
      return err
    }
  }
  return err
}

func ExecuteFiles( db *sqlitecloud.SQCloud, out *bufio.Writer, FilePathes []string, Settings *Parameter ) error {
  for _, file := range FilePathes {
    err := ExecuteFile( db, out, file, Settings )
    if( err != nil ) {
      return err
    }
  }
  return nil
}

func print( out *bufio.Writer, Message string, Settings *Parameter ) {
  println( "print:" + Message )
  if !Settings.Quiet {
    io.WriteString( out, Message )
  }
  out.Flush()
}

func bail( out *bufio.Writer, Message string, Settings *Parameter ) {
  print( out, fmt.Sprintf( "ERROR: %s\r\n", Message ), Settings )
  if Settings.Bail {
    os.Exit( 1 )
  }
}