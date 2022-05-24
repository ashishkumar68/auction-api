# auction-api

This is a learning golang project. It is composed of the RESTFul APIs that can be integrated with any client app.

The Complete API documentation for this is: [here](https://github.com/ashishkumar68/auction-api/blob/main/APIs-doc.md)

To run this locally.

1. `git clone git@github.com:ashishkumar68/auction-api.git`(Clone the repository)
2. `cd /user/local/path/to/auction-api`(move into directory)
3. (assuming that docker and docker-compose are installed on your machine.) Run `docker-compose up -d` and wait for about 5 minutes for containers to be ready to use. 
4. To verify that containers are ready, inside the same directory, run: `docker-compose logs -f -t`
and once you see this similar log message, it means that containers are ready to be used.

```
auction-api-app    | 2022-05-24T06:46:40.548095512Z [GIN-debug] Listening and serving HTTP on :8081
```
5. To run unit tests, from within container:
**a.** Login into app container via `docker container exec -it auction-api-app bash`
**b.** To run unit tests: `sh bash-scripts/run-unit-tests.sh`
6. To run unit tests, from outside container, run: `docker container exec auction-api-app bash-scripts/run-unit-tests.sh`