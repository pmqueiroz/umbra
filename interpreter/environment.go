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

func (env *Environment) Set(name string, value interface{}) bool {
	if _, exists := env.values[name]; exists {
		env.values[name] = value
		return true
	}
	if env.parent != nil {
		return env.parent.Set(name, value)
	}
	return false
}

func (env *Environment) Create(name string, value interface{}) bool {
	if _, exists := env.Get(name); exists {
		return false
	}
	env.values[name] = value
	return true
}

func (env *Environment) ListValues() map[string]interface{} {
	allValues := make(map[string]interface{})
	for key, value := range env.values {
		allValues[key] = value
	}
	if env.parent != nil {
		parentValues := env.parent.ListValues()
		for key, value := range parentValues {
			if _, exists := allValues[key]; !exists {
				allValues[key] = value
			}
		}
	}
	return allValues
}
