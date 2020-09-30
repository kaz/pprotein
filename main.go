package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/google/pprof/driver"
)

type (
	FlagSet struct {
		*flag.FlagSet

		input     []string
		usageMsgs []string
	}
)

func NewFlagSet(input []string) *FlagSet {
	return &FlagSet{
		flag.NewFlagSet("", flag.ContinueOnError),
		input,
		[]string{},
	}
}

func (f *FlagSet) StringList(o, d, c string) *[]*string {
	return &[]*string{f.String(o, d, c)}
}

func (f *FlagSet) ExtraUsage() string {
	return strings.Join(f.usageMsgs, "\n")
}
func (f *FlagSet) AddExtraUsage(eu string) {
	f.usageMsgs = append(f.usageMsgs, eu)
}

func (f *FlagSet) Parse(usage func()) []string {
	f.Usage = usage
	f.FlagSet.Parse(f.input)
	args := f.Args()
	if len(args) == 0 {
		usage()
	}
	return args
}

func main() {
	flagSet := NewFlagSet([]string{"-http", "0:9000", "http://localhost/pprof"})
	if err := driver.PProf(&driver.Options{Flagset: flagSet}); err != nil {
		fmt.Printf("%#v\n", err)
	}
}
