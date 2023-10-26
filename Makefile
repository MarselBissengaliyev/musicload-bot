run:
	docker run --name=go-app -p 80:8080 go-app

build:
	docker build -t go-app .