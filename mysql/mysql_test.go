// +build small

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
	"testing"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core/ctypes"
	. "github.com/smartystreets/goconvey/convey"
)

type foo int

func TestMySQLPlugin(t *testing.T) {
	Convey("Meta should return metadata for the plugin", t, func() {
		meta := Meta()
		So(meta.Name, ShouldResemble, name)
		So(meta.Version, ShouldResemble, version)
		So(meta.Type, ShouldResemble, plugin.PublisherPluginType)
	})

	Convey("Create MySQLPublisher", t, func() {
		mp := NewMySQLPublisher()
		Convey("So mp should not be nil", func() {
			So(mp, ShouldNotBeNil)
		})
		Convey("So mp should be of mysqlPublisher type", func() {
			So(mp, ShouldHaveSameTypeAs, &mysqlPublisher{})
		})
		configPolicy, err := mp.GetConfigPolicy()
		Convey("ip.GetConfigPolicy() should return a config policy", func() {
			Convey("So config policy should not be nil", func() {
				So(configPolicy, ShouldNotBeNil)
			})
			Convey("So getting a config policy should not return an error", func() {
				So(err, ShouldBeNil)
			})
			Convey("So config policy should be a cpolicy.ConfigPolicy", func() {
				So(configPolicy, ShouldHaveSameTypeAs, &cpolicy.ConfigPolicy{})
			})

			Convey("so processing configuration with all parameters configured", func() {
				testConfig := make(map[string]ctypes.ConfigValue)
				testConfig["username"] = ctypes.ConfigValueStr{Value: "root1"}
				testConfig["password"] = ctypes.ConfigValueStr{Value: "root1"}
				testConfig["hostname"] = ctypes.ConfigValueStr{Value: "localhost1"}
				testConfig["port"] = ctypes.ConfigValueStr{Value: "33061"}
				testConfig["database"] = ctypes.ConfigValueStr{Value: "SNAP_TEST1"}
				testConfig["tablename"] = ctypes.ConfigValueStr{Value: "info1"}

				cfg, errs := configPolicy.Get([]string{""}).Process(testConfig)

				Convey("So config policy should process testConfig and return a config", func() {
					So(cfg, ShouldNotBeNil)
				})

				Convey("so parameters should have correct values", func() {
					So((*cfg)["username"].(ctypes.ConfigValueStr).Value, ShouldEqual, "root1")
					So((*cfg)["password"].(ctypes.ConfigValueStr).Value, ShouldEqual, "root1")
					So((*cfg)["hostname"].(ctypes.ConfigValueStr).Value, ShouldEqual, "localhost1")
					So((*cfg)["port"].(ctypes.ConfigValueStr).Value, ShouldEqual, "33061")
					So((*cfg)["database"].(ctypes.ConfigValueStr).Value, ShouldEqual, "SNAP_TEST1")
					So((*cfg)["tablename"].(ctypes.ConfigValueStr).Value, ShouldEqual, "info1")
				})

				Convey("So testConfig processing should return no errors", func() {
					So(errs.HasErrors(), ShouldBeFalse)
				})
			})

			Convey("so processing configuration without parameters configured", func() {
				testConfig := make(map[string]ctypes.ConfigValue)

				cfg, errs := configPolicy.Get([]string{""}).Process(testConfig)

				Convey("So config policy should process testConfig and return a config", func() {
					So(cfg, ShouldNotBeNil)
				})

				Convey("so parameters should have correct values", func() {
					So((*cfg)["username"].(ctypes.ConfigValueStr).Value, ShouldEqual, "root")
					So((*cfg)["password"].(ctypes.ConfigValueStr).Value, ShouldEqual, "root")
					So((*cfg)["hostname"].(ctypes.ConfigValueStr).Value, ShouldEqual, "localhost")
					So((*cfg)["port"].(ctypes.ConfigValueStr).Value, ShouldEqual, "3306")
					So((*cfg)["database"].(ctypes.ConfigValueStr).Value, ShouldEqual, "SNAP_TEST")
					So((*cfg)["tablename"].(ctypes.ConfigValueStr).Value, ShouldEqual, "info")
				})

				Convey("So testConfig processing should return no errors", func() {
					So(errs.HasErrors(), ShouldBeFalse)
				})
			})
		})
	})
}

func TestInterfaceToString(t *testing.T) {
	Convey("Return properly formatted values", t, func() {
		Convey("Slice of strings should be formatted", func() {
			slice := []string{"foo", "bar", "baz"}
			Convey("So a slice with length greater than 1 should format", func() {
				str, err := interfaceToString(slice)
				So(err, ShouldBeNil)
				So(str, ShouldEqual, "foo, bar, baz")
			})
			Convey("So a slice with length of 1 should format", func() {
				str, err := interfaceToString(slice[:1])
				So(err, ShouldBeNil)
				So(str, ShouldEqual, "foo")
			})
			Convey("So an empty array should format", func() {
				str, err := interfaceToString([]string{})
				So(err, ShouldBeNil)
				So(str, ShouldBeBlank)
			})
		})
		Convey("So a string should be the same", func() {
			str, err := interfaceToString("foo")
			So(err, ShouldBeNil)
			So(str, ShouldEqual, "foo")
		})
		Convey("Slice of ints should be formatted", func() {
			slice := []int{1, 2, 3}
			Convey("So a slice with length greater than 1 should format", func() {
				str, err := interfaceToString(slice)
				So(err, ShouldBeNil)
				So(str, ShouldEqual, "1, 2, 3")
			})
			Convey("So a slice with length of 1 should format", func() {
				str, err := interfaceToString(slice[:1])
				So(err, ShouldBeNil)
				So(str, ShouldEqual, "1")
			})
			Convey("So an empty array should format", func() {
				str, err := interfaceToString([]int{})
				So(err, ShouldBeNil)
				So(str, ShouldBeBlank)
			})
		})
		Convey("So an int should be formatted", func() {
			str, err := interfaceToString(1)
			So(err, ShouldBeNil)
			So(str, ShouldEqual, "1")
		})
		Convey("Slice of uints should be formatted", func() {
			slice := []uint{1, 2, 3}
			Convey("So a slice with length greater than 1 should format", func() {
				str, err := interfaceToString(slice)
				So(err, ShouldBeNil)
				So(str, ShouldEqual, "1, 2, 3")
			})
			Convey("So a slice with length of 1 should format", func() {
				str, err := interfaceToString(slice[:1])
				So(err, ShouldBeNil)
				So(str, ShouldEqual, "1")
			})
			Convey("So an empty array should format", func() {
				str, err := interfaceToString([]uint{})
				So(err, ShouldBeNil)
				So(str, ShouldBeBlank)
			})
		})
		Convey("So an uint should be formatted", func() {
			str, err := interfaceToString(uint(1))
			So(err, ShouldBeNil)
			So(str, ShouldEqual, "1")
		})
		Convey("Slice of uint64s should be formatted", func() {
			slice := []uint64{1, 2, 3}
			Convey("So a slice with length greater than 1 should format", func() {
				str, err := interfaceToString(slice)
				So(err, ShouldBeNil)
				So(str, ShouldEqual, "1, 2, 3")
			})
			Convey("So a slice with length of 1 should format", func() {
				str, err := interfaceToString(slice[:1])
				So(err, ShouldBeNil)
				So(str, ShouldEqual, "1")
			})
			Convey("So an empty array should format", func() {
				str, err := interfaceToString([]uint64{})
				So(err, ShouldBeNil)
				So(str, ShouldBeBlank)
			})
		})
		Convey("So an uint64 should be formatted", func() {
			str, err := interfaceToString(uint64(1))
			So(err, ShouldBeNil)
			So(str, ShouldEqual, "1")
		})
		Convey("Slice of float64s should be formatted", func() {
			slice := []float64{1.123456789, 2.123456789}
			Convey("So a slice with length greater than 1 should format", func() {
				str, err := interfaceToString(slice)
				So(err, ShouldBeNil)
				So(str, ShouldEqual, "1.123456789, 2.123456789")
			})
			Convey("So a slice with length of 1 should format", func() {
				str, err := interfaceToString(slice[:1])
				So(err, ShouldBeNil)
				So(str, ShouldEqual, "1.123456789")
			})
			Convey("So an empty array should format", func() {
				str, err := interfaceToString([]int{})
				So(err, ShouldBeNil)
				So(str, ShouldBeBlank)
			})
		})
		Convey("So an float64 should be formatted", func() {
			str, err := interfaceToString(float64(1.123456789))
			So(err, ShouldBeNil)
			So(str, ShouldEqual, "1.123456789")
		})
		Convey("So nil should be formatted", func() {
			str, err := interfaceToString(nil)
			So(err, ShouldBeNil)
			So(str, ShouldEqual, "nil")
		})
		Convey("So an unsupported type should return an error", func() {
			str, err := interfaceToString(foo(1))
			So(str, ShouldBeBlank)
			So(err.Error(), ShouldStartWith, "Unsupported type")
		})
	})
}
