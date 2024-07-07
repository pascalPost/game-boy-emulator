# game-boy-emulator

This repository contains the WIP of a game boy emulator written in go.

## Steps

1) Create a cli app that dumps the content of ROMs in a format (XXXX : XX XX XX XX  ....) containing the address, the
byte data and its ASCII representation. You can test, e.g. with this
[homebrew snake ROM](https://hh.gbdev.io/game/snake-gb). You may use `cmd/hexDump/hexDump.go` as a reference.

## Resources

Main references:
- A very well written introduction: https://www.inspiredpython.com/course/game-boy-emulator/let-s-write-a-game-boy-emulator-in-python
- Single most comprehensive technical reference to Game Boy: https://gbdev.io/pandocs/
- Homebrew Game Boy games: https://hh.gbdev.io/

Others (not yet checked):
- A go example for the Game Boy Advance: https://dev.to/aurelievache/learning-go-by-examples-part-5-create-a-game-boy-advance-gba-game-in-go-5944
- Rust examples:
  - https://jeremybanks.github.io/0dmg/
  - https://rylev.github.io/DMG-01/public/book/introduction.html
  - https://read.cv/mehdi/uNGQ7pgWb2CO1QfJkb1n
- C/C++:
  - https://cturt.github.io/cinoop.html
  - http://www.codeslinger.co.uk/pages/projects/gameboy.html



