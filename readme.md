Solididty Disassembler
======================

Geth includes a disassembler but annoyingly, like most solidity disassemblers, the addresses are in decimal.

This prompted me to fix it - which was surprisingly easy as long as you have geth installed because geth does most of the work for you.

I had been intending to remove the necessity for geth but since I intend to allow it to slurp contracts from the blockchain, I may still need it.
Decide later. Maybe we need two versions :-)

Why?
---

When I found a popular token was made up of 42 contracts of which only about 5 were verified, I got curious. 

I then started thinking that there must be a better way to at least find the entry points.

myDasm is designed 
This is intended to help people

* analyse gas costs
* examine EVM code
* look at unverified contracts listing entry points

Usage (so far), assuming that you build it

`myDasm -input <hex file> -gas <yes/no> -opcode <yes/no>`

`<hex file>`

The hex file is an ascii encoded file of hex characters something like

```
606060405234156200001057600080fd5b6040516020806200247b833981016040528080519060200190919050506200006d7f693a746f6b656e0000000000000000000000000000000000000000000000000082620000806401000000000262001f08176401000000009004565b15156200007957600080fd5b50620002df565b6000808273ffffffffffffffffffffffffffffffffffffffff1663cf309012
```

Vanilla output looks like..

```
000000:  PUSH1 0x60 
000002:  PUSH1 0x40 
000004:  MSTORE 
000005:  CALLVALUE 
000006:  ISZERO 
000007:  PUSH3 0x000010 
00000B:  JUMPI 
```

`-gas`

If you answer yes to this, you will see the gas uasage...

```
000000: (3) PUSH1 0x60 
000002: (3) PUSH1 0x40 
000004: (*) MSTORE 
000005: (2) CALLVALUE 
000006: (3) ISZERO 
000007: (3) PUSH3 0x000010 
00000B: (10) JUMPI 
00000C: (3) PUSH1 0x00 
```

`-opcode`

If you answer yes to this, you will see every line explained (over and over) except PUSH4 and PUSH32

```
000000:  PUSH1 0x60 // Place 1 byte item on stack
000002:  PUSH1 0x40 // Place 1 byte item on stack
000004:  MSTORE // Save word to memory
000005:  CALLVALUE // Get deposited value by the instruction/transaction responsible for this execution
000006:  ISZERO // Simple not operator
000007:  PUSH3 0x000010 // Place 3-byte item on stack
00000B:  JUMPI // Conditionally alter the program counter
00000C:  PUSH1 0x00 // Place 1 byte item on stack
```
Probably a bit over the top?

Building it and installing it
---

You need go installed and you need the go-ethereum sourcecode at
`$GOPATH/src/github.com/ethereum/go-ethereum`

`$ go install myDasm.go`

bingo - you are ready to go

Acknowledgements

* The go-ethereum team for an amazing product
* Howard Yeah for his great deep-dive into the EVM series
* Trailofbits for their opcode-gas table

