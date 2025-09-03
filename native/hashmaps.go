package native

func del(args []interface{}) (interface{}, error) {
	hashmap := args[0].(map[interface{}]interface{})
	prop := args[1].(string)

	delete(hashmap, prop)

	return nil, nil
}

func keys(args []interface{}) (interface{}, error) {
	hashmap := args[0].(map[interface{}]interface{})
	keys := make([]interface{}, 0, len(hashmap))
	for k := range hashmap {
		keys = append(keys, k)
	}
	return keys, nil
}

var HashmapModule = InternalModule{
	symbols: map[string]InternalModuleFn{
		"delete": del,
		"keys":   keys,
	},
}
