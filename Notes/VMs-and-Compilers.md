## VMs and Compilers

- At the base level, VMs and Compilers are both essentially programming patterns and can have multiple implementations.
- Compilers general mean something that produce/outputs an executable file. 
    - Technically, compilers are just software that translate code from one language into another
    - Executing to an executable is basically translation from a high level language into native machine language
- Compilers and interpreters are similar in many ways, especially in terms of the lexer and the parser(together called the frontend). Both can make use of the generated AST
-But, while an interpreter directly implements whatever is in the AST, the compiler translates the AST into a language understood by the computer and the computer executes it
- The relevant points here are 
    - What is the target language of the source code generation?
    - How does the compiler translate the language?
    - What exactly gets generated?
- As usual, the answer is dependent on factors such as the architecture of the system, how the output of the compiler is used, performance needs, resource allocation and so on.
- General life cycle of the code being compiled into machine code is as such:  

![LifeCycle](/Notes/assets/lifecycle.png)  

- The AST is turned into an IR (Internal Representation) by an optimizer/compiler because an IR might be better optimized to be translated into the target language. The optimizations can be
    - Removal of dead code
    - Pre-calculations
    - Code restructuring (removing what doesn't need to be in a loop, for example)

### Workings of a CPU
- In regards to compilers, we specifically care about how the CPU and the memory interact
- How does the CPU know where to store and retrieve things that are in memory?
    - Program counter, part of CPU, keeps track of next instruction fetch, aka, numbers
- Computer memory is segmented into **words** [click here](https://en.wikipedia.org/wiki/Word_(computer_architecture))
    - A word is a base unit, smallest addressable region of memory
    - Standard is either 32 or 64 bits
- Addressing a single byte of a word requires offsets
- Data storage in memory is also optimized. One type of data is stored in one region and another type in another
    - A heap and a stack, for example
    - The stack is technically a *call stack*, basically to keep track of information to execute a program. Current execution, next execution, returns, arguments, variables et cetera.
- The CPU can also store data in process registers, which are faster to access. 
    - An x86-64 has only 16 registers each holding 64 bits of data
    - Registers are used to store small but frequently accessed data
    - The most common example is to store the mem address pointing to the top of the stack. This has a designated register called stack pointers

### VMs
- VMs are useful with compilers because we can create our own target language and then create a virtual machine that executes the target language. This is exactly what java does which is why it is platform independent
- *Side Note* : Interesting link about the [Von Neumann Arch](https://en.wikipedia.org/wiki/Von_Neumann_architecture) which is how computers are built  

![VonNeu](/Notes/assets/VonNeuArch.png)  

- Computers execute in **fetch-decode-execute cycle** steps
    - Fetch instruction from memory
    - Decode instruction
    - Execute instruction
- VMs have their own cycles, own stacks and counters, all built using software. 
- A VM in Javascript
```
let virtualMachine = function(program) {
    let programCounter = 0;
    let stack = [];
    let stackPointer = 0;
    while (programCounter < program.length) { // Run loop to execute instructions
        let currentInstruction = program[programCounter]; // Fetch current instruction
        switch (currentInstruction) { // This is called dispatching - selecting an implementation of an instruction.
        case PUSH:
            stack[stackPointer] = program[programCounter+1];
            stackPointer++;
            programCounter++;
            break;
        case ADD:
            right = stack[stackPointer-1]
            stackPointer--;
            left = stack[stackPointer-1]
            stackPointer--;
            stack[stackPointer] = left + right;
            stackPointer++;
            break;
        case MINUS:
            right = stack[stackPointer-1]
            stackPointer--;
            left = stack[stackPointer-1]
            stackPointer--;
            stack[stackPointer] = left - right;
            stackPointer++;
            break;
        }
        
        programCounter++;
    }

    console.log("stacktop: ", stack[stackPointer-1]);
}

let program = [
PUSH, 3,
PUSH, 4,
ADD,
PUSH, 5,
MINUS
];

virtualMachine(program); // Function call
```
- One of the main design choices of a VM is whether to make is a *stack machine* or a *register machine*.
- The stack machine is easier to build since everything happens on the stack
    - But this is also restrictive since everything *must* happen on the stack. Everything must be pushed and popped. 
    - So the point of stack machines/code is to be as efficient as possible with as little instructions.
- Register machines are harder to build since register are an additional thing to be taken care of, alongside a stack.
    - Instead of constantly popping and pushing on to the stack, registers can help speed up this process

## Bytecode
- Bytecode is called that way because the opcodes contained in each instruction are one byte in size
    - An “opcode” is the “operator” part of an instruction, sometimes also called “op”.
    - The PUSH for example is such an opcode, except that in the example it was a multi-byte string and not just one byte. 
    - In a proper implementation PUSH would just be the name that refers to an opcode, which itself is one byte wide. 
    - These names, like PUSH or POP, are called mnemonics.
- Assembly language is readable versions of bytecode

![opcode](/Notes/assets/opcode.png)
- Bytecodes are binary formatted.
- Operands don't have to be just one byte in size  
    - There are two representation of operands - *little endian* and *big endian*
    - Little endian - least significant bit first, stored in lowest mem address
    - Big endian - most significant bit first
- Representing the above opcode in actual bytecode, with PUSH as 1, ADD as 2 and ints stored in big endian  

![bytecode](/Notes/assets/bytecode.png)  

- Bytecode is always domain specific for the domain specific machine. Bytecode for Java built to run in a JVM wouldn't work for other languages, for example.