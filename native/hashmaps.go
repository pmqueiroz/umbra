package native

func del(args []interface{}) (interface{}, error) {
	hashmap := args[0].(map[interface{}]interface{})
	prop := args[1].(string)

	delete(hashmap, prop)

	return nil, nil
}

var HashmapModule = InternalModule{
	symbols: map[string]InternalModuleFn{
		"delete": del,
	},
}
