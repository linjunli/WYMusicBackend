WYMusicBackend = ./bin/WYMusicBackend

all: $(WYMusicBackend)
$(WYMusicBackend) : $(shell find src/main.go ./bin/) 
	go vet src/main.go
	go build -o $@ src/main.go
clean : 
	@rm -rf bin/*
