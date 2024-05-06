package flag

type stringFlag struct {
	name         string
	defaultValue string
	dest         *string
	seted        bool
}

// StringVar defines a string flag with specified name and default value.
func (f *FlagSet) StringVar(dest *string, name string, defaultValue string) {
	f.flags = append(f.flags, &stringFlag{name, defaultValue, dest, false})
}

func (f *stringFlag) getName() string {
	return f.name
}

func (f *stringFlag) set(value string) bool {
	*f.dest = value
	f.seted = true
	return true
}

func (f *stringFlag) setDefault() {
	*f.dest = f.defaultValue
}

func (f *stringFlag) wantValue() bool {
	return true
}

func (f *stringFlag) gotValue() bool {
	return f.seted
}
