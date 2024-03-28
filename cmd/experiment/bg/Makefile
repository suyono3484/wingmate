all:
	git describe > version.txt
	go build -v

clean:
	echo "dev" > version.txt
	go clean -i -cache -testcache

install:
	install bg ${DESTDIR}/wmbg