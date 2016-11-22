# Snap publisher plugin - MySQL
This plugin publishes data to MySQL database. 

It's used in the [Snap framework](http://github.com:intelsdi-x/snap).

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Installation](#installation)
  * [Configuration and Usage](#configuration-and-usage)
2. [Documentation](#documentation)
  * [Task Manifest Config](#task-manifest-config)
  * [Database schema](#database-schema)
  * [Examples](#examples)
  * [Roadmap](#roadmap)
3. [Community Support](#community-support)
4. [Contributing](#contributing)
5. [License](#license-and-authors)
6. [Acknowledgements](#acknowledgements)

## Getting Started
### System Requirements
* [golang 1.6+](https://golang.org/dl/) (needed only for building)

### Operating systems
All OSs currently supported by Snap:
* Linux/amd64
* Darwin/amd64

### Installation
#### Download mysql plugin binary:
You can get the pre-built binaries for your OS and architecture at plugin's [GitHub Releases](https://github.com/intelsdi-x/snap-plugin-publisher-mysql/releases) page.

#### To build the plugin binary:
Fork https://github.com/intelsdi-x/snap-plugin-publisher-mysql

Clone repo into `$GOPATH/src/github.com/intelsdi-x/`:

```
$ git clone https://github.com/<yourGithubID>/snap-plugin-publisher-mysql.git
```

Build the plugin by running make within the cloned repo:
```
$ make
```
This builds the plugin in `./build`

### Configuration and Usage
* Set up the [Snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started)
* Load the plugin and create a task, see example in [Examples](#examples).

## Documentation
There are a number of other resources you can review to learn to use this plugin:

* [mysql](https://www.mysql.com/) 
* [go-mysql-driver](github.com/go-sql-driver/mysql)
* [Snap mysql integration test](https://github.com/intelsdi-x/snap-plugin-publisher-mysql/blob/master/mysql/mysql_integration_test.go)
* [Snap mysql unit test](https://github.com/intelsdi-x/snap-plugin-publisher-mysql/blob/master/mysql/mysql_test.go)

### Task Manifest Config

In task manifest, the config section of mysql publisher describes how to establish a connection to the MySQL server.

Name 	  	 | Data Type | Default       | Description
----------|-----------|---------------|-------------
hostname 	| string 	  | localhost     | the host of MySQL service
port 		   | string	 	 | 3306          | the port number of MySQL service
username  | string 	  | root          | the name of user
password 	| string 	  | root          | the password of user
database 	| string 	  | SNAP_TEST        | the name of database (use existed or create a new)
tablename | string 	  | info       | the name of table (use existed or create a new)

Each connection parameter has a default value, but it can be override by set a value in task manifest (see [exemplary task manifest](examples/tasks/psutil-mysql.json))

### Database schema

Metrics are saved in table with following schema:

```
+---------------+--------------+------+-----+---------+-------+
| Field         | Type         | Null | Key | Default | Extra |
+---------------+--------------+------+-----+---------+-------+
| timestamp     | varchar(200) | YES  |     | NULL    |       |
| source_column | varchar(200) | YES  |     | NULL    |       |
| key_column    | varchar(200) | YES  |     | NULL    |       |
| value_column  | varchar(200) | YES  |     | NULL    |       |
+---------------+--------------+------+-----+---------+-------+
```

### Examples

Example of running [psutil collector plugin](https://github.com/intelsdi-x/snap-plugin-collector-psutil) and publishing data to MySQL database.

Set up the [Snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started)

Ensure [Snap daemon is running](https://github.com/intelsdi-x/snap#running-snap):
* initd: `service snap-telemetry start`
* systemd: `systemctl start snap-telemetry`
* command line: `sudo snapteld -l 1 -t 0 &`


Download and load Snap plugins (paths to binary files for Linux/amd64):
```
$ wget http://snap.ci.snap-telemetry.io/plugins/snap-plugin-publisher-mysql/latest/linux/x86_64/snap-plugin-publisher-mysql
$ wget http://snap.ci.snap-telemetry.io/plugins/snap-plugin-collector-psutil/latest/linux/x86_64/snap-plugin-collector-psutil
$ snaptel plugin load snap-plugin-publisher-mysql
$ snaptel plugin load snap-plugin-collector-psutil
```

Create a [task manifest](https://github.com/intelsdi-x/snap/blob/master/docs/TASKS.md) (see [exemplary tasks](examples/tasks/)),
for example `psutil-mysql.json` with following content:
```json
{
  "version": 1,
  "schedule": {
    "type": "simple",
    "interval": "1s"
  },
  "max-failures": 10,
  "workflow": {
    "collect": {
      "metrics": {
        "/intel/psutil/load/load1": {},
        "/intel/psutil/load/load15": {},
        "/intel/psutil/load/load5": {},
        "/intel/psutil/vm/available": {},
        "/intel/psutil/vm/free": {},
        "/intel/psutil/vm/used": {}
      },
      "process": null,
      "publish": [
        {
          "plugin_name": "mysql",
          "config": {
            "username": "root",
            "password": "root",
            "hostname": "localhost",
            "port": "3306",
            "database": "mydb",
            "tablename": "snap_metrics"
          }
        }
      ]
    }
  }
}
```

Create a task:
```
$ snaptel task create -t psutil-mysql.json
```

Watch created task:
```
$ snaptel task watch <task_id>
```

To stop previously created task:
```
$ snaptel task stop <task_id>
```

The running task is collecting data and publishing them into MySQL database. The config section of publisher defining a database connection. If such database does not exist, a new one will be created (if a given user has privileges to create a db).

### Roadmap
There isn't a current roadmap for this plugin, but it is in active development. As we launch this plugin, we do not have any outstanding requirements for the next release. If you have a feature request, please add it as an [issue](https://github.com/intelsdi-x/snap-plugin-publisher-mysql/issues/new) and/or submit a [pull request](https://github.com/intelsdi-x/snap-plugin-publisher-mysql/pulls).

## Community Support
This repository is one of **many** plugins in **Snap**, a powerful telemetry framework. See the full project at http://github.com/intelsdi-x/snap To reach out to other users, head to the [main framework](https://github.com/intelsdi-x/snap#community-support)

## Contributing
We love contributions!

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

## License
[Snap](http://github.com:intelsdi-x/snap), along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements
* Author: [@Lactem](https://github.com/Lactem/)
* Author: [@geauxvirtual](https://github.com/geauxvirtual/)

And **thank you!** Your contribution, through code and participation, is incredibly important to us.
