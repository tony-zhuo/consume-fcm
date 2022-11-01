# Off-site Test

## Setup
1. cp .env.example .env.
2. set variable in .env file, specially FIREBASE_PROJECT_ID.( if you want to use own key file )
3. add firebase service account key in root of project as name `test-firebase-account-key.json`
4. run migrate `make migrate` in project root
5. run `docker-compose up` in project root