# cowboyolith

### getting started

#### local ssl

install https://github.com/FiloSottile/mkcer

Run the following commands:

```shell
mkcert -install
mkcert localhost
```

This will by default create two files:
```shell
localhost.pem
localhost-key.pem
```

#### database

Install postgres (developed against PostgreSQL 16.3 but likely many versions work)

```sql
drop database if exists cowboyolith;
create database cowboyolith;
\c cowboyolith
```

Then copy and paste the setup.sql file to create the tables

Insert your own admin using this sql

```sql
insert into users (id, email, is_admin) values (gen_random_uuid(), 'youremailhere@gmail.com', true);
```

#### aws ses

Set up an email only user in AWS so that it can only hit SES.


### running

Create your own `env-dev.sh` file following the `env-example.sh` file for environmental variables.

Run the following commands

```shell
source env-dev.sh
go run web_server.go
```
