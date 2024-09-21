<img src=".github/logo.svg" width="150px" align="right"/>

# Welcome to source code of Umbra language interpreter
made in go!

![code](https://img.shields.io/github/languages/code-size/umbra-lang/umbra)
[![test-ci](https://github.com/pmqueiroz/umbra/actions/workflows/ci.yml/badge.svg)](https://github.com/pmqueiroz/umbra/actions/workflows/ci.yml)

TODO:
- Improve tokens structure

```go
type ColumnRange struct {
	from int
	to   int
}

type Loc struct {
	Line  int
	Range ColumnRange
}

type Token struct {
	Lexeme string
	Type   TokenType
	Loc
}
```
