// MIMA Assembler

package mima

import (
	"errors"
)

type Bytecode struct {
	Start uint32
	Mem []uint32
}

// Assembles the calling program.
func (program *Program) Assemble() (*Bytecode, error) {
	bytecode := new(Bytecode)
	// Find the program entry point.
	if start, ok := program.Marks["START"]; ok {
		bytecode.Start = start
	} else {
		return nil, errors.New("No START label given.")
	}
	// Initialize the memory.
	bytecode.Mem = make([]uint32, 0xFFFFF)

	// Write all instructions to their positions.
	for position, instruction := range program.Instructions {
		var result uint32 = 0
		var arg uint32 = 0
		// Try to resolve the argument.
		if val, ok := program.Marks[instruction.Argument]; ok {
			arg = val
		} else if val, err := instruction.ParseArgument(); err == nil {
			arg = val
		}
		// Save opcode/argument.
		saveArg := func(opcode uint32) {
			result = opcode << 20
			result += arg
		}
		// Translate the instructions.
		switch instruction.Op {
		case "DS":   result = arg
		case "LDC":  saveArg(0)
		case "LDV":  saveArg(1)
		case "STV":  saveArg(2)
		case "ADD":  saveArg(3)
		case "AND":  saveArg(4)
		case "OR":   saveArg(5)
		case "XOR":  saveArg(6)
		case "EQL":  saveArg(7)
		case "JMP":  saveArg(8)
		case "JMN":  saveArg(9)
		case "HALT": result = 0xF00000
		case "NOT":  result = 0xF10000
		case "RAR":  result = 0xF20000

		default:
			return nil, errors.New("Invalid instruction: " + instruction.Op)
		}
		bytecode.Mem[position] = result
	}
	return bytecode, nil
}
