# 6502 THE SPREADSHEET (The Repository)

**Under Construction**

## Introduction

All 6502 instructions with their actions completely unrolled so you don't have to pick through detailed documentation to understand what they do or how to emulate them accurately.

This is used to code generate emulators in a variety of languages directly in the spreadsheet as columns which can be copy/pasted as source code.

The spreadsheet is under construction and is at:
https://docs.google.com/spreadsheets/d/1rCwO6CoIXiCqkJUaZ5Q9W5QZNqmhLpXaAhSthJPuV8I/edit?gid=0#gid=0

![image](https://github.com/user-attachments/assets/31bceb69-3d3f-4adc-9b40-a0d8ea07790d)

## STATUS: Alpha / under construction

* DONE: Get the list of instructions in place
* DONE: Get the first pass of the functionality of instructions in place
* DOING: Checking the functionality of each instruction
* DOING: Adding languages to codegen and fixing bugs
* DOING: Refining the reversibility information
* TODO: Make test cases and put it through its paces

## How it works

For each instruction there is a set of expressions for each of the necessary memory reads, register updates, flag updates, and memory writes.

These expressions are provided in a simple portable form, just using the symbols + - << >> & | () and are thus valid expressions in many different programming languages from JavaScript to Go - enabling easy code generation using basic string manipulation of the spreadsheet rows (as TSV, for instance).

## Code generation

The spreadsheet is capable of generating emulator code directly, which comes as columns which can be copy-pasted into your own code to have a full self contained 6502 instruction emulator in a function with no expternal dependencies (other than passing in the registers and memory array to the function).

![image](https://github.com/user-attachments/assets/b79542f3-ce62-4c48-83af-fbdb047df711)

Zooming in a bit, the form of the emulator code is visible; a generated if statement which checks for and executes each instruction:

![image](https://github.com/user-attachments/assets/154d616b-c8e1-46b8-92f9-aebcd0f36dec)

## Making an Emulator or Disassembler

Each row can be turned into a small block of code which processes that particular instruction.

The basic flow of the generated code for an emulator would be something like this. After setting up the initial memory, register variables, and PC pointing to the start of the program:

* Get byte at PC; this is the instruction

* Use a set of generated if statements which check for each possible byte value, and if that's the current byte then do the following:

* Do memory fetches for this instruction (if any)

* Calculate updated values for registers and flags (if any)

* Do memory writes (if any)

* Update register and flag variables with the new values (including PC)

* Loop around to get the next instruction or halt if PC has not changed value

## Reversibility

The spreadsheet also details for each instruction what data needs to be saved and what register, flag and memory changes need to be made to roll back the instruction.

This is the information which would otherwise be trashed by the instruction.

If you incorporate the recording of these values into your emulator as it executes, then you gain the ability to step backwards through the previously executed instructions when debugging.

## Keep It Simple

If you don't want to mess around with code generation, your 6502-enabled program can just have the spreadsheet itself in memory - probably broken down into an array of lines - and then do the string processing of the row on the fly to perform the operation or disassemble the instruction.

The actual distinct expressions for the updates to registers et al. are few enough that you can check for each one as a string specifically and perform the equivalent calculation in code; you don't need to do any expression parsing.

## Roadmap:

1) Complete basic list of instructions and first pass of the actions of each
2) Run through each instruction making sure its expressions are correct
3) Make codegen examples for JavaScript and Go for disassemblers and emulators
4) Make test cases and put it through its paces

## TECHNICALS

### Spreadsheet Columns

Here are the columns of the spreadsheet:

| Index | Position | Header Title                 | Description                                                                 |
|-------|----------|------------------------------|-----------------------------------------------------------------------------|
| 0     | A        | Bytes                        | Number of bytes used by the instruction                                     |
| 1     | B        | Name                         | Instruction mnemonic, e.g., BRK, ORA, etc.                                  |
| 2     | C        | Cycles                       | Number of CPU cycles the instruction takes                                  |
| 3     | D        | Args                         | The arguments used by the instruction, e.g., immediate, zero-page, etc.     |
| 4     | E        | Description of the instruction| A short description of what the instruction does                            |
| 5     | F        | zeroPageWord(Address) =        | Address used to perform zero-page word memory access, results in readByte   |
| 6     | G        | readByte(Address) =            | Address used to perform a byte read from memory, results in readByte        |
| 7     | H        | readWord(Address) =            | Address used to perform a word read from memory, results in readWord        |
| 8     | I        | result (word) =              | The result of the operation, affects flags and registers                    |
| 9     | J        | A =                          | The new value of the accumulator after the operation                        |
| 10    | K        | X =                          | The new value of the X register after the operation                         |
| 11    | L        | Y =                          | The new value of the Y register after the operation                         |
| 12    | M        | SP =                         | The new value of the stack pointer after the operation                      |
| 13    | N        | PC =                         | The new value of the program counter after the operation                    |
| 14    | O        | carry =                      | The new value of the carry flag                                             |
| 15    | P        | zero =                       | The new value of the zero flag                                              |
| 16    | Q        | interrupt =                  | The new value of the interrupt disable flag                                 |
| 17    | R        | decimal =                    | The new value of the decimal mode flag                                      |
| 18    | S        | break =                      | The new value of the break command flag                                     |
| 19    | T        | overflow =                   | The new value of the overflow flag                                          |
| 20    | U        | negative =                   | The new value of the negative flag                                          |
| 21    | V        | writeByteAddress             | Address where a byte is written in memory                                   |
| 22    | W        | writeByteValue               | The value that is written to the memory byte                                |
| 23    | X        | writeWordAddress             | Address where a word is written in memory                                   |
| 24    | Y        | writeWordValue               | The value that is written to the memory word                                |

### Example column values

Here with some example values from the `ORA ($LL, X)` instruction (opcode `0x01`):

| Header Title                 | Example Value (ORA ($LL, X))            | Notes |
|------------------------------|-----------------------------------------|-------|
| Bytes                        | 01 LL                                       | Opcode and one byte of argument (LL) |
| Name                         | ORA                                     | Assembly langugaue instruction name |
| Cycles                       | 6                                       | Number of cycles |
| Args                         | ($LL, X)                                | What the arguments to this assembly language instruction look like |
| Description of the instruction| Bitwise OR between A and byte at address (LL + X) | |
| zeroPageWord(Address) =        | LL + X                                  | This causes `zeroPageWord` to be set to `memory[zeroPageWordAddress&0xFF` | 
| readByte(Address) =            | zeroPageWord                            | This causes `readByte` to be set to `memory[readByteAddress&0xFFFF` |
| readWord(Address) =            |                                         | This would set `readWord` if present |
| result (byte?) =              | A | readByte                            | A convenient variable to store the result |
| A =                          | result                                  | Puts the result into the new value of A register|
| X =                          |                                         | No change to X |
| Y =                          |                                         | No change to Y register |
| SP =                         |                                         | No change to Stack Pointer |
| PC =                         |                                   | Two byte instruction; moving Program Counter on by 2 is the default behaviour |
| carry =                      |                                         | |
| zero =                       | result == 0                             | Zero flag gets updated |
| interrupt =                  |                                 | |
| decimal =                    |                                  | |
| break =                      |                                         | |
| overflow =                   |                                 | |
| negative =                   | result >> 7                             | Negative flag gets updated |
| writeByteAddress             |                                         | No memory writes |
| writeByteValue               |                                         | |
| writeWordAddress             |                                         | |
| writeWordValue               |                                         | |

### Go Language Support

#### func Emulate6502( A *uint8, X *unit8, Y *unit8, SP *unit8, PC *uint16, SR *uint8, memory *byte[]) bool

Column BE in the spreadsheet contains Go code generated from the instruction information.

This whole column can simply be copy-pasted into your source code to have a function called `Emulate6502` which emulates 6502 processing an instruction - given the current state updates that state based on the instruction at PC, moving PC on to the next instruction.

The generated code column starts with some fixed text:

```
func Emulate6502( A *uint8, X *unit8, Y *unit8, SP *unit8, PC *uint16, SR *uint8, memory *byte[]) bool {
  instruction := memory[PC & 0xFFFF]
```

Then has an if statement for each row which encodes the execution of the instruction in that row. Here is ORA (opcode 0x01) as an example:

```
if instruction == 0x01 { LL := memory[PC + 1]; zeroPageWord := memory[(LL + X)&0xFF] + (memory[((LL + X)+1)&0xFF]<<8); readByte := memory[(zeroPageWord)&0xFFFF]; result := A | readByte; newA := result; *A = newA; carry := SR&1; zero := result == 0;interrupt := (SR>>2)&1; decimal := (SR>>3)&1; overflow := (SR>>6)&1; negative := (SR>>7)&1; *SR = carry + (zero<<1) + (interrupt<<2) + (decimal<<3) + (overflow<<6) + (negative<<7); *PC = PC + 2; return true }
```

And finally after all the instruction handling there is the end of the `Emulate6502` function:

```
  return false
}
```

Although the expression looks complex, it breaks down into a set of easy string replacements in the spreadsheet formula:

```
="  if instruction == 0x" & LEFT(A9,2) & " { " &

JOIN("",

IF(LEN(A9) > 3, "LL := memory[PC + 1]; ", ""),

IF(I9 <> "", "zeroPageWord := memory[("&I9&")&0xFF] + (memory[(("&I9&")+1)&0xFF]<<8); ", ""),

IF(J9 <> "", "readByte := memory[("&J9&")&0xFFFF]; ", ""),
IF(K9 <> "", "readWord := memory[("&K9&")&0xFFFF] + (memory[(("&K9&")+1)&0xFFFF]<<8); ", ""),

IF(L9 <> "", "result := "&L9&"; ", ""),

IF(M9 <> "", "newA := "&M9&"; ", ""),
IF(N9 <> "", "newX := "&N9&"; ", ""),
IF(O9 <> "", "newY := "&O9&"; ", ""),
IF(P9 <> "", "newSP := "&P9&"; ", ""),
IF(Q9 <> "", "newPC := "&Q9&"; ", ""),

IF(M9 <> "", "*A = newA; ", ""),
IF(N9 <> "", "*X = newX; ", ""),
IF(O9 <> "", "*Y = newY; ", ""),
IF(P9 <> "", "*SP = newSP; ", ""),
IF(Q9 <> "", "*PC = newPC; ", ""),

IF( (R9 & S9 & T9 & U9 & V9 & W9) <> "", JOIN("",

  IF( R9 <> "", "carry := " & R9 & ";", "carry := SR&1; "),
  IF( S9 <> "", "zero := " & S9 & ";", "zero := (SR>>1)&1; "),
  IF( T9 <> "", "interrupt := " & T9 & ";", "interrupt := (SR>>2)&1; "),
  IF( U9 <> "", "decimal := " & U9 & ";", "decimal := (SR>>3)&1; "),
  IF( V9 <> "", "overflow := " & V9 & ";", "overflow := (SR>>6)&1; "),
  IF( W9 <> "", "negative := " & W9 & ";", "negative := (SR>>7)&1; ")

), ""),

IF( (R9 & S9 & T9 & U9 & V9 & W9) <> "", "*SR = carry + (zero<<1) + (interrupt<<2) + (decimal<<3) + (overflow<<6) + (negative<<7); ", ""),

IF( Q9 <> "", "*PC = " & Q9 & ";", "*PC = PC + " & (LEN(A9&" ")/3) & "; " ),

IF( Y9 <> "", "memory[(" & Y9 & ")&0xFFFF] = " & Z9 & "; ", ""),

IF( AA9 <> "", "memory[(" & AA9 & ")&0xFFFF] = (" & AB9 & ")&0xFF; memory[((" & AA9 & ")+1)&0xFFFF] = ((" & AB9 & ")>>8&0xFF; ", ""),
"return true"
)

& " }"
```
