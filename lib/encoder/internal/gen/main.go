// Package main provides utilities for encoder.
package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/artpar/rclone/lib/encoder"
)

const (
	edgeLeft = iota
	edgeRight
)

type mapping struct {
	mask     encoder.MultiEncoder
	src, dst []rune
}
type stringPair struct {
	a, b string
}

const header = `// Code generated by ./internal/gen/main.go. DO NOT EDIT.

` + `//go:generate go run ./internal/gen/main.go

package encoder

`

var maskBits = []struct {
	mask encoder.MultiEncoder
	name string
}{
	{encoder.EncodeZero, "EncodeZero"},
	{encoder.EncodeSlash, "EncodeSlash"},
	{encoder.EncodeSingleQuote, "EncodeSingleQuote"},
	{encoder.EncodeBackQuote, "EncodeBackQuote"},
	{encoder.EncodeLtGt, "EncodeLtGt"},
	{encoder.EncodeSquareBracket, "EncodeSquareBracket"},
	{encoder.EncodeSemicolon, "EncodeSemicolon"},
	{encoder.EncodeDollar, "EncodeDollar"},
	{encoder.EncodeDoubleQuote, "EncodeDoubleQuote"},
	{encoder.EncodeColon, "EncodeColon"},
	{encoder.EncodeQuestion, "EncodeQuestion"},
	{encoder.EncodeAsterisk, "EncodeAsterisk"},
	{encoder.EncodePipe, "EncodePipe"},
	{encoder.EncodeHash, "EncodeHash"},
	{encoder.EncodePercent, "EncodePercent"},
	{encoder.EncodeBackSlash, "EncodeBackSlash"},
	{encoder.EncodeCrLf, "EncodeCrLf"},
	{encoder.EncodeDel, "EncodeDel"},
	{encoder.EncodeCtl, "EncodeCtl"},
	{encoder.EncodeLeftSpace, "EncodeLeftSpace"},
	{encoder.EncodeLeftPeriod, "EncodeLeftPeriod"},
	{encoder.EncodeLeftTilde, "EncodeLeftTilde"},
	{encoder.EncodeLeftCrLfHtVt, "EncodeLeftCrLfHtVt"},
	{encoder.EncodeRightSpace, "EncodeRightSpace"},
	{encoder.EncodeRightPeriod, "EncodeRightPeriod"},
	{encoder.EncodeRightCrLfHtVt, "EncodeRightCrLfHtVt"},
	{encoder.EncodeInvalidUtf8, "EncodeInvalidUtf8"},
	{encoder.EncodeDot, "EncodeDot"},
}

type edge struct {
	mask    encoder.MultiEncoder
	name    string
	edge    int
	orig    []rune
	replace []rune
}

var allEdges = []edge{
	{encoder.EncodeLeftSpace, "EncodeLeftSpace", edgeLeft, []rune{' '}, []rune{'␠'}},
	{encoder.EncodeLeftPeriod, "EncodeLeftPeriod", edgeLeft, []rune{'.'}, []rune{'．'}},
	{encoder.EncodeLeftTilde, "EncodeLeftTilde", edgeLeft, []rune{'~'}, []rune{'～'}},
	{encoder.EncodeLeftCrLfHtVt, "EncodeLeftCrLfHtVt", edgeLeft,
		[]rune{'\t', '\n', '\v', '\r'},
		[]rune{'␀' + '\t', '␀' + '\n', '␀' + '\v', '␀' + '\r'},
	},
	{encoder.EncodeRightSpace, "EncodeRightSpace", edgeRight, []rune{' '}, []rune{'␠'}},
	{encoder.EncodeRightPeriod, "EncodeRightPeriod", edgeRight, []rune{'.'}, []rune{'．'}},
	{encoder.EncodeRightCrLfHtVt, "EncodeRightCrLfHtVt", edgeRight,
		[]rune{'\t', '\n', '\v', '\r'},
		[]rune{'␀' + '\t', '␀' + '\n', '␀' + '\v', '␀' + '\r'},
	},
}

var allMappings = []mapping{{
	encoder.EncodeZero, []rune{
		0,
	}, []rune{
		'␀',
	}}, {
	encoder.EncodeSlash, []rune{
		'/',
	}, []rune{
		'／',
	}}, {
	encoder.EncodeLtGt, []rune{
		'<', '>',
	}, []rune{
		'＜', '＞',
	}}, {
	encoder.EncodeSquareBracket, []rune{
		'[', ']',
	}, []rune{
		'［', '］',
	}}, {
	encoder.EncodeSemicolon, []rune{
		';',
	}, []rune{
		'；',
	}}, {
	encoder.EncodeDoubleQuote, []rune{
		'"',
	}, []rune{
		'＂',
	}}, {
	encoder.EncodeSingleQuote, []rune{
		'\'',
	}, []rune{
		'＇',
	}}, {
	encoder.EncodeBackQuote, []rune{
		'`',
	}, []rune{
		'｀',
	}}, {
	encoder.EncodeDollar, []rune{
		'$',
	}, []rune{
		'＄',
	}}, {
	encoder.EncodeColon, []rune{
		':',
	}, []rune{
		'：',
	}}, {
	encoder.EncodeQuestion, []rune{
		'?',
	}, []rune{
		'？',
	}}, {
	encoder.EncodeAsterisk, []rune{
		'*',
	}, []rune{
		'＊',
	}}, {
	encoder.EncodePipe, []rune{
		'|',
	}, []rune{
		'｜',
	}}, {
	encoder.EncodeHash, []rune{
		'#',
	}, []rune{
		'＃',
	}}, {
	encoder.EncodePercent, []rune{
		'%',
	}, []rune{
		'％',
	}}, {
	encoder.EncodeSlash, []rune{
		'/',
	}, []rune{
		'／',
	}}, {
	encoder.EncodeBackSlash, []rune{
		'\\',
	}, []rune{
		'＼',
	}}, {
	encoder.EncodeCrLf, []rune{
		rune(0x0D), rune(0x0A),
	}, []rune{
		'␍', '␊',
	}}, {
	encoder.EncodeDel, []rune{
		0x7F,
	}, []rune{
		'␡',
	}}, {
	encoder.EncodeCtl,
	runeRange(0x01, 0x1F),
	runeRange('␁', '␟'),
}}

var (
	rng *rand.Rand

	printables          = runeRange(0x20, 0x7E)
	fullwidthPrintables = runeRange(0xFF00, 0xFF5E)
	encodables          = collectEncodables(allMappings)
	encoded             = collectEncoded(allMappings)
	greek               = runeRange(0x03B1, 0x03C9)
)

func main() {
	seed := flag.Int64("s", 42, "random seed")
	flag.Parse()
	rng = rand.New(rand.NewSource(*seed))

	fd, err := os.Create("encoder_cases_test.go")
	fatal(err, "Unable to open encoder_cases_test.go:")
	defer func() {
		fatal(fd.Close(), "Failed to close encoder_cases_test.go:")
	}()
	fatalW(fd.WriteString(header))("Failed to write header:")

	fatalW(fd.WriteString("var testCasesSingle = []testCase{\n\t"))("Write:")
	_i := 0
	i := func() (r int) {
		r, _i = _i, _i+1
		return
	}
	for _, m := range maskBits {
		if len(getMapping(m.mask).src) == 0 {
			continue
		}
		if _i != 0 {
			fatalW(fd.WriteString(" "))("Write:")
		}
		in, out := buildTestString(
			[]mapping{getMapping(m.mask)},                               // pick
			[]mapping{getMapping(0)},                                    // quote
			printables, fullwidthPrintables, encodables, encoded, greek) // fill
		fatalW(fmt.Fprintf(fd, `{ // %d
		mask: %s,
		in:   %s,
		out:  %s,
	},`, i(), m.name, strconv.Quote(in), strconv.Quote(out)))("Error writing test case:")
	}
	fatalW(fd.WriteString(`
}

var testCasesSingleEdge = []testCase{
	`))("Write:")
	_i = 0
	for _, e := range allEdges {
		for idx, orig := range e.orig {
			if _i != 0 {
				fatalW(fd.WriteString(" "))("Write:")
			}
			fatalW(fmt.Fprintf(fd, `{ // %d
		mask: %s,
		in:   %s,
		out:  %s,
	},`, i(), e.name, strconv.Quote(string(orig)), strconv.Quote(string(e.replace[idx]))))("Error writing test case:")
		}
		for _, m := range maskBits {
			if len(getMapping(m.mask).src) == 0 || invalidMask(e.mask|m.mask) {
				continue
			}
			for idx, orig := range e.orig {
				replace := e.replace[idx]
				pairs := buildEdgeTestString(
					[]edge{e}, []mapping{getMapping(0), getMapping(m.mask)}, // quote
					[][]rune{printables, fullwidthPrintables, encodables, encoded, greek}, // fill
					func(rIn, rOut []rune, quoteOut []bool, testMappings []mapping) (out []stringPair) {
						testL := len(rIn)
						skipOrig := false
						for _, m := range testMappings {
							if runePos(orig, m.src) != -1 || runePos(orig, m.dst) != -1 {
								skipOrig = true
								break
							}
						}
						if !skipOrig {
							rIn[10], rOut[10], quoteOut[10] = orig, orig, false
						}

						out = append(out, stringPair{string(rIn), quotedToString(rOut, quoteOut)})
						for _, i := range []int{0, 1, testL - 2, testL - 1} {
							for _, j := range []int{1, testL - 2, testL - 1} {
								if j < i {
									continue
								}
								rIn := append([]rune{}, rIn...)
								rOut := append([]rune{}, rOut...)
								quoteOut := append([]bool{}, quoteOut...)

								for _, in := range []rune{orig, replace} {
									expect, quote := in, false
									if i == 0 && e.edge == edgeLeft ||
										i == testL-1 && e.edge == edgeRight {
										expect, quote = replace, in == replace
									}
									rIn[i], rOut[i], quoteOut[i] = in, expect, quote

									if i != j {
										for _, in := range []rune{orig, replace} {
											expect, quote = in, false
											if j == testL-1 && e.edge == edgeRight {
												expect, quote = replace, in == replace
											}
											rIn[j], rOut[j], quoteOut[j] = in, expect, quote
										}
									}
									out = append(out, stringPair{string(rIn), quotedToString(rOut, quoteOut)})
								}
							}
						}
						return
					})
				for _, p := range pairs {
					fatalW(fmt.Fprintf(fd, ` { // %d
		mask: %s | %s,
		in:   %s,
		out:  %s,
	},`, i(), m.name, e.name, strconv.Quote(p.a), strconv.Quote(p.b)))("Error writing test case:")
				}
			}
		}
	}
	fatalW(fmt.Fprintf(fd, ` { // %d
		mask: EncodeLeftSpace,
		in:   "  ",
		out:  "␠ ",
	}, { // %d
		mask: EncodeLeftPeriod,
		in:   "..",
		out:  "．.",
	}, { // %d
		mask: EncodeLeftTilde,
		in:   "~~",
		out:  "～~",
	}, { // %d
		mask: EncodeRightSpace,
		in:   "  ",
		out:  " ␠",
	}, { // %d
		mask: EncodeRightPeriod,
		in:   "..",
		out:  ".．",
	}, { // %d
		mask: EncodeLeftSpace | EncodeRightPeriod,
		in:   " .",
		out:  "␠．",
	}, { // %d
		mask: EncodeLeftSpace | EncodeRightSpace,
		in:   " ",
		out:  "␠",
	}, { // %d
		mask: EncodeLeftSpace | EncodeRightSpace,
		in:   "  ",
		out:  "␠␠",
	}, { // %d
		mask: EncodeLeftSpace | EncodeRightSpace,
		in:   "   ",
		out:  "␠ ␠",
	}, { // %d
		mask: EncodeLeftPeriod | EncodeRightPeriod,
		in:   "...",
		out:  "．.．",
	}, { // %d
		mask: EncodeRightPeriod | EncodeRightSpace,
		in:   "a. ",
		out:  "a.␠",
	}, { // %d
		mask: EncodeRightPeriod | EncodeRightSpace,
		in:   "a .",
		out:  "a ．",
	},
}

var testCasesDoubleEdge = []testCase{
	`, i(), i(), i(), i(), i(), i(), i(), i(), i(), i(), i(), i()))("Error writing test case:")
	_i = 0
	for _, e1 := range allEdges {
		for _, e2 := range allEdges {
			if e1.mask == e2.mask {
				continue
			}
			for _, m := range maskBits {
				if len(getMapping(m.mask).src) == 0 || invalidMask(m.mask|e1.mask|e2.mask) {
					continue
				}
				orig, replace := e1.orig[0], e1.replace[0]
				edges := []edge{e1, e2}
				pairs := buildEdgeTestString(
					edges, []mapping{getMapping(0), getMapping(m.mask)}, // quote
					[][]rune{printables, fullwidthPrintables, encodables, encoded, greek}, // fill
					func(rIn, rOut []rune, quoteOut []bool, testMappings []mapping) (out []stringPair) {
						testL := len(rIn)
						for _, i := range []int{0, testL - 1} {
							for _, secondOrig := range e2.orig {
								rIn := append([]rune{}, rIn...)
								rOut := append([]rune{}, rOut...)
								quoteOut := append([]bool{}, quoteOut...)

								rIn[1], rOut[1], quoteOut[1] = secondOrig, secondOrig, false
								rIn[testL-2], rOut[testL-2], quoteOut[testL-2] = secondOrig, secondOrig, false

								for _, in := range []rune{orig, replace} {
									rIn[i], rOut[i], quoteOut[i] = in, in, false
									fixEdges(rIn, rOut, quoteOut, edges)

									out = append(out, stringPair{string(rIn), quotedToString(rOut, quoteOut)})
								}
							}
						}
						return
					})

				for _, p := range pairs {
					if _i != 0 {
						fatalW(fd.WriteString(" "))("Write:")
					}
					fatalW(fmt.Fprintf(fd, `{ // %d
		mask: %s | %s | %s,
		in:   %s,
		out:  %s,
	},`, i(), m.name, e1.name, e2.name, strconv.Quote(p.a), strconv.Quote(p.b)))("Error writing test case:")
				}
			}
		}
	}
	fatalW(fmt.Fprint(fd, "\n}\n"))("Error writing test case:")
}

func fatal(err error, s ...interface{}) {
	if err != nil {
		log.Println(append(s, err))
	}
}
func fatalW(_ int, err error) func(...interface{}) {
	if err != nil {
		return func(s ...interface{}) {
			log.Println(append(s, err))
		}
	}
	return func(s ...interface{}) {}
}

func invalidMask(mask encoder.MultiEncoder) bool {
	return mask&(encoder.EncodeCtl|encoder.EncodeCrLf) != 0 && mask&(encoder.EncodeLeftCrLfHtVt|encoder.EncodeRightCrLfHtVt) != 0
}

// construct a slice containing the runes between (l)ow (inclusive) and (h)igh (inclusive)
func runeRange(l, h rune) []rune {
	if h < l {
		panic("invalid range")
	}
	out := make([]rune, h-l+1)
	for i := range out {
		out[i] = l + rune(i)
	}
	return out
}

func getMapping(mask encoder.MultiEncoder) mapping {
	for _, m := range allMappings {
		if m.mask == mask {
			return m
		}
	}
	return mapping{}
}
func collectEncodables(m []mapping) (out []rune) {
	for _, s := range m {
		out = append(out, s.src...)
	}
	return
}
func collectEncoded(m []mapping) (out []rune) {
	for _, s := range m {
		out = append(out, s.dst...)
	}
	return
}

func buildTestString(mappings, testMappings []mapping, fill ...[]rune) (string, string) {
	combinedMappings := append(mappings, testMappings...)
	var (
		rIn  []rune
		rOut []rune
	)
	for _, m := range mappings {
		if len(m.src) == 0 || len(m.src) != len(m.dst) {
			panic("invalid length")
		}
		rIn = append(rIn, m.src...)
		rOut = append(rOut, m.dst...)
	}
	inL := len(rIn)
	testL := inL * 3
	if testL < 30 {
		testL = 30
	}
	rIn = append(rIn, make([]rune, testL-inL)...)
	rOut = append(rOut, make([]rune, testL-inL)...)
	quoteOut := make([]bool, testL)
	set := func(i int, in, out rune, quote bool) {
		rIn[i] = in
		rOut[i] = out
		quoteOut[i] = quote
	}
	for i, r := range rOut[:inL] {
		set(inL+i, r, r, true)
	}

outer:
	for pos := inL * 2; pos < testL; pos++ {
		m := pos % len(fill)
		i := rng.Intn(len(fill[m]))
		r := fill[m][i]
		for _, m := range combinedMappings {
			if pSrc := runePos(r, m.src); pSrc != -1 {
				set(pos, r, m.dst[pSrc], false)
				continue outer
			} else if pDst := runePos(r, m.dst); pDst != -1 {
				set(pos, r, r, true)
				continue outer
			}
		}
		set(pos, r, r, false)
	}

	rng.Shuffle(testL, func(i, j int) {
		rIn[i], rIn[j] = rIn[j], rIn[i]
		rOut[i], rOut[j] = rOut[j], rOut[i]
		quoteOut[i], quoteOut[j] = quoteOut[j], quoteOut[i]
	})

	var bOut strings.Builder
	bOut.Grow(testL)
	for i, r := range rOut {
		if quoteOut[i] {
			bOut.WriteRune(encoder.QuoteRune)
		}
		bOut.WriteRune(r)
	}
	return string(rIn), bOut.String()
}

func buildEdgeTestString(edges []edge, testMappings []mapping, fill [][]rune,
	gen func(rIn, rOut []rune, quoteOut []bool, testMappings []mapping) []stringPair,
) []stringPair {
	testL := 30
	rIn := make([]rune, testL)      // test input string
	rOut := make([]rune, testL)     // test output string without quote runes
	quoteOut := make([]bool, testL) // if true insert quote rune before the output rune

	set := func(i int, in, out rune, quote bool) {
		rIn[i] = in
		rOut[i] = out
		quoteOut[i] = quote
	}

	// populate test strings with values from the `fill` set
outer:
	for pos := 0; pos < testL; pos++ {
		m := pos % len(fill)
		i := rng.Intn(len(fill[m]))
		r := fill[m][i]
		for _, m := range testMappings {
			if pSrc := runePos(r, m.src); pSrc != -1 {
				set(pos, r, m.dst[pSrc], false)
				continue outer
			} else if pDst := runePos(r, m.dst); pDst != -1 {
				set(pos, r, r, true)
				continue outer
			}
		}
		set(pos, r, r, false)
	}

	rng.Shuffle(testL, func(i, j int) {
		rIn[i], rIn[j] = rIn[j], rIn[i]
		rOut[i], rOut[j] = rOut[j], rOut[i]
		quoteOut[i], quoteOut[j] = quoteOut[j], quoteOut[i]
	})
	fixEdges(rIn, rOut, quoteOut, edges)
	return gen(rIn, rOut, quoteOut, testMappings)
}

func fixEdges(rIn, rOut []rune, quoteOut []bool, edges []edge) {
	testL := len(rIn)
	for _, e := range edges {
		for idx, o := range e.orig {
			r := e.replace[idx]
			if e.edge == edgeLeft && rIn[0] == o {
				rOut[0], quoteOut[0] = r, false
			} else if e.edge == edgeLeft && rIn[0] == r {
				quoteOut[0] = true
			} else if e.edge == edgeRight && rIn[testL-1] == o {
				rOut[testL-1], quoteOut[testL-1] = r, false
			} else if e.edge == edgeRight && rIn[testL-1] == r {
				quoteOut[testL-1] = true
			}
		}
	}
}

func runePos(r rune, s []rune) int {
	for i, c := range s {
		if c == r {
			return i
		}
	}
	return -1
}

// quotedToString returns a string for the chars slice where an encoder.QuoteRune is
// inserted before a char[i] when quoted[i] is true.
func quotedToString(chars []rune, quoted []bool) string {
	var out strings.Builder
	out.Grow(len(chars))
	for i, r := range chars {
		if quoted[i] {
			out.WriteRune(encoder.QuoteRune)
		}
		out.WriteRune(r)
	}
	return out.String()
}
