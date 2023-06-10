package main

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type formatFunc func(path string, info FileInfo) string

type Formatter struct {
	fs  []formatFunc
	Sep string
}

var formatMap = map[string]formatFunc{
	"path": relPath,

	"dir":  dirName,
	"base": baseName,

	// "iso":  isoTime,
	"unix": unixTime,
	"time": lsTime,

	"size": size,
	// "unit": unitSize,
}

func ParseFormat(s string) (*Formatter, error) {
	fs := make([]formatFunc, 0)
	for _, k := range strings.Split(s, ",") {
		f, ok := formatMap[k]
		if !ok {
			return nil, fmt.Errorf("parse format: unknown keyword: %s", k)
		}
		fs = append(fs, f)
	}
	return &Formatter{fs: fs}, nil
}

func (f *Formatter) Execute(path string, info FileInfo) string {
	list := make([]string, 0, len(f.fs))
	for _, fn := range f.fs {
		list = append(list, fn(path, info))
	}
	return strings.Join(list, f.Sep)
}

func relPath(path string, _ FileInfo) string {
	return path
}

func dirName(path string, _ FileInfo) string {
	return filepath.Dir(path)
}

func baseName(_ string, info FileInfo) string {
	return info.Name()
}

func unixTime(_ string, info FileInfo) string {
	return strconv.FormatInt(info.ModTime().Unix(), 10)
}

var sixMonthsAgo = time.Now().AddDate(0, -6, 0)

func lsTime(_ string, info FileInfo) string {
	t := info.ModTime()

	if t.Before(sixMonthsAgo) {
		return t.Format("Jan _2  2006")
	}
	return t.Format("Jan _2 15:04")
}

func size(_ string, info FileInfo) string {
	return strconv.FormatInt(info.Size(), 10)
}

func userName(_ string, info FileInfo) string {
	return info.User().Username
}

func groupName(_ string, info FileInfo) string {
	return info.Group().Name
}