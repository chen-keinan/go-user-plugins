package uplugin

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"plugin"
	"reflect"
)

const (
	//SourceExt const
	SourceExt = ".go"
	//CompiledExt const
	CompiledExt = ".so"
)

// PluginLoader keeps the context needed to find where ObjPlugins and
// objects are stored.
type PluginLoader struct {
	pluginsDir string
	objectsDir string
}

//NewPluginLoader return new plugin loader object with src and compiled folders
func NewPluginLoader(sourcePath, objPath string) *PluginLoader {
	return &PluginLoader{objectsDir: objPath, pluginsDir: sourcePath}
}

//Compile the go plugin in a given path and hook name and return it symbol
func (l *PluginLoader) Compile(name string) (string, error) {
	return l.compile(name)
}

// compile compiles the code in the given path, builds a
// plugin, and returns its path.
//nolint:gosec
func (l *PluginLoader) compile(name string) (string, error) {
	// Copy the file to the objects directory with a different name
	// each time, to avoid retrieving the cached version.
	// Apparently the cache key is the path of the file compiled and
	// there's no way to invalidate it.
	fullPath := filepath.Join(l.pluginsDir, name)
	f, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return "", fmt.Errorf("could not read %s: %v", fullPath, err)
	}

	name = fmt.Sprintf("%d.go", rand.Int())
	srcPath := filepath.Join(l.objectsDir, name)
	fileCreated, err := os.Create(srcPath)
	if err != nil {
		return "", fmt.Errorf("could not write %s: %v", name, err)
	}
	defer func() {
		err = fileCreated.Close()
		if err != nil {
			fmt.Print(err.Error())
		}
		err = os.Remove(path.Join(l.pluginsDir, name))
		if err != nil {
			fmt.Print(err.Error())
		}
	}()
	_, err = fileCreated.WriteString(string(f))
	if err != nil {
		return "", fmt.Errorf("could not write %s: %v", name, err)
	}
	objectPath := srcPath[:len(srcPath)-3] + ".so"
	cmd := exec.Command("go", "build", "-buildmode=plugin", fmt.Sprintf("-o=%s", objectPath), srcPath)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("could not compile %s: %v", name, err)
	}
	return objectPath, nil
}

// Load loads the plugin object in the given path and runs the Run
// function.
func (l *PluginLoader) Load(object string, hookName string) (plugin.Symbol, error) {
	fullPath := path.Join(l.objectsDir, object)
	p, err := plugin.Open(fullPath)
	if err != nil {
		return nil, err
	}
	return p.Lookup(hookName)
}

//LoadAndInvoke invoke plugin function with params and return it results
func (l *PluginLoader) LoadAndInvoke(plugin, method string, params ...interface{}) ([]interface{}, error) {
	sym, err := l.Load(plugin, method)
	if err != nil {
		return nil, err
	}
	return l.Invoke(sym, params...)
}

//Invoke invoke plugin function with params and return it results
func (l *PluginLoader) Invoke(sym plugin.Symbol, params ...interface{}) ([]interface{}, error) {
	results := make([]interface{}, 0)
	f := reflect.ValueOf(sym)
	if len(params) != f.Type().NumIn() {
		return nil, fmt.Errorf("the number of params %d is out of index", len(params))
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	res := f.Call(in)
	if len(res) > 0 {
		for _, r := range res {
			results = append(results, r.Interface())
		}
	}
	return results, nil
}

//nolint gosec
//Plugins lists all the files in the ObjPlugins
func (l *PluginLoader) Plugins(ext string) ([]string, error) {
	pluginFolder := l.pluginsDir
	if ext == CompiledExt {
		pluginFolder = l.objectsDir
	}
	dir, err := os.Open(pluginFolder)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = dir.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}()
	names, err := dir.Readdirnames(-1)
	if err != nil {
		return nil, err
	}

	var res []string
	for _, name := range names {
		if filepath.Ext(name) == ext {
			res = append(res, name)
		}
	}
	return res, nil
}
