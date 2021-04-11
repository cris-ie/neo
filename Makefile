.PHONY: build
build:
	docker build -t neo .

.PHONY: run
run: build
	docker run -it --rm neo 

.PHONY: deploy
deploy:
	helm dep update chart
	helm upgrade -n neo -i neo chart --atomic --wait --timeout 200s --create-namespace --debug
