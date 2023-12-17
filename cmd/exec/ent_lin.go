//go:build cgo && linux

package main

import (
	"errors"
)

func getUid(user string) (uint64, error) {
	return 0, errors.New("not implemented")
}

func getGid(group string) (uint64, error) {
	return 0, errors.New("not implemented")
}
