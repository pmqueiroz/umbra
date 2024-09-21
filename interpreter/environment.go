package interpreter

type Environment struct {
	values map[string]interface{}
	parent *Environment
}

func NewEnvironment(parent *Environment) *Environment {
	return &Environment{
		values: make(map[string]interface{}),
		parent: parent,
	}
}

func (env *Environment) Get(name string) (interface{}, bool) {
	value, exists := env.values[name]
	if !exists && env.parent != nil {
		return env.parent.Get(name)
	}
	return value, exists
}

func (env *Environment) Set(name string, value interface{}) {
	env.values[name] = value
}
