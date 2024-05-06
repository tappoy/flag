package flag

import (
	"reflect"
	"strings"
	"testing"
)

var (
	s1 = "-111"
	i1 = -111
	b1 = false
	s2 = "-222"
	i2 = -222
	b2 = false
)

type stringTest struct {
	name         string
	dest         *string
	defaultValue string
	want         string
}

type intTest struct {
	name         string
	dest         *int
	defaultValue int
	want         int
}

type boolTest struct {
	name string
	dest *bool
	want bool
}

type setting struct {
	title       string
	args        []string
	stringFlags []stringTest
	intFlags    []intTest
	boolFlags   []boolTest
	wantArgs    []string
	wantError   ParseError
}

func split(s string) []string {
	return strings.Split(s, " ")
}

func makeFlagSet(s *setting) *FlagSet {
	fs := NewFlagSet(s.args)
	for _, f := range s.stringFlags {
		fs.StringVar(f.dest, f.name, f.defaultValue)
	}
	for _, f := range s.intFlags {
		fs.IntVar(f.dest, f.name, f.defaultValue)
	}
	for _, f := range s.boolFlags {
		fs.BoolVar(f.dest, f.name)
	}
	return fs
}

func checkResult(t *testing.T, s *setting, args []string, err error) {
	if s.wantError.Err != nil {
		if err == nil {
			t.Errorf("NG. want error: %v, got: nil", s.wantError)
		} else if err != nil && err.(ParseError).Err != s.wantError.Err {
			t.Errorf("NG. want error: %v, got: %v", s.wantError, err)
		}
		return
	}
	if err != nil {
		t.Errorf("NG. unexpected error: %v", err)
		return
	}
	if !reflect.DeepEqual(args, s.wantArgs) {
		t.Errorf("NG. args: got: %v, want: %v", args, s.wantArgs)
	}
	for _, f := range s.stringFlags {
		if *f.dest != f.want {
			t.Errorf("NG. name: %v, got: %v, want: %v", f.name, *f.dest, f.want)
		}
	}
	for _, f := range s.intFlags {
		if *f.dest != f.want {
			t.Errorf("NG. name: %v, got: %v, want: %v", f.name, *f.dest, f.want)
		}
	}
	for _, f := range s.boolFlags {
		if *f.dest != f.want {
			t.Errorf("NG. name: %v, got: %v, want: %v", f.name, *f.dest, f.want)
		}
	}
}

func doTest(t *testing.T, s *setting) {
	t.Logf("Test: %v", s.title)
	fs := makeFlagSet(s)
	err := fs.Parse()
	checkResult(t, s, fs.Args(), err)
}

func TestFlagSet(t *testing.T) {
	doTest(t, &setting{
		title: "normal 1",
		args:  split("command -boolflag -intflag 42 -stringflag answer args1 args2"),
		stringFlags: []stringTest{
			{"stringflag", &s1, "default", "answer"},
		},
		intFlags: []intTest{
			{"intflag", &i1, 0, 42},
		},
		boolFlags: []boolTest{
			{"boolflag", &b1, true},
		},
		wantArgs: split("command args1 args2"),
	})
	doTest(t, &setting{
		title:       "normal 2",
		args:        split("command args1 args2"),
		stringFlags: []stringTest{},
		intFlags:    []intTest{},
		boolFlags:   []boolTest{},
		wantArgs:    split("command args1 args2"),
	})
	doTest(t, &setting{
		title: "normal 3",
		args:  split("command args1 -boolflag args2 -intflag 42 -stringflag answer"),
		stringFlags: []stringTest{
			{"stringflag", &s1, "default", "answer"},
		},
		intFlags: []intTest{
			{"intflag", &i1, 0, 42},
		},
		boolFlags: []boolTest{
			{"boolflag", &b1, true},
		},
		wantArgs: split("command args1 args2"),
	})
	doTest(t, &setting{
		title: "default",
		args:  split("command args1 args2"),
		stringFlags: []stringTest{
			{"stringflag", &s1, "default", "default"},
		},
		intFlags: []intTest{
			{"intflag", &i1, 42, 42},
		},
		boolFlags: []boolTest{
			{"boolflag", &b1, false},
		},
		wantArgs: split("command args1 args2"),
	})
	doTest(t, &setting{
		title: "double hyphen",
		args:  split("command args1 --boolflag args2 -intflag 42 -stringflag answer"),
		stringFlags: []stringTest{
			{"stringflag", &s1, "default", "answer"},
		},
		intFlags: []intTest{
			{"intflag", &i1, 0, 42},
		},
		boolFlags: []boolTest{
			{"boolflag", &b1, true},
		},
		wantArgs:  split("command args1 args2"),
		wantError: ParseError{"--boolflag", ErrDoubleHyphen, ""},
	})
	doTest(t, &setting{
		title: "unknown flag",
		args:  split("command args1 -unknown --boolflag args2 -intflag 42 -stringflag answer"),
		stringFlags: []stringTest{
			{"stringflag", &s1, "default", "answer"},
		},
		intFlags: []intTest{
			{"intflag", &i1, 0, 42},
		},
		boolFlags: []boolTest{
			{"boolflag", &b1, true},
		},
		wantArgs:  split("command args1 args2"),
		wantError: ParseError{"-unknown", ErrUnknownFlag, ""},
	})
	doTest(t, &setting{
		title: "missing value: string",
		args:  split("command args1 args2 -intflag 42 -stringflag"),
		stringFlags: []stringTest{
			{"stringflag", &s1, "default", "answer"},
		},
		intFlags: []intTest{
			{"intflag", &i1, 0, 42},
		},
		boolFlags: []boolTest{
			{"boolflag", &b1, true},
		},
		wantArgs:  split("command args1 args2"),
		wantError: ParseError{"-stringflag", ErrMissingValue, ""},
	})
	doTest(t, &setting{
		title: "missing value: int",
		args:  split("command args1 args2 -stringflag answer -intflag"),
		stringFlags: []stringTest{
			{"stringflag", &s1, "default", "answer"},
		},
		intFlags: []intTest{
			{"intflag", &i1, 0, 42},
		},
		boolFlags: []boolTest{
			{"boolflag", &b1, true},
		},
		wantArgs:  split("command args1 args2"),
		wantError: ParseError{"-intflag", ErrMissingValue, ""},
	})
	doTest(t, &setting{
		title: "invalid value: int",
		args:  split("command args1 args2 -stringflag answer -intflag notint"),
		stringFlags: []stringTest{
			{"stringflag", &s1, "default", "answer"},
		},
		intFlags: []intTest{
			{"intflag", &i1, 0, 42},
		},
		boolFlags: []boolTest{
			{"boolflag", &b1, true},
		},
		wantArgs:  split("command args1 args2"),
		wantError: ParseError{"-intflag", ErrInvalidValue, "notint"},
	})
}

func TestAlreadyParsed(t *testing.T) {
	fs := NewFlagSet(split("command args1 args2"))
	err := fs.Parse()
	if err != nil {
		t.Errorf("NG. unexpected error: %v", err)
	}
	err = fs.Parse()
	if err == nil {
		t.Errorf("NG. want error: %v, got: nil", ErrAlreadyParsed)
	} else if err != nil {
		if err.(ParseError).Err != ErrAlreadyParsed {
			t.Errorf("NG. want error: %v, got: %v", ErrAlreadyParsed, err)
		}
		if err.Error() != "Err:Already parsed\tFlag:\tValue:" {
			t.Errorf("NG. want error string: Err:Already parsed\tFlag:\tValue:, got: %v", err.Error())
		}
	}
}
