all:
	go build server.go

run:
	./server

clean:
	rm server
