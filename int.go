package flag

import (
	"strconv"
)

type intFlag struct {
	name         string
	defaultValue int
	dest         *int
	seted        bool
}

// IntVar defines a int flag with specified name and default value.
func (f *FlagSet) IntVar(dest *int, name string, defaultValue int) {
	f.flags = append(f.flags, &intFlag{name, defaultValue, dest, false})
}

func (f *intFlag) getName() string {
	return f.name
}

func (f *intFlag) set(value string) bool {
	v, err := strconv.Atoi(value)
	if err != nil {
		return false
	}
	*f.dest = v
	f.seted = true
	return true
}

func (f *intFlag) setDefault() {
	*f.dest = f.defaultValue
}

func (f *intFlag) wantValue() bool {
	return true
}

func (f *intFlag) gotValue() bool {
	return f.seted
}
