# Critch Backend
**_Critch_**: is a desktop application that aims to increase work productivity of teams
through proper, simple, and efficient communication process.

## Installation
1. download the project to your machine
```bash
git clone https://github.com/critch-app/critch-backend.git
```

2. change your directory to the project folder
```bash
cd ./critch-backend
```

3. download the required dependencies
```bash
make get-pkgs
```

4. create `.env` file with the following variables to connect to a `PostgreSQL` database and generate tokens.
```
DB_HOST=
DB_PORT=
DB_USER=
DB_PASS=
DB_NAME=

JWT_KEY=
```

5. run migrations after you setup connection variables for the database
```bash
make migrate
```

6. explore `MakeFile` to test, build, or run the project. e.g.
```bash
make run
```