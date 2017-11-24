
res: cleanRes
	# Building resources...
	go-bindata -o res/res.go -pkg res -prefix res/ res/...

cleanRes:
	rm res/res.go 2> /dev/null || true

clean: cleanRes
	rm objects 2> /dev/null || true

build: clean res
	# Building application...
	go build -o objects .

run: build
	./objects
