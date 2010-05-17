# Copyright 2009 The Go Authors. All rights reserved.
# Copyright 2010 Peter Williams. All rights reserved.
# Use of this source code is governed by the new BSD license
# license that can be found in the LICENSE file.

# After editing the DIRS= list or adding imports to any Go files
# in any of those directories, run:
#
#	$GOROOT/src/pkg/deps.bash
#
# to rebuild the dependency information in Make.deps.

nullstring :=
space := $(nullstring) # a space at the end
ifndef GOBIN
QUOTED_HOME=$(subst $(space),\ ,$(HOME))
GOBIN=$(QUOTED_HOME)/bin
endif
QUOTED_GOBIN=$(subst $(space),\ ,$(GOBIN))

all: install

DIRS=\
	mudlark/tree/llrb_tree

NOTEST=

NOBENCH=\
	mudlark/tree/llrb_tree\

TEST=\
	$(filter-out $(NOTEST),$(DIRS))

BENCH=\
	$(filter-out $(NOBENCH),$(TEST))

clean.dirs: $(addsuffix .clean, $(DIRS))
install.dirs: $(addsuffix .install, $(DIRS))
nuke.dirs: $(addsuffix .nuke, $(DIRS))
test.dirs: $(addsuffix .test, $(TEST))
bench.dirs: $(addsuffix .bench, $(BENCH))

%.clean:
	+cd $* && $(QUOTED_GOBIN)/gomake clean

%.install:
	+cd $* && $(QUOTED_GOBIN)/gomake install

%.nuke:
	+cd $* && $(QUOTED_GOBIN)/gomake nuke

%.test:
	+cd $* && $(QUOTED_GOBIN)/gomake test

%.bench:
	+cd $* && $(QUOTED_GOBIN)/gomake bench

clean: clean.dirs

install: install.dirs

test:	test.dirs

bench:	bench.dirs ../../test/garbage.bench

nuke: nuke.dirs
	rm -rf "$(GOROOT)"/pkg/$(GOOS)_$(GOARCH)/mudlark

deps:
	$(GOROOT)/src/pkg/deps.bash

-include Make.deps
