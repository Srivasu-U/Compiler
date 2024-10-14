## C-Monkey (Compiled Monkey) v2
- As a reminder, Monkey has three types of statements: `let`, `return` and `expressions`
    - `let` and `return` statements both use the value that is produced
    - `expression` statement is just a wrapper around other expressions so they can occur on their own. The value produced is it not reused
- So, if we end up having a lot of expression statements in a Monkey code snippet like `1; 2; 3;` then, all three values end up on the stack without ever being purged
    - This can cause memory problems since these values must be used
- To solve this, we can explicity tell the VM to pop elements when we estimate that we won't need them in the future
    - Just to make it clear, the way we "tell" the VM anything is through the compiler instructions.
    - This step must also be emitted to ensure proper tracking
    - The opcode for this is `OpPop`
- In order to check if the stack state is correct on the VM, we have to essentially do "Brother VM, this *should* have been on the stack right before this got popped off"

### Stack arithmetic and comparisons
- Stack arithmetic: `+, -, *, /`. `OpCodes` are
    - `+`: `OpAdd`
    - `-`: `OpSub`
    - `*`: `OpMul`
    - `/`: `OpDiv`
- Comparisons are:
    - `==`: `OpEqual`
    - `!=`: `OpNotEqual`
    - `>`: `OpGreaterThan` (This is also used for lesser than by just reordering the operands. For example, `3 > 5` or `5 < 3` produces the same result)
    - `>=`: `OpGreaterThanOrEqual` (Same as above with reordering)
- This reordering is only possible with compilation and not interpretation because of the formation of a tree structure during interpretation
    - Our compiler will essentially take any `<` expression and reorder the operands to give us a `>` expression
- Essentially, neither `<` nor `<=` even exist for our VM, it doesn't know the meaning of that symbol because we have not defined it
- Just as a reminder: ***This doesn't have to be executed like this***. 
    - We can just make `<` and `<=` opcodes like `OpLessThan` and `OpLessThanOrEqual` and we would be fine
- `vm.Run()` takes care of the actual arithmetic. Check `TestIntegerArithmetic` from `vm_test.go` to see the full range of execution

### Booleans
- With our compiler and VM, when encountering a boolean value (`true` or `false`), this must also be loaded on to the stack
- Boolean values are not treated as `OpConstants` to save on resources since these require lesser processing
    - So instead we have `OpTrue` and `OpFalse`
    - Both of these are like `OpPop` where they have no operands, just the opcode tells what is to be loaded on to the stack