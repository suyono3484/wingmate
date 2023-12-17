//go:build cgo && linux

package main

/*
#include<errno.h>
#include<string.h>
#include<sys/types.h>
#include<pwd.h>
#include<grp.h>

static uid_t getuid(const char* username) {
	struct passwd local, *rv;
	errno = 0;
	rv = getpwnam(username);
	if (errno != 0) {
		return 0;
	}

	memcpy(&local, rv, sizeof(struct passwd));
	return local.pw_uid;
}
static gid_t getgid(const char* groupname) {
	struct group local, *rv;
	errno = 0;
	rv = getgrnam(groupname);
	if (errno != 0) {
		return 0;
	}

	memcpy(&local, rv, sizeof(struct group));
	return local.gr_gid;
}
*/
import "C"

func getUid(user string) (uint64, error) {
	u, err := C.getuid(C.CString(user))
	if err != nil {
		return 0, err
	}
	return uint64(u), nil
}

func getGid(group string) (uint64, error) {
	g, err := C.getgid(C.CString(group))
	if err != nil {
		return 0, err
	}
	return uint64(g), nil
}
