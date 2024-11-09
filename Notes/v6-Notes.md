## Functions
- ***NOTE TO SELF***: Go through the entire function flow again to fully digest it. It is suprisingly complicated
- Part of the challenge is that functions in Monkey are both a series of statements and first class citizens that can be passed around and returned
- Other issues include
    - Passing params to methods
    - Return values to original control position
- Functions also evaluate to a value like every other literal. 
    - Functions are also set as global bindings using `OpSetGlobal` and using `OpGetGlobal`
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
    - `object.CompiledFunction` => hold the instructions of a compiled function and to pass them from the compiler to the VM as part of the bytecode, as a constant. Also holds the number of local bindings used by the function
    - `code.OpCall` => tell the VM to start executing the *object.CompiledFunction sitting on top of the stack.
    - `code.OpReturnValue` => tell the VM to return the value on top of the stack to the calling context and to resume execution there.
    - `code.OpReturn` => similar to code.OpReturnValue, except that there is no explicit value to return but an implicit vm.Null.

- The way we depict function scopes on the stack is by using the `scopes` value in the `Compiler` struct
    - `scopes` is just an array of `CompilationScope`s, which is a new struct, and we push new scopes into that "stack" when we start compiling a function
        - `CompilationScope` is a struct of instruction that are moved out from the `Compiler` struct
    - After compilation, it is popped off `scopes` and put into a new `*object.CompiledFunction`

- The start of the function call is that we put the function we want to call on the stack
    - `OpCall` is the function call instruction
    - The VM executed the function instructions and then pops the function of the stack, to replace it with the return value
        - Since functions are treated as literals as well `OpConstant` automatically pushed the function on to the stack. This is why the function needs to be implicitly popped to keep the stack clean 
        - If there is no return value, only the function is popped
        - That is, the popping is implicitly built into the VM


### Frames
- With how the VM handles execution of functions, non-linear execution has already been tried out once with jump instructions
    - The additional challenge here is that after first "jumping" to the function execution, we also need to "jump back" to the original location of the function call instruction to maintain order
- This is a ***`temporary storage`*** that lives for as long as a function call
- We can use something called a `frame` or also called a `call frame` or `stack frame`
    - This is the data structure that hold execution relevant information (whatever that means)
    - Frames are part of the stack itself, not separate
- Frames are where data such as the return address, arguments to the current function and local variables are stored
- As part of the stack, frames are easy to pop off after the function is done executing
- If assembly language was actually used to build a machine, then we would have to think about `memory addresses` in a much more real way
    - But since we are just building a virtual machine, we have more freedom in terms of how to create and store frames
- Hence, a `frame` for us is a struct built as such
```
type Frame struct {
    fn *object.CompiledFunction
    ip int // The instruction pointer within this particular frame
    basePointer // The value of the stack pointer before the execution of a function begins
}
```
- With the addition of frames, we have two options
    - Change the entire VM to use only frames when calling/executing functions
    - Change the VM with treats the `main` function like a function as well, which is what we will do
        - This is much simpler to learn and generally a more elegant solution since we already have so much of the VM already built
- This is good because our test suite can actually validate that none of the preexisting behaviours change with this new addition
- Essentially, the function is pushed on to both the stack and the frames
    - I am not fully certain of the nuances of the relationship between frames and stack

### Local bindings
- Multiple local bindings must be kept track of across multiple functions, alongside the global bindings
    - For this, another store is used
- The opcodes for local bindings will be similar to global opcodes: `OpSetLocal` and `OpGetLocal`
    - We'll also keep the width of this smaller at `1` instead, sort of an inherent and implicit limitations that a local binding set much never be greater than or equal to the global binding set
- The symbol table in the compiler is also extended to understand and take care of the different scopes. This is the "new store" that was previously alluded to in point 1
    - Essentially, we just tell the symbol when to enter or leave a scope
    - At the end of the day, both `let` statements and `identifiers` are the same in global or local scope, they are both framed as `ast.LetStatement` or `ast.Indentifier` from the AST. 
    - Within the new scope, the index restarts from 0
- The symbol table struct now has a field `Outer` which is another `*SymbolTable`
    - This `Outer` is the "parent" symbol table
    - For the global scope, `Outer` is nil
    - `Resolve()` is called recursively to check for a binding both within the current scope and any number of `Outer` scopes

### Within the VM
- For storing local bindings within the VM, we have a couple of options
    - Dynamically allocate a new slice with a function call that is used to store and retrieve local values with operands from the instruction as the index
    - Use the stack itself to store execution relevant data
- This project takes the second option. It is more complicated to understand but results in a lot of learning
    - Also saves memory allocation
    - This is the more commmon real world way of implementation as well, the common practice
- Working:
    - When an `OpCall` instruction is encountered, the current value of the stack pointer is stored for later use
    - The stack pointer is then increased by the number of locals used by the function to be executed
    - This results in a large empty space on the stack since we have just increased the stack pointer without actually pushing any values
        - Below this space is all the previously pushed values, and above is the function's workspace
    - The index of a local binding is the index to one of the holes on the stack 
        - The index essentially serves as the offset
    ![stackHoles](/Notes/assets/stackHoles.png)
- Since we also stored the original position of the stack pointer, the restoration/reset of the stack after execution is easy
- The number of local bindings used by a function is calculated at compiler time and passed on to the VM as part of the `*object.CompiledFunction` struct, as the field `NumLocals`
- The original location of the stack pointer is part of the `Frame` as field `basePointer`
    - `basePointer` is the conventional name for given for such a thing, ie, the pointer that points to the bottom of the stack of the current frame
    - Conventionally, it can also be called the `frame pointer`

### Dealing with arguments/params
- This is just a specialized way to declare new bindings in the local scope instead of using `let` statements
- The `OpCall` opcode is given an operand that indicates how many arguments are passed
- We just push the arguments on to the stack right after the function also has been pushed. So the offset for the function is also easily calculable, ie, just the number of arguments