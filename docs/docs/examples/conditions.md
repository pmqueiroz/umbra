In Umbra, conditional logic is handled using the `if`, `else if`, and `else` constructs, which allow branching based on `bool` expressions. The basic structure begins with `if`, followed by a condition enclosed in parentheses and the block of code to execute if the condition is `true`:

```u title="conditions.md"
if true {
  # Code executes if the condition is true
}
```

You can add an `else if` block to check additional conditions if the first one is `false`:

```u title="conditions.md"
else if false {
  # Code executes if the first condition is false, but this one is true
}
```

Finally, an optional `else` block can be added to handle all cases where none of the previous conditions are met:

```u title="conditions.md"
else {
  # Code executes if all previous conditions are false
}
```

This structure allows for clean and readable conditional branching, enabling complex decision-making within the code.
