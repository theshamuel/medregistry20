# medregistry20

### DESCRIPTION
The backend API version 2.0 for [MedRegistry](https://github.com/theshamuel/medregistry) open-source project provided CRM system in a medical domain.


#### Database
The first start of db container import medregdb_starter.tar.gz dumb of database with collection and demo admin-user for getting JWT token and manage `users` (admin/admin). The database is ready deploy on production and has all necessary indexes for optimisation in queries. Also for container has scripts for backup on host instanse and to `AWS S3`. If you need this setup specal enviroment variables.
Important after first start change `MONGO_RESTORE` variable to false.

There is deploy one container of mongo `medregdb`.

#### API 1.0
The directory `api` has the first version of jar api (the version which will start into docker container is setting up as environment variable (VERSION) in medregapi-v1 service in docker-compose.yml)
Logs collects to volume which mapped into docker container
This docker container based on my own openjdk docker images version 15 for [building](https://hub.docker.com/repository/docker/theshamuel/baseing-java-build) and [launching](https://hub.docker.com/repository/docker/theshamuel/baseing-java-app).


#### API 2.0
The new [API 2.0](https://github.com/theshamuel/medregistry20) [`https://github.com/theshamuel/medregistry20`]

## DEPLOYMENT
1. Download any release
1. Define all necessary environment variables into docker-compose.yml.
   - MEDREG_API_V1_URL - URL to API 1.0
   - REPORT_PATH - Filepath to folder with report`s templates
1. Execute `docker-compose up --build`