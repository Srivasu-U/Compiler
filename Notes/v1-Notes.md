## C-Monkey (Compiled Monkey) V1

- Still uses the frontend (lexer and parser) from Interpreted Monkey
- But the AST from the parser is compiled into bytecode and then executes on a VM
    - Involves building a compiler and a VM
- The architecture for the VM is a *stack machine*
    - Mostly for easier beginner understand/less performance concerns
- This implies that for an instruction to be correctly executed, the instruction ordering/stack arithmetic should be correct
- Opcodes for instructions must be decided for proper encoding-decoding