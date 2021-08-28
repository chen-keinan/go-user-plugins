package userplugin

import (
	"fmt"
	"testing"
)

func TestUserPlugin(t *testing.T) {
	var path = "/vagrant/test_plugin"
	pl,err:=NewPluginLoader(path)
	if err != nil {
		t.Errorf(err.Error())
	}
	s,err:=pl.Load("test.so","Test")
	if err != nil {
		fmt.Println(err.Error())
	}
	a,err:=pl.InvokeFunc(s)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(a.(string))
}
