
test:
	go test -cover

cover:
	go test -cover -coverprofile=c.out
	go tool cover -html=c.out -o coverage.html