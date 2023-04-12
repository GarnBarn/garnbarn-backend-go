run:
	ENV=dev go run .

run-compose:
	docker-compose up -d

start-dev-db-linux:
	docker start mysql-dev || docker run --name mysql-dev -d \
        -p 3306:3306 \
        -e MYSQL_ROOT_PASSWORD=P@ssw0rd \
        -e MYSQL_DATABASE=garnbarn \
        --restart unless-stopped \
        mysql:8.0.32

start-dev-php-admin:
	docker start phpmyadmin || docker run --name phpmyadmin -d --link mysql-dev:db -p 8080:80 phpmyadmin