.PHONY: build
build:
	docker build -t neo .

.PHONY: run
run: build
	docker run -it --rm neo 
