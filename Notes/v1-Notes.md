## C-Monkey (Compiled Monkey) V1

- Still uses the frontend (lexer and parser) from Interpreted Monkey
- But the AST from the parser is compiled into bytecode and then executes on a VM
    - Involves building a compiler and a VM
- The architecture for the VM is a *stack machine*
    - Mostly for easier beginner understand/less performance concerns
- This implies that for an instruction to be correctly executed, the instruction ordering/stack arithmetic should be correct
- Opcodes for instructions must be decided for proper encoding-decoding
    - Our opcodes will be constants and this is determined at compile time
    - Compile time look up allows the compiler to just refer to the constants in the instructions
- ***This is the flow of operations***
    - When we come across an integer literal (a constant expression) while compiling, we’ll evaluate it and keep track of resulting `*object.Integer` by storing it in memory and assigning it a number. 
    - In the bytecode instructions we’ll refer to the `*object.Integer` by this number.
    - After we’re done compiling and pass the instructions to the VM for execution, we’ll also hand over all the constants we’ve found by putting them in a data structure – our *constant pool* – where the number that has been assigned to each constant can be used as an index to retrieve it.

### Basic compiler
- The first version of the compiler only needs to produce two `OpConstant` instructions to load values `1` and `2` on to the stack
    - First traverse the AST
    - Find `*ast.IntegerLiteral`
    - Eval into `*object.Integer`
    - Add to constant pool
    - Emit `OpConstant` instructions to reference constants
- The compiler must also be able to emit Bytecode instructions in human readable lang, instead of bytes for easier testing and debugging.


### VM
- Basic VM to deal with the `Bytecode` produced by the compiler
    - Fetch, decode and execute `OpConstant` instructions
    - At the end, the numbers should be pushed on to the VM's stack
- Lexed -> Parsed -> Compiled -> Passed into new instance of VM
- 