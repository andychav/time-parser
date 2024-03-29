package main

// Code generated by C:\Users\Andres\go\bin\peg.exe -inline grammar.peg DO NOT EDIT.

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

const endSymbol rune = 1114112

/* The rule types inferred from the grammar are below. */
type pegRule uint8

const (
	ruleUnknown pegRule = iota
	ruletime
	rulemilitary
	rulehour
	ruleminute
	rulehTens
	rulemTens
	ruledigit
	rulemeridian
	ruleEND
	rulePegText
	ruleAction0
	ruleAction1
	ruleAction2
	ruleAction3
)

var rul3s = [...]string{
	"Unknown",
	"time",
	"military",
	"hour",
	"minute",
	"hTens",
	"mTens",
	"digit",
	"meridian",
	"END",
	"PegText",
	"Action0",
	"Action1",
	"Action2",
	"Action3",
}

type token32 struct {
	pegRule
	begin, end uint32
}

func (t *token32) String() string {
	return fmt.Sprintf("\x1B[34m%v\x1B[m %v %v", rul3s[t.pegRule], t.begin, t.end)
}

type node32 struct {
	token32
	up, next *node32
}

func (node *node32) print(w io.Writer, pretty bool, buffer string) {
	var print func(node *node32, depth int)
	print = func(node *node32, depth int) {
		for node != nil {
			for c := 0; c < depth; c++ {
				fmt.Fprintf(w, " ")
			}
			rule := rul3s[node.pegRule]
			quote := strconv.Quote(string(([]rune(buffer)[node.begin:node.end])))
			if !pretty {
				fmt.Fprintf(w, "%v %v\n", rule, quote)
			} else {
				fmt.Fprintf(w, "\x1B[36m%v\x1B[m %v\n", rule, quote)
			}
			if node.up != nil {
				print(node.up, depth+1)
			}
			node = node.next
		}
	}
	print(node, 0)
}

func (node *node32) Print(w io.Writer, buffer string) {
	node.print(w, false, buffer)
}

func (node *node32) PrettyPrint(w io.Writer, buffer string) {
	node.print(w, true, buffer)
}

type tokens32 struct {
	tree []token32
}

func (t *tokens32) Trim(length uint32) {
	t.tree = t.tree[:length]
}

func (t *tokens32) Print() {
	for _, token := range t.tree {
		fmt.Println(token.String())
	}
}

func (t *tokens32) AST() *node32 {
	type element struct {
		node *node32
		down *element
	}
	tokens := t.Tokens()
	var stack *element
	for _, token := range tokens {
		if token.begin == token.end {
			continue
		}
		node := &node32{token32: token}
		for stack != nil && stack.node.begin >= token.begin && stack.node.end <= token.end {
			stack.node.next = node.up
			node.up = stack.node
			stack = stack.down
		}
		stack = &element{node: node, down: stack}
	}
	if stack != nil {
		return stack.node
	}
	return nil
}

func (t *tokens32) PrintSyntaxTree(buffer string) {
	t.AST().Print(os.Stdout, buffer)
}

func (t *tokens32) WriteSyntaxTree(w io.Writer, buffer string) {
	t.AST().Print(w, buffer)
}

func (t *tokens32) PrettyPrintSyntaxTree(buffer string) {
	t.AST().PrettyPrint(os.Stdout, buffer)
}

func (t *tokens32) Add(rule pegRule, begin, end, index uint32) {
	tree, i := t.tree, int(index)
	if i >= len(tree) {
		t.tree = append(tree, token32{pegRule: rule, begin: begin, end: end})
		return
	}
	tree[i] = token32{pegRule: rule, begin: begin, end: end}
}

func (t *tokens32) Tokens() []token32 {
	return t.tree
}

type parser struct {
	min      int
	hour     int
	afterMid int

	Buffer string
	buffer []rune
	rules  [15]func() bool
	parse  func(rule ...int) error
	reset  func()
	Pretty bool
	tokens32
}

func (p *parser) Parse(rule ...int) error {
	return p.parse(rule...)
}

func (p *parser) Reset() {
	p.reset()
}

type textPosition struct {
	line, symbol int
}

type textPositionMap map[int]textPosition

func translatePositions(buffer []rune, positions []int) textPositionMap {
	length, translations, j, line, symbol := len(positions), make(textPositionMap, len(positions)), 0, 1, 0
	sort.Ints(positions)

search:
	for i, c := range buffer {
		if c == '\n' {
			line, symbol = line+1, 0
		} else {
			symbol++
		}
		if i == positions[j] {
			translations[positions[j]] = textPosition{line, symbol}
			for j++; j < length; j++ {
				if i != positions[j] {
					continue search
				}
			}
			break search
		}
	}

	return translations
}

type parseError struct {
	p   *parser
	max token32
}

func (e *parseError) Error() string {
	tokens, err := []token32{e.max}, "\n"
	positions, p := make([]int, 2*len(tokens)), 0
	for _, token := range tokens {
		positions[p], p = int(token.begin), p+1
		positions[p], p = int(token.end), p+1
	}
	translations := translatePositions(e.p.buffer, positions)
	format := "parse error near %v (line %v symbol %v - line %v symbol %v):\n%v\n"
	if e.p.Pretty {
		format = "parse error near \x1B[34m%v\x1B[m (line %v symbol %v - line %v symbol %v):\n%v\n"
	}
	for _, token := range tokens {
		begin, end := int(token.begin), int(token.end)
		err += fmt.Sprintf(format,
			rul3s[token.pegRule],
			translations[begin].line, translations[begin].symbol,
			translations[end].line, translations[end].symbol,
			strconv.Quote(string(e.p.buffer[begin:end])))
	}

	return err
}

func (p *parser) PrintSyntaxTree() {
	if p.Pretty {
		p.tokens32.PrettyPrintSyntaxTree(p.Buffer)
	} else {
		p.tokens32.PrintSyntaxTree(p.Buffer)
	}
}

func (p *parser) WriteSyntaxTree(w io.Writer) {
	p.tokens32.WriteSyntaxTree(w, p.Buffer)
}

func (p *parser) SprintSyntaxTree() string {
	var bldr strings.Builder
	p.WriteSyntaxTree(&bldr)
	return bldr.String()
}

func (p *parser) Execute() {
	buffer, _buffer, text, begin, end := p.Buffer, p.buffer, "", 0, 0
	for _, token := range p.Tokens() {
		switch token.pegRule {

		case rulePegText:
			begin, end = int(token.begin), int(token.end)
			text = string(_buffer[begin:end])

		case ruleAction0:

			if p.hour != 12 {
				p.afterMid += p.hour * 60
			}
			p.afterMid += p.min
			fmt.Println(text, "is", p.afterMid, "minutes after midnight")

		case ruleAction1:

			i, _ := strconv.Atoi(text)
			p.hour = i

		case ruleAction2:

			i, _ := strconv.Atoi(text)
			p.min = i

		case ruleAction3:

			if text == "pm" {
				p.afterMid += 12 * 60
			}

		}
	}
	_, _, _, _, _ = buffer, _buffer, text, begin, end
}

func Pretty(pretty bool) func(*parser) error {
	return func(p *parser) error {
		p.Pretty = pretty
		return nil
	}
}

func Size(size int) func(*parser) error {
	return func(p *parser) error {
		p.tokens32 = tokens32{tree: make([]token32, 0, size)}
		return nil
	}
}
func (p *parser) Init(options ...func(*parser) error) error {
	var (
		max                  token32
		position, tokenIndex uint32
		buffer               []rune
	)
	for _, option := range options {
		err := option(p)
		if err != nil {
			return err
		}
	}
	p.reset = func() {
		max = token32{}
		position, tokenIndex = 0, 0

		p.buffer = []rune(p.Buffer)
		if len(p.buffer) == 0 || p.buffer[len(p.buffer)-1] != endSymbol {
			p.buffer = append(p.buffer, endSymbol)
		}
		buffer = p.buffer
	}
	p.reset()

	_rules := p.rules
	tree := p.tokens32
	p.parse = func(rule ...int) error {
		r := 1
		if len(rule) > 0 {
			r = rule[0]
		}
		matches := p.rules[r]()
		p.tokens32 = tree
		if matches {
			p.Trim(tokenIndex)
			return nil
		}
		return &parseError{p, max}
	}

	add := func(rule pegRule, begin uint32) {
		tree.Add(rule, begin, position, tokenIndex)
		tokenIndex++
		if begin != position && position > max.end {
			max = token32{rule, begin, position}
		}
	}

	matchDot := func() bool {
		if buffer[position] != endSymbol {
			position++
			return true
		}
		return false
	}

	/*matchChar := func(c byte) bool {
		if buffer[position] == c {
			position++
			return true
		}
		return false
	}*/

	/*matchRange := func(lower byte, upper byte) bool {
		if c := buffer[position]; c >= lower && c <= upper {
			position++
			return true
		}
		return false
	}*/

	_rules = [...]func() bool{
		nil,
		/* 0 time <- <(<((hour meridian) / (hour ':' minute meridian) / military)+> END Action0)> */
		func() bool {
			position0, tokenIndex0 := position, tokenIndex
			{
				position1 := position
				{
					position2 := position
					{
						position5, tokenIndex5 := position, tokenIndex
						if !_rules[rulehour]() {
							goto l6
						}
						if !_rules[rulemeridian]() {
							goto l6
						}
						goto l5
					l6:
						position, tokenIndex = position5, tokenIndex5
						if !_rules[rulehour]() {
							goto l7
						}
						if buffer[position] != rune(':') {
							goto l7
						}
						position++
						if !_rules[ruleminute]() {
							goto l7
						}
						if !_rules[rulemeridian]() {
							goto l7
						}
						goto l5
					l7:
						position, tokenIndex = position5, tokenIndex5
						{
							position8 := position
							if !_rules[rulehour]() {
								goto l0
							}
							if buffer[position] != rune(':') {
								goto l0
							}
							position++
							if !_rules[ruleminute]() {
								goto l0
							}
							add(rulemilitary, position8)
						}
					}
				l5:
				l3:
					{
						position4, tokenIndex4 := position, tokenIndex
						{
							position9, tokenIndex9 := position, tokenIndex
							if !_rules[rulehour]() {
								goto l10
							}
							if !_rules[rulemeridian]() {
								goto l10
							}
							goto l9
						l10:
							position, tokenIndex = position9, tokenIndex9
							if !_rules[rulehour]() {
								goto l11
							}
							if buffer[position] != rune(':') {
								goto l11
							}
							position++
							if !_rules[ruleminute]() {
								goto l11
							}
							if !_rules[rulemeridian]() {
								goto l11
							}
							goto l9
						l11:
							position, tokenIndex = position9, tokenIndex9
							{
								position12 := position
								if !_rules[rulehour]() {
									goto l4
								}
								if buffer[position] != rune(':') {
									goto l4
								}
								position++
								if !_rules[ruleminute]() {
									goto l4
								}
								add(rulemilitary, position12)
							}
						}
					l9:
						goto l3
					l4:
						position, tokenIndex = position4, tokenIndex4
					}
					add(rulePegText, position2)
				}
				{
					position13 := position
					{
						position14, tokenIndex14 := position, tokenIndex
						if !matchDot() {
							goto l14
						}
						goto l0
					l14:
						position, tokenIndex = position14, tokenIndex14
					}
					add(ruleEND, position13)
				}
				{
					add(ruleAction0, position)
				}
				add(ruletime, position1)
			}
			return true
		l0:
			position, tokenIndex = position0, tokenIndex0
			return false
		},
		/* 1 military <- <(hour ':' minute)> */
		nil,
		/* 2 hour <- <(<(('2' [0-3]) / (hTens digit) / digit)> Action1)> */
		func() bool {
			position17, tokenIndex17 := position, tokenIndex
			{
				position18 := position
				{
					position19 := position
					{
						position20, tokenIndex20 := position, tokenIndex
						if buffer[position] != rune('2') {
							goto l21
						}
						position++
						if c := buffer[position]; c < rune('0') || c > rune('3') {
							goto l21
						}
						position++
						goto l20
					l21:
						position, tokenIndex = position20, tokenIndex20
						{
							position23 := position
							{
								position24, tokenIndex24 := position, tokenIndex
								if buffer[position] != rune('0') {
									goto l25
								}
								position++
								goto l24
							l25:
								position, tokenIndex = position24, tokenIndex24
								if buffer[position] != rune('1') {
									goto l22
								}
								position++
							}
						l24:
							add(rulehTens, position23)
						}
						if !_rules[ruledigit]() {
							goto l22
						}
						goto l20
					l22:
						position, tokenIndex = position20, tokenIndex20
						if !_rules[ruledigit]() {
							goto l17
						}
					}
				l20:
					add(rulePegText, position19)
				}
				{
					add(ruleAction1, position)
				}
				add(rulehour, position18)
			}
			return true
		l17:
			position, tokenIndex = position17, tokenIndex17
			return false
		},
		/* 3 minute <- <(<(mTens digit)> Action2)> */
		func() bool {
			position27, tokenIndex27 := position, tokenIndex
			{
				position28 := position
				{
					position29 := position
					{
						position30 := position
						if c := buffer[position]; c < rune('0') || c > rune('5') {
							goto l27
						}
						position++
						add(rulemTens, position30)
					}
					if !_rules[ruledigit]() {
						goto l27
					}
					add(rulePegText, position29)
				}
				{
					add(ruleAction2, position)
				}
				add(ruleminute, position28)
			}
			return true
		l27:
			position, tokenIndex = position27, tokenIndex27
			return false
		},
		/* 4 hTens <- <('0' / '1')> */
		nil,
		/* 5 mTens <- <[0-5]> */
		nil,
		/* 6 digit <- <[0-9]> */
		func() bool {
			position34, tokenIndex34 := position, tokenIndex
			{
				position35 := position
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l34
				}
				position++
				add(ruledigit, position35)
			}
			return true
		l34:
			position, tokenIndex = position34, tokenIndex34
			return false
		},
		/* 7 meridian <- <(<((('a' / 'A') ('m' / 'M')) / (('p' / 'P') ('m' / 'M')))> Action3)> */
		func() bool {
			position36, tokenIndex36 := position, tokenIndex
			{
				position37 := position
				{
					position38 := position
					{
						position39, tokenIndex39 := position, tokenIndex
						{
							position41, tokenIndex41 := position, tokenIndex
							if buffer[position] != rune('a') {
								goto l42
							}
							position++
							goto l41
						l42:
							position, tokenIndex = position41, tokenIndex41
							if buffer[position] != rune('A') {
								goto l40
							}
							position++
						}
					l41:
						{
							position43, tokenIndex43 := position, tokenIndex
							if buffer[position] != rune('m') {
								goto l44
							}
							position++
							goto l43
						l44:
							position, tokenIndex = position43, tokenIndex43
							if buffer[position] != rune('M') {
								goto l40
							}
							position++
						}
					l43:
						goto l39
					l40:
						position, tokenIndex = position39, tokenIndex39
						{
							position45, tokenIndex45 := position, tokenIndex
							if buffer[position] != rune('p') {
								goto l46
							}
							position++
							goto l45
						l46:
							position, tokenIndex = position45, tokenIndex45
							if buffer[position] != rune('P') {
								goto l36
							}
							position++
						}
					l45:
						{
							position47, tokenIndex47 := position, tokenIndex
							if buffer[position] != rune('m') {
								goto l48
							}
							position++
							goto l47
						l48:
							position, tokenIndex = position47, tokenIndex47
							if buffer[position] != rune('M') {
								goto l36
							}
							position++
						}
					l47:
					}
				l39:
					add(rulePegText, position38)
				}
				{
					add(ruleAction3, position)
				}
				add(rulemeridian, position37)
			}
			return true
		l36:
			position, tokenIndex = position36, tokenIndex36
			return false
		},
		/* 8 END <- <!.> */
		nil,
		nil,
		/* 11 Action0 <- <{
		    if p.hour != 12 {
		        p.afterMid += p.hour * 60
		    }
		    p.afterMid += p.min
		    fmt.Println(text, "is",p.afterMid, "minutes after midnight")
		}> */
		nil,
		/* 12 Action1 <- <{
		    i, _ := strconv.Atoi(text)
		       p.hour = i
		}> */
		nil,
		/* 13 Action2 <- <{
		    i, _ := strconv.Atoi(text)
		       p.min = i
		}> */
		nil,
		/* 14 Action3 <- <{
		    if text == "pm" {
		        p.afterMid += 12 * 60
		    }
		}> */
		nil,
	}
	p.rules = _rules
	return nil
}
