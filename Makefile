

all: wingmate dummy oneshot spawner starter pidproxy

wingmate:
	$(MAKE) -C cmd/wingmate all

pidproxy:
	$(MAKE) -C cmd/pidproxy all

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
	$(MAKE) -C cmd/experiment/dummy clean
	$(MAKE) -C cmd/experiment/oneshot clean
	$(MAKE) -C cmd/experiment/spawner clean
	$(MAKE) -C cmd/experiment/starter clean
