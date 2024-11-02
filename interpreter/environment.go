package interpreter

import (
	"fmt"
	"os"

	"github.com/pmqueiroz/umbra/exception"
)

type Value struct {
	data      interface{}
	isPrivate bool
}

type Namespace struct {
	env Environment
}

type Environment struct {
	values     map[string]Value
	namespaces map[string]Namespace
	parent     *Environment
}

func NewEnvironment(parent *Environment) *Environment {
	return &Environment{
		values:     make(map[string]Value),
		namespaces: make(map[string]Namespace),
		parent:     parent,
	}
}

func (env *Environment) Get(name string, allowPrivate bool) (interface{}, bool) {
	value, exists := env.values[name]
	if exists {
		if value.isPrivate && !allowPrivate {
			return nil, false
		}
		return value.data, true
	}

	if env.parent != nil {
		return env.parent.Get(name, allowPrivate)
	}
	return nil, false
}

func (env *Environment) Set(name string, value interface{}) bool {
	if val, exists := env.values[name]; exists {
		env.values[name] = Value{data: value, isPrivate: val.isPrivate}
		return true
	}
	if env.parent != nil {
		return env.parent.Set(name, value)
	}
	return false
}

func (env *Environment) Create(name string, value interface{}) bool {
	if _, exists := env.Get(name, true); exists {
		fmt.Println(exception.NewRuntimeError(fmt.Sprintf("variable %s already exists", name)))
		os.Exit(1)
		return false
	}
	env.values[name] = Value{data: value, isPrivate: true}
	return true
}

func (env *Environment) ListValues(includePrivate bool) map[string]interface{} {
	allValues := make(map[string]interface{})
	for key, value := range env.values {
		if !value.isPrivate || includePrivate {
			allValues[key] = value.data
		}
	}
	if env.parent != nil {
		parentValues := env.parent.ListValues(includePrivate)
		for key, value := range parentValues {
			if _, exists := allValues[key]; !exists {
				allValues[key] = value
			}
		}
	}
	return allValues
}

func (env *Environment) MakePublic(name string) bool {
	if value, exists := env.values[name]; exists {
		if value.isPrivate {
			value.isPrivate = false
			env.values[name] = value
		}
		return true
	}
	if env.parent != nil {
		return env.parent.MakePublic(name)
	}
	return false
}

func (env *Environment) GetNamespace(name string) (Environment, bool) {
	namespace, exists := env.namespaces[name]
	if exists {
		return namespace.env, true
	}

	if env.parent != nil {
		return env.parent.GetNamespace(name)
	}

	return Environment{}, false
}

func (env *Environment) CreateNamespace(name string, namespace *Environment) bool {
	if _, exists := env.GetNamespace(name); exists {
		fmt.Println(exception.NewRuntimeError(fmt.Sprintf("namespace %s already exists", name)))
		os.Exit(1)
		return false
	}
	env.namespaces[name] = Namespace{env: *namespace}
	return true
}

func (env *Environment) ListNamespaces() map[string]interface{} {
	allValues := make(map[string]interface{})
	for key, value := range env.namespaces {
		allValues[key] = value.env
	}

	return allValues
}
