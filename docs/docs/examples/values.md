Umbra has various value types including strings (`str`), numbers (`num`), booleans (`bool`), etc. Here are a few basic examples.

Strings, which can be added together with +.

```u title="values.u"
io::println("umbra" + "lang")
```
Numbers.

```u title="values.u"
io::println("1+1 =", 1+1)
io::println("7/3 =", 7/3)
```

Booleans, with boolean operators as you’d expect.

```u title="values.u"
io::println(true and false)
io::println(true or false)
io::println(!true)
```

```sh
$ umbra values.u
# umbralang
# 1+1 = 2
# 7/3 = 2.3333333333333335
# false
# true
# false
```

Next example: [Variables](/examples/variables)
