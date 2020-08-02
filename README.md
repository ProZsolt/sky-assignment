# Sky Assignment

## Install

```
git clone git@github.com:ProZsolt/sky-assignment.git
cd sky-assignment
go install ./...
```

## Usage

Create a MySQL database

Create a database table with the following command:

```
CREATE TABLE metrics (
  timestamp INT NOT NULL,
  cpuLoad FLOAT NULL DEFAULT NULL,
  concurrency INT NULL DEFAULT NULL,
  PRIMARY KEY (timestamp)
)
```

Set the following environment variables (modify it if needed):

```
export SKY_DB_HOST=127.0.0.1:3306
export SKY_DB_USERNAME=sky
export SKY_DB_PASSWORD=password
export SKY_DB_DATABASE=sky
```

`skyingestion` command generates metric entries for the past 5 minutes

`skyapi` command will launch an API server

You can access the API via `http://localhost:8080/api?from=1500000000&to=1600000000` where `from` and `to` are unix timestamps

## Testing

Set the following environment variables for the integration test (modify it if needed):

```
export SKY_TEST_DB_HOST=127.0.0.1:3306
export SKY_TEST_DB_USERNAME=sky
export SKY_TEST_DB_PASSWORD=password
export SKY_TEST_DB_DATABASE=sky
```

Run tests:

```
go test ./...
```

After running the test you have to manually clean the database.

## Dockerized version

Run `docker-compose up -d` at the application's root directory

It will launch the service and generate data for the past 5 minutes

You can access the API via `http://localhost:8080/api?from=1500000000&to=1600000000` where `from` and `to` are unix timestamps

To stop the service run `docker-compose down -v`

