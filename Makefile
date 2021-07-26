.PHONY: serve build clean

serve: clean
	docker-compose run --rm mkdocs serve

build: clean
	docker-compose run --rm mkdocs build -c

generator:
	go run scripts/main.go

clean:
	@rm -rf site
