all:
	git describe > version.txt
	go build -v

clean:
	go clean -i -cache -testcache

install:
	install exec ${DESTDIR}/wmexec