# game-boy-emulator

[![Go build & test](https://github.com/pascalPost/game-boy-emulator/actions/workflows/go.yml/badge.svg)](https://github.com/pascalPost/game-boy-emulator/actions/workflows/go.yml)

This repository contains the WIP of a game boy emulator written in go.

## Steps

1) Create a cli app that dumps the content of ROMs in a format (XXXX : XX XX XX XX  ....) containing the address, the
byte data and its ASCII representation. You can test, e.g. with this
[homebrew snake ROM](https://hh.gbdev.io/game/snake-gb). You may use `cmd/hexDump/hexDump.go` as a reference.
2) Parse the header and crate unit tests w.r.t., e.g., the Title, the Nintendo Logo and the Cartridge type.
3) Parse the json file with all opcodes (https://gbdev.io/gb-opcodes/Opcodes.json).
4) Write a disassembler. You may test with the snake ROM.
5) Begin programming the emulator by adding instructions for the load sequence of snake.
6) Add the graphics.

## Resources

Main references:
- A very well written introduction: https://www.inspiredpython.com/course/game-boy-emulator/let-s-write-a-game-boy-emulator-in-python
- Single most comprehensive technical reference to Game Boy: https://gbdev.io/pandocs/
- Homebrew Game Boy games: https://hh.gbdev.io/

Others (not yet checked):
- A go example for the Game Boy Advance: https://dev.to/aurelievache/learning-go-by-examples-part-5-create-a-game-boy-advance-gba-game-in-go-5944
- General detailed document about emulation:
  - http://www.codeslinger.co.uk/files/emu.pdf
- explanation of the opcodes:
  - https://gist.github.com/SakiiR/62661e45ee8b2ab13f0dc8203a7dfbd9
  - https://rgbds.gbdev.io/docs/v0.8.0/gbz80.7#LD_SP,n16
- Rust examples:
  - https://jeremybanks.github.io/0dmg/
  - https://rylev.github.io/DMG-01/public/book/introduction.html
  - https://read.cv/mehdi/uNGQ7pgWb2CO1QfJkb1n
- C/C++:
  - https://cturt.github.io/cinoop.html
  - http://www.codeslinger.co.uk/pages/projects/gameboy.html
