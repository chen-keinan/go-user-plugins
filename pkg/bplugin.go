package bplugin

import (
	"fmt"
  	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"plugin"
)

// PluginLoader keeps the context needed to find where ObjPlugins and
// objects are stored.
type PluginLoader struct {
	pluginsDir string
	objectsDir string
}

//NewPluginLoader return new plugin loader object with src and compiled folders
func NewPluginLoader(obgPath string) (*PluginLoader, error) {
	return &PluginLoader{objectsDir: obgPath}, nil
}

//Compile the go plugin in a given path and hook name and return it symbol
func (l *PluginLoader) Compile(name string, hookName string) (plugin.Symbol, error) {
	obj, err := l.compile(name)
	if err != nil {
		return nil, fmt.Errorf("could not compile %s: %v", name, err)
	}

	var sym plugin.Symbol
	if sym, err = l.Load(obj, hookName); err != nil {
		return nil, fmt.Errorf("could not compile %s: %v", name, err)
	}
	return sym, nil
}

// compile compiles the code in the given path, builds a
// plugin, and returns its path.
//nolint:gosec
func (l *PluginLoader) compile(name string) (string, error) {
	// Copy the file to the objects directory with a different name
	// each time, to avoid retrieving the cached version.
	// Apparently the cache key is the path of the file compiled and
	// there's no way to invalidate it.
	f, err := ioutil.ReadFile(filepath.Join(l.pluginsDir, name))
	if err != nil {
		return "", fmt.Errorf("could not read %s: %v", name, err)
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
	p, err := plugin.Open(object)
	if err != nil {
		return nil, fmt.Errorf("could not open %s: %v", object, err)
	}
	return p.Lookup(hookName)
}

/*//ExecuteK8sAuditResults execute on K8s api call hook
func ExecuteK8sAuditResults(sym plugin.Symbol, auditEvt models.KubeAuditResults) error {
	runFunc, ok := sym.(func(netEvent models.KubeAuditResults) error)
	if !ok {
		return fmt.Errorf("found Run but type is %T instead of func() error", sym)
	}
	return runFunc(auditEvt)
}*/

//Plugins lists all the files in the ObjPlugins
func (l *PluginLoader) Plugins() ([]string, error) {
	dir, err := os.Open(l.pluginsDir)
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
		if filepath.Ext(name) == ".go" {
			res = append(res, name)
		}
	}
	return res, nil
}

//ObjPlugins lists all the files in the ObjPlugins
func (l *PluginLoader) ObjPlugins() ([]string, error) {
	dir, err := os.Open(l.objectsDir)
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
		if filepath.Ext(name) == ".so" {
			res = append(res, path.Join(l.objectsDir, name))
		}
	}
	return res, nil
}
