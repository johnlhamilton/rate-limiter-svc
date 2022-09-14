# rate-limiter-svc
This is a simple POC of a rate-limiting microservice using redis.  The full design can be found here: 
[Distributed Rate Limiter Design Doc](https://docs.google.com/document/d/1Hz1ziW7_--5UtQnglY6eqRNfIcWrvpWAG1I_-0JOOf4/edit?usp=sharing)

## Running Locally
This POC was built to be run locally in a docker container.  A docker-compose file is included
that will bring up both the rate-limiter-svc container and a redis container.

### Requirements
* docker-compose
* Port 8080 available

### Steps:
Starting the service:
```shell script
docker-compose up -d
```

Viewing the logs:
```shell script
docker-compose logs rate-limiter-svc
```

Testing the service:
```shell script
for i in {1..30}; do curl -X POST http://localhost:8080/request/ip -H 'Content-Type: application/json' -d '{"key": "foo"}'; done
```
