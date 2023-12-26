DESTDIR = /usr/local/bin

installs = install-dir
programs = wingmate pidproxy exec
ifdef TEST_BUILD
	programs += oneshot spawner starter dummy
	installs += install-test
endif

all: ${programs}

wingmate:
	$(MAKE) -C cmd/wingmate all

pidproxy:
	$(MAKE) -C cmd/pidproxy all

exec:
	$(MAKE) -C cmd/exec all

dummy:
	$(MAKE) -C cmd/experiment/dummy all

oneshot:
	$(MAKE) -C cmd/experiment/oneshot all

spawner:
	$(MAKE) -C cmd/experiment/spawner all

starter:
	$(MAKE) -C cmd/experiment/starter all

clean:
	$(MAKE) -C cmd/wingmate clean
	$(MAKE) -C cmd/pidproxy clean
	$(MAKE) -C cmd/exec clean
	$(MAKE) -C cmd/experiment/dummy clean
	$(MAKE) -C cmd/experiment/oneshot clean
	$(MAKE) -C cmd/experiment/spawner clean
	$(MAKE) -C cmd/experiment/starter clean

install: ${installs}
	$(MAKE) -C cmd/wingmate DESTDIR=${DESTDIR} install
	$(MAKE) -C cmd/pidproxy DESTDIR=${DESTDIR} install
	$(MAKE) -C cmd/exec DESTDIR=${DESTDIR} install

install-test:
	$(MAKE) -C cmd/experiment/dummy DESTDIR=${DESTDIR} install
	$(MAKE) -C cmd/experiment/oneshot DESTDIR=${DESTDIR} install
	$(MAKE) -C cmd/experiment/spawner DESTDIR=${DESTDIR} install
	$(MAKE) -C cmd/experiment/starter DESTDIR=${DESTDIR} install

install-dir:
	install -d ${DESTDIR}
