# vm

Overview
--------
VM is intended to be a virtual machine for the person, namely me, wanting to
learn the ins and outs of a CPU. More effect has been made to make it mimic
an actual CPU rather than produce something efficient. If you're looking
for something fast and good for every day use, better look elsewhere.

The VM comes with its own assembly language. At the time of writing, it has
only 20 instructions but expect to see that number rise, if only slightly.

The virtual machine is intended to be part of an entire tool chain spanning
from a simple, easily parsed high level language ("Simple C" is the working
title), a custom object format, an assembler and a linker.

The assembler, linker and virtual machine are in working order but no work
has been done on the high level language at this time.

Limitations
-----------
* You will not be able pre-declare data or variables and, concequently, to 
load said data. The object format currently supports a data section but the 
instructions to load the data and the corresponding assembly has not been
implemented. It is planned, however.
* Branching is somewhat limited. Conditionals within the virtual machine are
limited to being zero or non-zero. No logical tests like greater or less than
are planned at this time but are easily added.
* No dynamic loading or linking. It is beyond the scope of this project.

CPU Specification
-----------------
* Data BUS: 8 bit
* Address BUS: 16 bit
* Working Registers: Accumulator, two general purpose named B and C.
* Non-Accessible Registers: Stack Pointer, Instruction, Temporary, Data,
and Address
* Instructions: 20

Inspirations
------------
* CPU: Intel 4004, Intel 8008, MOS6502 and the Intel 8080.
* Object Format: ELF

