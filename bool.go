package flag

type boolFlag struct {
	name  string
	dest  *bool
	seted bool
}

// BoolVar defines a bool flag with specified name.
func (f *FlagSet) BoolVar(dest *bool, name string) {
	f.flags = append(f.flags, &boolFlag{name, dest, false})
}

func (f *boolFlag) getName() string {
	return f.name
}

func (f *boolFlag) set(value string) bool {
	*f.dest = true
	f.seted = true
	return true
}

func (f *boolFlag) setDefault() {
	*f.dest = false
}

func (f *boolFlag) wantValue() bool {
	return false
}

func (f *boolFlag) gotValue() bool {
	return f.seted
}
