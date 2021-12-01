NAME=gitlab-vars

build:
	go build -o ${NAME}

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${NAME}

win:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ${NAME}.exe

mac:
  CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ${NAME}

clean:
	rm -f ${NAME}
	rm -r ${NAME}.exe
