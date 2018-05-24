build: clean ftmpl
	go build -v -i -o anki2dict ./...
	./anki2dict
test: ftmpl
	go test ./...
ftmpl:
	ftmpl -package=main -targetgo anki2dictionary/ftmpl_generated.go templates/ftmpl/
clean:
	rm -Rf output
	mkdir -p output
