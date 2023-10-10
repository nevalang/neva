# Components

## Interpreter

Runs source code without separate step of saving IR to the disc. Usage example `neva run cmd/server`

## Compiler

Reads source code, analyzes it and, if no errors found, generates IR and writes it to the disc as a file so it later can be used by VM

## Virtual Machine

Reads IR file from the disc and passes the data to the runtime

## Runtime

Executes IR
