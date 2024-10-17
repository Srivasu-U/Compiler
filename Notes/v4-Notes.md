## Bindings
- Adding support for let statements and identifiers
    - Variable resolves to a value, essentially
- We use the stack we have to bind a value to an identifier
- The opcodes to be used are `OpSetGlobal` and `OpGetGlobal` for variables in the global environment (non-functional)
    - Let statement emits `OpSetGlobal`
    - Identifier compilation emits `OpGetGlobal`, aka, retrieval
- For ``` let x = 33;
let y = 66;
let z = x + y;```, diagrammatically, we have

![letStack](/Notes/assets/letStack.png)

- Within the VM, a slice is used as a "global store" with the operands from `OpSetGlobal` and `OpGetGlobal` as indices
    - Execution of `OpSetGlobal` -> pop topmost elem, save in global store at the index of the operand
    - Execution of `OpGetGlobal` -> Retrieve using index operand and push on to stack

### Symbol tables
- When compiling say `let x = 33;`, we assign a unique number to `x` because the operands to our opcodes can only be numbers and we cannot just pass `x` to be converted into bytecode
    - When we come across `x` again, we just use the previously assigned number
    - To keep it simple, we can use increasing number values, ie, the `x` would be `0` and every consecutive identifier after would be `1, 2, 3...`
    - We are using something called `symbol tables` (a struct), to choose what number we have to for our identifier (Look at `symbol_table.go`)
        - Symbol tables are data structures used in compilers to associate identifiers with information
        - Information can include location, scope, datatype, or anything else that is useful
        - In our case, we will use the `symbol tables` for scope and unique numbers
            - The table should associate identifiers in the global scope with a unique number (`define`)
            - Get previously associated number for an identifier (`resolve`)
            - An identifier is associated with a `symbol` and the symbol itself contains the information (a `symbol` is a struct for us, basically)
