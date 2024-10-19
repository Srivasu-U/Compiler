## String support
- Includes support for string concatenation
- String literals are constants, ie, the value doesn't change between compile and runtime
    - This means we can just turn string values into `*object.String` at compile time and add them to constant pool in `compiler.Bytecode`

## Arrays and Hashes
- Includes support for indexing
- Array values are not static between compile time and run time, ie, they can change, since array values can be expressions during compile time that evaluate into values at run time.
    - This means that only at runtime, we can determine what the array values are
    - An optimized compiler could pre-compute these values
- So the approach to be taken is that the bytecode instructions we send to the VM when compiling an array, tells the VM how to build its own array
    - That is, we do not compile the values and add it to the constant pool
- Done by defining a new opcode: `OpArray`, with one operand stating the number of elements in an array
    - When we then compile an `*ast.ArrayLiteral`, we first compile all of its elements. 
    - Since these are `ast.Expressions`, compiling them results in instructions that leave N values on the VM’s stack, where N is the number of elements in the array literal. 
    - Then, we’re going to emit an OpArray instruction with the operand being N, the number of elements. This is the end of the compilation.
    - When the VM then executes the OpArray instruction it takes the N elements off the stack, builds an `*object.Array` out of them, and pushes that on to the stack. 
- The approach is similar for hash, but the underlying structure is a map instead of a slice
- As a refresher, hashes work like this:
    - `object.Hash` has a Pairs field that contains a `map[HashKey]HashPair`. 
    - A HashKey can be created by calling the `HashKey` method of an `object.Hashable`, an interface that `*object.String`, `*object.Boolean` and `*object.Integer` implement. 
    - A `HashPair` then has a Key and a Value field, both containing an `object.Object`. This is where the real key and the real value are stored.
- For indexing, we have the `OpIndex` opcode, that takes the top two values of the stack: topmost being the index and top-1 being the object to be indexed
    - `OpIndex` pops both, performs the retrieval and pushes the new value on to the stack
