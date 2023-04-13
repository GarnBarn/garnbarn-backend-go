run:
	ENV=dev go run .

run-compose:
	docker-compose up -d

start-dev-php-admin:
	docker start phpmyadmin || docker run --name phpmyadmin -d --link mysql-dev:db -p 8080:80 phpmyadmin