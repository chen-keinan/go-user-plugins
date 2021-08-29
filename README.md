[![Go Report Card](https://goreportcard.com/badge/github.com/chen-keinan/go-user-plugins)](https://goreportcard.com/report/github.com/chen-keinan/go-user-plugins)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/chen-keinan/go-user-plugins/blob/master/LICENSE)
<img src="./pkg/img/coverage_badge.png" alt="test coverage badge">
<br><img src="./pkg/img/golang-plugins.png" width="300" alt="golang plugin logo"><br>
# go-user-plugins

go-user-plugins is an open source lib who perform action on user plugin ,plugin generation from source ,load and invoke a plugin method


* [Installation](#installation)
* [Usage](#usage)
* [Contribution](#Contribution)


## Installation

```
go get github.com/chen-keinan/go-user-plugins
```

## Usage
### Generate plugin from source

```go
// init plugin folder
userPlugin := NewPluginLoader("./soureFolder", "./objFolder")
sourcePluginName:="test.go"
// compile plugin from source
compiledPlugin, err:=userPlugin.Compile(sourcePluginName)
if err != nil {
    fmt.Print(err.Error())
}
// print compiled plugin name
fmt.Println(fmt.Sprintf("compiled plugin name %s",compiledPlugin))
```

### Load and invoke a compiled plugin

```go
// init plugin folder
userPlugin:= NewPluginLoader("./soureFolder", "./objFolder")
compiledPluginName:="test.so"
pluginMethodName:="Test"
// load plugin method symbol
methodsymbol, err:=userPlugin.Load(compiledPluginName,pluginMethodName)
if err != nil {
    fmt.Print(err.Error())
}
// invoke plugin method symbol
results,err:=userPlugin.InvokeFunc(methodSymbol)
if err != nil {
  fmt.Print(err.Error())
}
res:=results[0].(string)
fmt.Println(fmt.Sprintf(res)
```

### Full Example 
#### test.go (source plugin)
```go
package main

//Test this plugin
func Test() string {
	return "test string"
}
```
#### Compile plugin
```shell
 go build -buildmode=plugin -o ./objFolder/test.so ./soureFolder/test.go
```

#### test.go (source plugin)
```go
//Test this plugin
userPlugin:= NewPluginLoader("./soureFolder", "./objFolder")
compiledPluginName:="test.so"
pluginMethodName:="Test"
// load plugin method symbol
methodsymbol, err:=userPlugin.Load(compiledPluginName,pluginMethodName)
if err != nil {
    fmt.Print(err.Error())
}
// invoke plugin method symbol
results,err:=userPlugin.InvokeFunc(methodSymbol)
if err != nil {
    fmt.Print(err.Error())
}
res:=results[0].(string)
fmt.Println(fmt.Sprintf(res)
```

## Contribution
code contribution is welcome !! , contribution with passing tests and linter is more than welcome :)