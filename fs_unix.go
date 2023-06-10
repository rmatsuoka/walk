//go:build unix

package main

import (
	"log"
	"os"
	"os/user"
	"strconv"
	"syscall"
)

func readDir(name string) ([]DirEntry, error) {
	osEnts, err := os.ReadDir(name)
	if err != nil {
		return nil, err
	}

	ents := make([]DirEntry, 0, len(osEnts))
	for _, e := range osEnts {
		ents = append(ents, &dirEntry{e})
	}
	return ents, nil
}

type dirEntry struct {
	os.DirEntry
}

var _ DirEntry = &dirEntry{}

func (d *dirEntry) Info() (FileInfo, error) {
	info, err := d.DirEntry.Info()
	if err != nil {
		return nil, err
	}
	return FileInfoFromOSPkg(info), nil
}

type fileInfo struct {
	os.FileInfo
	stat *syscall.Stat_t
}

var _ FileInfo = &fileInfo{}

func FileInfoFromOSPkg(info os.FileInfo) FileInfo {
	return &fileInfo{info, info.Sys().(*syscall.Stat_t)}
}

func (i *fileInfo) NLink() int64 {
	return int64(i.stat.Nlink)
}

func (i *fileInfo) User() *user.User {
	u, err := user.LookupId(strconv.FormatInt(int64(i.stat.Uid), 10))
	if err != nil {
		log.Print(err)
		return &user.User{}
	}
	return u
}

func (i *fileInfo) Group() *user.Group {
	g, err := user.LookupGroupId(strconv.FormatInt(int64(i.stat.Uid), 10))
	if err != nil {
		log.Print(err)
		return &user.Group{}
	}
	return g
}