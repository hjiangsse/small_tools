PREFIX=/usr/local
BINDIR=${PREFIX}/bin
DESTDIR=
BLDDIR = build
BLDFLAGS=
EXT=
ifeq (${GOOS},windows)
    EXT=.exe
endif

APPS = myarch
all: $(APPS)

$(BLDDIR)/myarch:        $(wildcard apps/myarch/*.go myarch/*.go)

$(BLDDIR)/%:
	@mkdir -p $(dir $@)
	go build ${BLDFLAGS} -o $@ ./apps/$*

$(APPS): %: $(BLDDIR)/%

clean:
	rm -fr $(BLDDIR)

.PHONY: install clean all
.PHONY: $(APPS)

install: $(APPS)
	install -m 755 -d ${DESTDIR}${BINDIR}
	for APP in $^ ; do install -m 755 ${BLDDIR}/$$APP ${DESTDIR}${BINDIR}/$$APP${EXT} ; done
