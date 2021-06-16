docker build --tag golang-docker-tutorial:test .
docker run -p 8080:8080 -v $(pwd):/app golang-docker-tutorial:test