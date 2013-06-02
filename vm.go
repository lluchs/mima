// MIMA Virtual Machine for executing bytecode.

package mima

import (
	"errors"
	"fmt"
)

// Runs the bytcode, returning the resulting memory.
func (bytecode *Bytecode) Run() ([]uint32, error) {
	// Copy the bytecode as execution will probably change the memory.
	mem := bytecode.Mem[:]
	var Akku, IAR, IR uint32
	// Execution position.
	IAR = bytecode.Start
	// Akkumulator
	Akku = 0

	// Execute!
	for {
		// Fetch
		IR = mem[IAR]
		IAR++

		// Decode + Execute
		op := IR & 0xF00000
		op >>= 20
		arg := IR & 0x0FFFFF
		switch op {
		case 0: // LDC
			Akku = arg
		case 1: // LDV
			Akku = mem[arg]
		case 2: // STV
			mem[arg] = Akku
		case 3: // ADD
			Akku += mem[arg]
			Akku %= 0x1000000
		case 4: // AND
			Akku &= mem[arg]
		case 5: // OR
			Akku |= mem[arg]
		case 6: // XOR
			Akku ^= mem[arg]
		case 7: // EQL
			if Akku == mem[arg] {
				Akku = 0xFFFFFF
			} else {
				Akku = 0
			}
		case 8: // JMP
			IAR = arg
		case 9: // JMN
			if (Akku & 0x800000) != 0 {
				IAR = arg
			}
		case 0xF: // Additional commands
			op2 := (arg & 0xF0000) >> 16
			switch op2 {
			case 0: // HALT
				goto Out
			case 1: // NOT
				Akku = ^Akku
				Akku &= 0xFFFFFF
			case 2: // RAR
				rot := Akku & 1
				Akku >>= 1
				Akku += rot << 23
			default:
				return nil, errors.New(fmt.Sprintf("Invalid special OpCode F%X at 0x%06X", op2, IAR))
			}
		default:
			return nil, errors.New(fmt.Sprintf("Invalid OpCode %X at 0x%06X", op, IAR))
		}
	}

Out:
	return mem, nil
}
