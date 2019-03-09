[![Build Status](https://travis-ci.org/rebel-l/schema.svg?branch=master)](https://travis-ci.org/rebel-l/schema) 
[![codecov](https://codecov.io/gh/rebel-l/schema/branch/master/graph/badge.svg)](https://codecov.io/gh/rebel-l/schema)
[![Go Report Card](https://goreportcard.com/badge/github.com/rebel-l/schema)](https://goreportcard.com/report/github.com/rebel-l/schema)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

[![Release](https://img.shields.io/github/release/rebel-l/schema.svg?label=Release)](https://github.com/rebel-l/schema/releases)
[![GitHub issues](https://img.shields.io/github/issues/rebel-l/schema.svg)](https://github.com/rebel-l/schema/issues)
[![Documentation](https://godoc.org/github.com/rebel-l/schema?status.svg)](https://godoc.org/github.com/rebel-l/schema)

# Schema Package
This library written in [go](https://golang.org) helps you to manage your database schema. It executes SQL scripts by a 
given folder. The following operations are provided:
- **upgrade**: applies new scripts which has not been executed successfully yet or executes everything from scratch.
- **revert**: reverts latest script. If you execute more than once, it takes the second latest, third latest and so on.
- **recreate**: resets the database by reverting all executed scripts and recreates it from scratch by using upgrade.

It requires Go 1.11 or higher. Lower versions might work but weren't tested.

## Supported Databases
Every SQL database is supported which has a driver for [go](https://golang.org) and which is compatible with the built in
package `"database/sql"` or the package `github.com/jmoiron/sqlx`. Originally it was developed for the following databases 
so far:
- sqlite3

## Write SQL Schema Script
Each script must have at least an `up` and a `down` command represented by the following SQL comments: `-- up` / `-- down`.
You can skip the down command but beware that you `revert` and `recreate` are not working. As an example a schema script
can look like the following (name: _001_example.sql_):

```sql
-- up
CREATE TABLE IF NOT EXISTS example(id INTEGER);
CREATE TABLE IF NOT EXISTS something(id INTEGER);

-- down
DROP TABLE IF EXISTS example;
DROP TABLE IF EXISTS something;

```

Remember the following restrictions to be a valid schema script:
- file ending must be **.sql**
- all files need to be in the same folder, sub folders are not executed
- files are executed in ascending (descending for _revert_) order of their filenames. I recommend to prefix files with 
three or more digits (001, 002, 003, ...) or timestamps like [yyyymmdd] (20190224).

## Usage of the Library

### Install as Project Dependency
I recommend to include this library to your project with [golang dep](https://github.com/golang/dep). If you have _dep_
installed just do

```bash
dep ensure -add github.com/rebel-l/schema
```

Alternately you can get it by go directly by calling

```bash
go get -u github.com/rebel-l/schema
```

### Usage: Upgrade
The library makes the usage as simple as possible. There is only a method called _Execute_ which receives the follwoing 
parameters:
- **path**: the path to your SQL scripts
- **command**: the command can be `upgrade`, `revert` or `recreate`
- **application version**: if you want you can set this value to version of your application to which the not applied 
scripts belong to. If you want to skip this, use a blank string `""` 
 
To `upgrade` (includes also **creation**) your database you can do

```go
package main

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rebel-l/schema"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	db, err := sqlx.Open("sqlite3", "database.db")
	if err != nil {
		log.Fatal(err)
	}

	s := schema.New(db)
	if err = s.Execute("./path_to_your_scripts", schema.CommandUpgrade, "Application Version"); err != nil {
		log.Fatal(err)
	}
}
``` 

The interesting part happens in the last 4 lines where we get the schema struct from `schema.New` and then apply the scripts
with `s.Execute("./path_to_your_scripts", schema.CommandUpgrade, "Application Version")`. 

As logger you can use of course different ones, but it is expected that they follow the `logrus.FieldLogger` interface.
Instead of `sqlx` you can also the internal go `sql` or any other which follows the `store.DatabaseConnector` interface
delivered with this _package_. 

### Usage: Revert
Regarding the example from the chapter before to `revert` the latest changes is very similar

```go
package main

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rebel-l/schema"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	db, err := sqlx.Open("sqlite3", "database.db")
	if err != nil {
		log.Fatal(err)
	}

	s := schema.New(db)
	if err = s.Execute("./path_to_your_scripts", schema.CommandRevert, ""); err != nil {
		log.Fatal(err)
	}
}
```

The only line which has changed is `s.Execute("./path_to_your_scripts", schema.CommandRevert, "")` where we just use 
the command `schema.CommandRevert`.

### Usage: Recreate
As you can imagine from the examples above `recreate` the database is no big deal

```go
package main

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rebel-l/schema"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	db, err := sqlx.Open("sqlite3", "database.db")
	if err != nil {
		log.Fatal(err)
	}

	s := schema.New(db)
	if err = s.Execute("./path_to_your_scripts", schema.CommandRecreate, ""); err != nil {
		log.Fatal(err)
	}
}
```

Just exchange the command in the call of `Execute`.

### Usage with Progress Bar
Optional you can show a progress bar on the command line. All you need to do is calling the method `WithProgressBar()`
before executing anything

```go
package main

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rebel-l/schema"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	db, err := sqlx.Open("sqlite3", "database.db")
	if err != nil {
		log.Fatal(err)
	}

	s := schema.New(db)
	s.WithProgressBar()
	if err = s.Execute("./path_to_your_scripts", schema.CommandUpgrade, "Application Version"); err != nil {
		log.Fatal(err)
	}
}
``` 

# Contributing to this Package
You are welcome to contribute to this repository. Please ensure that you created an issue and push your changes in a
feature branch.

## Setup Local Environment
At first you need to clone this repository and have [go](https://golang.org) installed. To get all the other necessary
stuff just run 

````bash
./scripts/tools/setup.sh
````

NOTE: works also with Windows using `Git Bash`.

The script installs:
- dep (and all necessary go packages)
- gometalinter (including all linters)
- goconvey
- git hooks
