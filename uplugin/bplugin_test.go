package uplugin

import (
	"strings"
	"testing"
)

func TestInvokeFunc(t *testing.T) {
	tests := []struct {
		name         string
		pluginPath   string
		pluginMethod string
		wantParam    string
		wantResult   string
	}{
		{name: "plugin no input and no output", pluginPath: "test_no_input_no_output.so", pluginMethod: "Test", wantParam: "", wantResult: ""},
		{name: "plugin and no output", pluginPath: "test_no_output.so", pluginMethod: "Test", wantParam: "param test", wantResult: ""},
		{name: "plugin test", pluginPath: "test.so", pluginMethod: "Test", wantParam: "param test", wantResult: "return from test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pl, err := NewPluginLoader("./fixture")
			if err != nil {
				t.Errorf("failed to open plugin %s", tt.pluginPath)
			}
			sym, err := pl.Load(tt.pluginPath, tt.pluginMethod)
			if err != nil {
				t.Errorf("failed to open plugin %s", tt.pluginPath)
			}
			var got []interface{}
			if len(tt.wantParam) == 0 {
				if got, err = pl.InvokeFunc(sym); err != nil {
					t.Errorf("TestInvokeFunc() failed to invoke function %s  error:%s", tt.pluginMethod, err.Error())
				}
			} else {
				if got, err = pl.InvokeFunc(sym, tt.wantParam); err != nil {
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
		wantResult string
	}{
		{name: "compile and invoke", sourceName: "test.go", funcName: "Test", wantResult: "test string"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pl, err := NewPluginLoader("./fixture")
			if err != nil {
				t.Errorf("failed to open plugin %s", tt.sourceName)
			}
			var objName string
			if objName, err = pl.Compile(tt.sourceName); err != nil {
				t.Errorf("TestCompile() failed to compile source %s error:%s", tt.sourceName, err.Error())
			}
			if !strings.HasSuffix(objName, ".so") {
				t.Errorf("TestCompile() want %s got %s ", "so", objName)
			}
		})
	}
}
