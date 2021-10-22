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
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"runtime"
	"sqlitecloud"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
	"unicode"

	"github.com/PaesslerAG/gval"
)

// ----------------------------- Debugger -----------------------------

var debug = true // true // false
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
	name    string
	line    int
	scanner *bufio.Scanner
	file    *os.File
}

type worker struct {
	// scannerc  chan *bufio.Scanner
	// filec     chan *os.File
	id        int
	conn      *sqlitecloud.SQCloud
	res       *sqlitecloud.Result
	taskc     chan *task
	stopc     chan struct{}
	stopped   bool
	wg        *sync.WaitGroup // WaitGroup used for --wait all
	completec chan struct{}   // channel used to wake up a worker waiting for this one to complete
	t         *testing.T
	env       map[string]interface{}
}

var workers map[int]*worker // = make(map[int]*worker)

func newWorker(id int, connstring string, wg *sync.WaitGroup, stopc chan struct{}, t *testing.T) (*worker, error) {
	Debugf("w%d created with connstring: %s", id, connstring)
	db, err := sqlitecloud.Connect(connstring) // "sqlitecloud://dev1.sqlitecloud.io/X"
	if err != nil {
		return nil, err
	}
	w := worker{id: id, conn: db, wg: wg}
	// t.scannerc = make(chan *bufio.Scanner, 100)
	// t.filec = make(chan *os.File, 100)
	w.taskc = make(chan *task, 100)
	w.completec = make(chan struct{}, 1)
	w.stopc = stopc
	w.t = t
	w.env = make(map[string]interface{})
	return &w, nil
}

// ----------------------------- COMMANDS -----------------------------

type commandFunc func(w *worker, args []string, opts []bool, t *task) error // scanner *bufio.Scanner

type command struct {
	template string
	f        commandFunc
	tokens   []Token
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
		{"--sleep %ms", commandSleep, nil},
		{"--wait %workerid_or_all [%timeout_ms]", commandWait, nil},
		{"--task %workerid [name %name] [%connectionstring]", commandTask, nil},
		{"--match type is %type", commandMatchType, nil},
		{"--match buffer %string [%row %col]", commandMatchBuffer, nil},
		{"--match value %value [%row %col]", commandMatchValue, nil},
		{"--loop %var=%val; &expr; %var=&expr;", commandLoop, nil},
	}
	commands = append(commands, cs[:]...)
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
		w.wg.Wait()
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
		} else if stacklevel == 0 && commandIsEnd(trimmedbytes) {
			break
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
	workerid, err := strconv.Atoi(args[iargs])
	iargs++
	if err != nil {
		err = fmt.Errorf("the first argument is invalid: %v\n", err)
		return err
	}

	// get the initial line for the code of this task
	tline := t.line

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
		if len(args) < 2 {
			err := fmt.Errorf("missing arguments\n")
			return err
		}

		connstring := args[iargs]
		taskworker, err = newWorker(workerid, connstring, w.wg, w.stopc, w.t)
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

	w.wg.Add(1)
	taskworker.taskc <- &task{name: tname, line: tline, scanner: taskscanner, file: t.file}

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
	match := false
	switch strings.ToUpper(expected) {
	case "ERROR":
		match = w.res.IsError()
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
		return fmt.Errorf("expected: %s, got: %v", expected, w.res.GetType())
	}

	return nil
}

type resvalue interface {
	IsOK() bool
	IsNULL() bool
	IsString() bool
	IsJSON() bool
	IsInteger() bool
	IsFloat() bool
	IsBLOB() bool
	IsPSUB() bool
	IsCommand() bool
	IsReconnect() bool
	IsError() bool
	IsRowSet() bool
	IsArray() bool

	GetBuffer() []byte
	// GetString() string
	//GetString() (string, error)
	GetInt32() (int32, error)
	GetInt64() (int64, error)
	GetFloat32() (float32, error)
	GetFloat64() (float64, error)
	GetError() (int, string, error)
}

// helper method to get the Value of a result or the Value at ROW COL with GetValue(row, col)
func getResValue(w *worker, args []string, opts []bool) (*resvalue, error) {
	wantsrowsetval := opts[0]
	var resvalue resvalue
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

		resvalue = v
	} else {
		// TODO should get the Value from res and parse with the same code as before
		if w.res == nil {
			return nil, fmt.Errorf("invalid result")
		}

		resvalue = w.res
		// resbuffer = w.res.GetBuffer()
	}

	return &resvalue, nil
}

func commandMatchBuffer(w *worker, args []string, opts []bool, t *task) error {
	if len(args) < 1 {
		return fmt.Errorf("missing arguments")
	}

	if len(opts) < 1 {
		return fmt.Errorf("invalid arguments")
	}

	expected := args[0]
	rv, err := getResValue(w, args, opts)
	if err != nil {
		return err
	}
	resvalue := *rv

	if res := bytes.Compare(resvalue.GetBuffer(), []byte(expected)); res != 0 {
		return fmt.Errorf("expected: %s, got: %v", expected, resvalue.GetBuffer())
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
	rv, err := getResValue(w, args, opts)
	if err != nil {
		return err
	}
	resvalue := *rv

	if resvalue.IsNULL() {
		if strings.ToUpper(expected) != "NULL" {
			return fmt.Errorf("expected: %s, got: NULL", expected)
		}
	} else if resvalue.IsString() || resvalue.IsJSON() {
		// TODO: GetString() has different prototypes in Value and Result
		// s, err := resvalue.GetString()
		// if err != nil {
		// 	return nil
		// }

		// if s != strings.ToUpper(expected) {
		// 	return fmt.Errorf("Result value is %s, expected %s", s, expected)
		// }
	} else if resvalue.IsInteger() {
		v, err := resvalue.GetInt32()
		if err != nil {
			return err
		}
		if expv, _ := strconv.Atoi(expected); v != int32(expv) {
			return fmt.Errorf("expected: %s, got: %d", expected, v)
		}
	} else if resvalue.IsFloat() {
		v, err := resvalue.GetFloat64()
		if err != nil {
			return err
		}

		expv, err := strconv.ParseFloat(expected, 64)
		result := big.NewFloat(v).Cmp(big.NewFloat(expv))
		if result != 0 {
			return fmt.Errorf("expected: %s, got: %f", expected, v)
		}
	} else if resvalue.IsBLOB() {
		data, err := base64.StdEncoding.DecodeString(expected)
		if err != nil {
			return err
		}
		if res := bytes.Compare(resvalue.GetBuffer(), data); res != 0 {
			return fmt.Errorf("expected: %v, got: %v", data, resvalue.GetBuffer())
		}
	}

	/*	TODO:
		func (this *Value) IsPSUB()      bool { return this.GetType() == '|' }
		func (this *Value) IsCommand()   bool { return this.GetType() == '^' }
		func (this *Value) IsReconnect() bool { return this.GetType() == '@' }
		func (this *Value) IsError()     bool { return this.GetType() == '-' }
		func (this *Value) IsRowSet()    bool { return this.GetType() == '*' }
		func (this *Value) IsArray()
	*/

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

func boolExpr(w *worker, t *task, expression string, env map[string]interface{}) bool {
	value, err := gval.Evaluate(expression, env)
	if err != nil {
		w.t.Logf("WARNING: w%d %s:%d error evaluating expression:%s with env:%v", w.id, t.name, t.line, expression, env)
		return false
	}
	return value.(bool)
}

func assignExpr(w *worker, t *task, key string, expression string, env map[string]interface{}) {
	value, err := gval.Evaluate(expression, env)
	if err != nil {
		w.t.Logf("WARNING: w%d %s:%d error evaluating expression:%s with env:%v", w.id, t.name, t.line, expression, env)
		return
	}
	env[key] = value
	Debugf("assignExpr env:%v", env)
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

	// copy the env variable and possibly override variables in the loop scope
	loopenv := copyMap(w.env)

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
	for assignExpr(w, t, initvar, initexpr, loopenv); boolExpr(w, t, condexpr, loopenv); assignExpr(w, t, finalvar, finalexpr, loopenv) {
		loopreader.Seek(0, io.SeekStart)
		t.scanner = bufio.NewScanner(loopreader)
		t.line = loopscannerstartline
		w.processTask(t)
	}

	// restore the task with the saved scanner and line to restart the execution after the loop
	t.scanner = stackedscanner
	t.line = stackedline

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
	EQUALS

	// Literals
	IDENT // arguments %xxx
	EXPR  // expressions $xx xx xx;

	// Keyword
	KEYWORD
)

var eof = rune(0)
var identch = '%'
var exprch = '&'
var varch = '$'
var optch = '['
var optendch = ']'

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}

func isSemicolon(ch rune) bool {
	return ch == ';'
}

func isEquals(ch rune) bool {
	return ch == '='
}

func isExprCh(ch rune) bool {
	return ch == exprch
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isValidExprCh(ch rune) bool {
	if isSemicolon(ch) || isStartOptional(ch) {
		return false
	}
	return unicode.IsLetter(ch) || unicode.IsDigit(ch) || unicode.IsSymbol(ch) || unicode.IsPunct(ch) || ch == '_' || ch == '-' || ch == identch || ch == exprch
}

func isValidIdentCh(ch rune) bool {
	if isSemicolon(ch) || isEquals(ch) || isStartOptional(ch) {
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
	} else if isEquals(ch) {
		s.unread()
		return s.scanEquals()
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
func (s *Scanner) scanEquals() Token {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	return Token{EQUALS, buf.String(), false, false}
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

// Parse a builtin command
func (p *Parser) ParseCommand() ([]Token, error) {
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
					for nextparsedtok.id == IDENT {
						expr.WriteString(nextparsedtok.lit)
						expr.WriteString(" ")
						if idx < len(tokens)-1 {
							idx += 1
							nextparsedtok = &tokens[idx]
						} else {
							// end of tokens, so stop
							break
						}
					}

					// if the SEMICOLON has been reached, append the expr to args
					if nextparsedtok.id == SEMICOLON {
						if ctok.islastopt {
							args = append(args, expr.String())
							opts = append(opts, true)
						} else {
							args = append(args, expr.String())
						}
						// unread the SEMICOLON token, it will match the next ctok
						idx -= 1
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
				// for example SEMICOLON or EQUALS token ids
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
				w.t.Errorf("w%d, task:%s(%d), %v", w.id, task.name, task.line, err)
				iserror = true
				Debugf("w%d failed task %s, error: %v", w.id, task.name, err)
			} else {
				Debugf("w%d completed task %s", w.id, task.name)
			}

		case <-w.stopc:
			Debugf("w%d stopped", w.id)
			return
		}

		// nonblocking write, skip if nobody is waiting or channel is full
		select {
		case w.completec <- struct{}{}:
		default:
		}

		w.wg.Done()

		if iserror {
			close(w.stopc)
		}
	}
}

func (w *worker) processTask(task *task) error { // scanner *bufio.Scanner
	if task.scanner == nil {
		if task.file == nil {
			return fmt.Errorf("invalid task")
		}

		Debugf("w%d processing file %s", w.id, filepath.Base(task.file.Name()))
		defer task.file.Close()

		task.scanner = bufio.NewScanner(task.file)
	}

	Debugf("w%d processing task %s (%s:%d)", w.id, task.name, filepath.Base(task.file.Name()), task.line)

	task.scanner.Split(bufio.ScanLines)
	for task.scanner.Scan() {
		Debugf("w%d processing line: %s", w.id, task.scanner.Text())

		// non-blocking test on stopc
		select {
		case <-w.stopc:
			Debugf("w%d stopping processTask %s", w.id, task.name)
			w.stopped = true
			return nil
		default:
		}

		// check if the line starts with "--" using scanner.Bytes() instead of
		// scanner.Text() to avoid allocation
		if bytes.HasPrefix(bytes.TrimLeft(task.scanner.Bytes(), " \t"), []byte("--")) {
			// is a comment, maybe a test command
			Debugf("Parsing command %s ...", task.scanner.Text())
			linereader := bytes.NewReader(task.scanner.Bytes())
			p := NewParser(linereader)
			c, args, opts, err := p.Parse()
			if err != nil {
				// command not found, it's just a comment
				Debugf("... parse error: %v", err)
			} else {
				// command found
				Debugf("w%d parsed command with args: %v", w.id, args)
				err = c.f(w, args, opts, task)
				if err != nil {
					return err
				}
			}
		} else {
			// not a command, pass it to sqlitecloud
			if w.res != nil {
				w.res.Free()
			}
			w.res, _ = w.conn.Select(task.scanner.Text())

			if debug {
				Debugf("w%d executing sqlite statement: %s...", w.id, task.scanner.Text())
				w.res.Dump()
			}
		}
		task.line++
	}

	return nil
}

func init() {
	initCommands()
}

var ppath = flag.String("path", "tester-scripts", "File or Directory containing the test scripts")
var pconnstring = flag.String("connstring", "sqlitecloud://dev1.sqlitecloud.io/X", "Connection string for the main worker")

func processFile(t *testing.T, path string, connstring string) {
	file, err := os.Open(path)
	if err != nil {
		t.Errorf("w0, %v", err)
		t.SkipNow()
	}
	// start default task
	var wg sync.WaitGroup
	stopc := make(chan struct{})
	workers = make(map[int]*worker)

	stopped := false
	task := task{name: filepath.Base(file.Name()), line: 1, file: file}
	t.Run(task.name, func(t *testing.T) {
		w, err := newWorker(0, connstring, &wg, stopc, t)
		if err != nil {
			w.t.Errorf("w%d, task:%s(l:%d), %v", w.id, task.name, task.line, err)
			t.SkipNow()
		}
		workers[0] = w

		err = w.processTask(&task)
		if err != nil {
			w.t.Errorf("w%d, task:%s(l:%d), %v", w.id, task.name, task.line, err)
			t.SkipNow()
			// Errorf("w%d, task:%s(%d), %v", w.id, task.name, task.line, err)
		}

		stopped = w.stopped
	})

	Debugf("w0 stopping ...")
	if !stopped {
		close(stopc)
	}

	// give time to workers loops to close
	time.Sleep(time.Duration(100) * time.Millisecond)
	Debugf("w0 stopped")
}

// func main() {
func TestTester(t *testing.T) {
	// path := "tester-scripts"
	// connstring := "sqlitecloud://dev1.sqlitecloud.io/X"

	// parse flags
	flag.Parse()
	path := *ppath
	connstring := *pconnstring

	Debugf("Parsing commands ...")
	for i, c := range commands {
		r := strings.NewReader(c.template)
		p := NewParser(r)
		tokens, err := p.ParseCommand()
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
