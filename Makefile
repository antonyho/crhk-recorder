BIN=crhkrecorder

.PHONY: dep dep-update dep-download build

all: dep-download build

dep:
	go mod init
	go mod tidy

dep-update:
	go test ./...
	go list -m all
	go mod tidy

dep-download:
	go mod download

clean:
	go clean
	if [ -f ${BIN} ]; then rm ${BIN}; fi

build:
	go build -o ${BIN}