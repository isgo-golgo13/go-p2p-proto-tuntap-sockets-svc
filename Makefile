MAKE=make

DEP=dep ensure -v
all:
	+$(DEP)
	+$(MAKE) -C p2pc/
	+$(MAKE) -C p2ps/

clean:
	rm -f p2pc/p2pc_svc
	rm -f p2ps/p2ps_svc