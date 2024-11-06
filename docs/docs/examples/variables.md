In Umbra, variables must always be explicitly typed, with the option to declare them as either `const` or `mut`, providing control over immutability and mutability in code.

### Mutable and constants

```u title="types.u"
const name str = "Peam"
mut age num = 24
const fool bool = true

mut person hashmap = {
  name: name,
  age: age,
}

const people arr = [person]
```

### Nullable

Additionally, by adding a `?` after the type, a variable can be made nullable, allowing it to hold either a value of the specified type or `null`.

```u title="types.u"
mut name str? = "Peam"
name = null
```

!!! info "Not initialized"
    Variables not initialized should be set as nullable when declaring
    ```u title="types.u"
    mut name str?
    ```

Next example: [Loops](/examples/loops)
