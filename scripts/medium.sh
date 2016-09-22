#!/bin/bash

set -e
set -u
set -o pipefail

mysql -e "create database IF NOT EXISTS snap_test;" -uroot
SNAP_MYSQL_HOST=127.0.0.1
UNIT_TEST="go_test"
test_unit
