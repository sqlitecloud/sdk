//
//                    ////              SQLite Cloud
//        ////////////  ///
//      ///             ///  ///        Product     : SQLite Cloud GO SDK
//     ///             ///  ///         Version     : 1.0.1
//     //             ///   ///  ///    Date        : 2021/10/08
//    ///             ///   ///  ///    Author      : Andrea Donetti
//   ///             ///   ///  ///
//   ///     //////////   ///  ///      Description : SQLite Cloud server test
//   ////                ///  ///                     run test from source file
//     ////     //////////   ///
//        ////            ////
//          ////     /////
//             ///                      Copyright   : 2021 by SQLite Cloud Inc.
//
// -----------------------------------------------------------------------TAB=2

package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sqlitecloud"
	"strconv"
	"strings"
	"sync"
	"testing"
	"text/template"
	"time"
	"unicode"

	"github.com/PaesslerAG/gval"
)

// ----------------------------- Debugger -----------------------------

var logger = log.New(os.Stderr, "tester_test.go: ", log.LstdFlags|log.Lmicroseconds)

func header(lvl, msg string) string {
	_, file, line, _ := runtime.Caller(2)
	logger.SetPrefix(fmt.Sprintf("%s:%05d: ", filepath.Base(file), line))
	return fmt.Sprintf("%s: %s", lvl, msg)
}

func Debug(v ...interface{}) {
	if debug {
		logger.Output(2, header("DEBUG", fmt.Sprint(v...)))
	}
}

func Debugf(format string, v ...interface{}) {
	if debug {
		logger.Output(2, header("DEBUG", fmt.Sprintf(format, v...)))
	}
}

// func Error(v ...interface{}) {
// 	logger.Output(2, header("ERROR", fmt.Sprint(v...)))
// }

// func Errorf(format string, v ...interface{}) {
// 	logger.Output(2, header("ERROR", fmt.Sprintf(format, v...)))
// }

// ----------------------------- Workers -----------------------------

type task struct {
	name       string
	line       int
	scanner    *bufio.Scanner
	connstring string
	file       *os.File
	env        map[string]interface{} // local copy environment variables
}

type statistics struct {
	failed        bool
	nchecks       int
	nfailedchecks int
	nretries      int
}

type worker struct {
	// scannerc  chan *bufio.Scanner
	// filec     chan *os.File
	id        int
	conn      *sqlitecloud.SQCloud
	res       *sqlitecloud.Result // last recived sqlitecloud result
	reserr    error               // err returned by the execution of the last sqlite command
	taskc     chan *task          // channel used to feed new tasks to the worker
	completec chan struct{}       // channel used to wake up a worker waiting for this one to complete
	exiting   bool                // flag set if the worker is exiting
	t         *testing.T          //
	stats     statistics          // thread safe, local for the worker
}

var (
	workers  map[int]*worker        // = make(map[int]*worker)
	wg       sync.WaitGroup         // WaitGroup used for --wait all
	stopmu   sync.Mutex             // mutex used to protect the access to stopc and tstopped var
	stopc    chan struct{}          // channel used to stop the execution of all the workers by calling close(stopc), for example in case of error
	stopped  bool                   // flag to avoid calling multiple times close(stopc)
	envMutex sync.RWMutex           // global envMutex, it's a pointer because mutexes cannot be copied
	env      map[string]interface{} // global environment variables
)

func tryStopc(force bool) {
	stopmu.Lock()
	if !stopped && (skip || force) {
		stopped = true
		close(stopc)
		Debugf("Execution Stopped")
	}
	stopmu.Unlock()
}

func newWorker(id int, connstring string, t *testing.T) (*worker, error) {
	Debugf("w%d created with connstring: %s", id, connstring)
	w := worker{
		id:        id,
		conn:      nil,
		res:       &sqlitecloud.Result{},
		reserr:    nil,
		taskc:     make(chan *task, 100),
		completec: make(chan struct{}, 1),
		exiting:   false,
		t:         t,
		stats:     statistics{},
	}

	// add the first task to connect to the server
	Debugf("w%d connection task %s", w.id, connstring)
	w.taskc <- &task{name: "-", connstring: connstring}
	return &w, nil
}

// ----------------------------- COMMANDS -----------------------------

type commandFunc func(w *worker, args []string, opts []bool, t *task) error // scanner *bufio.Scanner

type command struct {
	template     string
	desctription string
	f            commandFunc
	tokens       []Token
}

var commands = make([]command, 0, 50)
var commandsRequiringEnd = [][]byte{
	[]byte("--task "),
	[]byte("--loop "),
}
var endcommand = "--end"
var endcommandb = []byte(endcommand)

func initCommands() {
	cs := [...]command{
		{"--list commands", "List the tester builtin commands.", commandListCommands, nil},
		{"--sleep %ms", "Pause for N milliseconds.", commandSleep, nil},
		{"--wait %workerid_or_all [%timeout_ms]", "Wait until all tasks complete for the given client.  If WORKERID_OR_ALL is \"all\" then wait for all clients to complete, otherwise only for the specified WORKERID.  Wait no longer than TIMEOUT_MS milliseconds (default 10,000).", commandWait, nil},
		{"--dump", "Dump the content of the last received response from the server.", commandDump, nil}, // dump the last result
		{"--exit [%rc]", "Exit this worker.  If N>0 then exit without explicitly disconnecting (In other words, simulate a crash.).", commandExit, nil},
		{"--task %workerid_or_expr [name %name] [%connectionstring]", "Assign work to a worker. Start the worker if it is not running already using the CONNECTIONSTRING (or the default connectionstring argument of the tester, if not specified)", commandTask, nil},
		{"--match type is %type", "Check to see if last response type matches TYPE.  Report an error if not.", commandMatchType, nil},
		{"--match buffer %string [%row %col]", "Check to see if last response buffer matches STRING. If the last response is an Array or a Rowset, the ROW and COL options can be use to check only the specified subvalue. Report an error if not.", commandMatchBuffer, nil},
		{"--match value %value [%row %col]", "Check to see if last response value matches VALUE. If the last response is an Array or a Rowset, the ROW and COL options can be use to check only the specified subvalue. Report an error if not.", commandMatchValue, nil},
		{"--match [rows %nrows] [cols %ncols]", "Check to see if last response has the specified number of ROW and/or COL (only works for Array and Rowset response). Report an error if not.", commandMatchRowsCols, nil},
		{"--loop %var=%val; &expr; %var=&expr;", "Repeat the following snipped until the next matching --end using the specified init-clause, cond-expression and iteration-expression", commandLoop, nil},
		{"--set %var=&expr", "Assign the EXPR result value to the variable VAR. Each task uses a local copy of the environment variables created during the creation of the task, so if a task is created inside a loop it will see the expected values of the variables involved in the iteration. On assignments, both the local copy of the env vars and the global env vars are updated.", commandSet, nil},
	}
	commands = append(commands, cs[:]...)
}

// command helpers

// check if the input command line is the start of a statement which requires an endcommand
// used to calculate the stack of nested statements
func commandRequiresEnd(linebytes []byte) bool {
	for _, commandprefix := range commandsRequiringEnd {
		if bytes.HasPrefix(linebytes, commandprefix) {
			return true
		}
	}
	return false
}

// check if the input bytes represents and end command (--end)
// TODO: improve the parsing to check also the next bytes, not only the 5 bytes of "--end" because it can be followed by \n, by a space and other chars or directly by other characters (letters, digit, punctuation, etc.)
func commandIsEnd(linebytes []byte) bool {
	return bytes.HasPrefix(linebytes, endcommandb)
}

func newReaderUntilEnd(t *task) (*bytes.Reader, error) { // bufio.Scanner
	// get all the code for this task until the endcommand (--end) of this statement
	// consider nested statments requiring the endcommand (stack of nested statements)
	// consume this code from the caller scanner
	var b bytes.Buffer
	stacklevel := 0
	for t.scanner.Scan() {
		t.line += 1

		trimmedbytes := bytes.TrimLeft(t.scanner.Bytes(), " \t")
		if commandRequiresEnd(trimmedbytes) {
			stacklevel += 1
		} else if commandIsEnd(trimmedbytes) {
			if stacklevel == 0 {
				break
			} else {
				stacklevel -= 1
			}
		}

		_, err := b.Write(t.scanner.Bytes())
		if err != nil {
			return nil, err
		}
		b.Write([]byte("\n"))
	}

	// create a new reader with the code for the new task
	return bytes.NewReader(b.Bytes()), nil
	// return bufio.NewScanner(bytes.NewReader(b.Bytes())), nil
}

func interfaceToInt(interfaceval interface{}) (int, error) {
	intval, ok := interfaceval.(int)
	if !ok {
		floatval, ok := interfaceval.(float64)
		if !ok {
			return 0, fmt.Errorf("can't convert %v to int", interfaceval)
		}
		intval = int(floatval)
	}
	return intval, nil
}

// commands functions

// list the tester builtin commands
func commandListCommands(w *worker, args []string, opts []bool, t *task) error {
	for _, c := range commands {
		w.t.Logf("%s:", c.template)
		w.t.Logf("\t\t%s:", c.desctription)
	}
	return nil
}

func commandSleep(w *worker, args []string, opts []bool, t *task) error {
	if len(args) < 1 {
		err := fmt.Errorf("missing arguments\n")
		return err
	}

	intms, err := strconv.Atoi(args[0])
	if err != nil {
		err = fmt.Errorf("the first argument is invalid: %v\n", err)
		return err
	}

	Debugf("w%d executing SLEEP %d", w.id, intms)
	var ms = time.Duration(intms)
	time.Sleep(ms * time.Millisecond)
	return nil
}

func commandWait(w *worker, args []string, opts []bool, t *task) error {
	if len(args) < 1 {
		err := fmt.Errorf("missing arguments")
		return err
	}

	timeoutms := 10000
	if len(opts) > 0 && opts[0] && len(args) > 1 {
		var err error
		timeoutms, err = strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid timeout argument in wait command: %v", err)
		}
	}

	Debugf("w%d executing WAIT %s TIMEOUT %dms", w.id, args[0], timeoutms)

	switch strings.ToUpper(args[0]) {
	case "ALL":
		if w.id != 0 {
			return fmt.Errorf("only the main worker can use '--wait all'")
		}
		wg.Wait()
	default:
		waitwid, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("the first argument is invalid: %v", err)
		}
		waitw, exists := workers[waitwid]
		if !exists {
			return fmt.Errorf("WAIT can't find worker with id %d", waitwid)
		}
		// wait for waitw to complete
		select {
		case <-waitw.completec:
		case <-time.After(time.Duration(timeoutms) * time.Millisecond):
			return fmt.Errorf("expected: wait %s, got: timeout in %dms", args[0], timeoutms)
		}
	}

	Debugf("w%d completed WAIT %s ", w.id, args[0])

	return nil
}

func commandDump(w *worker, args []string, opts []bool, t *task) error {
	if w.res == nil {
		fmt.Printf("w%d error: %v\n", w.id, w.reserr)
	} else {
		w.res.Dump()
	}
	return nil
}

func commandExit(w *worker, args []string, opts []bool, t *task) error {
	if len(opts) < 1 {
		return fmt.Errorf("missing arguments")
	}

	mustclose := true
	if opts[0] {
		if len(args) < 1 {
			return fmt.Errorf("missing arguments")
		}

		strval := args[0]
		interfaceval, err := gval.Evaluate(strval, t.env)
		if err != nil {
			return fmt.Errorf("expected int arg, got %s (%v)", strval, err)
		}

		intval, err := interfaceToInt(interfaceval)
		if err != nil {
			return fmt.Errorf("expected int arg, got %s (%v)", strval, err)
		}
		mustclose = intval == 0
	}

	if mustclose && w.conn != nil {
		// Debugf("w%d closing connection", w.id)
		if err := w.conn.Close(); err != nil {
			w.t.Logf("w%d conn close error: %v", w.id, err)
		}
		w.conn = nil
	}

	return nil
}

func commandTask(w *worker, args []string, opts []bool, t *task) error {
	if len(args) < 1 {
		return fmt.Errorf("missing arguments")
	}

	if len(opts) < 2 {
		return fmt.Errorf("missing arguments")
	}

	if t.scanner == nil {
		return fmt.Errorf("missing scanner")
	}

	iargs := 0
	// workerid, err := strconv.Atoi(args[iargs])
	interfaceval, err := gval.Evaluate(args[iargs], t.env)
	iargs++
	if err != nil {
		return fmt.Errorf("expected int arg, got %s (%v)", args[iargs], err)
	}

	// transform the interface{} result of Evaluate to int
	workerid, err := interfaceToInt(interfaceval)
	if err != nil {
		return err
	}

	// get the initial line for the code of this task
	tline := t.line + 1

	// get the task name from the first optional arg or set a default one
	var tname string
	if opts[0] {
		tname = args[iargs]
		iargs++
	} else {
		tname = fmt.Sprintf("%s:%d", filepath.Base(t.file.Name()), tline)
	}

	taskworker, exists := workers[workerid]
	if !exists {
		var connstring string
		if opts[1] {
			connstring = args[iargs]
			iargs++
		} else {
			connstring = *pconnstring
		}

		taskworker, err = newWorker(workerid, connstring, w.t)
		if err != nil {
			return err
		}
		workers[workerid] = taskworker
		go taskworker.runLoop()
	}

	// get all the code for this task until the endcommand (--end) of this statement
	// consider nested statments requiring the endcommand (stack of nested statements)
	// consume this code from the caller scanner
	// create a new scanner with the code for the new task and add it to the t.scannerch
	taskreader, err := newReaderUntilEnd(t)
	if err != nil {
		return err
	}

	taskscanner := bufio.NewScanner(taskreader)

	Debugf("w%d new task %s code(%s:%d-%d):\n", taskworker.id, tname, filepath.Base(t.file.Name()), tline, t.line)

	// clear completec before starting new task
	select {
	case <-w.completec:
	default:
	}

	wg.Add(1)
	Debugf("w%d add task %s", taskworker.id, tname)

	envMutex.RLock()
	taskenv := copyMap(env)
	envMutex.RUnlock()

	taskworker.taskc <- &task{name: tname, line: tline, scanner: taskscanner, file: t.file, env: taskenv}

	return nil
}

// func commandMatchIsError(w *worker, args []string, opts []bool, t *task) error {
// 	if !w.res.IsError() {
// 		return fmt.Errorf("expected: error, got:  %v, expected error type", w.res.GetType())
// 	}
// 	return nil
// }

func commandMatchType(w *worker, args []string, opts []bool, t *task) error {
	if len(args) < 1 {
		return fmt.Errorf("missing arguments")
	}

	expected := args[0]

	w.stats.nchecks++

	// check result error case (res == nil, reserr != nil)
	if w.res == nil && strings.ToUpper(expected) != "ERROR" {
		w.stats.nfailedchecks++
		return fmt.Errorf("expected: %s, got error %v", expected, w.reserr)
	}

	match := false
	switch strings.ToUpper(expected) {
	case "OK":
		match = w.res.IsOK()
	case "ERROR":
		if w.res != nil {
			match = w.res.IsError()
		} else {
			match = true
		}
	case "NULL":
		match = w.res.IsNULL()
	case "JSON":
		match = w.res.IsJSON()
	case "STRING":
		match = w.res.IsString()
	case "INTEGER", "INT":
		match = w.res.IsInteger()
	case "FLOAT":
		match = w.res.IsFloat()
	case "PSUB":
		match = w.res.IsPSUB()
	case "COMMAND":
		match = w.res.IsCommand()
	case "RECONNECT":
		match = w.res.IsReconnect()
	case "BLOB":
		match = w.res.IsBLOB()
	case "ARRAY":
		match = w.res.IsArray()
	case "TEXT":
		match = w.res.IsText()
	case "ROWSET":
		match = w.res.IsRowSet()
	case "LITERAL":
		match = w.res.IsLiteral()
	}

	if !match {
		w.stats.nfailedchecks++
		return fmt.Errorf("expected: %s, got: %v", expected, w.res.GetType())
	}

	return nil
}

// helper method to get the string value of a result or its value at ROW COL
func getString(w *worker, args []string, opts []bool) (string, error) {
	wantsrowsetval := opts[0]
	var s string
	if wantsrowsetval {
		if len(args) < 3 {
			return "", fmt.Errorf("missing arguments")
		}
		r, err := strconv.Atoi(args[1])
		if err != nil {
			return "", fmt.Errorf("invalid ROW argument %s", args[1])
		}
		c, err := strconv.Atoi(args[2])
		if err != nil {
			return "", fmt.Errorf("invalid COL argument %s", args[1])
		}

		str, err := w.res.GetStringValue(uint64(r), uint64(c))
		if err != nil {
			return "", err
		}

		s = str

	} else {
		if w.res == nil {
			return "", fmt.Errorf("invalid result")
		}

		var str string
		var err error
		if w.res.IsArray() || w.res.IsRowSet() {
			// it is an array or a rowset but the row/col aren't specified, use default row/col 0 0
			str, err = w.res.GetStringValue(0, 0)
		} else {
			// resbuffer = w.res.GetBuffer()
			str, err = w.res.GetString()
		}

		if err != nil {
			return "", err
		}

		s = str
	}

	return s, nil
}

func getBuffer(w *worker, args []string, opts []bool) ([]byte, error) {
	wantsrowsetval := opts[0]
	var b []byte
	if wantsrowsetval {
		if len(args) < 3 {
			return nil, fmt.Errorf("missing arguments")
		}
		r, err := strconv.Atoi(args[1])
		if err != nil {
			return nil, fmt.Errorf("invalid ROW argument %s", args[1])
		}
		c, err := strconv.Atoi(args[2])
		if err != nil {
			return nil, fmt.Errorf("invalid COL argument %s", args[1])
		}

		v, err := w.res.GetValue(uint64(r), uint64(c))
		if err != nil {
			return nil, err
		}

		b = v.GetBuffer()

	} else {
		if w.res == nil {
			return nil, fmt.Errorf("invalid result")
		}

		b = w.res.GetBuffer()
	}

	return b, nil
}

func commandMatchBuffer(w *worker, args []string, opts []bool, t *task) error {
	if len(args) < 1 {
		return fmt.Errorf("missing arguments")
	}

	if len(opts) < 1 {
		return fmt.Errorf("invalid arguments")
	}

	expected := args[0]

	w.stats.nchecks++

	if w.res == nil {
		w.stats.nfailedchecks++
		return fmt.Errorf("expected: %s, got error %v", expected, w.reserr)
	}

	b, err := getBuffer(w, args, opts)
	if err != nil {
		return err
	}

	if res := bytes.Compare(b, []byte(expected)); res != 0 {
		w.stats.nfailedchecks++
		return fmt.Errorf("expected: %s, got: %v", expected, b)
	}

	return nil
}

func commandMatchValue(w *worker, args []string, opts []bool, t *task) error {
	if len(args) < 1 {
		return fmt.Errorf("missing arguments")
	}

	if len(opts) < 1 {
		return fmt.Errorf("invalid arguments")
	}

	expected := args[0]
	expectedeval, err := gval.Evaluate(expected, t.env)
	expectedevalstring := fmt.Sprint(expectedeval)

	w.stats.nchecks++

	// check result error case (res == nil, reserr != nil)
	if w.res == nil {
		w.stats.nfailedchecks++
		return fmt.Errorf("expected: %s, got error %v", expected, w.reserr)
	}

	s, err := getString(w, args, opts)
	if err != nil {
		w.stats.nfailedchecks++
		return err
	}

	if s != expected && s != expectedevalstring {
		w.stats.nfailedchecks++
		return fmt.Errorf("expected: %s, got: %s", expected, s)
	}

	return nil
}

func commandMatchRowsCols(w *worker, args []string, opts []bool, t *task) error {
	if len(args) < 1 {
		return fmt.Errorf("missing arguments")
	}

	w.stats.nchecks++

	if !w.res.IsRowSet() && !w.res.IsArray() {
		// TODO: print the human-readable representation of the type
		w.stats.nfailedchecks++
		return fmt.Errorf("expected RowSet or Array, got %v", w.res.GetType())
	}

	idx := 0
	optsnames := [2]string{"rows", "cols"}
	for i, optname := range optsnames {
		if opts[i] {
			expected := args[idx]
			idx++
			expectedinterface, _ := gval.Evaluate(expected, t.env)

			expectedint, err := interfaceToInt(expectedinterface)
			if err != nil {
				w.stats.nfailedchecks++
				return fmt.Errorf("expected int arg, got %s (%v)", expected, err)
			}

			val := 0
			if i == 0 {
				val = int(w.res.GetNumberOfRows())
			} else {
				val = int(w.res.GetNumberOfColumns())
			}

			if val != expectedint {
				w.stats.nfailedchecks++
				return fmt.Errorf("expected %d %s, got %d", expectedint, optname, val)
			}
		}
	}

	return nil
}

func copyMap(m map[string]interface{}) map[string]interface{} {
	cp := make(map[string]interface{})
	for k, v := range m {
		vm, ok := v.(map[string]interface{})
		if ok {
			cp[k] = copyMap(vm)
		} else {
			cp[k] = v
		}
	}

	return cp
}

func boolExpr(w *worker, t *task, expression string) bool {
	value, err := gval.Evaluate(expression, t.env)
	if err != nil {
		w.t.Logf("WARNING: w%d %s:%d error evaluating expression:%s with env:%v", w.id, t.name, t.line, expression, t.env)
		return false
	}
	return value.(bool)
}

func assignExpr(w *worker, t *task, key string, expression string) {
	value, err := gval.Evaluate(expression, t.env)
	if err != nil {
		// w.t.Logf("WARNING: w%d %s:%d error evaluating expression:%s with env:%v", w.id, t.name, t.line, expression, w.env)
		value = expression
	}

	t.env[key] = value
	envMutex.Lock()
	env[key] = value
	envMutex.Unlock()

	if debug {
		envMutex.RLock()
		Debugf("assignExpr env:%v", env)
		envMutex.RUnlock()
	}
}

func commandLoop(w *worker, args []string, opts []bool, t *task) error {
	if len(args) < 5 {
		return fmt.Errorf("missing arguments")
	}

	// [initialization]; [condition]; [final-expression];
	// initialization: vartype varname=expression;
	// condition: expression
	// final-expression: varname2=expression
	initvar := strings.Trim(args[0], " ")
	initexpr := args[1]
	condexpr := args[2]
	finalvar := strings.Trim(args[3], " ")
	finalexpr := args[4]

	Debugf("loop initialization: var:%s, initvalue:%s", initvar, initexpr)
	Debugf("loop condition:%s", condexpr)
	Debugf("loop finalexpression: var:%s, expr:%s", finalexpr, finalexpr)

	loopscannerstartline := t.line + 1

	// get new scanner with content from current line to end line
	// consume loop lines in the current scanner
	loopreader, err := newReaderUntilEnd(t)
	if err != nil {
		return err
	}
	Debugf("loop lines:%d-%d", loopscannerstartline-1, t.line)

	// save the scanner and the line of the task to restart from there after the loop ends
	stackedscanner := t.scanner
	stackedline := t.line

	// setup the current task with the scanner of the loop
	// perform the loop using [initialization]; [condition]; [final-expression];
	// parse expressions with github.com/PaesslerAG/gval
	// initial initiza
	for assignExpr(w, t, initvar, initexpr); boolExpr(w, t, condexpr); assignExpr(w, t, finalvar, finalexpr) {
		loopreader.Seek(0, io.SeekStart)
		t.scanner = bufio.NewScanner(loopreader)
		t.line = loopscannerstartline
		err := w.processTask(t)
		if err != nil {
			return err
		}
		if w.exiting {
			break
		}
	}

	// restore the task with the saved scanner and line to restart the execution after the loop
	t.scanner = stackedscanner
	t.line = stackedline

	return nil
}

func commandSet(w *worker, args []string, opts []bool, t *task) error {
	if len(args) < 2 {
		return fmt.Errorf("missing arguments")
	}

	varname := strings.Trim(args[0], " ")
	expr := args[1]
	assignExpr(w, t, varname, expr)

	if debug {
		Debugf("set var:%s, expr:%s, value:%v", varname, expr, t.env[varname])
	}

	return nil
}

// ----------------------------- Parser & Lexer -----------------------------
// https://blog.gopheracademy.com/advent-2014/parsers-lexers/

// Token represents a lexical token.
type Token struct {
	id        TokenID
	lit       string
	isopt     bool
	islastopt bool
}

type TokenID int

const (
	// Special tokens
	ILLEGAL TokenID = iota
	EOF
	WS
	SEMICOLON
	MATH

	// Literals
	IDENT // arguments %xxx
	EXPR  // expressions &xx xx xx;

	// Keyword
	KEYWORD
)

var eof = rune(0)
var identch = '%'
var exprch = '&'
var optch = '['
var optendch = ']'

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}

func isSemicolon(ch rune) bool {
	return ch == ';'
}

func isExprCh(ch rune) bool {
	return ch == exprch
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isMath(ch rune) bool {
	return unicode.Is(unicode.Sm, ch)
}

func isValidExprCh(ch rune) bool {
	if isSemicolon(ch) || isStartOptional(ch) {
		return false
	}
	return unicode.IsLetter(ch) || unicode.IsDigit(ch) || unicode.IsSymbol(ch) || unicode.IsPunct(ch) || ch == '_' || ch == '-' || ch == identch || ch == exprch
}

func isValidIdentCh(ch rune) bool {
	if isSemicolon(ch) || isMath(ch) || isStartOptional(ch) {
		return false
	}
	return unicode.IsLetter(ch) || unicode.IsDigit(ch) || unicode.IsSymbol(ch) || unicode.IsPunct(ch) || ch == '_' || ch == '-' || ch == identch || ch == exprch
}

func isStartOptional(ch rune) bool {
	return (ch == optch)
}

func isEndOptional(ch rune) bool {
	return (ch == optendch)
}

func isDigit(ch rune) bool {
	return (ch >= '0' && ch <= '9')
}

// Scanner represents a lexical scanner.
type Scanner struct {
	r *bufio.Reader
}

// NewScanner returns a new instance of Scanner.
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

// read reads the next rune from the bufferred reader.
// Returns the rune(0) if an error occurs (or io.EOF is returned).
func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

// unread places the previously read rune back on the reader.
func (s *Scanner) unread() { _ = s.r.UnreadRune() }

// Scan returns the next token and literal value.
func (s *Scanner) Scan(istemplate bool, isopt bool) Token {
	// Read the next rune.
	ch := s.read()

	// If we see whitespace then consume all contiguous whitespace.
	// If we see a letter then consume as an ident or reserved word.
	if isWhitespace(ch) {
		s.unread()
		return s.scanWhitespace()
	} else if isSemicolon(ch) {
		s.unread()
		return s.scanSemicolon()
	} else if isMath(ch) {
		s.unread()
		return s.scanMath()
	} else if isExprCh(ch) {
		s.unread()
		return s.scanExpr(istemplate, isopt && istemplate)
	} else if isValidIdentCh(ch) || (!istemplate && isStartOptional(ch)) {
		s.unread()
		return s.scanIdent(istemplate, isopt && istemplate)
	} else if istemplate && isStartOptional(ch) {
		return s.scanIdent(istemplate, true)
	}

	// Otherwise read the individual character.
	switch ch {
	case eof:
		return Token{EOF, "", false, false}
	}

	return Token{ILLEGAL, string(ch), false, false}
}

// scanWhitespace consumes the current rune and all contiguous whitespace.
func (s *Scanner) scanWhitespace() Token {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent whitespace character into the buffer.
	// Non-whitespace characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return Token{WS, buf.String(), false, false}
}

// scanSemicolon consumes the current rune
func (s *Scanner) scanSemicolon() Token {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	return Token{SEMICOLON, buf.String(), false, false}
}

// scanEquals consumes the current rune
func (s *Scanner) scanMath() Token {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent whitespace character into the buffer.
	// Non-whitespace characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isMath(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return Token{MATH, buf.String(), false, false}
}

// scanExpr consumes the current rune and all contiguous valid runes for expressions.
func (s *Scanner) scanExpr(istemplate bool, isopt bool) Token {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent ident character into the buffer.
	// Non-ident characters and EOF will cause the loop to exit.
	islastopt := false
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isValidExprCh(ch) && (!isopt || !isEndOptional(ch)) {
			s.unread()
			break
		} else if !isEndOptional(ch) {
			_, _ = buf.WriteRune(ch)
		} else if isopt && isEndOptional(ch) {
			islastopt = true
		}
	}

	return Token{EXPR, buf.String(), isopt, islastopt}
}

// scanIdent consumes the current rune and all contiguous ident runes.
func (s *Scanner) scanIdent(istemplate bool, isopt bool) Token {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent ident character into the buffer.
	// Non-ident characters and EOF will cause the loop to exit.
	islastopt := false
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isValidIdentCh(ch) && (!isopt || !isEndOptional(ch)) {
			s.unread()
			break
		} else if !isEndOptional(ch) {
			_, _ = buf.WriteRune(ch)
		} else if isopt && isEndOptional(ch) {
			islastopt = true
		}
	}

	if istemplate && !strings.HasPrefix(buf.String(), string(identch)) {
		return Token{KEYWORD, strings.ToUpper(buf.String()), isopt, islastopt}
	}

	// Otherwise return as a regular identifier.
	return Token{IDENT, buf.String(), isopt, islastopt}
}

// Parser represents a parser.
type Parser struct {
	s   *Scanner
	buf struct {
		tok Token // last read token
		n   int   // buffer size (max=1)
	}
}

// NewParser returns a new instance of Parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

// scan returns the next token from the underlying scanner.
// If a token has been unscanned then read that instead.
func (p *Parser) scan(istemplate bool, isopt bool) (tok Token) {
	// If we have a token on the buffer, then return it.
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok
	}

	// Otherwise read the next token from the scanner.
	tok = p.s.Scan(istemplate, isopt)

	// Save it to the buffer in case we unscan later.
	p.buf.tok = tok

	return
}

// unscan pushes the previously read token back onto the buffer.
func (p *Parser) unscan() { p.buf.n = 1 }

// scanIgnoreWhitespace scans the next non-whitespace token.
func (p *Parser) scanIgnoreWhitespace(istemplate bool, isopt bool) (tok Token) {
	tok = p.scan(istemplate, isopt)
	if tok.id == WS {
		tok = p.scan(istemplate, isopt)
	}
	return
}

// Parse a builtin template
func (p *Parser) ParseTemplate() ([]Token, error) {
	var tokens []Token
	tok := Token{}
	for tok.id != EOF {
		tok = p.scanIgnoreWhitespace(true, tok.isopt && !tok.islastopt)
		if tok.id != EOF {
			tokens = append(tokens, tok)
		}
	}
	return tokens, nil
}

func (p *Parser) Parse() (*command, []string, []bool, error) {
	var tokens = make([]Token, 0, 20)
	tok := Token{}
	for tok.id != EOF {
		tok = p.scanIgnoreWhitespace(false, false)
		if tok.id != EOF {
			tokens = append(tokens, tok)
		}
	}

	// must not be empty to go on
	if len(tokens) == 0 {
		return nil, nil, nil, fmt.Errorf("found 0 tokens")
	}

	// lookup the tokens of the parsed command in the builtin commands array
	args := make([]string, 0, 20)
	opts := make([]bool, 0, 20)
	var expr bytes.Buffer
	for _, c := range commands {
		found := false
		idx := 0
		skipopttokens := false
		for _, ctok := range c.tokens {
			if idx < len(tokens) {
				if skipopttokens && ctok.isopt {
					if ctok.islastopt {
						skipopttokens = false
					}
					continue
				}

				// get the next token of the parsed line, if any
				var nextparsedtok = &tokens[idx]
				// check if the next token of the parsed line match
				// a constant keyword of the template
				if ctok.id == KEYWORD && ctok.lit == strings.ToUpper(nextparsedtok.lit) {
					if ctok.isopt && ctok.islastopt {
						opts = append(opts, true)
					}
					idx += 1
					found = true
					continue
				}

				// check if the next token of the parsed line match
				// a variable argument of the template
				if ctok.id == IDENT {
					if ctok.islastopt {
						args = append(args, nextparsedtok.lit)
						opts = append(opts, true)
					} else {
						args = append(args, nextparsedtok.lit)
					}
					idx += 1
					found = true
					continue
				}

				// check if the next token of the parsed line match
				// a variable argument of the template
				if ctok.id == EXPR {
					// an expression can be made by multiple IDENT tokens
					// create the expression string by concatenating all next IDENT tokens until SEMICOLON,
					// a space must be added between the tokens
					for nextparsedtok.id == IDENT || nextparsedtok.id == MATH {
						expr.WriteString(nextparsedtok.lit)
						if idx < len(tokens)-1 {
							expr.WriteString(" ")
							idx += 1
							nextparsedtok = &tokens[idx]
						} else {
							// end of tokens, so stop
							nextparsedtok = nil
							break
						}
					}

					// if the end of tokens or SEMICOLON has been reached, append the expr to args
					if nextparsedtok == nil || nextparsedtok.id == SEMICOLON {
						if ctok.islastopt {
							args = append(args, strings.Trim(expr.String(), " "))
							opts = append(opts, true)
						} else {
							args = append(args, strings.Trim(expr.String(), " "))
						}
						// unread the SEMICOLON token, it will match the next ctok
						if nextparsedtok != nil {
							idx -= 1
						}
					} else if !ctok.isopt {
						found = false
						break
					}

					// reset the expr buffer
					expr.Reset()

					idx += 1
					found = true
					continue
				}

				// check if id and lit are equals between ctok and nextparsedtok
				// for example SEMICOLON or MATH token ids
				if nextparsedtok.id == ctok.id && ctok.lit == strings.ToUpper(nextparsedtok.lit) {
					idx += 1
					found = true
					continue
				}
			} else if !ctok.isopt {
				// check if template token is not optional and
				// the current line has no more tokens -> not match
				found = false
				break
			}

			// the next token of the parsed line didn't match
			// the current token of the template, but it was optional
			if ctok.isopt {
				opts = append(opts, false)
				skipopttokens = true
				continue
			}

			// current token doesn't match and it is not optional,
			// the lookup failed for this builtin command,
			// so try with the next one
			found = false
			break
		}
		if found {
			return &c, args, opts, nil
		}
	}

	// didn't find any match
	return nil, nil, nil, fmt.Errorf("command not found")
}

// Worker Functions
func (w *worker) runLoop() {
	iserror := false
	for {
		select {
		case task := <-w.taskc:
			err := w.processTask(task)
			if err != nil {
				iserror = true
				w.stats.failed = true
				w.t.Errorf("w%d, task:%s(:%d), %v", w.id, task.name, task.line, err)
				Debugf("w%d failed task %s, error: %v", w.id, task.name, err)
			} else {
				Debugf("w%d completed task %s", w.id, task.name)
			}

			// only real tasks must notify the completec and the wg, not the helper connection tasks
			if task.connstring == "" {
				// nonblocking write, skip if nobody is waiting or channel is full
				select {
				case w.completec <- struct{}{}:
				default:
				}

				Debugf("w%d task %s done", w.id, task.name)
				wg.Done()
			} else {
				Debugf("w%d connection task %s connected", w.id, task.connstring)
			}

			if iserror {
				tryStopc(false)
			}

			if w.exiting {
				break
			}

		case <-stopc:
			Debugf("w%d stopped", w.id)
			if w.conn != nil {
				if err := w.conn.Close(); err != nil {
					w.t.Logf("w%d conn close error: %v", w.id, err)
				}
			}

			// consume all pending tasks and call wg.Done() for each
			end := false
			for !end {
				select {
				case t := <-w.taskc:
					select {
					case w.completec <- struct{}{}:
					default:
					}
					Debugf("w%d task %s done (stopc)", w.id, t.name)
					wg.Done()

				default:
					end = true
				}
			}

			return
		}
	}
}

func replaceEnvVarInSqlString(s string, w *worker, t *task) string {
	tmpl, err := template.New("sql").Parse(s)
	buf := &bytes.Buffer{}
	if err = tmpl.Execute(buf, t.env); err != nil {
		fmt.Println(err)
		return s
	}
	return buf.String()
}

func (w *worker) processTask(task *task) error {
	// check if this is a connect task
	if w.conn == nil {
		if task.connstring == "" {
			return fmt.Errorf("Worker not connected")
		}

		// it's only a connection task
		var err error
		Debugf("w%d connecting ...", w.id)
		if w.conn, err = sqlitecloud.Connect(task.connstring); err != nil {
			Debugf("w%d connection error %v", w.id, err)
			w.conn = nil
			return err
		}
		Debugf("w%d connected", w.id)
		return nil
	}

	// it's a normal task
	// the code can be from a scanner or from a file
	// prepare the scanner from file if needed
	if task.scanner == nil {
		if task.file == nil {
			return fmt.Errorf("invalid task")
		}

		Debugf("w%d processing file %s", w.id, filepath.Base(task.file.Name()))
		defer task.file.Close()

		task.scanner = bufio.NewScanner(task.file)
	}

	Debugf("w%d processing task %s (%s:%d)", w.id, task.name, filepath.Base(task.file.Name()), task.line)

	// process the scanner line by line
	task.scanner.Split(bufio.ScanLines)
	for task.scanner.Scan() {
		Debugf("w%d processing line: %s", w.id, task.scanner.Text())

		// non-blocking test on stopc
		select {
		case <-stopc:
			Debugf("w%d stopping processTask %s", w.id, task.name)
			return nil
		default:
		}

		// check if the line starts with "--" using scanner.Bytes() instead of
		// scanner.Text() to avoid allocation
		if bytes.HasPrefix(bytes.TrimLeft(task.scanner.Bytes(), " \t"), []byte("--")) {
			// is a comment, maybe a test command
			Debugf("w%d parsing command %s ...", w.id, task.scanner.Text())
			linereader := bytes.NewReader(task.scanner.Bytes())
			p := NewParser(linereader)
			c, args, opts, err := p.Parse()
			if err != nil {
				// command not found, it's just a comment
				Debugf("... parse error: %v", err)
			} else {
				// command found
				Debugf("w%d parsed command %s with args: %v", w.id, c.template, args)
				err = c.f(w, args, opts, task)
				if err != nil {
					return err
				}

				if w.exiting {
					break
				}
			}
		} else {
			// not a command, pass it to sqlitecloud
			sql := strings.Trim(task.scanner.Text(), " \t")
			sql = replaceEnvVarInSqlString(sql, w, task)

			// execute sql, but skip if it is empty
			if len(sql) > 0 {
				if w.res != nil {
					w.res.Free()
				}

				Debugf("w%d executing...: %s", w.id, sql)

				// TODO: change the name of the function from Select to Execute
				w.res, w.reserr = w.conn.Select(sql)
				
				// if err is busy then retry, no more than retrytimes
				retrytimes := 10
				for i := 1; w.reserr != nil && w.conn.ErrorCode == 5 && i < retrytimes; i++ { 
					// TODO: in case of explicit transaction, we should retry the full transaction instead of single commands
					Debugf("w%d busy error, retry(%d): %s", w.id, i, sql)
					w.stats.nretries++
					w.res, w.reserr = w.conn.Select(sql)
				}

				if debug {
					Debugf("w%d executed: %s", w.id, sql)
					if w.res != nil {
						w.res.Dump()
					}
					if w.reserr != nil {
						println(w.reserr.Error())
					}
				}
			}
		}
		task.line++
	}

	return nil
}

func processFile(t *testing.T, path string, connstring string) {
	file, err := os.Open(path)
	if err != nil {
		t.Errorf("w0, %v", err)
		trySkip(t)
	}
	// start default task
	wg = sync.WaitGroup{}
	stopc = make(chan struct{})
	workers = make(map[int]*worker)
	env = make(map[string]interface{})
	envMutex = sync.RWMutex{}

	envMutex.RLock()
	taskenv := copyMap(env)
	envMutex.RUnlock()

	task := task{name: filepath.Base(file.Name()), line: 1, file: file, env: taskenv}
	t.Run(task.name, func(t *testing.T) {
		w, err := newWorker(0, connstring, t)
		if err != nil {
			w.t.Errorf("w%d, task:%s(l:%d), %v", w.id, task.name, task.line, err)
			trySkip(t)
		}
		workers[0] = w

		// the connection task was added by newWorker()
		// the main worker (0) doens't execute the run loop like other workers
		// so get the conntask and run it synchronously before the main task from the script file
		conntask := <-w.taskc
		err = w.processTask(conntask)
		Debugf("w%d connection task %s connected", w.id, conntask.connstring)
		if err != nil {
			w.t.Errorf("w%d, task:%s(l:%d), %v", w.id, task.name, task.line, err)
			w.stats.failed = true
			trySkip(t)
		}

		err = w.processTask(&task)
		if err != nil {
			w.t.Errorf("w%d, task:%s(l:%d), %v", w.id, task.name, task.line, err)
			w.stats.failed = true
			trySkip(t)
			// Errorf("w%d, task:%s(%d), %v", w.id, task.name, task.line, err)
		}

		if !w.exiting {
			// Debugf("w%d closing connection", w.id)
			if w.conn != nil {
				if err := w.conn.Close(); err != nil {
					w.t.Logf("w%d conn close error: %v", w.id, err)
				}
				w.conn = nil
			}
		}

		printStats(t)
	})

	Debugf("w0 stopping ...")
	tryStopc(true)

	// give time to workers loops to close
	time.Sleep(time.Duration(100) * time.Millisecond)
	Debugf("w0 stopped")
}

func printStats(t *testing.T) {
	if !testing.Verbose() {
		return
	}

	nchecks := 0
	nfailedchecks := 0
	nfailedworkers := 0
	nretries := 0
	for _, w := range workers {
		nchecks += w.stats.nchecks
		nfailedchecks += w.stats.nfailedchecks
		if w.stats.failed {
			nfailedworkers += 1
		}
		nretries += w.stats.nretries
	}
	fmt.Printf("    --- STATS: %s workers:%d/%d, checks:%d/%d, retries:%d\n", t.Name(), len(workers)-nfailedworkers, len(workers), nchecks-nfailedchecks, nchecks, nretries)
}

func trySkip(t *testing.T) {
	if skip {
		t.SkipNow()
	}
}

func init() {
	initCommands()
}

var ppath = flag.String("path", "scripts", "File or Directory containing the test scripts")
var pconnstring = flag.String("connstring", "sqlitecloud://dev1.sqlitecloud.io/", "Connection string for the main worker")
var pdebug = flag.Bool("debug", false, "Enable debug logs")
var debug = false
var pskip = flag.Bool("skip", false, "Skip immediately in case of errors")
var skip = false

// func main() {
func TestTester(t *testing.T) {
	// parse flags
	flag.Parse()
	path := *ppath
	connstring := *pconnstring
	debug = *pdebug
	skip = *pskip

	Debugf("Parsing commands ...")
	for i, c := range commands {
		r := strings.NewReader(c.template)
		p := NewParser(r)
		tokens, err := p.ParseTemplate()
		Debugf("command:\"%s\", tokens: %v", c.template, tokens)
		if err != nil {
			fmt.Println(err)
			continue
		}
		commands[i].tokens = tokens
	}
	Debugf("Parsing commands end")

	// open source script file or dir
	var wd string
	if !filepath.IsAbs(path) {
		wd, _ = os.Getwd()
		path = filepath.Join(wd, path)

	}
	fi, err := os.Stat(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	switch mode := fi.Mode(); {
	case mode.IsDir():
		// do directory stuff
		files, err := ioutil.ReadDir(path)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		var dirpath string
		if len(wd) > 0 {
			dirpath = filepath.Join(wd, fi.Name())
		} else {
			dirpath = fi.Name()
		}

		for _, f := range files {
			if filepath.Ext(f.Name()) != ".test" {
				continue
			}

			path = filepath.Join(dirpath, f.Name())
			processFile(t, path, connstring)
		}
	case mode.IsRegular():
		processFile(t, path, connstring)
	}
}
