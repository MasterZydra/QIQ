# Internal workings

## Pass by value
The interpreter uses the structs that implement `IRuntimeValue` as pointers.
The interpreter creates new instances for each new runtime value of type `Boolean`, `Integer`, `Floating` or `String`.
This results in "pass by value" behavior.

The above behavior results in a "pass by ref" for runtime values of type `Array`.
To solve this problem, a `deepCopy` function is used.
It creates a deep copy for a given array runtime value.

There are two places in the interpreter logic where a deepCopy must be created:
- variable assignments: `variable = deepCopy(value)`
- function calls: `argument = deepCopy(argument)`