CC=go

TARGET=bin
LIBSRC=$(wildcard **/*.go)
LIBDIRS=$(patsubst %/,%, $(sort $(dir $(LIBSRC))))
EXECSRC=$(wildcard examples/**/*.go)
EXECDIRS=$(sort $(dir $(EXECSRC)))
EXEC=$(patsubst examples/%/,$(TARGET)/%,$(EXECDIRS))

.PHONY: clean test coverage

default: CHECK $(TARGET) $(EXEC)

test: CHECK
	@ for dir in $(LIBDIRS); do cd $$dir ; go test ; cd .. ; done

coverage: CHECK $(TARGET)
	@ for dir in $(LIBDIRS); do cd $$dir ; go test -coverprofile=../$(TARGET)/$$dir.html ; cd .. ; done
	@ for dir in $(LIBDIRS); do cd $$dir ; go tool cover -html=../$(TARGET)/$$dir.html ; cd .. ; done

install:
	@ go install

clean:
	@-rm -rf $(TARGET)

$(TARGET):
	@mkdir $(TARGET)

$(TARGET)/%: $(EXECSRC) $(LIBSRC)
	@cd $(patsubst bin/%,examples/%,$@) && $(CC) build -o ../../$@

CHECK:
ifndef GOPATH
	$(error GOPATH is undefined)
endif
