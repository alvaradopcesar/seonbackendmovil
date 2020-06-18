
# ProductService2
Proyecto Mobile para Inkafarma lado backend
2018

Compilation
-----------
* go get -d -v .
* go build *.go

Compilacion nativa Linux si estas en otro SO
--------------------------------------------
GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -v -o ProductService2 *.go

Containers
-----------
Se genera 1 solo container con 3 archios de configuracion para los 3 ambientes

Creacion de container 
---------------------
./build.sh

Ejecucion del container Produccion
----------------------------------
docker run -e ENVIRONMENT=prd -p 8080:8082 ProductService2

Ejecucion del container QA
---------------------------
docker run -e ENVIRONMENT=qa -p 8080:8082 ProductService2

Ejecucion del container CI
---------------------------
docker run -e ENVIRONMENT=ci -p 8080:8082 ProductService2
