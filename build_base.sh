# GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -v -o imb *.go
# docker build -t inkafarma/imb .
docker build -t imb_base -f Dockerfile_base .
