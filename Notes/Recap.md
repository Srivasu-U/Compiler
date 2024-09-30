## Recap

- Starts off right where this was left off from *Interpreter* ([click here to go to the repo](https://github.com/Srivasu-U/Interpreter.git))
- General structure of the interpreter (for recollection) is as follows:
    - Lexer to tokenize
    - Parser (Pratt parser) to convert the tokens into an AST (Abstract Syntax Tree)
    - Evaluate whatever is parsed using an object system (everything is an object through wrappers)
    - This entire process can be viewed through `repl.go`
- The main reason to extend whatever we had written is to deal with performance and langauge maturity.
    - This is achieved by turning the tree-walking interpreter into a ***bytecode compiler and a VM that executes this bytecode***, much like Java and the JVM
    - This system is modular and bytecode compilers are very fast
    - The expectation is that the performance of Monkey with be faster by 3 fold at the end ( 3x faster )
- The folder `interpreted-monkey` is the base of our changes