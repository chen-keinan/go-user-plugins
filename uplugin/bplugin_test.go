package uplugin

import (
	"fmt"
	"github.com/chen-keinan/lxd-probe/pkg/models"
	"os"
	"strings"
	"testing"
)

func TestInvokeFunc(t *testing.T) {
	tests := []struct {
		name         string
		pluginPath   string
		pluginMethod string
		wantParam    interface{}
		wantResult   string
	}{
		{name: "plugin no input and no output", pluginPath: "test_no_input_no_output.so", pluginMethod: "Test", wantParam: "", wantResult: ""},
		{name: "plugin and no output", pluginPath: "test_no_output.so", pluginMethod: "Test", wantParam: "param test", wantResult: ""},
		{name: "plugin test", pluginPath: "test.so", pluginMethod: "Test", wantParam: "param test", wantResult: "return from test"},
		{name: "plugin auditHook", pluginPath: "auditHook.so", pluginMethod: "AuditHook", wantParam: models.AuditBenchResult{Name: "AuditHook"}, wantResult: "AuditHook"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pl := NewPluginLoader("./fixture", "./fixture")
			var got []interface{}
			var err error
			v, ok := tt.wantParam.(string)
			if ok && len(v) == 0 {
				if got, err = pl.LoadAndInvoke(tt.pluginPath, tt.pluginMethod); err != nil {
					t.Errorf("TestInvokeFunc() failed to invoke function %s  error:%s", tt.pluginMethod, err.Error())
				}
			} else {
				if got, err = pl.LoadAndInvoke(tt.pluginPath, tt.pluginMethod, tt.wantParam); err != nil {
					t.Errorf("TestInvokeFunc() failed to invoke function %s error:%s ", tt.pluginMethod, err.Error())
				}
			}
			if len(tt.wantResult) > 0 {
				if tt.wantResult != got[0].(string) {
					t.Errorf("TestInvokeFunc() want %s got %s ", tt.wantResult, got[0].(string))
				}
			}

		})
	}
}

func TestCompile(t *testing.T) {
	tests := []struct {
		name       string
		sourceName string
		methodName string
		funcName   string
		suffix     string
		err        error
		wantResult string
	}{
		{name: "compile and invoke", sourceName: "test.go", funcName: "Test", wantResult: "test string", suffix: CompiledExt, err: nil},
		{name: "compile ano file exist", sourceName: "test1.go", funcName: "Test", wantResult: "test string", suffix: CompiledExt, err: fmt.Errorf("could not read fixture/test1.go: open fixture/test1.go: no such file or directory")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pl := NewPluginLoader("./fixture", "./fixture")
			var objName string
			var err error
			if objName, err = pl.Compile(tt.sourceName); err != nil {
				if err.Error() != tt.err.Error() {
					t.Errorf("TestCompile() want %s got %s ", err.Error(), tt.err.Error())
				}
			}
			if err == nil {
				if !strings.HasSuffix(objName, tt.suffix) {
					t.Errorf("TestCompile() want %s got %s ", tt.suffix, objName)
				}
				err = os.RemoveAll(objName)
				if err != nil {
					t.Errorf("failed to delete compiled plugin %s", objName)
				}
			}
		})
	}
}

func TestPlugins(t *testing.T) {
	tests := []struct {
		name      string
		ext       string
		pluginDir string
		err       error
		want      int
	}{
		{name: "source plugins", ext: SourceExt, want: 1, pluginDir: "./fixture", err: nil},
		{name: "plugin bad folder", ext: SourceExt, want: 1, pluginDir: "./fixture1", err: fmt.Errorf("open ./fixture1: no such file or directory")},
		{name: "compiled plugins", ext: CompiledExt, pluginDir: "./fixture", want: 4, err: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pl := NewPluginLoader(tt.pluginDir, tt.pluginDir)
			var plugins []string
			var err error
			if plugins, err = pl.Plugins(tt.ext); err != nil {
				if err.Error() != tt.err.Error() {
					t.Errorf("TestPlugins() want %s got %s ", tt.err.Error(), err.Error())
				}
			}
			if len(plugins) > 0 {
				if len(plugins) != tt.want {
					t.Errorf("TestPlugins() want %d got %d ", tt.want, len(plugins))
				}
			}
		})
	}
}
