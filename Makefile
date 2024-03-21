.PHONY: run clean kill

run:
	cd fixtures && docker-compose up -d && cd .. && go build -o collection	&& ./collection
clean:
	cd fixtures && docker-compose down -v && cd .. && rm -f collection
kill:
	kill -9 $$(lsof -t -i:8080)

