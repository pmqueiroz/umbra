In Umbra, for is the only looping construct, offering two distinct forms: the initialized `for` statement and the conditional for statement.

### Initialized

The initialized `for` loop has an initializer, stop condition, and an optional step value, structured as 

```u title="loops.u"
for mut i num = 0, 10, 2 {
  io::printLn(i)
}
```
where `i` starts at `0`, runs until `10`, and increments by `2` on each iteration, resulting in

```
0
2
4
6
8
10
```
### Conditional

The conditional `for` loop, on the other hand, is structured with a condition, executing as long as the condition remains true.

```u title="loops.u"
if i < 10 {
  
}
```

If no condition is specified, the loop becomes infinite, running until explicitly broken with `break`.

```u title="loops.u"
mut i num = 0

for {
  if i > 10 {
    break
  }

  i = i + 1
}
```
