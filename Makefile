VERSION=dev

build:
	CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags="-X github.com/almostmoore/gbquestion/vars.Varsion=$(VERSION)"

build_linux:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-X github.com/almostmoore/gbquestion/vars.Varsion=$(VERSION)"