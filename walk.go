package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func walk(path string) error {
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}
	walkDir(path, fs.FileInfoToDirEntry(info), 0)
	return nil
}

func walkFunc(path string, d fs.DirEntry) error {
	if *iflag && isDotfile(d.Name()) {
		return fs.SkipDir
	}
	if *lflag {
		printLongFormat(path, d)
	} else {
		fmt.Fprintln(stdout, path)
	}
	return nil
}

func isDotfile(name string) bool {
	if name == "." || name == ".." {
		return false
	}
	return strings.HasPrefix(name, ".")
}

func walkDir(path string, d fs.DirEntry, depth int) {
	if err := walkFunc(path, d); err != nil || !d.IsDir() {
		return
	}

	if *maxdepth >= 0 && depth+1 > *maxdepth {
		return
	}

	ents, err := os.ReadDir(path)
	if err != nil {
		log.Print(err)
		return
	}

	for _, ent := range ents {
		walkDir(filepath.Join(path, ent.Name()), ent, depth+1)
	}
}

func printLongFormat(path string, d fs.DirEntry) {
	info, err := d.Info()
	if err != nil {
		info = fakeInfo(0)
	}
	strtime := formatTime(info.ModTime())
	fmt.Fprintf(stdout, "%-13s %10d %s %s\n", info.Mode(), info.Size(), strtime, path)
}

var sixMonthsAgo = time.Now().AddDate(0, -6, 0)

func formatTime(t time.Time) string {
	if *uflag {
		return strconv.FormatInt(t.Unix(), 10)
	}

	if t.Before(sixMonthsAgo) {
		return t.Format("Jan _2  2006")
	}
	return t.Format("Jan _2 15:04")
}

type fakeInfo int

var _ fs.FileInfo = fakeInfo(0)

func (fakeInfo) Name() string       { return "?" }
func (fakeInfo) Size() int64        { return 0 }
func (fakeInfo) Mode() fs.FileMode  { return 0 }
func (fakeInfo) ModTime() time.Time { return time.Unix(0, 0) }
func (fakeInfo) IsDir() bool        { return false }
func (fakeInfo) Sys() any           { return nil }
