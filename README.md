#Microservice returning last currency exchange rates / Микросервис, отдающий свежие курсы валют.

Zen test 004

March 30, 2018

 - Take all rates for all pairs from [https://wex.nz](https://wex.nz) API ticker.last / Забирать курсы всех пар с биржи wex.nz (тикер last).

 - Save rates / Хранить курсы.

 - Return via API call average rates for last 10 mins / API метод возвращает среднее значение курсов за 10 минут.

 - Put service into Docker / Завернуть все в докер. 



Technologies choice is up to implementer / Выбор технологий, способов хранения и схемы API – на усмотрение соискателя.

##Implementation

###Environment variables

`export GOPATH=pwd && echo $GOPATH`

###Dependencies injection

`cd src/rates/ && dep ensure && cd ../..`

###Build

  - locally `go install rates/main`


  - with docker `docker build -f Dockerfile . -t rates`

###Run

  - locally
  
  `docker run --name=postgres --rm -d -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=changeit -e POSTGRES_DB=main -v pgdata:/var/lib/postgresql/data postgres:9.6`
  
  `docker inspect postgres | grep \"IPAddress\":` to know postgres container IP
  
  `bin/main --scrap-url https://wex.nz/api/3 --scrap-intervall 2s --listen-addr :8080 --api-path /rates/v0 --storage postgres --dsn "host=<IP_address> port=5432 user=postgres password=changeit dbname=main sslmode=disable"`


  - with docker `docker-compose -p rates up -d`

###Stop

  - locally `docker run postgres:9.6`
  
  - with docker `docker-compose -p rates down`

###Test

####Unit tests for packages

``

`go test -race .src/rates/api && go test -race .src/rates/storage && go test -race ./src/rates/scrapper/`

####Manually

`curl -iv -X GET http://localhost:8080/rates/v0/avg`
