//
// +build integration

package mysql

import (
	"bytes"
	"encoding/gob"
	"testing"
	"time"

	"github.com/intelsdi-x/pulse/control/plugin"
	"github.com/intelsdi-x/pulse/core/ctypes"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMySQLPublish(t *testing.T) {
	var buf bytes.Buffer
	metrics := []plugin.PluginMetricType{
		*plugin.NewPluginMetricType([]string{"test", "string"}, time.Now(), "127.0.0.1", "example_string"),
		*plugin.NewPluginMetricType([]string{"test", "int"}, time.Now(), "127.0.0.1", 1),
		*plugin.NewPluginMetricType([]string{"test", "string", "slice"}, time.Now(), "localhost", []string{"str1", "str2"}),
		*plugin.NewPluginMetricType([]string{"test", "string", "slice"}, time.Now(), "localhost", []int{1, 2}),
	}
	config := make(map[string]ctypes.ConfigValue)
	enc := gob.NewEncoder(&buf)
	enc.Encode(metrics)
	config["username"] = ctypes.ConfigValueStr{Value: "root"}
	config["password"] = ctypes.ConfigValueStr{Value: ""}
	config["database"] = ctypes.ConfigValueStr{Value: "pulse_test"}
	config["tablename"] = ctypes.ConfigValueStr{Value: "info"}
	sp := NewMySQLPublisher()
	Convey("Publish metrics to MySQL instance should succeed and not throw an error", t, func() {
		err := sp.Publish(plugin.PulseGOBContentType, buf.Bytes(), config)
		So(err, ShouldBeNil)
	})
}
