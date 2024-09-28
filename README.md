# 6502 THE SPREADSHEET (The Repository)

All 6502 instructions with their actions completely unrolled so you don't have to pick through detailed documentation to understand what they do or how to emulate them accurately.

The spreadsheet is under construction and is at:
https://docs.google.com/spreadsheets/d/1rCwO6CoIXiCqkJUaZ5Q9W5QZNqmhLpXaAhSthJPuV8I/edit?gid=0#gid=0

![image](https://github.com/user-attachments/assets/31bceb69-3d3f-4adc-9b40-a0d8ea07790d)

## How it works

For each instruction there is a set of expressions for each of the necessary memory reads, register updates, flag updates, and memory writes.

These expressions are provided in a simple portable form, just using the symbols + - << >> & | () and are thus valid expressions in many different programming languages from JavaScript to Go - enabling easy code generation using basic string manipulation of the spreadsheet rows (as TSV, for instance).

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

