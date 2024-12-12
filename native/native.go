package native

import (
	"github.com/pmqueiroz/umbra/environment"
)

type InternalModuleFn func([]interface{}) (interface{}, error)

type InternalModule struct {
	symbols map[string]InternalModuleFn
}

func (m InternalModule) Register(namespace *environment.Environment) (ok bool) {
	for name, symbol := range m.symbols {
		createOk := namespace.Create(name, symbol, "", false, true, false)
		pubOk := namespace.MakePublic(name)

		if !createOk || !pubOk {
			return false
		}
	}

	return true
}

func Register(name string, namespace *environment.Environment) (ok bool) {
	var module InternalModule

	switch name {
	case "os":
		module = OsModule
	case "path":
		module = PathModule
	case "hashmaps":
		module = HashmapModule
	default:
		ok = false
		return
	}

	ok = module.Register(namespace)
	return
}
