all:
	go install hello

run:
	$GOPATH/bin/hello

clean:
	rm $GOPATH/bin/*