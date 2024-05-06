// This package is simpler version of flag package in Go standard library.
//
// The differences from the standard package are as follows:
//   - Does not require an output stream of error messages. Immediately returns the offending argument in the event of an error.
//   - Parses arguments with a leading hyphen. If the argument requires a value, parse the next argument. Process any arguments given. Does not stop processing if there are arguments that are not covered. So, os.Args can be given as it is, starting from index 0.
//   - Args() can be used to retrieve arguments not subject to parsing.
//   - Double hyphens `--like_this` are not supported.
//   - The value specification by the equal sign `-like=this` is not supported.
//   - Supported types are string, int, and bool.
//   - The bool flag does not support value. If the flag is given, it is set to true. If not, it is set to false.
package flag

import (
	"errors"
	"fmt"
	"strings"
)

type FlagSet struct {
	args    []string
	flags   []flagDef
	remaing []string
	parsed  bool
}

type ParseError struct {
	// The flag that caused the error.
	Arg string

	// The error that occurred. It specifies in this package.
	Err error

	// The value that caused the error.
	Value string
}

// Error returns a string representation of the error.
func (e ParseError) Error() string {
	return fmt.Sprintf("Err:%s\tFlag:%s\tValue:%s", e.Err, e.Arg, e.Value)
}

var (
	// ErrUnknownFlag is returned when an unknown flag is given.
	ErrUnknownFlag = errors.New("Unknown flag")

	// ErrDoubleHyphen is returned when a double hyphen is given.
	ErrDoubleHyphen = errors.New("Double hyphen")

	// ErrMissingValue is returned when a value is missing.
	ErrMissingValue = errors.New("Missing value")

	// ErrInvalidValue is returned when an invalid value is given.
	ErrInvalidValue = errors.New("Invalid value")

	// ErrAlreadyParsed is returned when Parse is called twice.
	ErrAlreadyParsed = errors.New("Already parsed")
)

type flagDef interface {
	getName() string
	set(string) bool
	setDefault()
	wantValue() bool
	gotValue() bool
}

// NewFlagSet creates a new FlagSet.
func NewFlagSet(args []string) *FlagSet {
	return &FlagSet{
		args:    args,
		flags:   make([]flagDef, 0),
		remaing: make([]string, 0),
		parsed:  false,
	}
}

// Args returns the remaining arguments that are not subject to parsing.
func (f *FlagSet) Args() []string {
	return f.remaing
}

// Parse parses the arguments.
// If an error occurs, it returns a ParseError.
//
// ParseError.Err is one of the following:
//   - ErrUnknownFlag
//   - ErrDoubleHyphen
//   - ErrMissingValue
//   - ErrInvalidValue
//   - ErrAlreadyParsed
func (f *FlagSet) Parse() error {
	if f.parsed {
		return ParseError{"", ErrAlreadyParsed, ""}
	}

	for i := 0; i < len(f.args); i++ {
		arg := f.args[i]
		if !strings.HasPrefix(arg, "-") {
			f.remaing = append(f.remaing, arg)
			continue
		}

		if strings.HasPrefix(arg, "--") {
			return ParseError{arg, ErrDoubleHyphen, ""}
		}

		flagName := arg[1:]
		var flag flagDef
		for _, f := range f.flags {
			if f.getName() == flagName {
				flag = f
				break
			}
		}

		if flag == nil {
			return ParseError{arg, ErrUnknownFlag, ""}
		}

		if flag.wantValue() {
			if i == len(f.args)-1 {
				return ParseError{arg, ErrMissingValue, ""}
			}
			if !flag.set(f.args[i+1]) {
				return ParseError{arg, ErrInvalidValue, f.args[i+1]}
			}
			i++
		} else {
			flag.set("")
		}
	}

	for _, f := range f.flags {
		if !f.gotValue() {
			f.setDefault()
		}
	}

	f.parsed = true
	return nil
}
