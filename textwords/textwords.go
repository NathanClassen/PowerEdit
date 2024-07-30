package textwords

import (
	"fmt"
	"os"
	"strings"
	"unicode"
)

/*
represents a word in a block of text
W is the word, the first letter of the word is at index-S
in the block of text and the last letter at index-E
*/
type WordLoc struct {
	W string
	s int
	e int
	lws string
	rws string
}

/*
an abstraction over a large block of text
allows for iterating over and performing add, delete, edit
operations on the words of the text as an array of strings, while
maintaining the proper location of the words in the text

T is the text, Ws is an array of the words of the text, mapped to
wordLocs, and Offset is used to preserve the actual word indexes in the text
through various modifications to the text
*/
type TextWords struct {
	t      string
	ws     []WordLoc
	offset int
}

/*
	Create TextWords from content of a file, presumably a plain-txt file
*/
func FromFile(filename string) (*TextWords, error) {
	// read file
	if b, err := os.ReadFile(filename); err != nil {
		return &TextWords{}, err
	} else {
		return new(string(b)), nil
	}
	
}

/*
	Create TextWords from content of a string
*/
func FromString(txt string) *TextWords {
	return new(txt)
}

func new(txt string) *TextWords {
	wls := parsewordLocs(txt)
	return &TextWords{
		txt,
		wls,
		0,
	}
}



func (tw *TextWords) Insert(w WordLoc, at int) {
	if at >= len(tw.ws) {
		tw.ws = append(tw.ws, w)
	} else if at == 0 {
		tw.ws = append([]WordLoc{w}, tw.ws...)
	} else {
		tw.ws[at-1].rws = ""
		tw.ws[at].lws = ""

		tw.ws = append(tw.ws[0:at], append([]WordLoc{w}, tw.ws[at:]...)...)
	}
}

func (tw *TextWords) Edit(at int, newwl WordLoc) {
	//	TODO: update offset?
	tw.ws[at] = newwl
}

func (tw *TextWords) Delete(at int) {
	if at == 0 {
		tw.ws = tw.ws[1:]
	} else if at == len(tw.ws)-1 {
		tw.ws = tw.ws[0:at]
	} else {
		lw := &tw.ws[at-1]
		rw := &tw.ws[at+1]
	
		if lw.rws != " " || rw.lws != " " {			//	if the word on left or right of deleted word has significant wsp against deleted word
			if lw.rws != " " && rw.lws != " " {		//	check if both have significant wsp
				rw.lws = lw.rws + rw.lws			//	if so, ensure it's all retained by setting as rw.lws; this is what will be applied at text generation
			} else if lw.rws != " " {				//	else if it's the lw that has significant space
				rw.lws = lw.rws						//	ensure it's retained by trading rw's insignificant space
			}
		}
	
		tw.ws = append(tw.ws[0:at], tw.ws[at+1:]...)
	}
}

func (tw *TextWords) GetWord(at int) WordLoc {
	return tw.ws[at]
}

func (tw *TextWords) SurroundingText(at int) WordLoc {
	return tw.ws[at]
}

func (tw *TextWords) Text() string {
	return tw.getText(0, len(tw.ws))
}



func (tw *TextWords) getText(from, size int) string {
	txt := strings.Builder{}

	for _, wloc := range tw.ws[from:from+size] {
		txt.WriteString(wloc.lws)
		txt.WriteString(wloc.W)
	}

	return txt.String()
}

func parsewordLocs(txt string) []WordLoc {
	wls := []WordLoc{}
	o	:= false

	w := strings.Builder{}
	wspc := strings.Builder{}

	inWord := false
	ws := 0
	we := 0

	wl := WordLoc{}

	for i, char := range txt {

		if unicode.IsSpace(char) {

			inWord = false
			wspc.WriteRune(char)

			if w.Len() > 0 {

				wl.W = w.String()
				wl.s = ws
				wl.e = we
				
				wls = append(wls, wl)
				w.Reset()
				wl = WordLoc{}
			}

		} else {
			if inWord {
				w.WriteRune(char)
				we = i
			} else {
				inWord = true
				if char == 'i' {
					fmt.Println("")
				}
				w.WriteRune(char)
				ws = i
				we = i

				wl.lws = wspc.String()

				if o {
					wls[len(wls)-1].rws = wspc.String()
				}

				o = true
				wspc.Reset()
			}
		}
	}

	if w.Len() > 0 {
		wl.W = w.String()
		wl.s = ws
		wl.e = we

		wls = append(wls, wl)
		w.Reset()
		wspc.Reset()
	}


	return wls
}
