# GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -v -o imb *.go
docker build -t productservice2 .
