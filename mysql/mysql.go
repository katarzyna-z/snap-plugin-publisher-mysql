/*
http://www.apache.org/licenses/LICENSE-2.0.txt


Copyright 2015 Intel Corporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package mysql

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core"
	"github.com/intelsdi-x/snap/core/ctypes"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

const (
	name            = "mysql"
	version         = 8
	pluginType      = plugin.PublisherPluginType
	usernameDefault = "root"
	passwordDefault = "root"
	hostnameDefault = "localhost"
	tcpPortDefault  = "3306"

	databaseDefault = "SNAP_TEST"
	tableDefault    = "info"
	tableColumns    = "(timestamp VARCHAR(200), source_column VARCHAR(200), key_column VARCHAR(200), value_column VARCHAR(200))"
)

type mysqlPublisher struct {
	db           *sql.DB
	dbInsertStmt *sql.Stmt
}

func NewMySQLPublisher() *mysqlPublisher {
	return &mysqlPublisher{}
}

// Publish sends data to a MySQL server
func (s *mysqlPublisher) Publish(contentType string, content []byte, cfg map[string]ctypes.ConfigValue) error {
	logger := log.New()
	logger.Println("Publishing started")
	var metrics []plugin.MetricType

	switch contentType {
	case plugin.SnapGOBContentType:
		dec := gob.NewDecoder(bytes.NewBuffer(content))
		if err := dec.Decode(&metrics); err != nil {
			logger.Printf("Error decoding: error=%v content=%v", err, content)
			return err
		}
	default:
		logger.Printf("Error unknown content type '%v'", contentType)
		return errors.New(fmt.Sprintf("Unknown content type '%s'", contentType))
	}

	if err := s.init(cfg); err != nil {
		s.db.Close()
	}

	for _, m := range metrics {
		key := sliceToString(m.Namespace().Strings())
		value, err := interfaceToString(m.Data())
		if err != nil {
			logger.Printf("Error: Cannot convert incoming data to string, err=%v", err)
			return err
		}
		_, err = s.dbInsertStmt.Exec(m.Timestamp(), m.Tags()[core.STD_TAG_PLUGIN_RUNNING_ON], key, value)
		if err != nil {
			logger.Printf("Error: Cannot publish incoming metric to mysql db, err=%v", err)
			return err
		}
	}
	return nil
}

func Meta() *plugin.PluginMeta {
	return plugin.NewPluginMeta(name, version, pluginType, []string{plugin.SnapGOBContentType}, []string{plugin.SnapGOBContentType})
}

func (s *mysqlPublisher) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	cp := cpolicy.New()
	config := cpolicy.NewPolicyNode()

	username, err := cpolicy.NewStringRule("username", false, usernameDefault)
	handleErr(err)
	username.Description = "Username to login to the MySQL server"

	password, err := cpolicy.NewStringRule("password", false, passwordDefault)
	handleErr(err)
	password.Description = "Password to login to the MySQL server"

	hostName, err := cpolicy.NewStringRule("hostname", false, hostnameDefault)
	handleErr(err)
	password.Description = "The host of MySQL service"

	port, err := cpolicy.NewStringRule("port", false, tcpPortDefault)
	handleErr(err)
	password.Description = "The host of MySQL service"

	database, err := cpolicy.NewStringRule("database", false, databaseDefault)
	handleErr(err)
	database.Description = "The MySQL database that data will be pushed to"

	tableName, err := cpolicy.NewStringRule("tablename", false, tableDefault)
	handleErr(err)
	tableName.Description = "The MySQL table within the database where information will be stored"

	config.Add(username, password, hostName, port, database, tableName)

	cp.Add([]string{""}, config)
	return cp, nil
}

func (s *mysqlPublisher) init(cfg map[string]ctypes.ConfigValue) error {
	var err error
	logger := log.New()

	url := getMySQLConnectionURL(cfg["username"].(ctypes.ConfigValueStr).Value,
		cfg["password"].(ctypes.ConfigValueStr).Value,
		cfg["hostname"].(ctypes.ConfigValueStr).Value,
		cfg["port"].(ctypes.ConfigValueStr).Value)

	// Open connection and ping to make sure it works
	s.db, err = sql.Open("mysql", url)
	if err != nil {
		logger.Printf("Error: cannot open %v, err=%v", url, err)
		return err
	}

	// test connection
	err = s.db.Ping()
	if err != nil {
		logger.Printf("Error: cannot establish a connection, err=%v", err)
		return err
	}

	// switch the connection when SelectDB is defined in cfg
	if _, err = s.db.Exec("USE " + cfg["database"].(ctypes.ConfigValueStr).Value); err != nil {

		if _, err = s.db.Exec("CREATE DATABASE " + cfg["database"].(ctypes.ConfigValueStr).Value); err != nil {
			logger.Printf("Error: cannot create a new database `%v`, err=%v", cfg["database"].(ctypes.ConfigValueStr).Value, err)
			return err
		}

		// use this already created db
		if _, err = s.db.Exec("USE " + cfg["database"].(ctypes.ConfigValueStr).Value); err != nil {
			return err
		}
	}

	// Create the table if it's not already there
	_, err = s.db.Exec("CREATE TABLE IF NOT EXISTS" + " " + cfg["tablename"].(ctypes.ConfigValueStr).Value + " " + tableColumns)
	if err != nil {
		logger.Printf("Error: cannot create table %v, err=%v", cfg["tablename"].(ctypes.ConfigValueStr).Value, err)
		return err
	}

	// Put the values into the database with the current time
	s.dbInsertStmt, err = s.db.Prepare("INSERT INTO" + " " + cfg["tablename"].(ctypes.ConfigValueStr).Value + " VALUES( ?, ?, ?, ? )")
	if err != nil {
		fmt.Printf("Error: cannot prepare insert db statement, err=%v", err)
		return err
	}
	return nil
}

func getMySQLConnectionURL(user, passwd, host, port string) string {
	// formatting as `user:passwd@tcp(host:port)'
	mysqlConnectionURL := user + ":" + passwd + "@tcp(" + host + ":" + port + ")/"
	return mysqlConnectionURL
}

func handleErr(e error) {
	if e != nil {
		panic(e)
	}
}

func sliceToString(slice []string) string {
	return strings.Join(slice, ", ")
}

func interfaceToString(face interface{}) (string, error) {
	var (
		ret string
		err error
	)
	switch val := face.(type) {
	case []string:
		ret = sliceToString(val)
	case string:
		ret = val
	case []int:
		length := len(val)
		if length == 0 {
			return ret, err
		}
		ret = strconv.Itoa(val[0])
		if length == 1 {
			return ret, err
		}
		for i := 1; i < length; i++ {
			ret += ", "
			ret += strconv.Itoa(val[i])
		}
	case int:
		ret = strconv.Itoa(val)
	case []uint:
		length := len(val)
		if length == 0 {
			return ret, err
		}
		ret = strconv.FormatUint(uint64(val[0]), 10)
		if length == 1 {
			return ret, err
		}
		for i := 1; i < length; i++ {
			ret += ", "
			ret += strconv.FormatUint(uint64(val[i]), 10)
		}
	case []uint64:
		length := len(val)
		if length == 0 {
			return ret, err
		}
		ret = strconv.FormatUint(val[0], 10)
		if length == 1 {
			return ret, err
		}
		for i := 1; i < length; i++ {
			ret += ", "
			ret += strconv.FormatUint(val[i], 10)
		}
	case uint:
		ret = strconv.FormatUint(uint64(val), 10)
	case uint64:
		ret = strconv.FormatUint(val, 10)
	case float64:
		ret = strconv.FormatFloat(val, 'g', -1, 64)
	case []float64:
		length := len(val)
		if length == 0 {
			return ret, err
		}
		ret = strconv.FormatFloat(val[0], 'g', -1, 64)
		if length == 1 {
			return ret, err
		}
		for i := 1; i < length; i++ {
			ret += ", "
			ret += strconv.FormatFloat(val[i], 'g', -1, 64)
		}
	case nil:
		ret = "nil"
	default:
		err = fmt.Errorf("Unsupported type %v (currently supported data type: string, []string, int, []int, uint, []uint, uint64, []uint64, float64, []float64, or nil)", reflect.TypeOf(val))
	}
	return ret, err
}
