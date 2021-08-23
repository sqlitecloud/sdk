# SQLite Cloud Command Line Client

## Main Features
- Many command line arguments for good feature compatibility with other mayor database systems
- Connection Strings
- Batch and interactive mode
- Can read sql scripts from multible files and stdin simultaniously
- Many result output formats
- Many internal commands (.dot commands) for good feature compatibility with other mayor database systems
- Automatic line trunctions for nice terminal rendering
- Automatic rednering of numbers and size/time units
- Command history
- Static and dynamic query qutocomplete

## Compatibility with other database command line clients

✅ = Implemented, directly or indrectly available
🤔 = Not Implemented / Maybe later?
👎 = Decided agains implementation, not usefull
❌ = Does not apply / will never be implemented

### Command line compatibility with sqlite3

🤔   -A ARGS...             run ".archive ARGS" and exit
❌   -append                append the database to the end of the file
👎   -ascii                 set output mode to 'ascii'
✅   -bail                  stop after hitting an error
🤔   -batch                 force batch I/O
✅   -box                   set output mode to 'box'
✅   -column                set output mode to 'column'
✅   -cmd COMMAND           run "COMMAND" before reading stdin
✅   -csv                   set output mode to 'csv'
❌   -deserialize           open the database using sqlite3_deserialize()
✅   -echo                  print commands before execution
✅   -init FILENAME         read/process named file
✅   -[no]header            turn headers on or off
✅   -help                  show this message
✅   -html                  set output mode to HTML
👎   -interactive           force interactive I/O
✅   -json                  set output mode to 'json'
✅   -line                  set output mode to 'line'
✅   -list                  set output mode to 'list'
❌   -lookaside SIZE N      use N entries of SZ bytes for lookaside memory
✅   -markdown              set output mode to 'markdown'
❌   -maxsize N             maximum size for a --deserialize database
❌   -memtrace              trace all memory allocations and deallocations
❌   -mmap N                default mmap size set to N
✅   -newline SEP           set output row separator. Default: '\n'
❌   -nofollow              refuse to open symbolic links to database files
✅   -nullvalue TEXT        set text string for NULL values. Default ''
❌   -pagecache SIZE N      use N slots of SZ bytes each for page cache memory
✅   -quote                 set output mode to 'quote'
❌   -readonly              open the database read-only
✅   -separator SEP         set output column separator. Default: '|'
👎   -stats                 print memory stats before each finalize
✅   -table                 set output mode to 'table'
✅   -tabs                  set output mode to 'tabs'
✅   -version               show SQLite version
❌   -vfs NAME              use NAME as the default VFS
❌   -zip                   open the file as a ZIP Archive

### Internal command compatibility with sqlite3

🤔 .archive ...             Manage SQL archives
🤔 .auth ON|OFF             Show authorizer callbacks
🤔 .backup ?DB? FILE        Backup DB (default "main") to FILE
✅ .bail on|off             Stop after hitting an error.  Default OFF
👎 .binary on|off           Turn binary output on or off.  Default OFF
🤔 .cd DIRECTORY            Change the working directory to DIRECTORY
🤔 .changes on|off          Show number of rows changed by SQL
🤔 .check GLOB              Fail if output since .testcase does not match
🤔 .clone NEWDB             Clone data into NEWDB from the existing database
👎 .databases               List names and files of attached databases
👎 .dbconfig ?op? ?val?     List or change sqlite3_db_config() options
🤔 .dbinfo ?DB?             Show status information about the database
🤔 .dump ?OBJECTS?          Render database content as SQL
✅ .echo on|off             Turn command echo on or off
👎 .eqp on|off|full|...     Enable or disable automatic EXPLAIN QUERY PLAN
👎 .excel                   Display the output of next command in spreadsheet
✅ .exit ?CODE?             Exit this program with return-code CODE
🤔 .expert                  EXPERIMENTAL. Suggest indexes for queries
🤔 .explain ?on|off|auto?   Change the EXPLAIN formatting mode.  Default: auto
🤔 .filectrl CMD ...        Run various sqlite3_file_control() operations
🤔 .fullschema ?--indent?   Show schema and the content of sqlite_stat tables
✅ .headers on|off          Turn display of headers on or off
✅ .help ?-all? ?PATTERN?   Show help text for PATTERN
🤔 .import FILE TABLE       Import data from FILE into TABLE
🤔 .imposter INDEX TABLE    Create imposter table TABLE on index INDEX
🤔 .indexes ?TABLE?         Show names of indexes
🤔 .limit ?LIMIT? ?VAL?     Display or change the value of an SQLITE_LIMIT
🤔 .lint OPTIONS            Report potential schema issues.
❌ .load FILE ?ENTRY?       Load an extension library
🤔 .log FILE|off            Turn logging on or off.  FILE can be stderr/stdout
✅ .mode MODE ?TABLE?       Set output mode
✅ .nullvalue STRING        Use STRING in place of NULL values
👎 .once ?OPTIONS? ?FILE?   Output for the next SQL command only to FILE
🤔 .open ?OPTIONS? ?FILE?   Close existing database and reopen FILE
🤔 .output ?FILE?           Send output to FILE or stdout if FILE is omitted
🤔 .parameter CMD ...       Manage SQL parameter bindings
🤔 .print STRING...         Print literal STRING
🤔 .progress N              Invoke progress handler after every N opcodes
🤔 .prompt MAIN CONTINUE    Replace the standard prompts
✅ .quit                    Exit this program
🤔 .read FILE               Read input from FILE
❌ .recover                 Recover as much data as possible from corrupt db.
🤔 .restore ?DB? FILE       Restore content of DB (default "main") from FILE
❌ .save FILE               Write in-memory database into FILE
🤔 .scanstats on|off        Turn sqlite3_stmt_scanstatus() metrics on or off
🤔 .schema ?PATTERN?        Show the CREATE statements matching PATTERN
❌ .selftest ?OPTIONS?      Run tests defined in the SELFTEST table
✅ .separator COL ?ROW?     Change the column and row separators
👎 .session ?NAME? CMD ...  Create or control sessions
❌ .sha3sum ...             Compute a SHA3 hash of database content
🤔 .shell CMD ARGS...       Run CMD ARGS... in a system shell
🤔 .show                    Show the current values for various settings
🤔 .stats ?ARG?             Show stats or turn stats on or off
🤔 .system CMD ARGS...      Run CMD ARGS... in a system shell
👎 .tables ?TABLE?          List names of tables matching LIKE pattern TABLE
❌ .testcase NAME           Begin redirecting output to 'testcase-out.txt'
❌ .testctrl CMD ...        Run various sqlite3_test_control() operations
✅ .timeout MS              Try opening locked tables for MS milliseconds
🤔 .timer on|off            Turn SQL timer on or off
🤔 .trace ?OPTIONS?         Output each SQL statement as it is run
❌ .vfsinfo ?AUX?           Information about the top-level VFS
❌ .vfslist                 List all available VFSes
❌ .vfsname ?AUX?           Print the name of the VFS stack
🤔 .width NUM1 NUM2 ...     Set minimum column widths for columnar output

## Getting started

### Compile
```console
go env -w GO111MODULE=off
cd sdk/GO
export GOPATH=`pwd`
echo $GOPATH
male cli

```

### Usage
```console
./bin/sqlc --help
SQLite Cloud Command Line Application Command Line Interface.

Usage:
  sqlc [URL] [options] [<FILE>...]
  sqlc -?|--help|--version

Arguments:
  URL                      "sqlitecloud://user:pass@host.com:port/dbname?timeout=10&compress=NO"
  FILE...                  Execute SQL commands from FILE(s) after connecting to the SQLite Cloud database

Examples:
  sqlc "sqlitecloud://user:pass@host.com:8860/dbname?timeout=10&compress=lz4"
  sqlc --host ***REMOVED*** -u user --password=pass -d dbname -c LZ4
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
  --separator SEP          Set output column separator [default: "|"]
  --format (LIST|CSV|QUOTE|TABS|LINE|JSON|HTML|XML|MARKDOWN|TABLE|BOX)
                           Specify the Output mode [default: BOX]

Connection Options:
  -h, --host HOSTNAME      Connect to SQLite Cloud database server host name [default: localhost]
  -p, --port PORT          Use specified port to connect to SQLIte Cloud database server [default: 8860]
  -u, --user USERNAME      Use USERNAME for authentication 
  -w, --password PASSWORD  Use PASSWORD for authentication
  -t, --timeout SECS       Set Timeout for network operations to SECS seconds [default: 10]
  -c, --compress (NO|LZ4)  Use line compression [default: NO]

```

### Internal Commands
```console
./bin/sqlc --host=***REMOVED*** --dbname=X
   _____     
  /    /     SQLite Cloud Command Line Application, version 1.0.1
 / ___/ /    (c) 2021 by SQLite Cloud Inc.
 \  ___/ /   
  \_ ___/    Enter ".help" for usage hints.

***REMOVED***:X > .help

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

***REMOVED***:X > 

```

### Using the CLI


### Exiting the app
```console
***REMOVED***:X > .exit

```


### Testing
```console
./bin/sqlc -?
./bin/sqlc --help
./bin/sqlc --version

./bin/sqlc -> trying to connect
./bin/sqlc sqlitecloud://sqlitecloud://***REMOVED***/X -> trying to connect

./bin/sqlc sqlitecloud://***REMOVED***/X

./bin/sqlc sqlitecloud://***REMOVED***/X --list
./bin/sqlc sqlitecloud://***REMOVED***/X --list --format=xml
./bin/sqlc sqlitecloud://***REMOVED***/X --cmd "LIST DATABASES"
echo "LIST DATABASES" | ./bin/sqlc sqlitecloud://***REMOVED***/X --format=json
echo "LIST DATABASES" > script.sql; ./bin/sqlc sqlitecloud://***REMOVED***/X --format=json script.sql

./bin/sqlc sqlitecloud://***REMOVED***/X --list -o outputfile --quiet --format=xml

```