#
bindir:=./bin
target:=$(bindir)/go-http-server
buildopt:=-v -i
srcs:=$(shell find . -name '*.go') go.mod .go-version
port:=

default: bin

run: $(target)
	$(target) $(port)

bin: $(target)

$(target): $(srcs)
	go build $(buildopt) -o $(target) cmd/go-http-server/main.go
