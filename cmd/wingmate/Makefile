all:
	git describe > version.txt
	go build -v

clean:
	rm version.txt
	go clean -i -cache -testcache

install:
	install wingmate ${DESTDIR}/wingmate
	