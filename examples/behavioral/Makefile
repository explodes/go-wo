OUTFILE=game

res: cleanRes
	# Building resources...
	$(GOPATH)/bin/go-bindata -o res/res.go -pkg res -prefix res/ res/...

cleanRes:
	rm res/res.go 2> /dev/null || true

clean: cleanRes
	rm $(OUTFILE) 2> /dev/null || true

build: clean res
	# Building application...
	go build -o $(OUTFILE) .

run: build
	./$(OUTFILE)
