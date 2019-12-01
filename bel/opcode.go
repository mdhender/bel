package bel

// opcode is the enum for an op code.
type opcode int

// WARNING: The order of op codes is important.
// Changing the order may break the eval code.

// enums for op codes.
const (
	// load must be the zero value so that the interpreter will default to loading
	opcLoad opcode = iota
	opcTopLevel0	// top level of the interpreter, primary prompt
	opcTopLevel1	// top level of the interpreter, continuation prompt
	opcRead // top level read expression
	opcValuePrint
	opcError0	// print initial error string
	opcError1	// print remainder of error
	opcReadSExpr
	opcP0List
	opcP1List
	opcInvalid
)