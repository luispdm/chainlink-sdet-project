go-test:
	go test -v

compile:
	go test -c -o test

binary:
	./test -test.v

build:
	docker build -t chainlink-sdet-project:latest .

run:
	docker run --rm --env WSS=${WSS} --env PARALLEL=${PARALLEL} \
		--name chainlink-sdet-project chainlink-sdet-project:latest
