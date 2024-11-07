Umbra has various value types including strings (`str`), numbers (`num`), booleans (`bool`), etc. Here are a few basic examples.

Strings, which can be added together with +.

```u title="values.u"
io::printLn("umbra" + "lang")
```
Numbers.

```u title="values.u"
io::printLn("1+1 =", 1+1)
io::printLn("7/3 =", 7/3)
```

Booleans, with boolean operators as you’d expect.

```u title="values.u"
io::printLn(true and false)
io::printLn(true or false)
io::printLn(!true)
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
