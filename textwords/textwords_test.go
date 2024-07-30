package textwords

import (
	"fmt"
	"os"
	"strings"
	"testing"
)


func TestBasic(t *testing.T) {
	src := "Hello there you wonderful world"
	srcTextWords := FromString(src)

	if srcTextWords.GetWord(3).W != "wonderful" {
		t.Errorf("expected \"%s\" but got \"%s\"", "wonderful", srcTextWords.GetWord(3).W)
	}
}
func TestText(t *testing.T) {
	var tests = []struct {
		name string
		text string
	}{
		{"sentence with newlines","and?\n\nHow could you say that?\nReally, thats.. Pretty incredible."},
		{"single word","wow"},
		{"basic sentence","hello world you are looking round today"},
	}

	for _, tt := range tests {
		testname := tt.name
		t.Run(testname, func(t *testing.T) {
			txtWs := FromString(tt.text)
			res := txtWs.Text()

			if res != tt.text {
				t.Errorf("\ngot:\n   '%s'\nwant:\n   '%s'",res,tt.text)
			}
		})
	}
}
func TestTextFromFile(t *testing.T) {
	var tests = []struct {
		name string
		path string
	}{
		{
			"gutenberg-iliad.txt",
			"/Users/nathanclassen/go/src/power-edit/texts/gutenberg-iliad.txt",
		},
	}

	for _, tt := range tests {
		testname := tt.name
		t.Run(testname, func(t *testing.T) {
			txtWs, err := FromFile(tt.path)
			if err != nil {
				t.Errorf("\ngot error making TextWords from file:\n %v\n", err)
			}

			fileContents, err := os.ReadFile(tt.path)
			if err != nil {
				t.Errorf("\ngot error reading file:\n %v\n", err)
			}

			if string(fileContents) != txtWs.Text() {
				t.Error("file text did not match TextWords.FromFile.Text")
			}
		})
	}
}
func TestFromStringEndWord(t *testing.T) {
	var tests = []struct {
		name string
		text string
		want WordLoc
	}{
		{
			"sentence with newlines",
			"and?\n\nHow could you say that?\nReally, thats.. Pretty incredible.",
			WordLoc{W:"incredible.",lws: " ", rws: ""},
		},
		{
			"single word",
			"wow",
			WordLoc{W:"wow",lws: "", rws: ""},
		},
		{
			"basic sentence",
			"hello world you are looking round today",
			WordLoc{W:"today",lws: " ", rws: ""},
		},
	}

	for _, tt := range tests {
		testname := tt.name
		t.Run(testname, func(t *testing.T) {
			txtWs := FromString(tt.text)
			res := txtWs.GetWord(len(txtWs.ws)-1)

			if res.W != tt.want.W {
				t.Errorf("\ngot:\n'%s' word\nwant:\n'%s' word",res.W,tt.want.W)
			} else if res.lws != tt.want.lws {
				t.Errorf("\ngot:\n'%s' lws\nwant:\n'%s' lws\n", res.lws, tt.want.lws)
			} else if res.rws != tt.want.rws {
				t.Errorf("\ngot:\n'%s' rws\nwant:\n'%s' rws\n", res.rws, tt.want.rws)
			}
		})
	}
}
func TestAdd(t *testing.T) {
	var tests = []struct {
		wl WordLoc
		at int
		have string
		want string
	}{
		{
			wl: WordLoc{W:"\n\nHow "},
			at: 1,
			have: "and?\n\ncould you say that?",
			want: "and?\n\nHow could you say that?",
		},
		{
			wl: WordLoc{W:" so?\n\n"},
			at: 1,
			have: "and\n\nHow could you say that?",
			want: "and so?\n\nHow could you say that?",
		},
		{
			wl: WordLoc{W:" that?"},
			at: 4,
			have: "How could you say",
			want: "How could you say that?",
		},
		{
			wl: WordLoc{W:"and?\n\n"},
			at: 0,
			have: "How could you say that?",
			want: "and?\n\nHow could you say that?",
		},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("adding '%s' test", strings.Trim(tt.wl.W, " \n"))
		t.Run(testname, func(t *testing.T) {
			txtWs := FromString(tt.have)
			txtWs.Insert(tt.wl,tt.at)
			res := txtWs.Text()

			if res != tt.want {
				t.Errorf("\ngot:   %s\nwant:   %s",res,tt.want)
			}
		})
	}
}
func TestDelete(t *testing.T) {
	var tests = []struct {
		at int
		have string
		want string
	}{
		{
			at: 0,
			have: "and?\n\nHow could you say that?",
			want: "\n\nHow could you say that?",
		},
		{
			at: 1,
			have: "and?\n\nHow could you say that?",
			want: "and?\n\ncould you say that?",
		},
		{
			at: 2,
			have: "and?\n\nHow could you say that?",
			want: "and?\n\nHow you say that?",
		},
		{
			at: 3,
			have: "and?\n\nHow could you say that?",
			want: "and?\n\nHow could say that?",
		},
		{
			at: 4,
			have: "and?\n\nHow could you say that?",
			want: "and?\n\nHow could you that?",
		},
		{
			at: 5,
			have: "and?\n\nHow could you say that?",
			want: "and?\n\nHow could you say",
		},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("delete at %d", tt.at)
		t.Run(testname, func(t *testing.T) {
			txtWs := FromString(tt.have)
			txtWs.Delete(tt.at)
			res := txtWs.Text()

			if res != tt.want {
				t.Errorf("\ngot: %s\nwant: %s",res,tt.want)
			}
		})
	}
}
func TestSurroundingText(t *testing.T) {
	var tests = []struct {
		at int
		size int
		text string
		want string
	}{
		{
			at: 0,
			size: 2,
			text: "and?\n\nHow could you say that?",
			want: "*and?*\n\nHow could",
		},
		{
			at: 1,
			size: 2,
			text: "and?\n\nHow could you say that?",
			want: "and?\n\n*How* could you",
		},
		{
			at: 2,
			size: 2,
			text: "and?\n\nHow could you say that?",
			want: "and?\n\nHow *could* you say",
		},
		{
			at: 3,
			size: 2,
			text: "and?\n\nHow could you say that?",
			want: "\n\nHow could *you* say that?",
		},
		{
			at: 4,
			size: 2,
			text: "and?\n\nHow could you say that?",
			want: " could you *say* that?",
		},
		{
			at: 5,
			size: 2,
			text: "and?\n\nHow could you say that?",
			want: " you say *that?*",
		},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("surrounding %d at %d", tt.size, tt.at)
		t.Run(testname, func(t *testing.T) {
			txtWs := FromString(tt.text)
			res := txtWs.SurroundingText(tt.at,tt.size)

			if res != tt.want {
				t.Errorf("\ngot: %s\nwant: %s",res,tt.want)
			}
		})
	}
}