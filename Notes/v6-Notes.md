## Functions
- Part of the challenge is that functions in Monkey are both a series of statements and first class citizens that can be passed around and returned
- Other issues include
    - Passing params to methods
    - Return values to original control position
- Functions also evaluate to a value like every other literal. 
    - So we'll be treating functions as constants, at least from the point of view of a VM
    - That is, functions get compiled into a series of instructions and get added to the constant pool
    - Used in tandem with `OpConstant`
- In Interpred Monkey, we already have `object.Function` that holds function literals
    - That can now be an updated, new object that holds the functional bytecode called `CompiledFunction` -> this is the function literal
- For functions call, we introduce a new opcode called `OpCall` -> Mapped to `*ast.CallExpression`
- We get the function we want to call on the stack by using the `OpConstant` opcode and then issue an `OpCall` that tells the VM to execute the function on top of the stack
- Implicit and explicit return compiles to the same bytecode. These are the same
```
fn() { 5 + 10 }
fn() { return 5 + 10 }
```
- Opcode for return -> `OpReturnValue`
    - Returns the value on top of the stack 
- An edge case we have here: ***`functions that return no value`***, ie, everything is just wrapped up inside the function body
``` 
fn() {}
fn() { let a = 1; }
```
- Both are valid, and since function literals are also expressions they must produce a value. But these do not
    - So we can use `*object.Null` or `vm.Null` to represent this for us -> this would be the return value
    - Opcode for this is is `OpReturn` (maybe should have been `OpReturnNull`)

- ***Summarizing all this we have***
    - `object.CompiledFunction` => hold the instructions of a compiled function and to pass them from the compiler to the VM as part of the bytecode, as a constant.
    - `code.OpCall` => tell the VM to start executing the *object.CompiledFunction sitting on top of the stack.
    - `code.OpReturnValue` => tell the VM to return the value on top of the stack to the calling context and to resume execution there.
    - `code.OpReturn` => similar to code.OpReturnValue, except that there is no explicit value to return but an implicit vm.Null.