# Off-site Test

## 1. Env Setup
1. cp .env.example .env
2. set variable in .env file, specially FIREBASE_PROJECT_ID.( if you want to use own key file )
3. add firebase service account key in root of project as name `test-firebase-account-key.json`( if you want to use own key file )
4. run migrate `make migrate` in project root
5. run `docker-compose up` in project root

## 2. Run docker-compose
```shell
make run
```

## 3. Migration
```shell
make migrate
```

## RabbitMQ
when server is running, you can open the url `http://localhost:15672/` in browser.