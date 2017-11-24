
build: clean
	# Building resources...
	go-bindata -o res/res.go -pkg res -prefix res/ res/...
	# Building application...
	go build -o flappy .

clean:
	rm res/res.go 2> /dev/null || true
	rm flappy 2> /dev/null || true

run: build
	./flappy
