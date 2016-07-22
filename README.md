# snap publisher plugin - mysql
This plugin publishes data to MySQL database. 

It's used in the [snap framework](http://github.com:intelsdi-x/snap).

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Installation](#installation)
  * [Configuration and Usage](configuration-and-usage)
2. [Documentation](#documentation)
  * [Roadmap](#roadmap)
3. [Community Support](#community-support)
4. [Contributing](#contributing)
5. [License](#license-and-authors)
6. [Acknowledgements](#acknowledgements)

## Getting Started
### System Requirements
* [golang 1.5+](https://golang.org/dl/) (needed only for building)

### Operating systems
All OSs currently supported by snap:
* Linux/amd64
* Darwin/amd64

### Installation
#### Download mysql plugin binary:
You can get the pre-built binaries for your OS and architecture at snap's [GitHub Releases](https://github.com/intelsdi-x/snap/releases) page.

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
This builds the plugin in `/build/rootfs/`

### Configuration and Usage
* Set up the [snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started)
* Ensure `$SNAP_PATH` is exported  
`export SNAP_PATH=$GOPATH/src/github.com/intelsdi-x/snap/build`
* Load the plugin and create a task, see example in Examples.

## Documentation
There are a number of other resources you can review to learn to use this plugin:

* [mysql](https://www.mysql.com/) 
* [go-mysql-driver](github.com/go-sql-driver/mysql)
* [snap mysql integration test](https://github.com/intelsdi-x/snap-plugin-publisher-mysql/blob/master/mysql/mysql_integration_test.go)
* [snap mysql unit test](https://github.com/intelsdi-x/snap-plugin-publisher-mysql/blob/master/mysql/mysql_test.go)

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

Each connection parameter has a default value, but it can be override by set a value in task manifest (see [exemplary task manifest](examples/tasks/mock-mysql.json))

### Examples
Example of running snap mock collector and publishing data to mysql database.

Run the snap daemon:
```
$ $SNAP_PATH/bin/snapd -l 1 -t 0
```

In another terminal load mock collector plugin:
```
$ $SNAP_PATH/bin/snapctl plugin load $SNAP_PATH/plugin/snap-collector-mock1
Plugin loaded
Name: mock
Version: 1
Type: collector
Signed: false
Loaded Time: Wed, 22 Jun 2016 14:49:38 CEST
```

Load mysql publisher plugin:
```
$ $SNAP_PATH/bin/snapctl plugin load snap-plugin-publisher-mysql/build/rootfs/snap-plugin-publisher-mysql
Plugin loaded
Name: mysql
Version: 8
Type: publisher
Signed: false
Loaded Time: Wed, 22 Jun 2016 14:53:14 CEST
```
Create a task JSON file (exemplary files in [examples/tasks/] (examples/tasks/)):
```json
{
  "version": 1,
    "schedule": {
      "type": "simple",
      "interval": "1s"
    },
    "workflow": {
      "collect": {
        "metrics": {
          "/intel/mock/*": {}
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
$ $SNAP_PATH/bin/snapctl task create -t snap-plugin-publisher-mysql/examples/tasks/mock-mysql.json
Using task manifest to create task
Task created
ID: d4392f17-11f0-4f64-8701-2708f432b50a
Name: Task-d4392f17-11f0-4f64-8701-2708f432b50a
State: Running
```

The running task is collecting data and publishing them into mysql database. The config section of publisher defining a database connection. If such database does not exist, a new one will be created (if a given user has privileges to create a db).

### Roadmap
There isn't a current roadmap for this plugin, but it is in active development. As we launch this plugin, we do not have any outstanding requirements for the next release. If you have a feature request, please add it as an [issue](https://github.com/intelsdi-x/snap-plugin-publisher-mysql/issues/new) and/or submit a [pull request](https://github.com/intelsdi-x/snap-plugin-publisher-mysql/pulls).

## Community Support
This repository is one of **many** plugins in **snap**, a powerful telemetry framework. See the full project at http://github.com/intelsdi-x/snap To reach out to other users, head to the [main framework](https://github.com/intelsdi-x/snap#community-support)

## Contributing
We love contributions!

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

## License
[snap](http://github.com:intelsdi-x/snap), along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements
* Author: [@Lactem](https://github.com/Lactem/)
* Author: [@geauxvirtual](https://github.com/geauxvirtual/)

And **thank you!** Your contribution, through code and participation, is incredibly important to us.
