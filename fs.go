package main

import (
	"io/fs"
	"os/user"
	"time"
)

func ReadDir(name string) ([]DirEntry, error) {
	return readDir(name)
}

type DirEntry interface {
	Name() string
	IsDir() bool
	Info() (FileInfo, error)
}

type FileInfo interface {
	Name() string
	Size() int64
	Mode() fs.FileMode
	ModTime() time.Time
	IsDir() bool

	/*
		Device() string
	*/
	NLink() int64
	User() *user.User
	Group() *user.Group
}
