//go:build !(cgo && linux)

package main

func getUid(user string) (uint64, error) {
	panic("not implemented")
}

func getGid(group string) (uint64, error) {
	panic("not implemented")
}
