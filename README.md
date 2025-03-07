# Aggregator

  

Aggregator is a CLI tool for rss feed aggregation.

  

#### Requirements

  

##### Golang Toolchain

  

https://go.dev/dl/

  

##### PostgresSQL 17+

  

[https://www.postgresql.org/download/](https://www.postgresql.org/download/)

  

**Download and install both for your Operating System**


*During the install of PostgreSQL make sure to record the default password you will need it latter.*


#### For Linux Users

Run The fallowing commands.

##### PostgresSQL 17+

```
sudo apt update
sudo apt install postgresql postgresql-contrib
psql --version
sudo passwd postgres
```


##### Golang Toolchain 1.23+

One you have download  go%version%.linux-amd64.tar.gz https://go.dev/dl/

```
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.23.1.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
go version
```
    
  
  

### **Setup**

  
  

#### **PostgreSQL**

  

Use psql command line tool to connect to Postgres (in windows go to start and search for sql shell (psql))

  

Create database called gator

  

Command Line Is

  

>CREATE DATABASE gator;

  

Connect to gator.

  

>\c gator

  

Set DB user password 'postgres'

  

>ALTER USER postgres PASSWORD 'postgres';

  

If you are using defaults this is your database connection string.

  

>postgres://postgres:postgres@localhost:5432/gator?sslmode=disable

  

if you did not use defaults below is a string that you need to modify for your use case. put it in notepad you will need it later.

  

>postgres://%username%:%password%@%ipaddress%:%port%/%databasename%?sslmode=disable

  

Open the shell, navigate to the repository that you have cloned.

  

We need to configure the database. 
Install goose 
Enter the following into shell

  

>go install github.com/pressly/goose/v3/cmd/goose@latest

  

Once installed navigate to \sql\schema

  

Run this command

  

>goose postgres "%your database connection string" up

  

>goose postgres "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable" up

  

go back to root of the repository open and look at  .gatorconfig.json modify the "db_url" value if necessary. Copy .gatorconfig.json to home directly for windows its %userprofile%.

  

Now simply type

  

>go install

  
  
  
  

### Usage

  

1.      First register yourself as a user.

>       Aggregator register %name%

  

2.      Login as the user you just registered

>       Aggregator login %name%

  

3.      Add some RSS feeds.

>       Aggregator addfeed “name” “feed url”

  

4.      Aggregate the feeds: note that the feeds will be continually polled unit stopped “CTRL+C”. Polling interval valid arguments are 1s, 1m, 1h.

>       Aggregator agg 1s

  

5.      Browse your feed browse has an optional argument of number of feeds you’d like to see. It will default to 2 if no value is given.

>       Aggregator browse 5

  

Remember if you want to clear the database use the “Aggregator reset”, there is no going back after reset. The data in data base will be erased.

  
