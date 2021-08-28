package uplugin

import (
	"testing"
)

func TestUserPlugin(t *testing.T) {
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
		pl, err := NewPluginLoader("./fixture")
		if err != nil {
			t.Errorf("failed to open plugin %s", tt.pluginPath)
		}
		t.Run(tt.name, func(t *testing.T) {
			sym, err := pl.Load(tt.pluginPath, tt.pluginMethod)
			if err != nil {
				t.Errorf("failed to open plugin %s", tt.pluginPath)
			}
			var got []interface{}
			if len(tt.wantParam) == 0 {
				if got, err = pl.InvokeFunc(sym); err != nil {
					t.Errorf("TestUserPlugin() failed to invoke function %s ", tt.pluginMethod)
				}
			} else {
				if got, err = pl.InvokeFunc(sym, tt.wantParam); err != nil {
					t.Errorf("TestUserPlugin() failed to invoke function %s ", tt.pluginMethod)
				}
			}
			if len(tt.wantResult) > 0 {
				if tt.wantResult != got[0].(string) {
					t.Errorf("TestUserPlugin() want %s got %s ", tt.wantResult, got[0].(string))
				}
			}

		})
	}
}
