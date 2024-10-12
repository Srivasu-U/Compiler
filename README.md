# The Compiled Monkey Language

- This repo consists of incremental development stages of a compiler for the Monkey programming language, as created through *Writing a Compiler in Go* by Thorsten Ball
- This is a direct extension from the [Interpreter](https://github.com/Srivasu-U/Interpreter), where Monkey was an interpreted language
- Each subdirectory has its own incremental development
    - `interpreted-monkey` is the final stage from the `Interpreter` code. This is the base on which the compiler is started to be built
    - `c-monkey-v1` has the beginnings of a compiler and VM
    - `Notes` has the relevant notes for each subdirectory. Written as I went through each development stage

### Execution of code
- Each subdirectory can be executed by cloning the repo and running
```
> cd c-monkey-<version-number>/src
> go run main.go
```
- The test cases can be executed by running
```
> go test ./lexer
> go test ./parser
> go test ./ast
> go test ./evaluator
> go test ./object
> go test ./compiler
> go test ./code
> go test ./vm
```
