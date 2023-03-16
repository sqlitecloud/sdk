# SQLite Cloud Command Line Client

## Main Features
- Many command line arguments for good feature compatibility with other mayor database systems
- Connection Strings
- Batch and interactive mode
- Can read sql scripts from multible files and stdin simultaniously
- Many result output formats
- Many internal commands (.dot commands) for good feature compatibility with other mayor database systems
- Automatic line trunctions for nice terminal rendering
- Automatic rendering of numbers and size/time units
- Command history
- Static and dynamic SQL qutocomplete

## Compatibility with other database command line clients

<PRE>
✅ = Implemented, directly or indrectly available
🤔 = Not Implemented / Maybe later?
👎 = Decided agains implementation, not usefull
❌ = Does not apply / will never be implemented
</PRE>

### Command line compatibility with sqlite3

<PRE>
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
</PRE>

### Internal command compatibility with sqlite3

<PRE>
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
</PRE>

## Getting started

### Compile
```console
go env -w GO111MODULE=off
cd sdk/GO
export GOPATH=`pwd`
echo $GOPATH
make cli

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
  sqlc --host hostname -u user --password=pass -d dbname -c LZ4
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
  -h, --host HOSTNAME      Connect to SQLite Cloud database server host name [default::localhost]
  -p, --port PORT          Use specified port to connect to SQLIte Cloud database server [default::8860]
  -u, --user USERNAME      Use USERNAME for authentication
  -w, --password PASSWORD  Use PASSWORD for authentication
  -t, --timeout SECS       Set Timeout for network operations to SECS seconds [default::10]
  -c, --compress (NO|LZ4)  Use line compression [default::NO]
  --tls [YES|NO|INTERN|FILE] Encrypt the database connection using the host's root CA set (YES), a custom CA with a PEM from FILE (FILE), the internal SQLiteCloud CA (INTERN), or disable the encryption (NO) [default::YES]

```

### Internal Commands
```console
 hostname.com:X > .help

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

hostname:X > 

```

## Using the CLI

### Starting a new session
```console
./bin/sqlc --host=hostname --dbname=X --tls=INTERN
   _____     
  /    /     SQLite Cloud Command Line Application, version 1.0.1
 / ___/ /    (c) 2021 by SQLite Cloud Inc.
 \  ___/ /   
  \_ ___/    Enter ".help" for usage hints.

hostname:X >

```

### SELECT'ing some data
```console
hostname:X > SELECT * FROM Dummy;
┌─────┬───────────┬──────────┬───────┬──────────┬─────────────────────┐
│ ID  │ FirstName │ LastName │  ZIP  │   City   │       Address       │
├─────┼───────────┼──────────┼───────┼──────────┼─────────────────────┤
│ 369 │ Some      │ One      │ 96450 │ Coburg   │ Mohrenstraße 1      │
│ 370 │ Someone   │ Else     │ 96145 │ Sesslach │ Raiffeisenstraße 6  │
│ 371 │ One       │ More     │ 91099 │ Poxdorf  │ Langholzstr. 4      │
│ 372 │ Quotation │ Test     │ 12345 │ &"<>     │ 'Straße 0'          │
└─────┴───────────┴──────────┴───────┴──────────┴─────────────────────┘
Rows: 4 - Cols: 6: 282 Bytes Time: 86.43071ms

hostname:X > 

```
### DELETE'ing a row
```console
hostname:X > DELETE FROM Dummy WHERE ID = 372;
OK

hostname:X > 

```

### Changing the outformat
```console
hostname:X > .format xml
hostname:X > SELECT * FROM Dummy;
<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<resultset statement="SELECT * FROM Dummy;" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
  <row>
    <field name="ID">369</field>
    <field name="FirstName">Some</field>
    <field name="LastName">One</field>
    <field name="ZIP">96450</field>
    <field name="City">Coburg</field>
    <field name="Address">Mohrenstraße 1</field>
  </row>
  <row>
    <field name="ID">370</field>
    <field name="FirstName">Someone</field>
    <field name="LastName">Else</field>
    <field name="ZIP">96145</field>
    <field name="City">Sesslach</field>
    <field name="Address">Raiffeisenstraße 6</field>
  </row>
  <row>
    <field name="ID">371</field>
    <field name="FirstName">One</field>
    <field name="LastName">More</field>
    <field name="ZIP">91099</field>
    <field name="City">Poxdorf</field>
    <field name="Address">Langholzstr. 4</field>
  </row>
</resultset>
Rows: 3 - Cols: 6: 229 Bytes Time: 82.762014ms

hostname:X > 

```
### Changing the outformat back
One can enter `.format` without any argument to switch the output format back to its default format or one could enter an explicit format like `.format box`.

### Line Truncation explained:
Lets assume, that you have a narrow terminal window. If you have entered the following commmands, the output would look something like this:
```console
hostname:X > .format table
 hostname:X > SELECT * FROM Dummy;
+-----+-----------+----------+-------+----------+---------------
------+
| ID  | FirstName | LastName |  ZIP  |   City   |       Address 
      |
+-----+-----------+----------+-------+----------+---------------
------+
| 369 | Some      | One      | 96450 | Coburg   | Mohrenstraße 1
      |
| 370 | Someone   | Else     | 96145 | Sesslach | Raiffeisenstra
ße 6  |
| 371 | One       | More     | 91099 | Poxdorf  | Langholzstr. 4
      |
+-----+-----------+----------+-------+----------+---------------
------+
Rows: 3 - Cols: 6: 229 Bytes Time: 81.418646ms

hostname:X > .format json
hostname:X > SELECT * FROM Dummy;

[
  {"ID":369,"FirstName":"Some","LastName":"One","ZIP":96450,"Cit
y":"Coburg","Address":"Mohrenstraße 1",},
  {"ID":370,"FirstName":"Someone","LastName":"Else","ZIP":96145,
"City":"Sesslach","Address":"Raiffeisenstraße 6",},
  {"ID":371,"FirstName":"One","LastName":"More","ZIP":91099,"Cit
y":"Poxdorf","Address":"Langholzstr. 4",},
]
Rows: 3 - Cols: 6: 229 Bytes Time: 87.014696ms

hostname:X > 

```
You can see a nasty line break in the middle of the result line that can easily ruin the screen reading experience. To avoid this annoyance, sqlc build in line trucation mechanism trims its output line in a terminal session by default. The result looks like this:
```console
hostname:X > .format table
hostname:X > SELECT * FROM Dummy;
+-----+-----------+----------+-------+----------+--------------…
| ID  | FirstName | LastName |  ZIP  |   City   |       Address…
+-----+-----------+----------+-------+----------+--------------…
| 369 | Some      | One      | 96450 | Coburg   | Mohrenstraße …
| 370 | Someone   | Else     | 96145 | Sesslach | Raiffeisenstr…
| 371 | One       | More     | 91099 | Poxdorf  | Langholzstr. …
+-----+-----------+----------+-------+----------+--------------…
Rows: 3 - Cols: 6: 229 Bytes Time: 84.409225ms

hostname:X > .format json
hostname:X > SELECT * FROM Dummy;
[
  {"ID":369,"FirstName":"Some","LastName":"One","ZIP":96450,"Ci…
  {"ID":370,"FirstName":"Someone","LastName":"Else","ZIP":96145…
  {"ID":371,"FirstName":"One","LastName":"More","ZIP":91099,"Ci…
]
Rows: 3 - Cols: 6: 229 Bytes Time: 88.874433ms

hostname:X > 
```

If an output line was trimed to a certain width, the truncation can easily be spoted by the `…` character at the very end of a line. In batch mode, all output is sent to an output file, no line truncation will occure. You can switch off this autotrunction behaviour with a `.width 0` command. To switch back to auto truncation, use `.width -1`. Truncation to any other width is also possible with, for exampel a `.width 35` command. 

### Using Autocomplete
To use the build in autocomplete feature, use the [TAB] key. The [TAB] key will try to guess what SQL command you was trying to use and autocomplete this SQL command for you. If autocomplete has guessed the wrong command, keep pressing [TAB] until the right command shows up. The Autocomplete knows all available SQLite Cloud server and SQLite Cloud SQL commands and functions. If you have selected a database (`USE DATABASE ...`), autocomplete will also help you with the Table and Colum names. "UPDATING'ing some data" shows a simple example session:

### UPDATING'ing some data
```console
hostname:X > sel[TAB]
hostname:X > SELECT 
hostname:X > SELECT Fi[TAB]
hostname:X > SELECT FirstName
hostname:X > SELECT FirstName, Dum[TAB][TAB][TAB][TAB]
hostname:X > SELECT FirstName, Dummy.LastName
hostname:X > SELECT FirstName, Dummy.LastName Fr[TAB]
hostname:X > SELECT FirstName, Dummy.LastName FROM D[TAB]
hostname:X > SELECT FirstName, Dummy.LastName FROM Dummy[RETURN]
┌───────────┬──────────┐
│ FirstName │ LastName │
├───────────┼──────────┤
│ Some      │ One      │
│ Someone   │ Else     │
│ One       │ More     │
└───────────┴──────────┘
Rows: 3 - Cols: 2: 74 Bytes Time: 81.865386ms

hostname:X > up[TAB]
hostname:X > UPDATE D[TAB]
hostname:X > UPDATE Dummy SET La[TAB]
hostname:X > UPDATE Dummy SET LastName 
hostname:X > UPDATE Dummy SET LastName = "ONE" WH[TAB]
hostname:X > UPDATE Dummy SET LastName = "ONE" WHERE id=369[RETURN]
OK

hostname:X > SELECT * FROM Dummy;
┌─────┬───────────┬──────────┬───────┬──────────┬─────────────────────┐
│ ID  │ FirstName │ LastName │  ZIP  │   City   │       Address       │
├─────┼───────────┼──────────┼───────┼──────────┼─────────────────────┤
│ 369 │ Some      │ ONE      │ 96450 │ Coburg   │ Mohrenstraße 1      │
│ 370 │ Someone   │ Else     │ 96145 │ Sesslach │ Raiffeisenstraße 6  │
│ 371 │ One       │ More     │ 91099 │ Poxdorf  │ Langholzstr. 4      │
└─────┴───────────┴──────────┴───────┴──────────┴─────────────────────┘
Rows: 3 - Cols: 6: 229 Bytes Time: 82.797135ms

hostname:X > 
```

### Setting the prompt
One can set the user promt ether with the `--promt <format>` command line switch or within the app with a `.promt <format>`. The follwoing format strings will automatically be replaced by it's actual values:

(https://wiki.ubuntuusers.de/Bash/Prompt/)

| Format String | Meaning | Example replacement |
|:---|:---|:---|
| \H | Database Host name | hostname |
| \p | Database Port      | 8860 |
| \u | Username      | marco |
| \d | Database name | X |
| \ | Actual Date | 2021-10-08 |
| \t | Actual Time | 14:37:32 |
| \w | Local path | ~ |

The following prebuild patterns are available as a shortcut:

| Shortcut | Definition | Example |
|:---|:---|:---|
| default | ``$host$:$dbname$ >`` | hostname:X >_ |
| simple | ``sqlc >`` | sqlc >_ |
| full | ``$host$:$dbname$ $user$> `` | hostname:X username> _ |

### Using Synthax Highlighting
One can switch synthax highlighting on with ``.synthax color``. There is also a monocrome synthax highlighting mode available that is VT100 compatible and uses only bold and italic character representations. You can switch into this simple VT100 mode with ``.synthax vt100``. To switch synthay highlighting off, please enter ``.synthax off``. Synthax highlighting is set to color in interactive mode by default. Synthay highlighting is automatically switched off in batch mode.


### Exiting the app
```console
hostname:X > .exit

```

## Testing
```console
./bin/sqlc -?
./bin/sqlc --help
./bin/sqlc --version

./bin/sqlc -> trying to connect
./bin/sqlc sqlitecloud://hostname/X?tls=INTERN -> trying to connect

./bin/sqlc sqlitecloud://hostname/X?tls=INTERN

./bin/sqlc sqlitecloud://hostname/X?tls=INTERN --list
./bin/sqlc sqlitecloud://hostname/X?tls=INTERN --list --format=xml
./bin/sqlc sqlitecloud://hostname/X?tls=INTERN --cmd "LIST DATABASES"
echo "LIST DATABASES" | ./bin/sqlc sqlitecloud://hostname/X?tls=INTERN --format=json
echo "LIST DATABASES" > script.sql; ./bin/sqlc sqlitecloud://hostname/X?tls=INTERN --format=json script.sql

./bin/sqlc sqlitecloud://hostname/X?tls=INTERN --list -o outputfile --quiet --format=xml

```

## ToDos
- [ ] --promt
- [ ] --reconnect
- [x] When started without a host, localhost is assumed. If no database server is running on the localhost, a 10 second timeout will occure.
- [ ] Make internal .dot commands available in batch files
- [ ] Test for "empty commands"
- [ ] Add --log feature
- [x] Remove the table "sqlite_sequence" from dynamic autocomplete scanning
- [ ] Implement the Auth command to use the Password feature
- [ ] Add more Test example commands