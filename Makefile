include $(GOROOT)/src/Make.$(GOARCH)

TARG=container/heteroset
GOFILES=\
	ll_rb_tree.go \

include $(GOROOT)/src/Make.pkg

