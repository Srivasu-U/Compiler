## Dealing with conditional expressions
- The VM needs to execute different instructions based on the evaluation of a condition
- In the Interpreted Monkey, we could evaluate the condition and then choose to execute either the Consequence or the Alternative, based on the output
    - We could do this because we have the Consequence and Alternative nodes on our hand already from the AST
    - By with compilation, we don't have nodes or a tree structure because we flatten our bytecode
- The way we can achieve this conditional execution of code is by using *jumps*, also called *branch instructions*
    - Jumps are essentially instructions that tell the VM to change its instruction pointer
    ![Jumps](/Notes/assets/jumpsCompiler.png)  

    - `JUMP_IF_NOT_TRUE` is an instruction that tells the VM to jump to `OpConstant 4` if the boolean on the stack is not true, resulting in the execution of the alternative.
    - In case the consequence does get executed, then afterwards, `JUMP_NO_MATTER_WHAT` ensures that the consequence is never executed.
    - The location to jump to are just represented as numbers, as part of the operand for the jump instruction
        - The value of these is probably the index of the instruction the VM needs to jump to, aka, the offset.
        - This is an absolute offset, not a relative offset. With the offset, the diagram is as follows  

        ![jumpWithOffset](/Notes/assets/jumpWithOffset.png)

- The problem with jumps is that we don't know *where* to jump to because we wouldn't have compiled the consequence of alternative branch yet. So we don't know how many instructions to jump over. 
    - Figuring this out gives us the operand to the jump instruction
    - This figuring out is done by using `EmittedInstruction` struct in `compiler.go`, which keeps track of the `Position` of the last executed instruction.
        - We also have related methods with `EmittedInstruction` with `setLastInstruction(), lastInstructionIsPop() and removeLastPop()`
    - We can modify the operand to the jump instruction *after* compiling `node.Consequence` which gives us the right position to jump to
    - This is called ***back-patching***, common in single-pass compilers, ie, walks through the AST only once
    - Essentially, we will keep emitting `9999` as the operand to the jump instruction, until we figure out where to jump to. Then we'll go back and correct the offset
    - From `compiler.go`
    ```
    [...]
    case *ast.IfExpression:
		err := c.Compile(node.Condition)
		if err != nil {
			return err
		}

		jumpNotTruthyPos := c.emit(code.OpJumpNotTruthy, 9999)

		err = c.Compile(node.Consequence) // If block gets compiled regardless. Understand that this is compilation, and not actual execution
		if err != nil {
			return err
		}

		if c.lastInstructionIsPop() { // Last pop instruction of consequence expression is removed
			c.removeLastPop()
		}

		if node.Alternative == nil {
			// This basically gives us the offset after the execution of Consequence is done and where we want to jump to
			afterConsequencePos := len(c.instructions)

			// Going back and changing the operand, ie, the position where we jump to, instead of 9999
			c.changeOperand(jumpNotTruthyPos, afterConsequencePos)
		} else {
			// Emit an OpJump
			jumpPos := c.emit(code.OpJump, 9999)

			afterConsequencePos := len(c.instructions)
			c.changeOperand(jumpNotTruthyPos, afterConsequencePos)

			err := c.Compile(node.Alternative) // Alternative block, aka, else block is compiled
			if err != nil {
				return err
			}

			if c.lastInstructionIsPop() { // Last pop instruction of alternative expression is removed
				c.removeLastPop()
			}

			afterAlternativePos := len(c.instructions) // Replacing the operand of the jump instruction
			c.changeOperand(jumpPos, afterAlternativePos)
		}
    [...]
    ```
- We cannot `OpPop` an evaluated value after either the consequence or the alternative evaluates, because then we couldn't do something like `let x = if (5 > 3) { 5 } else { 3 }`.
    - If we evaluate to `5` and then `OpPop` is, `x` would be `nil`, not `5` which is not correct