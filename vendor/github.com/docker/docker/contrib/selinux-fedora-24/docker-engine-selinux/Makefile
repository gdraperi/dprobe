TARGETS?=docker
MODULES?=$***REMOVED***TARGETS:=.pp.bz2***REMOVED***
SHAREDIR?=/usr/share

all: $***REMOVED***TARGETS:=.pp.bz2***REMOVED***

%.pp.bz2: %.pp
	@echo Compressing $^ -\> $@
	bzip2 -9 $^

%.pp: %.te
	make -f $***REMOVED***SHAREDIR***REMOVED***/selinux/devel/Makefile $@

clean:
	rm -f *~  *.tc *.pp *.pp.bz2
	rm -rf tmp *.tar.gz

man: install
	sepolicy manpage --domain $***REMOVED***TARGETS***REMOVED***_t

install:
	semodule -i $***REMOVED***TARGETS***REMOVED***

