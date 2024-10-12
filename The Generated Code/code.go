package main

import (
	"fmt"
)

// Go

func Emulate6502(pA *uint8, pX *uint8, pY *uint8, pSP *uint8, pPC *uint16, pSR *uint8, memory *[65536]uint8) bool {
	A := *pA
	X := *pX
	Y := *pY
	SP := *pSP
	PC := *pPC
	SR := *pSR
	SPSP := uint16(SP)
	XXXX := uint16(X)
	YYYY := uint16(Y)
	carry := SR & 1
	zero := (SR >> 1) & 1 /*interrupt := (SR>>2)&1; decimal := (SR>>3)&1;*/
	overflow := (SR >> 6) & 1
	negative := (SR >> 7) & 1
	cccc := uint16(carry)
	zzzz := uint16(zero) & 1
	vvvv := uint16(overflow) & 1
	nnnn := uint16(negative) & 1
	instruction := memory[PC&0xFFFF]
	if instruction == 0x00 { /* 00 LL = BRK #$LL */
		readWordAddress := uint16(0xFFFE)
		readWord := uint16(memory[readWordAddress]) + (uint16(memory[(readWordAddress+1)&0xFFFF]) << 8)
		*pSP = SP - 3
		newInterrupt := uint8(0x1) & 1
		*pSR = (*pSR & uint8(0b11111011)) + (newInterrupt << 2)
		*pPC = readWord
		writeByteAddress := (0x200 + (SPSP-2)&0x00FF) & 0xFFFF
		writeByte := uint8(SR | 0x10)
		memory[writeByteAddress] = writeByte
		writeWordAddress := (0x200 + (SPSP-1)&0xFF) & 0xFFFF
		writeWord := uint16(PC + 2)
		memory[writeWordAddress] = uint8(writeWord)
		memory[(writeWordAddress+1)&0xFFFF] = uint8(writeWord >> 8)
		return true
	}
	if instruction == 0x01 { /* 01 LL = ORA ($LL, X) */
		LL := memory[PC+1]
		zeroPageWordAddress := uint8(LL + X)
		zeroPageWord := uint16(memory[zeroPageWordAddress]) + (uint16(memory[(zeroPageWordAddress+1)&0xFF]) << 8)
		readByteAddress := uint16(zeroPageWord)
		readByte := memory[readByteAddress]
		operand := readByte
		result := A | operand
		*pA = result
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 2
		return true
	}
	if instruction == 0x05 { /* 05 LL = ORA $LL */
		LL := memory[PC+1]
		LLLL := uint16(int8(LL))
		readByteAddress := uint16(LLLL)
		readByte := memory[readByteAddress]
		operand := readByte
		result := A | operand
		*pA = result
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 2
		return true
	}
	if instruction == 0x06 { /* 06 LL = ASL $LL */
		LL := memory[PC+1]
		LLLL := uint16(int8(LL))
		readByteAddress := uint16(LLLL)
		readByte := memory[readByteAddress]
		operand := readByte
		result := operand << 1
		newCarry := uint8(operand>>7) & 1
		*pSR = (*pSR & uint8(0b11111110)) + newCarry
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 2
		writeByteAddress := (readByteAddress) & 0xFFFF
		writeByte := uint8(result)
		memory[writeByteAddress] = writeByte
		return true
	}
	if instruction == 0x08 { /* 08 = PHP  */
		*pSP = SP - 1
		*pPC = PC + 1
		writeByteAddress := (0x200 + SPSP) & 0xFFFF
		writeByte := uint8(SR | 0x30)
		memory[writeByteAddress] = writeByte
		return true
	}
	if instruction == 0x09 { /* 09 LL = ORA #$LL */
		LL := memory[PC+1]
		operand := LL
		result := A | operand
		*pA = result
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 2
		return true
	}
	if instruction == 0x0A { /* 0A = ASL A */
		operand := A
		result := operand << 1
		*pA = result
		newCarry := uint8(operand>>7) & 1
		*pSR = (*pSR & uint8(0b11111110)) + newCarry
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 1
		return true
	}
	if instruction == 0x0D { /* 0D LL HH = ORA $HHLL */
		LL := memory[PC+1]
		HHLL := (uint16(memory[PC+2]) << 8) + uint16(int8(LL))
		readByteAddress := uint16(HHLL)
		readByte := memory[readByteAddress]
		operand := readByte
		result := A | operand
		*pA = result
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 3
		return true
	}
	if instruction == 0x0E { /* 0E LL HH = ASL $HHLL */
		LL := memory[PC+1]
		HHLL := (uint16(memory[PC+2]) << 8) + uint16(int8(LL))
		readByteAddress := uint16(HHLL)
		readByte := memory[readByteAddress]
		operand := readByte
		result := operand << 1
		newCarry := uint8(operand>>7) & 1
		*pSR = (*pSR & uint8(0b11111110)) + newCarry
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 3
		writeByteAddress := (readByteAddress) & 0xFFFF
		writeByte := uint8(result)
		memory[writeByteAddress] = writeByte
		return true
	}
	if instruction == 0x10 { /* 10 LL = BPL $LL */
		LL := memory[PC+1]
		LLLL := uint16(int8(LL))
		*pPC = PC + 2 + (LLLL & (nnnn - 1))
		return true
	}
	if instruction == 0x11 { /* 11 LL = ORA ($LL), Y */
		LL := memory[PC+1]
		zeroPageWordAddress := uint8(LL)
		zeroPageWord := uint16(memory[zeroPageWordAddress]) + (uint16(memory[(zeroPageWordAddress+1)&0xFF]) << 8)
		readByteAddress := uint16(zeroPageWord + YYYY)
		readByte := memory[readByteAddress]
		operand := readByte
		result := A | operand
		*pA = result
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 2
		return true
	}
	if instruction == 0x15 { /* 15 LL = ORA $LL, X */
		LL := memory[PC+1]
		LLLL := uint16(int8(LL))
		readByteAddress := uint16((LLLL + XXXX) & 0x00FF)
		readByte := memory[readByteAddress]
		operand := readByte
		result := A | operand
		*pA = result
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 2
		return true
	}
	if instruction == 0x16 { /* 16 LL = ASL ($LL, X) */
		LL := memory[PC+1]
		zeroPageWordAddress := uint8(LL + X)
		zeroPageWord := uint16(memory[zeroPageWordAddress]) + (uint16(memory[(zeroPageWordAddress+1)&0xFF]) << 8)
		readByteAddress := uint16(zeroPageWord)
		readByte := memory[readByteAddress]
		operand := readByte
		result := operand << 1
		newCarry := uint8(operand>>7) & 1
		*pSR = (*pSR & uint8(0b11111110)) + newCarry
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 2
		writeByteAddress := (readByteAddress) & 0xFFFF
		writeByte := uint8(result)
		memory[writeByteAddress] = writeByte
		return true
	}
	if instruction == 0x18 { /* 18 = CLC  */
		newCarry := uint8(0x0) & 1
		*pSR = (*pSR & uint8(0b11111110)) + newCarry
		*pPC = PC + 1
		return true
	}
	if instruction == 0x19 { /* 19 LL HH = ORA $HHLL, Y */
		LL := memory[PC+1]
		HHLL := (uint16(memory[PC+2]) << 8) + uint16(int8(LL))
		readByteAddress := uint16(HHLL&0xFF00 + (HHLL+YYYY)&0x00FF)
		readByte := memory[readByteAddress]
		operand := readByte
		result := A | operand
		*pA = result
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 3
		return true
	}
	if instruction == 0x1D { /* 1D LL HH = ORA $HHLL, X */
		LL := memory[PC+1]
		HHLL := (uint16(memory[PC+2]) << 8) + uint16(int8(LL))
		readByteAddress := uint16(HHLL&0xFF00 + (HHLL+XXXX)&0x00FF)
		readByte := memory[readByteAddress]
		operand := readByte
		result := A | operand
		*pA = result
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 3
		return true
	}
	if instruction == 0x1E { /* 1E LL HH = ASL $HHLL, X */
		LL := memory[PC+1]
		HHLL := (uint16(memory[PC+2]) << 8) + uint16(int8(LL))
		readByteAddress := uint16(HHLL&0xFF00 + (HHLL+XXXX)&0x00FF)
		readByte := memory[readByteAddress]
		operand := readByte
		result := operand << 1
		newCarry := uint8(operand>>7) & 1
		*pSR = (*pSR & uint8(0b11111110)) + newCarry
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 3
		writeByteAddress := (readByteAddress) & 0xFFFF
		writeByte := uint8(result)
		memory[writeByteAddress] = writeByte
		return true
	}
	if instruction == 0x20 { /* 20 LL HH = JSR $HHLL */
		LL := memory[PC+1]
		HHLL := (uint16(memory[PC+2]) << 8) + uint16(int8(LL))
		*pSP = SP - 2
		*pPC = HHLL
		writeWordAddress := (0x200 + (SPSP-1)&0xFF) & 0xFFFF
		writeWord := uint16(PC + 3)
		memory[writeWordAddress] = uint8(writeWord)
		memory[(writeWordAddress+1)&0xFFFF] = uint8(writeWord >> 8)
		return true
	}
	if instruction == 0x21 { /* 21 LL = AND ($LL, X) */
		LL := memory[PC+1]
		zeroPageWordAddress := uint8(LL + X)
		zeroPageWord := uint16(memory[zeroPageWordAddress]) + (uint16(memory[(zeroPageWordAddress+1)&0xFF]) << 8)
		readByteAddress := uint16(zeroPageWord)
		readByte := memory[readByteAddress]
		operand := readByte
		result := A & operand
		*pA = result
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 2
		return true
	}
	if instruction == 0x24 { /* 24 LL = BIT $LL */
		LL := memory[PC+1]
		LLLL := uint16(int8(LL))
		readByteAddress := uint16(LLLL)
		readByte := memory[readByteAddress]
		operand := readByte
		result := A & operand
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newOverflow := uint8(operand>>6) & 1
		*pSR = (*pSR & uint8(0b10111111)) + (newOverflow << 6)
		newNegative := uint8(readByte>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 2
		return true
	}
	if instruction == 0x25 { /* 25 LL = AND $LL */
		LL := memory[PC+1]
		LLLL := uint16(int8(LL))
		readByteAddress := uint16(LLLL)
		readByte := memory[readByteAddress]
		operand := readByte
		result := A & operand
		*pA = result
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 2
		return true
	}
	if instruction == 0x26 { /* 26 LL = ROL $LL */
		LL := memory[PC+1]
		LLLL := uint16(int8(LL))
		readByteAddress := uint16(LLLL)
		readByte := memory[readByteAddress]
		operand := readByte
		result := (operand << 1) + carry
		newCarry := uint8(operand>>7) & 1
		*pSR = (*pSR & uint8(0b11111110)) + newCarry
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 2
		writeByteAddress := (readByteAddress) & 0xFFFF
		writeByte := uint8(result)
		memory[writeByteAddress] = writeByte
		return true
	}
	if instruction == 0x28 { /* 28 = PLP  */
		readByteAddress := uint16(0x200 + ((SPSP + 1) & 0x00FF))
		readByte := memory[readByteAddress]
		operand := readByte
		result := operand
		*pSP = SP + 1
		newCarry := uint8(result) & 1
		*pSR = (*pSR & uint8(0b11111110)) + newCarry
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newInterrupt := uint8(result>>2) & 1
		*pSR = (*pSR & uint8(0b11111011)) + (newInterrupt << 2)
		newDecimal := uint8(result>>3) & 1
		*pSR = (*pSR & uint8(0b11110111)) + (newDecimal << 3)
		newBreak := uint8(0x0) & 1
		*pSR = (*pSR & uint8(0b11101111)) + (newBreak << 4)
		newOverflow := uint8(result>>6) & 1
		*pSR = (*pSR & uint8(0b10111111)) + (newOverflow << 6)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 1
		return true
	}
	if instruction == 0x29 { /* 29 LL = AND #$LL */
		LL := memory[PC+1]
		operand := LL
		result := A & operand
		*pA = result
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 2
		return true
	}
	if instruction == 0x2A { /* 2A = ROL A */
		operand := A
		result := (operand << 1) + carry
		*pA = result
		newCarry := uint8(operand>>7) & 1
		*pSR = (*pSR & uint8(0b11111110)) + newCarry
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 1
		return true
	}
	if instruction == 0x2C { /* 2C LL HH = BIT $HHLL */
		LL := memory[PC+1]
		HHLL := (uint16(memory[PC+2]) << 8) + uint16(int8(LL))
		readByteAddress := uint16(HHLL)
		readByte := memory[readByteAddress]
		operand := readByte
		result := A & operand
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newOverflow := uint8(operand>>6) & 1
		*pSR = (*pSR & uint8(0b10111111)) + (newOverflow << 6)
		newNegative := uint8(operand>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 3
		return true
	}
	if instruction == 0x2D { /* 2D LL HH = AND $HHLL */
		LL := memory[PC+1]
		HHLL := (uint16(memory[PC+2]) << 8) + uint16(int8(LL))
		readByteAddress := uint16(HHLL)
		readByte := memory[readByteAddress]
		operand := readByte
		result := A & operand
		*pA = result
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 3
		return true
	}
	if instruction == 0x2E { /* 2E LL HH = ROL $HHLL */
		LL := memory[PC+1]
		HHLL := (uint16(memory[PC+2]) << 8) + uint16(int8(LL))
		readByteAddress := uint16(HHLL)
		readByte := memory[readByteAddress]
		operand := readByte
		result := (operand << 1) + carry
		newCarry := uint8(operand>>7) & 1
		*pSR = (*pSR & uint8(0b11111110)) + newCarry
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 3
		writeByteAddress := (readByteAddress) & 0xFFFF
		writeByte := uint8(result)
		memory[writeByteAddress] = writeByte
		return true
	}
	if instruction == 0x30 { /* 30 LL = BMI $LL */
		LL := memory[PC+1]
		LLLL := uint16(int8(LL))
		*pPC = PC + 2 + (LLLL & (nnnn - 1))
		return true
	}
	if instruction == 0x31 { /* 31 LL = AND ($LL), Y */
		LL := memory[PC+1]
		zeroPageWordAddress := uint8(LL)
		zeroPageWord := uint16(memory[zeroPageWordAddress]) + (uint16(memory[(zeroPageWordAddress+1)&0xFF]) << 8)
		readByteAddress := uint16(zeroPageWord + YYYY)
		readByte := memory[readByteAddress]
		operand := readByte
		result := A & operand
		*pA = result
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 2
		return true
	}
	if instruction == 0x35 { /* 35 LL = AND ($LL, X) */
		LL := memory[PC+1]
		zeroPageWordAddress := uint8(LL + X)
		zeroPageWord := uint16(memory[zeroPageWordAddress]) + (uint16(memory[(zeroPageWordAddress+1)&0xFF]) << 8)
		readByteAddress := uint16(zeroPageWord)
		readByte := memory[readByteAddress]
		operand := readByte
		result := A & operand
		*pA = result
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 2
		return true
	}
	if instruction == 0x36 { /* 36 LL = ROL $LL, X */
		LL := memory[PC+1]
		LLLL := uint16(int8(LL))
		readByteAddress := uint16((LLLL + XXXX) & 0x00FF)
		readByte := memory[readByteAddress]
		operand := readByte
		result := (operand << 1) + carry
		newCarry := uint8(operand>>7) & 1
		*pSR = (*pSR & uint8(0b11111110)) + newCarry
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 2
		writeByteAddress := (readByteAddress) & 0xFFFF
		writeByte := uint8(result)
		memory[writeByteAddress] = writeByte
		return true
	}
	if instruction == 0x38 { /* 38 = SEC  */
		newCarry := uint8(0x1) & 1
		*pSR = (*pSR & uint8(0b11111110)) + newCarry
		newNegative := uint8(0x1) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 1
		return true
	}
	if instruction == 0x39 { /* 39 LL HH = AND $HHLL, Y */
		LL := memory[PC+1]
		HHLL := (uint16(memory[PC+2]) << 8) + uint16(int8(LL))
		readByteAddress := uint16(HHLL&0xFF00 + (HHLL+YYYY)&0x00FF)
		readByte := memory[readByteAddress]
		operand := readByte
		result := A & operand
		*pA = result
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 3
		return true
	}
	if instruction == 0x3D { /* 3D LL HH = AND $HHLL, X */
		LL := memory[PC+1]
		HHLL := (uint16(memory[PC+2]) << 8) + uint16(int8(LL))
		readByteAddress := uint16(HHLL&0xFF00 + (HHLL+XXXX)&0x00FF)
		readByte := memory[readByteAddress]
		operand := readByte
		result := A & operand
		*pA = result
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 3
		return true
	}
	if instruction == 0x3E { /* 3E LL HH = ROL $HHLL, X */
		LL := memory[PC+1]
		HHLL := (uint16(memory[PC+2]) << 8) + uint16(int8(LL))
		readByteAddress := uint16(HHLL&0xFF00 + (HHLL+XXXX)&0x00FF)
		readByte := memory[readByteAddress]
		operand := readByte
		result := (operand << 1) + carry
		newCarry := uint8(operand>>7) & 1
		*pSR = (*pSR & uint8(0b11111110)) + newCarry
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 3
		writeByteAddress := (readByteAddress) & 0xFFFF
		writeByte := uint8(result)
		memory[writeByteAddress] = writeByte
		return true
	}
	if instruction == 0x40 { /* 40 = RTI  */
		readByteAddress := uint16(0x200 + (SPSP+1)&0x00FF)
		readByte := memory[readByteAddress]
		readWordAddress := uint16(0x200 + (SPSP+2)&0x00FF)
		readWord := uint16(memory[readWordAddress]) + (uint16(memory[(readWordAddress+1)&0xFFFF]) << 8)
		operand := readByte
		result := operand
		*pSP = SP + 3
		newCarry := uint8(result) & 1
		*pSR = (*pSR & uint8(0b11111110)) + newCarry
		newZero := uint8(result>>1) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newInterrupt := uint8(result>>2) & 1
		*pSR = (*pSR & uint8(0b11111011)) + (newInterrupt << 2)
		newDecimal := uint8(result>>3) & 1
		*pSR = (*pSR & uint8(0b11110111)) + (newDecimal << 3)
		newBreak := uint8(0x0) & 1
		*pSR = (*pSR & uint8(0b11101111)) + (newBreak << 4)
		newOverflow := uint8(result>>6) & 1
		*pSR = (*pSR & uint8(0b10111111)) + (newOverflow << 6)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = readWord
		return true
	}
	if instruction == 0x41 { /* 41 LL = EOR ($LL, X) */
		LL := memory[PC+1]
		zeroPageWordAddress := uint8(LL + X)
		zeroPageWord := uint16(memory[zeroPageWordAddress]) + (uint16(memory[(zeroPageWordAddress+1)&0xFF]) << 8)
		readByteAddress := uint16(zeroPageWord)
		readByte := memory[readByteAddress]
		operand := readByte
		result := A ^ operand
		*pA = result
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 2
		return true
	}
	if instruction == 0x45 { /* 45 LL = EOR $LL */
		LL := memory[PC+1]
		LLLL := uint16(int8(LL))
		readByteAddress := uint16(LLLL)
		readByte := memory[readByteAddress]
		operand := readByte
		result := A ^ operand
		*pA = result
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 2
		return true
	}
	if instruction == 0x46 { /* 46 LL = LSR $LL */
		LL := memory[PC+1]
		LLLL := uint16(int8(LL))
		readByteAddress := uint16(LLLL)
		readByte := memory[readByteAddress]
		operand := readByte
		result := operand >> 1
		*pA = result
		newCarry := uint8(operand) & 1
		*pSR = (*pSR & uint8(0b11111110)) + newCarry
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 2
		writeByteAddress := (readByteAddress) & 0xFFFF
		writeByte := uint8(result)
		memory[writeByteAddress] = writeByte
		return true
	}
	if instruction == 0x48 { /* 48 = PHA  */
		operand := A
		*pSP = SP - 1
		*pPC = PC + 1
		writeByteAddress := (0x200 + SPSP) & 0xFFFF
		writeByte := uint8(operand)
		memory[writeByteAddress] = writeByte
		return true
	}
	if instruction == 0x49 { /* 49 LL = EOR #$LL */
		LL := memory[PC+1]
		operand := LL
		result := A ^ operand
		*pA = result
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 2
		return true
	}
	if instruction == 0x4A { /* 4A = LSR A */
		operand := A
		result := operand >> 1
		*pA = result
		newCarry := uint8(operand) & 1
		*pSR = (*pSR & uint8(0b11111110)) + newCarry
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 1
		return true
	}
	if instruction == 0x4C { /* 4C LL HH = JMP $HHLL */
		LL := memory[PC+1]
		HHLL := (uint16(memory[PC+2]) << 8) + uint16(int8(LL))
		*pPC = HHLL
		return true
	}
	if instruction == 0x4D { /* 4D LL HH = EOR $HHLL */
		LL := memory[PC+1]
		HHLL := (uint16(memory[PC+2]) << 8) + uint16(int8(LL))
		readByteAddress := uint16(HHLL)
		readByte := memory[readByteAddress]
		operand := readByte
		result := A ^ operand
		*pA = result
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 3
		return true
	}
	if instruction == 0x4E { /* 4E LL HH = LSR $HHLL */
		LL := memory[PC+1]
		HHLL := (uint16(memory[PC+2]) << 8) + uint16(int8(LL))
		readByteAddress := uint16(HHLL)
		readByte := memory[readByteAddress]
		operand := readByte
		result := operand >> 1
		*pA = result
		newCarry := uint8(operand) & 1
		*pSR = (*pSR & uint8(0b11111110)) + newCarry
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 3
		writeByteAddress := (readByteAddress) & 0xFFFF
		writeByte := uint8(result)
		memory[writeByteAddress] = writeByte
		return true
	}
	if instruction == 0x50 { /* 50 LL = BVC #$LL */
		LL := memory[PC+1]
		LLLL := uint16(int8(LL))
		*pPC = PC + 2 + (LLLL & (vvvv - 1))
		return true
	}
	if instruction == 0x51 { /* 51 LL = EOR ($LL), Y */
		LL := memory[PC+1]
		zeroPageWordAddress := uint8(LL)
		zeroPageWord := uint16(memory[zeroPageWordAddress]) + (uint16(memory[(zeroPageWordAddress+1)&0xFF]) << 8)
		readByteAddress := uint16(zeroPageWord + YYYY)
		readByte := memory[readByteAddress]
		operand := readByte
		result := A ^ operand
		*pA = result
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 2
		return true
	}
	if instruction == 0x55 { /* 55 LL = EOR $LL, X */
		LL := memory[PC+1]
		LLLL := uint16(int8(LL))
		readByteAddress := uint16((LLLL + XXXX) & 0x00FF)
		readByte := memory[readByteAddress]
		operand := readByte
		result := A ^ operand
		*pA = result
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 2
		return true
	}
	if instruction == 0x56 { /* 56 LL = LSR $LL, X */
		LL := memory[PC+1]
		LLLL := uint16(int8(LL))
		readByteAddress := uint16((LLLL + XXXX) & 0x00FF)
		readByte := memory[readByteAddress]
		operand := readByte
		result := operand >> 1
		*pA = result
		newCarry := uint8(operand) & 1
		*pSR = (*pSR & uint8(0b11111110)) + newCarry
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 2
		writeByteAddress := (readByteAddress) & 0xFFFF
		writeByte := uint8(result)
		memory[writeByteAddress] = writeByte
		return true
	}
	if instruction == 0x58 { /* 58 = CLI  */
		newInterrupt := uint8(0x0) & 1
		*pSR = (*pSR & uint8(0b11111011)) + (newInterrupt << 2)
		*pPC = PC + 1
		return true
	}
	if instruction == 0x59 { /* 59 LL HH = EOR $HHLL, Y */
		LL := memory[PC+1]
		HHLL := (uint16(memory[PC+2]) << 8) + uint16(int8(LL))
		readByteAddress := uint16(HHLL&0xFF00 + (HHLL+YYYY)&0x00FF)
		readByte := memory[readByteAddress]
		operand := readByte
		result := A ^ operand
		*pA = result
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 3
		return true
	}
	if instruction == 0x5D { /* 5D LL HH = EOR $HHLL, X */
		LL := memory[PC+1]
		HHLL := (uint16(memory[PC+2]) << 8) + uint16(int8(LL))
		readByteAddress := uint16(HHLL&0xFF00 + (HHLL+XXXX)&0x00FF)
		readByte := memory[readByteAddress]
		operand := readByte
		result := A ^ operand
		*pA = result
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 3
		return true
	}
	if instruction == 0x5E { /* 5E LL HH = LSR $HHLL, X */
		LL := memory[PC+1]
		HHLL := (uint16(memory[PC+2]) << 8) + uint16(int8(LL))
		readByteAddress := uint16(HHLL&0xFF00 + (HHLL+XXXX)&0x00FF)
		readByte := memory[readByteAddress]
		operand := readByte
		result := operand >> 1
		*pA = result
		newCarry := uint8(operand) & 1
		*pSR = (*pSR & uint8(0b11111110)) + newCarry
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 3
		writeByteAddress := (readByteAddress) & 0xFFFF
		writeByte := uint8(result)
		memory[writeByteAddress] = writeByte
		return true
	}
	if instruction == 0x60 { /* 60 = RTS  */
		readWordAddress := uint16(0x200 + (SPSP+1)&0x00FF)
		readWord := uint16(memory[readWordAddress]) + (uint16(memory[(readWordAddress+1)&0xFFFF]) << 8)
		*pSP = SP + 2
		*pPC = readWord
		return true
	}
	if instruction == 0x61 { /* 61 LL = ADC ($LL, X) */
		LL := memory[PC+1]
		zeroPageWordAddress := uint8(LL + X)
		zeroPageWord := uint16(memory[zeroPageWordAddress]) + (uint16(memory[(zeroPageWordAddress+1)&0xFF]) << 8)
		readByteAddress := uint16(zeroPageWord)
		readByte := memory[readByteAddress]
		operand := readByte
		result := A + operand + carry
		*pA = result
		newCarry := uint8((result-A)>>7) & 1
		*pSR = (*pSR & uint8(0b11111110)) + newCarry
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newOverflow := uint8(((A^result)&(operand^result))>>7) & 1
		*pSR = (*pSR & uint8(0b10111111)) + (newOverflow << 6)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 2
		return true
	}
	if instruction == 0x65 { /* 65 LL = ADC $LL */
		LL := memory[PC+1]
		LLLL := uint16(int8(LL))
		readByteAddress := uint16(LLLL)
		readByte := memory[readByteAddress]
		operand := readByte
		result := A + operand + carry
		*pA = result
		newCarry := uint8((result-A)>>7) & 1
		*pSR = (*pSR & uint8(0b11111110)) + newCarry
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newOverflow := uint8(((A^result)&(operand^result))>>7) & 1
		*pSR = (*pSR & uint8(0b10111111)) + (newOverflow << 6)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 2
		return true
	}
	if instruction == 0x66 { /* 66 LL = ROR $LL */
		LL := memory[PC+1]
		LLLL := uint16(int8(LL))
		readByteAddress := uint16(LLLL)
		readByte := memory[readByteAddress]
		operand := readByte
		result := (operand >> 1) + (carry << 7)
		newCarry := uint8(operand) & 1
		*pSR = (*pSR & uint8(0b11111110)) + newCarry
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 2
		writeByteAddress := (readByteAddress) & 0xFFFF
		writeByte := uint8(result)
		memory[writeByteAddress] = writeByte
		return true
	}
	if instruction == 0x68 { /* 68 = PLA  */
		readByteAddress := uint16(0x200 + (SPSP+1)&0x00FF)
		readByte := memory[readByteAddress]
		operand := readByte
		result := operand
		*pA = result
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 1
		return true
	}
	if instruction == 0x69 { /* 69 LL = ADC #$LL */
		LL := memory[PC+1]
		operand := LL
		result := A + operand + carry
		*pA = result
		newCarry := uint8((result-A)>>7) & 1
		*pSR = (*pSR & uint8(0b11111110)) + newCarry
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newOverflow := uint8(((A^result)&(operand^result))>>7) & 1
		*pSR = (*pSR & uint8(0b10111111)) + (newOverflow << 6)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 2
		return true
	}
	if instruction == 0x6A { /* 6A = ROR A */
		operand := A
		result := (operand >> 1) + (carry << 7)
		*pA = result
		newCarry := uint8(operand) & 1
		*pSR = (*pSR & uint8(0b11111110)) + newCarry
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 1
		return true
	}
	if instruction == 0x6C { /* 6C LL HH = JMP ($HHLL) */
		LL := memory[PC+1]
		HHLL := (uint16(memory[PC+2]) << 8) + uint16(int8(LL))
		readWordAddress := uint16(HHLL)
		readWord := uint16(memory[readWordAddress]) + (uint16(memory[(readWordAddress+1)&0xFFFF]) << 8)
		*pPC = readWord
		return true
	}
	if instruction == 0x6D { /* 6D LL HH = ADC $HHLL */
		LL := memory[PC+1]
		HHLL := (uint16(memory[PC+2]) << 8) + uint16(int8(LL))
		readByteAddress := uint16(HHLL)
		readByte := memory[readByteAddress]
		operand := readByte
		result := A + operand + carry
		*pA = result
		newCarry := uint8((result-A)>>7) & 1
		*pSR = (*pSR & uint8(0b11111110)) + newCarry
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newOverflow := uint8(((A^result)&(operand^result))>>7) & 1
		*pSR = (*pSR & uint8(0b10111111)) + (newOverflow << 6)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 3
		return true
	}
	if instruction == 0x6E { /* 6E LL HH = ROR $HHLL */
		LL := memory[PC+1]
		HHLL := (uint16(memory[PC+2]) << 8) + uint16(int8(LL))
		readByteAddress := uint16(HHLL)
		readByte := memory[readByteAddress]
		operand := readByte
		result := (operand >> 1) + (carry << 7)
		newCarry := uint8(operand) & 1
		*pSR = (*pSR & uint8(0b11111110)) + newCarry
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 3
		writeByteAddress := (readByteAddress) & 0xFFFF
		writeByte := uint8(result)
		memory[writeByteAddress] = writeByte
		return true
	}
	if instruction == 0x70 { /* 70 LL = BVS $LL */
		LL := memory[PC+1]
		LLLL := uint16(int8(LL))
		*pPC = PC + 2 + (LLLL & (-vvvv))
		return true
	}
	if instruction == 0x71 { /* 71 LL = ADC ($LL),Y */
		LL := memory[PC+1]
		zeroPageWordAddress := uint8(LL)
		zeroPageWord := uint16(memory[zeroPageWordAddress]) + (uint16(memory[(zeroPageWordAddress+1)&0xFF]) << 8)
		readByteAddress := uint16(zeroPageWord + YYYY)
		readByte := memory[readByteAddress]
		operand := readByte
		result := A + operand + carry
		*pA = result
		newCarry := uint8((result-A)>>7) & 1
		*pSR = (*pSR & uint8(0b11111110)) + newCarry
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newOverflow := uint8(((A^result)&(operand^result))>>7) & 1
		*pSR = (*pSR & uint8(0b10111111)) + (newOverflow << 6)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 2
		return true
	}
	if instruction == 0x75 { /* 75 LL = ADC $LL, X */
		LL := memory[PC+1]
		LLLL := uint16(int8(LL))
		readByteAddress := uint16((LLLL + XXXX) & 0x00FF)
		readByte := memory[readByteAddress]
		operand := readByte
		result := A + operand + carry
		*pA = result
		newCarry := uint8((result-A)>>7) & 1
		*pSR = (*pSR & uint8(0b11111110)) + newCarry
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newOverflow := uint8(((A^result)&(operand^result))>>7) & 1
		*pSR = (*pSR & uint8(0b10111111)) + (newOverflow << 6)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 2
		return true
	}
	if instruction == 0x76 { /* 76 LL = ROR $LL, X */
		LL := memory[PC+1]
		LLLL := uint16(int8(LL))
		readByteAddress := uint16((LLLL + XXXX) & 0x00FF)
		readByte := memory[readByteAddress]
		operand := readByte
		result := (operand >> 1) + (carry << 7)
		newCarry := uint8(operand) & 1
		*pSR = (*pSR & uint8(0b11111110)) + newCarry
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 2
		writeByteAddress := (readByteAddress) & 0xFFFF
		writeByte := uint8(result)
		memory[writeByteAddress] = writeByte
		return true
	}
	if instruction == 0x78 { /* 78 = SEI  */
		newInterrupt := uint8(0x1) & 1
		*pSR = (*pSR & uint8(0b11111011)) + (newInterrupt << 2)
		*pPC = PC + 1
		return true
	}
	if instruction == 0x79 { /* 79 LL HH = ADC $HHLL, Y */
		LL := memory[PC+1]
		HHLL := (uint16(memory[PC+2]) << 8) + uint16(int8(LL))
		readByteAddress := uint16(HHLL&0xFF00 + (HHLL+YYYY)&0x00FF)
		readByte := memory[readByteAddress]
		operand := readByte
		result := A + operand + carry
		*pA = result
		newCarry := uint8((result-A)>>7) & 1
		*pSR = (*pSR & uint8(0b11111110)) + newCarry
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newOverflow := uint8(((A^result)&(operand^result))>>7) & 1
		*pSR = (*pSR & uint8(0b10111111)) + (newOverflow << 6)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 3
		return true
	}
	if instruction == 0x7D { /* 7D LL HH = ADC $HHLL, X */
		LL := memory[PC+1]
		HHLL := (uint16(memory[PC+2]) << 8) + uint16(int8(LL))
		readByteAddress := uint16(HHLL&0xFF00 + (HHLL+XXXX)&0x00FF)
		readByte := memory[readByteAddress]
		operand := readByte
		result := A + operand + carry
		*pA = result
		newCarry := uint8((result-A)>>7) & 1
		*pSR = (*pSR & uint8(0b11111110)) + newCarry
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newOverflow := uint8(((A^result)&(operand^result))>>7) & 1
		*pSR = (*pSR & uint8(0b10111111)) + (newOverflow << 6)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 3
		return true
	}
	if instruction == 0x7E { /* 7E LL HH = ROR $HHLL, X */
		LL := memory[PC+1]
		HHLL := (uint16(memory[PC+2]) << 8) + uint16(int8(LL))
		readByteAddress := uint16(HHLL&0xFF00 + (HHLL+XXXX)&0x00FF)
		readByte := memory[readByteAddress]
		operand := readByte
		result := (operand >> 1) + (carry << 7)
		newCarry := uint8(operand) & 1
		*pSR = (*pSR & uint8(0b11111110)) + newCarry
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 3
		writeByteAddress := (readByteAddress) & 0xFFFF
		writeByte := uint8(result)
		memory[writeByteAddress] = writeByte
		return true
	}
	if instruction == 0x81 { /* 81 LL = STA ($LL, X) */
		LL := memory[PC+1]
		zeroPageWordAddress := uint8(LL + X)
		zeroPageWord := uint16(memory[zeroPageWordAddress]) + (uint16(memory[(zeroPageWordAddress+1)&0xFF]) << 8)
		operand := A
		*pPC = PC + 2
		writeByteAddress := (zeroPageWord) & 0xFFFF
		writeByte := uint8(operand)
		memory[writeByteAddress] = writeByte
		return true
	}
	if instruction == 0x84 { /* 84 LL = STY $LL */
		LL := memory[PC+1]
		LLLL := uint16(int8(LL))
		operand := Y
		*pPC = PC + 2
		writeByteAddress := (LLLL) & 0xFFFF
		writeByte := uint8(operand)
		memory[writeByteAddress] = writeByte
		return true
	}
	if instruction == 0x85 { /* 85 LL = STA $LL */
		LL := memory[PC+1]
		LLLL := uint16(int8(LL))
		operand := A
		*pPC = PC + 2
		writeByteAddress := (LLLL) & 0xFFFF
		writeByte := uint8(operand)
		memory[writeByteAddress] = writeByte
		return true
	}
	if instruction == 0x86 { /* 86 LL = STX $LL */
		LL := memory[PC+1]
		LLLL := uint16(int8(LL))
		operand := X
		*pPC = PC + 2
		writeByteAddress := (LLLL) & 0xFFFF
		writeByte := uint8(operand)
		memory[writeByteAddress] = writeByte
		return true
	}
	if instruction == 0x88 { /* 88 = DEY  */
		operand := Y
		result := operand - 1
		*pY = result
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 1
		return true
	}
	if instruction == 0x8A { /* 8A = TXA  */
		operand := X
		result := operand
		*pA = result
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 1
		return true
	}
	if instruction == 0x8C { /* 8C LL HH = STY $HHLL */
		LL := memory[PC+1]
		HHLL := (uint16(memory[PC+2]) << 8) + uint16(int8(LL))
		operand := Y
		*pPC = PC + 3
		writeByteAddress := (HHLL) & 0xFFFF
		writeByte := uint8(operand)
		memory[writeByteAddress] = writeByte
		return true
	}
	if instruction == 0x8D { /* 8D LL HH = STA $HHLL */
		LL := memory[PC+1]
		HHLL := (uint16(memory[PC+2]) << 8) + uint16(int8(LL))
		operand := A
		*pPC = PC + 3
		writeByteAddress := (HHLL) & 0xFFFF
		writeByte := uint8(operand)
		memory[writeByteAddress] = writeByte
		return true
	}
	if instruction == 0x8E { /* 8E LL HH = STX $HHLL */
		LL := memory[PC+1]
		HHLL := (uint16(memory[PC+2]) << 8) + uint16(int8(LL))
		operand := X
		*pPC = PC + 3
		writeByteAddress := (HHLL) & 0xFFFF
		writeByte := uint8(operand)
		memory[writeByteAddress] = writeByte
		return true
	}
	if instruction == 0x90 { /* 90 LL = BCC $LL */
		LL := memory[PC+1]
		LLLL := uint16(int8(LL))
		*pPC = PC + 2 + (LLLL & (cccc - 1))
		return true
	}
	if instruction == 0x91 { /* 91 LL = STA ($LL), Y */
		LL := memory[PC+1]
		zeroPageWordAddress := uint8(LL)
		zeroPageWord := uint16(memory[zeroPageWordAddress]) + (uint16(memory[(zeroPageWordAddress+1)&0xFF]) << 8)
		operand := A
		*pPC = PC + 2
		writeByteAddress := (zeroPageWord + YYYY) & 0xFFFF
		writeByte := uint8(operand)
		memory[writeByteAddress] = writeByte
		return true
	}
	if instruction == 0x94 { /* 94 LL = STY $LL, X */
		LL := memory[PC+1]
		LLLL := uint16(int8(LL))
		operand := Y
		*pPC = PC + 2
		writeByteAddress := ((LLLL + XXXX) & 0x00FF) & 0xFFFF
		writeByte := uint8(operand)
		memory[writeByteAddress] = writeByte
		return true
	}
	if instruction == 0x95 { /* 95 LL = STA $LL, X */
		LL := memory[PC+1]
		LLLL := uint16(int8(LL))
		operand := A
		*pPC = PC + 2
		writeByteAddress := ((LLLL + XXXX) & 0x00FF) & 0xFFFF
		writeByte := uint8(operand)
		memory[writeByteAddress] = writeByte
		return true
	}
	if instruction == 0x96 { /* 96 LL = STX $LL, Y */
		LL := memory[PC+1]
		zeroPageWordAddress := uint8(LL)
		zeroPageWord := uint16(memory[zeroPageWordAddress]) + (uint16(memory[(zeroPageWordAddress+1)&0xFF]) << 8)
		operand := X
		*pPC = PC + 2
		writeByteAddress := (zeroPageWord + YYYY) & 0xFFFF
		writeByte := uint8(operand)
		memory[writeByteAddress] = writeByte
		return true
	}
	if instruction == 0x98 { /* 98 = TYA  */
		operand := Y
		result := operand
		*pA = result
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 1
		return true
	}
	if instruction == 0x99 { /* 99 LL HH = STA $HHLL, Y */
		LL := memory[PC+1]
		HHLL := (uint16(memory[PC+2]) << 8) + uint16(int8(LL))
		operand := A
		*pPC = PC + 3
		writeByteAddress := (HHLL&0xFF00 + (HHLL+YYYY)&0x00FF) & 0xFFFF
		writeByte := uint8(operand)
		memory[writeByteAddress] = writeByte
		return true
	}
	if instruction == 0x9A { /* 9A = TXS  */
		operand := X
		result := operand
		*pSP = result
		*pPC = PC + 1
		return true
	}
	if instruction == 0x9D { /* 9D LL HH = STA $HHLL, X */
		LL := memory[PC+1]
		HHLL := (uint16(memory[PC+2]) << 8) + uint16(int8(LL))
		operand := A
		*pPC = PC + 3
		writeByteAddress := (HHLL&0xFF00 + (HHLL+XXXX)&0x00FF) & 0xFFFF
		writeByte := uint8(operand)
		memory[writeByteAddress] = writeByte
		return true
	}
	if instruction == 0xA0 { /* A0 LL = LDY #$LL */
		LL := memory[PC+1]
		operand := LL
		result := operand
		*pY = result
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 2
		return true
	}
	if instruction == 0xA1 { /* A1 LL = LDA ($LL, X) */
		LL := memory[PC+1]
		zeroPageWordAddress := uint8(LL + X)
		zeroPageWord := uint16(memory[zeroPageWordAddress]) + (uint16(memory[(zeroPageWordAddress+1)&0xFF]) << 8)
		readByteAddress := uint16(zeroPageWord)
		readByte := memory[readByteAddress]
		operand := readByte
		result := operand
		*pA = result
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 2
		return true
	}
	if instruction == 0xA2 { /* A2 LL = LDX #$LL */
		LL := memory[PC+1]
		operand := LL
		result := operand
		*pX = result
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 2
		return true
	}
	if instruction == 0xA4 { /* A4 LL = LDY $LL */
		LL := memory[PC+1]
		LLLL := uint16(int8(LL))
		readByteAddress := uint16(LLLL)
		readByte := memory[readByteAddress]
		operand := readByte
		result := operand
		*pY = result
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 2
		return true
	}
	if instruction == 0xA5 { /* A5 LL = LDA $LL */
		LL := memory[PC+1]
		LLLL := uint16(int8(LL))
		readByteAddress := uint16(LLLL)
		readByte := memory[readByteAddress]
		operand := readByte
		result := operand
		*pA = result
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 2
		return true
	}
	if instruction == 0xA6 { /* A6 LL = LDX $LL */
		LL := memory[PC+1]
		LLLL := uint16(int8(LL))
		readByteAddress := uint16(LLLL)
		readByte := memory[readByteAddress]
		operand := readByte
		result := operand
		*pX = result
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 2
		return true
	}
	if instruction == 0xA8 { /* A8 = TAY  */
		operand := A
		result := operand
		*pY = result
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 1
		return true
	}
	if instruction == 0xA9 { /* A9 LL = LDA #$LL */
		LL := memory[PC+1]
		operand := LL
		result := operand
		*pA = result
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 2
		return true
	}
	if instruction == 0xAA { /* AA = TAX  */
		operand := A
		result := operand
		*pX = result
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 1
		return true
	}
	if instruction == 0xAC { /* AC LL HH = LDY $HHLL */
		LL := memory[PC+1]
		HHLL := (uint16(memory[PC+2]) << 8) + uint16(int8(LL))
		readByteAddress := uint16(HHLL)
		readByte := memory[readByteAddress]
		operand := readByte
		result := operand
		*pY = result
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 3
		return true
	}
	if instruction == 0xAD { /* AD LL HH = LDA $HHLL */
		LL := memory[PC+1]
		HHLL := (uint16(memory[PC+2]) << 8) + uint16(int8(LL))
		readByteAddress := uint16(HHLL)
		readByte := memory[readByteAddress]
		operand := readByte
		result := operand
		*pA = result
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 3
		return true
	}
	if instruction == 0xAE { /* AE LL HH = LDX $HHLL */
		LL := memory[PC+1]
		HHLL := (uint16(memory[PC+2]) << 8) + uint16(int8(LL))
		readByteAddress := uint16(HHLL)
		readByte := memory[readByteAddress]
		operand := readByte
		result := operand
		*pX = result
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 3
		return true
	}
	if instruction == 0xB0 { /* B0 LL = BCS $LL */
		LL := memory[PC+1]
		LLLL := uint16(int8(LL))
		*pPC = PC + 2 + (LLLL & (-cccc))
		return true
	}
	if instruction == 0xB1 { /* B1 LL = LDA ($LL), Y */
		LL := memory[PC+1]
		zeroPageWordAddress := uint8(LL)
		zeroPageWord := uint16(memory[zeroPageWordAddress]) + (uint16(memory[(zeroPageWordAddress+1)&0xFF]) << 8)
		readByteAddress := uint16(zeroPageWord + YYYY)
		readByte := memory[readByteAddress]
		operand := readByte
		result := operand
		*pA = result
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 2
		return true
	}
	if instruction == 0xB4 { /* B4 LL = LDY $LL, X */
		LL := memory[PC+1]
		LLLL := uint16(int8(LL))
		readByteAddress := uint16((LLLL + XXXX) & 0x00FF)
		readByte := memory[readByteAddress]
		operand := readByte
		result := operand
		*pY = result
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 2
		return true
	}
	if instruction == 0xB5 { /* B5 LL = LDA $LL, X */
		LL := memory[PC+1]
		LLLL := uint16(int8(LL))
		readByteAddress := uint16((LLLL + XXXX) & 0x00FF)
		readByte := memory[readByteAddress]
		operand := readByte
		result := operand
		*pA = result
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 2
		return true
	}
	if instruction == 0xB6 { /* B6 LL = LDX $LL, Y */
		LL := memory[PC+1]
		LLLL := uint16(int8(LL))
		readByteAddress := uint16((LLLL + XXXX) & 0x00FF)
		readByte := memory[readByteAddress]
		operand := readByte
		result := operand
		*pX = result
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 2
		return true
	}
	if instruction == 0xB8 { /* B8 = CLV  */
		newOverflow := uint8(0x0) & 1
		*pSR = (*pSR & uint8(0b10111111)) + (newOverflow << 6)
		*pPC = PC + 1
		return true
	}
	if instruction == 0xB9 { /* B9 LL HH = LDA $HHLL, Y */
		LL := memory[PC+1]
		HHLL := (uint16(memory[PC+2]) << 8) + uint16(int8(LL))
		readByteAddress := uint16(HHLL&0xFF00 + (HHLL+YYYY)&0x00FF)
		readByte := memory[readByteAddress]
		operand := readByte
		result := operand
		*pA = result
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 3
		return true
	}
	if instruction == 0xBA { /* BA = TSX  */
		operand := SP
		result := operand
		*pX = result
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 1
		return true
	}
	if instruction == 0xBC { /* BC LL HH = LDY $HHLL, X */
		LL := memory[PC+1]
		HHLL := (uint16(memory[PC+2]) << 8) + uint16(int8(LL))
		readByteAddress := uint16(HHLL&0xFF00 + (HHLL+XXXX)&0x00FF)
		readByte := memory[readByteAddress]
		operand := readByte
		result := operand
		*pY = result
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 3
		return true
	}
	if instruction == 0xBD { /* BD LL HH = LDA $HHLL, X */
		LL := memory[PC+1]
		HHLL := (uint16(memory[PC+2]) << 8) + uint16(int8(LL))
		readByteAddress := uint16(HHLL&0xFF00 + (HHLL+XXXX)&0x00FF)
		readByte := memory[readByteAddress]
		operand := readByte
		result := operand
		*pA = result
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 3
		return true
	}
	if instruction == 0xBE { /* BE LL HH = LDX $HHLL, Y */
		LL := memory[PC+1]
		HHLL := (uint16(memory[PC+2]) << 8) + uint16(int8(LL))
		readByteAddress := uint16(HHLL&0xFF00 + (HHLL+YYYY)&0x00FF)
		readByte := memory[readByteAddress]
		operand := readByte
		result := operand
		*pX = result
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 3
		return true
	}
	if instruction == 0xC0 { /* C0 LL = CPY #$LL */
		LL := memory[PC+1]
		operand := LL
		result := Y - operand
		newCarry := uint8(1-(result>>7)) & 1
		*pSR = (*pSR & uint8(0b11111110)) + newCarry
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 2
		return true
	}
	if instruction == 0xC1 { /* C1 LL = CMP ($LL, X) */
		LL := memory[PC+1]
		zeroPageWordAddress := uint8(LL + X)
		zeroPageWord := uint16(memory[zeroPageWordAddress]) + (uint16(memory[(zeroPageWordAddress+1)&0xFF]) << 8)
		readByteAddress := uint16(zeroPageWord)
		readByte := memory[readByteAddress]
		operand := readByte
		result := A - operand
		newCarry := uint8(1-(result>>7)) & 1
		*pSR = (*pSR & uint8(0b11111110)) + newCarry
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 2
		return true
	}
	if instruction == 0xC4 { /* C4 LL = CPY $LL */
		LL := memory[PC+1]
		LLLL := uint16(int8(LL))
		readByteAddress := uint16(LLLL)
		readByte := memory[readByteAddress]
		operand := readByte
		result := Y - operand
		newCarry := uint8(1-(result>>7)) & 1
		*pSR = (*pSR & uint8(0b11111110)) + newCarry
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 2
		return true
	}
	if instruction == 0xC5 { /* C5 LL = CMP $LL */
		LL := memory[PC+1]
		LLLL := uint16(int8(LL))
		readByteAddress := uint16(LLLL)
		readByte := memory[readByteAddress]
		operand := readByte
		result := A - operand
		newCarry := uint8(1-(result>>7)) & 1
		*pSR = (*pSR & uint8(0b11111110)) + newCarry
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 2
		return true
	}
	if instruction == 0xC6 { /* C6 LL = DEC $LL */
		LL := memory[PC+1]
		LLLL := uint16(int8(LL))
		readByteAddress := uint16(LLLL)
		readByte := memory[readByteAddress]
		operand := readByte
		result := operand - 1
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 2
		writeByteAddress := (readByteAddress) & 0xFFFF
		writeByte := uint8(result)
		memory[writeByteAddress] = writeByte
		return true
	}
	if instruction == 0xC8 { /* C8 = INY  */
		operand := Y
		result := operand + 1
		*pY = result
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 1
		return true
	}
	if instruction == 0xC9 { /* C9 LL = CMP #$LL */
		LL := memory[PC+1]
		operand := LL
		result := A - operand
		newCarry := uint8(1-(result>>7)) & 1
		*pSR = (*pSR & uint8(0b11111110)) + newCarry
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 2
		return true
	}
	if instruction == 0xCA { /* CA = DEX  */
		operand := X
		result := operand - 1
		*pX = result
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 1
		return true
	}
	if instruction == 0xCC { /* CC LL HH = CPY $HHLL */
		LL := memory[PC+1]
		HHLL := (uint16(memory[PC+2]) << 8) + uint16(int8(LL))
		readByteAddress := uint16(HHLL)
		readByte := memory[readByteAddress]
		operand := readByte
		result := Y - operand
		newCarry := uint8(1-(result>>7)) & 1
		*pSR = (*pSR & uint8(0b11111110)) + newCarry
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 3
		return true
	}
	if instruction == 0xCD { /* CD LL HH = CMP $HHLL */
		LL := memory[PC+1]
		HHLL := (uint16(memory[PC+2]) << 8) + uint16(int8(LL))
		readByteAddress := uint16(HHLL)
		readByte := memory[readByteAddress]
		operand := readByte
		result := A - operand
		newCarry := uint8(1-(result>>7)) & 1
		*pSR = (*pSR & uint8(0b11111110)) + newCarry
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 3
		return true
	}
	if instruction == 0xCE { /* CE LL HH = DEC $HHLL */
		LL := memory[PC+1]
		HHLL := (uint16(memory[PC+2]) << 8) + uint16(int8(LL))
		readByteAddress := uint16(HHLL)
		readByte := memory[readByteAddress]
		operand := readByte
		result := operand - 1
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 3
		writeByteAddress := (readByteAddress) & 0xFFFF
		writeByte := uint8(result)
		memory[writeByteAddress] = writeByte
		return true
	}
	if instruction == 0xD0 { /* D0 LL = BNE $LL */
		LL := memory[PC+1]
		LLLL := uint16(int8(LL))
		*pPC = PC + 2 + (LLLL & (zzzz - 1))
		return true
	}
	if instruction == 0xD1 { /* D1 LL = CMP ($LL), Y */
		LL := memory[PC+1]
		zeroPageWordAddress := uint8(LL)
		zeroPageWord := uint16(memory[zeroPageWordAddress]) + (uint16(memory[(zeroPageWordAddress+1)&0xFF]) << 8)
		readByteAddress := uint16(zeroPageWord + YYYY)
		readByte := memory[readByteAddress]
		operand := readByte
		result := A - operand
		newCarry := uint8(1-(result>>7)) & 1
		*pSR = (*pSR & uint8(0b11111110)) + newCarry
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 2
		return true
	}
	if instruction == 0xD5 { /* D5 LL = CMP $LL, X */
		LL := memory[PC+1]
		LLLL := uint16(int8(LL))
		readByteAddress := uint16((LLLL + XXXX) & 0x00FF)
		readByte := memory[readByteAddress]
		operand := readByte
		result := A - operand
		newCarry := uint8(1-(result>>7)) & 1
		*pSR = (*pSR & uint8(0b11111110)) + newCarry
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 2
		return true
	}
	if instruction == 0xD6 { /* D6 LL = DEC $LL, X */
		LL := memory[PC+1]
		LLLL := uint16(int8(LL))
		readByteAddress := uint16((LLLL + XXXX) & 0x00FF)
		readByte := memory[readByteAddress]
		operand := readByte
		result := operand - 1
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 2
		writeByteAddress := (readByteAddress) & 0xFFFF
		writeByte := uint8(result)
		memory[writeByteAddress] = writeByte
		return true
	}
	if instruction == 0xD8 { /* D8 = CLD  */
		newDecimal := uint8(0x0) & 1
		*pSR = (*pSR & uint8(0b11110111)) + (newDecimal << 3)
		*pPC = PC + 1
		return true
	}
	if instruction == 0xD9 { /* D9 LL HH = CMP $HHLL, Y */
		LL := memory[PC+1]
		HHLL := (uint16(memory[PC+2]) << 8) + uint16(int8(LL))
		readByteAddress := uint16(HHLL&0xFF00 + (HHLL+YYYY)&0x00FF)
		readByte := memory[readByteAddress]
		operand := readByte
		result := A - operand
		newCarry := uint8(1-(result>>7)) & 1
		*pSR = (*pSR & uint8(0b11111110)) + newCarry
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 3
		return true
	}
	if instruction == 0xDD { /* DD LL HH = CMP $HHLL, X */
		LL := memory[PC+1]
		HHLL := (uint16(memory[PC+2]) << 8) + uint16(int8(LL))
		readByteAddress := uint16(HHLL&0xFF00 + (HHLL+XXXX)&0x00FF)
		readByte := memory[readByteAddress]
		operand := readByte
		result := A - operand
		newCarry := uint8(1-(result>>7)) & 1
		*pSR = (*pSR & uint8(0b11111110)) + newCarry
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 3
		return true
	}
	if instruction == 0xDE { /* DE LL HH = DEC $HHLL, X */
		LL := memory[PC+1]
		HHLL := (uint16(memory[PC+2]) << 8) + uint16(int8(LL))
		readByteAddress := uint16((HHLL + XXXX) & 0x00FF)
		readByte := memory[readByteAddress]
		operand := readByte
		result := operand - 1
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 3
		writeByteAddress := (readByteAddress) & 0xFFFF
		writeByte := uint8(result)
		memory[writeByteAddress] = writeByte
		return true
	}
	if instruction == 0xE0 { /* E0 LL = CPX #$LL */
		LL := memory[PC+1]
		operand := LL
		result := X - operand
		newCarry := uint8(1-(result>>7)) & 1
		*pSR = (*pSR & uint8(0b11111110)) + newCarry
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 2
		return true
	}
	if instruction == 0xE1 { /* E1 LL = SBC ($LL, X) */
		LL := memory[PC+1]
		zeroPageWordAddress := uint8(LL + X)
		zeroPageWord := uint16(memory[zeroPageWordAddress]) + (uint16(memory[(zeroPageWordAddress+1)&0xFF]) << 8)
		readByteAddress := uint16(zeroPageWord)
		readByte := memory[readByteAddress]
		operand := readByte
		result := A - operand - (carry ^ 1)
		*pA = result
		newCarry := uint8(1-(result>>7)) & 1
		*pSR = (*pSR & uint8(0b11111110)) + newCarry
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newOverflow := uint8(((A^result)&(operand^result))>>7) & 1
		*pSR = (*pSR & uint8(0b10111111)) + (newOverflow << 6)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 2
		return true
	}
	if instruction == 0xE4 { /* E4 LL = CPX $LL */
		LL := memory[PC+1]
		LLLL := uint16(int8(LL))
		readByteAddress := uint16(LLLL)
		readByte := memory[readByteAddress]
		operand := readByte
		result := X - operand
		newCarry := uint8(1-(result>>7)) & 1
		*pSR = (*pSR & uint8(0b11111110)) + newCarry
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 2
		return true
	}
	if instruction == 0xE5 { /* E5 LL = SBC $LL */
		LL := memory[PC+1]
		LLLL := uint16(int8(LL))
		readByteAddress := uint16(LLLL)
		readByte := memory[readByteAddress]
		operand := readByte
		result := A - operand - (carry ^ 1)
		*pA = result
		newCarry := uint8(1-(result>>7)) & 1
		*pSR = (*pSR & uint8(0b11111110)) + newCarry
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newOverflow := uint8(((A^result)&(operand^result))>>7) & 1
		*pSR = (*pSR & uint8(0b10111111)) + (newOverflow << 6)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 2
		return true
	}
	if instruction == 0xE6 { /* E6 LL = INC $LL */
		LL := memory[PC+1]
		LLLL := uint16(int8(LL))
		readByteAddress := uint16(LLLL)
		readByte := memory[readByteAddress]
		operand := readByte
		result := operand + 1
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 2
		writeByteAddress := (readByteAddress) & 0xFFFF
		writeByte := uint8(result)
		memory[writeByteAddress] = writeByte
		return true
	}
	if instruction == 0xE8 { /* E8 = INX  */
		operand := X
		result := operand + 1
		*pX = result
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 1
		return true
	}
	if instruction == 0xE9 { /* E9 LL = SBC #$LL */
		LL := memory[PC+1]
		operand := LL
		result := A - operand - (carry ^ 1)
		*pA = result
		newCarry := uint8(1-(result>>7)) & 1
		*pSR = (*pSR & uint8(0b11111110)) + newCarry
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newOverflow := uint8(((A^result)&(operand^result))>>7) & 1
		*pSR = (*pSR & uint8(0b10111111)) + (newOverflow << 6)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 2
		return true
	}
	if instruction == 0xEA { /* EA = NOP  */
		*pPC = PC + 1
		return true
	}
	if instruction == 0xEC { /* EC LL HH = CPX $HHLL */
		LL := memory[PC+1]
		HHLL := (uint16(memory[PC+2]) << 8) + uint16(int8(LL))
		readByteAddress := uint16(HHLL)
		readByte := memory[readByteAddress]
		operand := readByte
		result := X - operand
		newCarry := uint8(1-(result>>7)) & 1
		*pSR = (*pSR & uint8(0b11111110)) + newCarry
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 3
		return true
	}
	if instruction == 0xED { /* ED LL HH = SBC $HHLL */
		LL := memory[PC+1]
		HHLL := (uint16(memory[PC+2]) << 8) + uint16(int8(LL))
		readByteAddress := uint16(HHLL)
		readByte := memory[readByteAddress]
		operand := readByte
		result := A - operand - (carry ^ 1)
		*pA = result
		newCarry := uint8(1-(result>>7)) & 1
		*pSR = (*pSR & uint8(0b11111110)) + newCarry
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newOverflow := uint8(((A^result)&(operand^result))>>7) & 1
		*pSR = (*pSR & uint8(0b10111111)) + (newOverflow << 6)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 3
		return true
	}
	if instruction == 0xEE { /* EE LL HH = INC $HHLL */
		LL := memory[PC+1]
		HHLL := (uint16(memory[PC+2]) << 8) + uint16(int8(LL))
		readByteAddress := uint16(HHLL)
		readByte := memory[readByteAddress]
		operand := readByte
		result := operand + 1
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 3
		writeByteAddress := (readByteAddress) & 0xFFFF
		writeByte := uint8(result)
		memory[writeByteAddress] = writeByte
		return true
	}
	if instruction == 0xF0 { /* F0 LL = BEQ $LL */
		LL := memory[PC+1]
		LLLL := uint16(int8(LL))
		*pPC = PC + 2 + (LLLL & (-zzzz))
		return true
	}
	if instruction == 0xF1 { /* F1 LL = SBC ($LL), Y */
		LL := memory[PC+1]
		zeroPageWordAddress := uint8(LL + X)
		zeroPageWord := uint16(memory[zeroPageWordAddress]) + (uint16(memory[(zeroPageWordAddress+1)&0xFF]) << 8)
		readByteAddress := uint16(zeroPageWord + YYYY)
		readByte := memory[readByteAddress]
		operand := readByte
		result := A - operand - (carry ^ 1)
		*pA = result
		newCarry := uint8(1-(result>>7)) & 1
		*pSR = (*pSR & uint8(0b11111110)) + newCarry
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newOverflow := uint8(((A^result)&(operand^result))>>7) & 1
		*pSR = (*pSR & uint8(0b10111111)) + (newOverflow << 6)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 2
		return true
	}
	if instruction == 0xF5 { /* F5 LL = SBC $LL, X */
		LL := memory[PC+1]
		LLLL := uint16(int8(LL))
		readByteAddress := uint16((LLLL + XXXX) & 0x00FF)
		readByte := memory[readByteAddress]
		operand := readByte
		result := A - operand - (carry ^ 1)
		*pA = result
		newCarry := uint8(1-(result>>7)) & 1
		*pSR = (*pSR & uint8(0b11111110)) + newCarry
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newOverflow := uint8(((A^result)&(operand^result))>>7) & 1
		*pSR = (*pSR & uint8(0b10111111)) + (newOverflow << 6)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 2
		return true
	}
	if instruction == 0xF6 { /* F6 LL = INC $LL, X */
		LL := memory[PC+1]
		LLLL := uint16(int8(LL))
		readByteAddress := uint16((LLLL + XXXX) & 0x00FF)
		readByte := memory[readByteAddress]
		operand := readByte
		result := operand + 1
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 2
		writeByteAddress := (readByteAddress) & 0xFFFF
		writeByte := uint8(result)
		memory[writeByteAddress] = writeByte
		return true
	}
	if instruction == 0xF8 { /* F8 = SED  */
		newDecimal := uint8(0x1) & 1
		*pSR = (*pSR & uint8(0b11110111)) + (newDecimal << 3)
		*pPC = PC + 1
		return true
	}
	if instruction == 0xF9 { /* F9 LL HH = SBC $HHLL, Y */
		LL := memory[PC+1]
		HHLL := (uint16(memory[PC+2]) << 8) + uint16(int8(LL))
		readByteAddress := uint16(HHLL&0xFF00 + (HHLL+YYYY)&0x00FF)
		readByte := memory[readByteAddress]
		operand := readByte
		result := A - operand - (carry ^ 1)
		*pA = result
		newCarry := uint8(1-(result>>7)) & 1
		*pSR = (*pSR & uint8(0b11111110)) + newCarry
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newOverflow := uint8(((A^result)&(operand^result))>>7) & 1
		*pSR = (*pSR & uint8(0b10111111)) + (newOverflow << 6)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 3
		return true
	}
	if instruction == 0xFD { /* FD LL HH = SBC $HHLL, X */
		LL := memory[PC+1]
		HHLL := (uint16(memory[PC+2]) << 8) + uint16(int8(LL))
		readByteAddress := uint16(HHLL&0xFF00 + (HHLL+XXXX)&0x00FF)
		readByte := memory[readByteAddress]
		operand := readByte
		result := A - operand - (carry ^ 1)
		*pA = result
		newCarry := uint8(1-(result>>7)) & 1
		*pSR = (*pSR & uint8(0b11111110)) + newCarry
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newOverflow := uint8(((A^result)&(operand^result))>>7) & 1
		*pSR = (*pSR & uint8(0b10111111)) + (newOverflow << 6)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 3
		return true
	}
	if instruction == 0xFE { /* FE LL HH = INC $HHLL, X */
		LL := memory[PC+1]
		HHLL := (uint16(memory[PC+2]) << 8) + uint16(int8(LL))
		readByteAddress := uint16(HHLL&0xFF00 + (HHLL+XXXX)&0x00FF)
		readByte := memory[readByteAddress]
		operand := readByte
		result := operand + 1
		newZero := uint8(1-((result|-result)>>7)) & 1
		*pSR = (*pSR & uint8(0b11111101)) + (newZero << 1)
		newNegative := uint8(result>>7) & 1
		*pSR = (*pSR & uint8(0b01111111)) + (newNegative << 7)
		*pPC = PC + 3
		writeByteAddress := (readByteAddress) & 0xFFFF
		writeByte := uint8(result)
		memory[writeByteAddress] = writeByte
		return true
	}
	return false
}

//
//func Assemble6502( memory uint8[], address uint16, code string) bool {
//
//	//  Break code into lines, remove all whitespace from lines, filter out non-blank lines to get final result
//	lines := strings.FieldsFunc(strings.ReplaceAll(code, " ", ""), func(r rune) bool { return r == '\n' && strings.TrimSpace(string(r)) == "" })
//
//	labelMap := make(map[string]uint16)
//
//	//  Loop through each line, determining whether it is a comment
//	for i, line := range lines {
//		if strings.HasPrefix(line, ".") {
//			labelMap[line] = address
//		}
//		else {
//			//  It is an instruction; deal with here
//
//			//  If text is just BRK then it is instruction 0
//
//
//
//			address = address + 1
//		}
//	}
//
//	return true
//}

//  This prints out a line with the status of the chip
func Status6502(A uint8, X uint8, Y uint8, SP uint8, PC uint16, SR uint8) {

	//  Fish out the flags
	carry := SR & 1
	zero := (SR >> 1) & 1
	interrupt := (SR >> 2) & 1
	decimal := (SR >> 3) & 1
	overflow := (SR >> 6) & 1
	negative := (SR >> 7) & 1

	// Print the registers as hex and the carry flag
	fmt.Printf("A: $%02X  X: $%02X  Y: $%02X  SP: $%02X  PC: $%04X  SR: $%02X  |  C:%d Z:%d I:%d D:%d V:%d N:%d\n",
		A, X, Y, SP, PC, SR, carry, zero, interrupt, decimal, overflow, negative)
}

// Function to print a memory page in hex and ASCII, starting at startAddress
func MemoryPage6502(startAddress uint16, memory [65536]uint8) {
	for i := 0; i < 256; i += 16 {
		// Calculate the current line's address
		address := startAddress + uint16(i)

		// Print the address
		fmt.Printf("%04X: ", address)

		// Print 16 bytes in hexadecimal
		for j := 0; j < 16; j++ {
			fmt.Printf("%02X ", memory[address+uint16(j)])
		}

		// Print a separator for the ASCII section
		fmt.Print(" | ")

		// Print 16 bytes as ASCII characters, replacing non-printable ones with '.'
		for j := 0; j < 16; j++ {
			byteValue := memory[address+uint16(j)]
			if byteValue >= 32 && byteValue <= 126 {
				fmt.Printf("%c", byteValue)
			} else {
				fmt.Print(".")
			}
		}

		// Print a newline at the end of the line
		fmt.Println()
	}
}

func main() {

	//  Prepare the CPU
	var A uint8 = 0
	var X uint8 = 0
	var Y uint8 = 0
	var SP uint8 = 0
	var PC uint16 = 0x400
	var SR uint8 = 0
	var memory [65536]uint8

	//  Write the program
	program := []uint8{
		0xa9, 0x00, // LDA #$00
		0xa2, 0x00, // LDX #$00
		0x9d, 0x00, 0xa0, // loop: STA $A000, X
		0x18,       // CLC
		0x69, 0x01, // ADC #$01
		0xe8,       // INX
		0xd0, 0xf7} // BNE loop
	copy(memory[PC:], program)

	//  Display the initial state
	fmt.Println("")
	fmt.Println("INITIALISING")
	fmt.Println("=\n")
	fmt.Println("Status:")
	fmt.Println("-\n")
	Status6502(A, X, Y, SP, PC, SR)
	fmt.Println("")
	fmt.Println("Program:")
	fmt.Println("-\n")
	MemoryPage6502(PC, memory)

	//  Emulate
	fmt.Println("")
	fmt.Println("EMULATING...")
	fmt.Println("=\n")
	var emulateOk = true
	var counter = 0
	steps := 1000000
	for emulateOk && counter < steps && PC != 0x0000 {
		emulateOk = Emulate6502(&A, &X, &Y, &SP, &PC, &SR, &memory)
		counter++
		if emulateOk {
			Status6502(A, X, Y, SP, PC, SR)
		}
	}

	//  Display result
	if emulateOk {
		fmt.Println("")
		fmt.Println("COMPLETE")
		fmt.Println("=\n")
		fmt.Println("Status:")
		fmt.Println("-\n")
		Status6502(A, X, Y, SP, PC, SR)
		fmt.Println("")
		fmt.Println("Page 0xA0:")
		fmt.Println("-\n")
		MemoryPage6502(0xA000, memory)

	} else {
		fmt.Println("")
		fmt.Println("ERROR: Emulation halted unexpectedly.")
	}
}
