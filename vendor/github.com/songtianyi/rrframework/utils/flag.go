// Copyright (C) 2013-2015 by Maxim Bublis <b@codemonkey.ru>
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// Package uuid provides implementation of Universally Unique Identifier (UUID).
// Supported versions are 1, 3, 4 and 5 (as specified in RFC 4122) and
// version 2 (as specified in DCE 1.1).
package rrutils

import (
	"flag"
	"fmt"
	"strconv"
)

func check() {
	if !flag.Parsed() {
		flag.Parse()
	}
}

func FlagHelp() {
	flag.PrintDefaults()
}

func FlagDump() string {
	check()
	var ret string
	fn := func(f *flag.Flag) {
		ret += f.Name + "=" + f.Value.String() + "\n"
	}
	flag.Visit(fn)
	return ret
}

func FlagIsSet(name string) bool {
	check()
	ret := false
	fn := func(f *flag.Flag) {
		if f.Name == name {
			ret = true
			return
		}
	}
	flag.Visit(fn)
	return ret
}

func FlagGetInt(option string) (int, error) {
	check()
	if op := flag.Lookup(option); op != nil {
		v := op.Value.String()
		i, err := strconv.Atoi(v)
		if err != nil {
			return -1, err
		}
		return i, nil
	} else {
		return -1, fmt.Errorf("no option %s", option)
	}
}

func FlagGetString(option string) (string, error) {
	check()
	if op := flag.Lookup(option); op != nil {
		return op.Value.String(), nil
	} else {
		return "", fmt.Errorf("no option %s", option)
	}
}

func FlagGetBool(option string) (bool, error) {
	check()
	if op := flag.Lookup(option); op != nil {
		v := op.Value.String()
		if v == "true" {
			return true, nil
		}
		return false, nil
	} else {
		return false, fmt.Errorf("no option %s", option)
	}
}
