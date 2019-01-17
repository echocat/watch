package main

import (
	"fmt"
	"io"
	"os"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func mustFprintf(w io.Writer, format string, args ...interface{}) {
	mustFprint(w, fmt.Sprintf(format, args...))
}

func mustFprint(w io.Writer, args ...interface{}) {
	aa := make([][]byte, len(args))
	for i, arg := range args {
		aa[i] = []byte(fmt.Sprint(arg))
	}
	_, err := fmt.Fprint(w, args...)
	must(err)
}

func warn(pattern string, args ...interface{}) {
	mustFprint(os.Stderr, "[WARNING] "+fmt.Sprintf(pattern, args...))
}

func fatal(pattern string, args ...interface{}) {
	mustFprint(os.Stderr, "[FATAL]   "+fmt.Sprintf(pattern, args...))
	os.Exit(3)
}

func intSliceContains(in []int, what int) bool {
	for _, candidate := range in {
		if candidate == what {
			return true
		}
	}
	return false
}
