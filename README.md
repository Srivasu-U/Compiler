# The Compiled Monkey Language

- This repo consists of incremental development stages of a compiler for the Monkey programming language, as created through *Writing a Compiler in Go* by Thorsten Ball
- This is a direct extension from the [Interpreter](https://github.com/Srivasu-U/Interpreter), where Monkey was an interpreted language
- Each subdirectory has its own incremental development (ie, it has everything from the previous iterations and adds something new as mentioned)
    - The latest version always has all the work in progress and may break. Right now, v5 works perfectly fine and the version supports as follows:
    - `interpreted-monkey` is the final stage from the `Interpreter` code. This is the base on which the compiler is started to be built
    - `c-monkey-v1` has the beginnings of a compiler and VM. 
        - Support for addition operation only, and the foundation of bytecode.
    - `c-monkey-v2` enables the compilation of execution of expressions, infix and postfix.
        - Support for other arithmetic operations
    - `c-monkey-v3` includes supports for conditionals, ie, `if...else...` blocks, and how to execute the consequence and alternatives
    - `c-monkey-v4` consists of support for let statements and accessing identifiers during execution (only in global scope)
    - `c-monkey-v5` supports strings, arrays and hashes
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
