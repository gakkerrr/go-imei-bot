up:
	docker-compose up -d

build:
	docker-compose build

down:
	docker-compose down

reload:
	docker-compose down && \
	docker-compose up --build -d

restart:
	docker-compose down && \
	docker-compose up -d

launch:
	docker-compose up --build -d