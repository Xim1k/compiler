package main

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
	// "os"
)

type stateFn func(*lexer) stateFn
type doubleStateFn func(*lexer, itemType) stateFn

type Pos int

func (p Pos) Position() Pos {
	return p
}


type item struct {
	typ  itemType
	pos  Pos
	val  string
	line int
}

func (i item) String() string {
	switch {
	case i.typ == itemEOF:
		return "EOF"
	case i.typ == itemError:
		return i.val
	case i.typ > itemKeyword:
		return fmt.Sprintf("<%s>", i.val)
	case len(i.val) > 10:
		return fmt.Sprintf("%.10q...", i.val)
	}
	return fmt.Sprintf("%q", i.val)
}

type itemType int

var keyWords = map[string]itemType{
	".":        itemDot,
	"block":    itemBlock,
	"define":   itemDefine,
	"else":     itemElse,
	"end":      itemEnd,
	"if":       itemIf,
	"range":    itemRange,
	"nil":      itemNil,
	"template": itemTemplate,
	"with":     itemWith,
	"for": itemFor,
	"package": itemPackage,
	"import": itemImport,
  "func": itemFunctionDefine,
  "map": itemMap,
  "var": itemVar,
  "int": itemIntType,
  "byte": itemByteType,
  "string": itemStringType,
}

type lexer struct {
	input          string
	pos            Pos
	start          Pos
	width          Pos
	items          chan item
	line           int
	startLine      int
	currentPosOnLine Pos
	currentStartOnLine Pos
	previousUnknown bool
}

const (
	itemError        itemType = iota
	itemBool
	itemChar
	itemCharConstant
	itemAssign
	itemDeclare
	itemEOF
	itemFunction
	itemField
	itemIdentifier
	itemLeftDelim
	itemLeftParen
	itemNumber
	itemPipe
	itemRawString
	itemRightDelim
	itemRightParen
	itemSpace
	itemString
	itemText
	itemVariable
	itemKeyword
	itemBlock
	itemDot
	itemDefine
	itemElse
	itemEnd
	itemIf
	itemNil
	itemRange
	itemTemplate
	itemWith
	itemFor
	itemDoublePlus
	itemDoubleMinus
	itemMinus
	itemPlus
	itemNewLine
	// header types
	itemPackage
	itemPackageValue
	itemImport
	itemImportValue
	itemFunctionDefine
	itemNotEqual
	itemFunctionName
	itemVariableType
	itemMap
	itemVar
	itemColon
	itemByteType
	itemStringType
	itemIntType
	itemUnknownToken
	itemSemiColon
	itemComment
	itemNode
)

const eof = -1

func (l *lexer) next() rune {
	if int(l.pos) >= len(l.input) {
		l.width = 0

		return eof
	}

	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = Pos(w)
	l.pos += l.width
	l.currentPosOnLine += l.width

	if r == '\n' {
		l.line++
	}

	return r
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()

	return r
}

func (l *lexer) backup() {
	l.pos -= l.width
	l.currentPosOnLine -= l.width

	if l.width == 1 && l.input[l.pos] == '\n' {
		l.line--
	}
}

func (l *lexer) emit(t itemType) {
	l.items <- item{t, l.currentStartOnLine, l.input[l.start:l.pos], l.startLine}
	l.start = l.pos
	l.currentStartOnLine = l.currentPosOnLine
	l.startLine = l.line
}

func (l *lexer) accept(valid string) bool {
	if strings.ContainsRune(valid, l.next()) {
		return true
	}

	l.backup()

	return false
}

func (l *lexer) acceptRun(valid string) {
	for strings.ContainsRune(valid, l.next()) {
	}

	l.backup()
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- item{itemError, l.currentStartOnLine, fmt.Sprintf(format, args...), l.startLine}

	return nil
}

func lex(input string) *lexer {
	l := &lexer{
		input:          input,
		items:          make(chan item),
		line:           1,
		startLine:      1,
		currentPosOnLine: 1,
		previousUnknown : false,
	}

	go l.run()

	return l
}

func (l *lexer) run() {
	for state := lexAction; state != nil; {
		state = state(l)
	}

	close(l.items)
}

const (
	leftDelim    = "{"
	rightDelim   = "}"
)

func lexAction(l *lexer) stateFn {
	switch r := l.next(); {
	  case r == '\n':
			l.currentPosOnLine = Pos(0)

			return lexWithUnknownConditionAndDoubleArguments(l, lexDefaultToken, itemNewLine)
		case r == eof:
			return nil;
		case isSpace(r):
			l.backup()

			return lexWithUnknownCondition(l, lexSpace)
		case r == '=':
			return lexWithUnknownConditionAndDoubleArguments(l, lexDefaultToken, itemAssign)
		case r == ':':
			return lexWithUnknownCondition(l, lexDeclareOrColon)
		case r == '|':
			return lexWithUnknownConditionAndDoubleArguments(l, lexDefaultToken, itemPipe)
		case r == '"':
			return lexWithUnknownCondition(l, lexQuote)
		case r == '`':
			return lexWithUnknownCondition(l, lexRawQuote)
		case r == '$':
			return lexWithUnknownCondition(l, lexVariable)
		case r == '/':
			nextRune := l.peek()

			if nextRune == '*' {
				return lexWithUnknownCondition(l, lexMultilineComment)
			}

			if nextRune == '/' {
				return lexWithUnknownCondition(l, lexOneLineComment)
			}
		case r == '\'':
			return lexWithUnknownCondition(l, lexChar)
		case r == '.':
			if l.pos < Pos(len(l.input)) {
				r := l.input[l.pos]

				if r < '0' || '9' < r {
					return lexWithUnknownCondition(l, lexField)
				}
			}

			fallthrough
	  case r == '+' || r == '-':
	    nextRune := l.peek()

	    if nextRune == '+' && r == '+' {
				return lexWithUnknownConditionAndDoubleArguments(l, lexDoubleSign, itemDoublePlus)
	    }

	    if nextRune == '-' && r == '-' {
				return lexWithUnknownConditionAndDoubleArguments(l, lexDoubleSign, itemDoubleMinus)
	    }

	    if nextRune == ' ' && r == '+' {
				return lexWithUnknownConditionAndDoubleArguments(l, lexDefaultToken, itemPlus)
	    }

			if nextRune == ' ' && r == '-' {
				return lexWithUnknownConditionAndDoubleArguments(l, lexDefaultToken, itemMinus)
			}
	  case r == '{':
			return lexWithUnknownConditionAndDoubleArguments(l, lexDefaultToken, itemLeftDelim)
	  case r == '}':
			return lexWithUnknownConditionAndDoubleArguments(l, lexDefaultToken, itemRightDelim)
		case r == '!':
			if l.previousUnknown {
				return lexUnknownToken
			}

			subRune := l.peek()

			if subRune == '=' {
				l.pos += Pos(1)
				l.emit(itemNotEqual)

				return lexAction
			}
		case ('0' <= r && r <= '9'):
			return lexWithUnknownCondition(l, lexNumber)
		case isAlphaNumeric(r):
			return lexWithUnknownCondition(l, lexIdentifier)
		case r == '(':
			return lexWithUnknownConditionAndDoubleArguments(l, lexDefaultToken, itemLeftParen)
		case r == ')':
			return lexWithUnknownConditionAndDoubleArguments(l, lexDefaultToken, itemRightParen)
		case r == ';':
			return lexWithUnknownConditionAndDoubleArguments(l, lexDefaultToken, itemSemiColon)
		case r == ',' || r == '[' || r == ']' || r == '<' || r == '>':
			return lexWithUnknownConditionAndDoubleArguments(l, lexDefaultToken, itemChar)
		default:
			l.previousUnknown = true
	}

	return lexAction
}

func lexWithUnknownCondition(l *lexer, lexState stateFn) stateFn {
	if l.previousUnknown {
		return lexUnknownToken
	}

	return lexState
}

func lexWithUnknownConditionAndDoubleArguments(l *lexer, lexState doubleStateFn, expectedType itemType) stateFn {
	if l.previousUnknown {
		return lexUnknownToken
	}

	return lexState(l, expectedType)
}

func lexMultilineComment(l *lexer) stateFn {
	if n := strings.Index(string(l.input[l.start:]), "*/"); n != -1 {
		l.pos = l.start + Pos(n + 2)

		if countNewLines := strings.Count(string(l.input[l.start:l.pos]), "\n"); n > 0 {
			l.line +=  countNewLines
			l.currentPosOnLine = Pos(len(l.input[l.start: l.pos]) - strings.LastIndex(string(l.input[l.start: l.pos]), "\n"))
		}

		l.emit(itemComment)
	}

	return lexAction
}

func lexOneLineComment(l *lexer) stateFn {
	if n := strings.Index(string(l.input[l.start:]), "\n"); n != -1 {
		l.pos = l.start + Pos(n)
		l.currentPosOnLine = Pos(len(l.input[l.start:l.pos]))

		l.emit(itemComment)
	}

	return lexAction
}

func lexUnknownToken(l *lexer) stateFn {
	if l.start - (l.pos - Pos(1)) != Pos(0) {
		l.pos -= Pos(1)
	}

	l.emit(itemUnknownToken)
	l.previousUnknown = false

	return lexAction
}

func lexDefaultToken(l *lexer, expectedType itemType) stateFn {
	l.emit(expectedType)

	return lexAction
}

func lexDeclareOrColon(l *lexer) stateFn {
	if l.peek() == '=' {
		l.next()
		l.emit(itemDeclare)
	} else {
		l.emit(itemColon)
	}

	return lexAction
}

func lexSpace(l *lexer) stateFn {
	var r rune
	var numSpaces int

	for {
		r = l.peek()

		if !isSpace(r) {
			break
		}

		l.next()
		numSpaces++
	}

	l.emit(itemSpace)

	return lexAction
}

func lexIdentifier(l *lexer) stateFn {
	l.backup()

	Loop:
		for {
			switch r := l.next(); {
			case isAlphaNumeric(r):
				// absorb.
			default:
				l.backup()
				word := l.input[l.start:l.pos]

				if !l.atTerminator() {
					return l.errorf("bad character %s %#U", word, r)
				}

				switch {
				case keyWords[word] > itemKeyword:
					l.emit(keyWords[word])

	        if keyWords[word] == itemPackage {
						l.currentPosOnLine -= Pos(1)

	          return lexPackageValue
	        }

	        if keyWords[word] == itemFunctionDefine {
	          return lexFunctionDefine
	        }
				case l.peek() == '(':
	        l.emit(itemFunction)
				case word == "true", word == "false":
					l.emit(itemBool)
				default:
					l.emit(itemIdentifier)
				}

				break Loop
			}
		}

	return lexAction
}

func lexFunctionDefine(l *lexer) stateFn {
  for {
    switch r := l.next(); {
      case isSpace(r):
        l.emit(itemSpace)

        break
      case isAlphaNumeric(r):
      	// absorb.
      default:
      	l.backup()
      	word := l.input[l.start:l.pos]

      	if !l.atTerminator() {
      		return l.errorf("bad character %s %#U", word, r)
      	}

        l.emit(itemFunctionName)

      	return lexAction
    }
  }
}

func lexPackageValue(l *lexer) stateFn {
	for {
		switch r := l.next(); {
			case isSpace(r):
				l.emit(itemSpace)

				break
			case isAlphaNumeric(r):
				// absorb.
			default:
				l.backup()
				word := l.input[l.start:l.pos]

				if !l.atTerminator() {
					return l.errorf("bad character %s %#U", word, r)
				}

				l.emit(itemPackageValue)

				return lexAction
		}
	}

	return l.errorf("expected package value")
}

func lexField(l *lexer) stateFn {
	return lexFieldOrVariable(l, itemField)
}

func lexVariable(l *lexer) stateFn {
	if l.atTerminator() {
		l.emit(itemVariable)
		return lexAction
	}

	return lexFieldOrVariable(l, itemVariable)
}

func lexFieldOrVariable(l *lexer, expectedType itemType) stateFn {
	if l.atTerminator() {
		if expectedType == itemVariable {
			l.emit(itemVariable)
		} else {
			l.emit(itemDot)
		}

		return lexAction
	}

	var r rune

	for {
		r = l.next()

		if !isAlphaNumeric(r) {
			l.backup()
			break
		}
	}

	if !l.atTerminator() {
		return l.errorf("bad character %#U", r)
	}

  if expectedType == itemField && l.peek() == '(' {
    l.emit(itemFunction)
  } else {
    l.emit(expectedType)
  }

	return lexAction
}

func (l *lexer) atTerminator() bool {
	r := l.peek()

	if isSpace(r) || isEndOfLine(r) || isPlus(r) {
		return true
	}

	switch r {
  case eof, '.', ',', '|', ':', ')', '(', '[', ']', '{', '}':
    	return true
	}

	return false
}

func lexChar(l *lexer) stateFn {
	Loop:
		for {
			switch l.next() {
			case '\\':
				if r := l.next(); r != eof && r != '\n' {
					break
				}

				fallthrough
			case eof, '\n':
				return l.errorf("unterminated character constant")
			case '\'':
				break Loop
			}
		}

	l.emit(itemCharConstant)

	return lexAction
}

func lexNumber(l *lexer) stateFn {
	l.backup()

	if !l.scanNumber() {
		return l.errorf("bad number syntax: %q", l.input[l.start:l.pos])
	}

	l.emit(itemNumber)

	return lexAction
}

func lexDoubleSign(l *lexer, expectedType itemType) stateFn {
  l.pos += Pos(1)
  l.currentPosOnLine += Pos(1)
  l.emit(expectedType)
  l.next()

  return lexAction
}

func (l *lexer) scanNumber() bool {
	digits := "0123456789_"
	if l.accept("0") {
		if l.accept("xX") {
			digits = "0123456789abcdefABCDEF_"
		} else if l.accept("oO") {
			digits = "01234567_"
		} else if l.accept("bB") {
			digits = "01_"
		}
	}

	l.acceptRun(digits)

	if l.accept(".") {
		l.acceptRun(digits)
	}

	if len(digits) == 10+1 && l.accept("eE") {
		l.accept("+-")
		l.acceptRun("0123456789_")
	}

	if len(digits) == 16+6+1 && l.accept("pP") {
		l.accept("+-")
		l.acceptRun("0123456789_")
	}

	l.accept("i")

	if isAlphaNumeric(l.peek()) {
		l.next()

		return false
	}

	return true
}

func lexQuote(l *lexer) stateFn {
	Loop:
		for {
			switch l.next() {
			case '\\':
				if r := l.next(); r != eof && r != '\n' {
					break
				}

				fallthrough
			case eof, '\n':
				return l.errorf("unterminated quoted string")
			case '"':
				break Loop
			}
		}

	l.emit(itemString)

	return lexAction
}

func lexRawQuote(l *lexer) stateFn {
	Loop:
		for {
			switch l.next() {
			case eof:
				return l.errorf("unterminated raw quoted string")
			case '`':
				break Loop
			}
		}

	l.emit(itemRawString)

	return lexAction
}

func isSpace(r rune) bool {
	return r == ' ' || r == '\t'
}

func isEndOfLine(r rune) bool {
	return r == '\r' || r == '\n'
}

func isPlus(r rune) bool {
  return r == '-' || r == '+'
}

func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

func (l *lexer) backPeek() rune {
  l.backup()
  r, _ := utf8.DecodeRuneInString(l.input[l.pos:])
	l.next()

	return r
}
