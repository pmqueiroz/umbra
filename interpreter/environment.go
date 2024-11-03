package interpreter

import (
	"fmt"
	"os"

	"github.com/pmqueiroz/umbra/exception"
	"github.com/pmqueiroz/umbra/tokens"
)

type Variable struct {
	data     interface{}
	dataType tokens.TokenType
	private  bool
	nullable bool
}

type Namespace struct {
	env Environment
}

type Environment struct {
	values     map[string]Variable
	namespaces map[string]Namespace
	parent     *Environment
}

func NewEnvironment(parent *Environment) *Environment {
	return &Environment{
		values:     make(map[string]Variable),
		namespaces: make(map[string]Namespace),
		parent:     parent,
	}
}

func (env *Environment) Get(name string, allowPrivate bool) (Variable, bool) {
	value, exists := env.values[name]
	if exists {
		if value.private && !allowPrivate {
			return Variable{}, false
		}
		return value, true
	}

	if env.parent != nil {
		return env.parent.Get(name, allowPrivate)
	}
	return Variable{}, false
}

func (env *Environment) Set(name string, value interface{}) bool {
	if val, exists := env.values[name]; exists {
		env.values[name] = Variable{data: value, dataType: val.dataType, private: val.private, nullable: val.nullable}
		return true
	}
	if env.parent != nil {
		return env.parent.Set(name, value)
	}
	return false
}

func (env *Environment) Create(name string, value interface{}, dataType tokens.TokenType, nullable bool) bool {
	if _, exists := env.Get(name, true); exists {
		fmt.Println(exception.NewRuntimeError(fmt.Sprintf("variable %s already exists", name)))
		os.Exit(1)
		return false
	}
	env.values[name] = Variable{data: value, dataType: dataType, private: true, nullable: nullable}
	return true
}

func (env *Environment) ListValues(includePrivate bool) map[string]interface{} {
	allValues := make(map[string]interface{})
	for key, value := range env.values {
		if !value.private || includePrivate {
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
		if value.private {
			value.private = false
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
