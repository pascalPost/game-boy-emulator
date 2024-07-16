package internal

import (
	"bytes"
	"encoding/binary"
)

type Header struct {
	Raw struct {
		NotUsed [64 * 4]byte

		// [0x0100:0x0104]  https://gbdev.io/pandocs/The_Cartridge_Header.html#0100-0103--entry-point
		EntryPoint [4]byte

		// [0x0104:0x0134] https://gbdev.io/pandocs/The_Cartridge_Header.html#0104-0133--nintendo-logo
		NintendoLogo [4 * 12]byte

		// [0x0134:0x0144]
		// Title (https://gbdev.io/pandocs/The_Cartridge_Header.html#0134-0143--title)
		// might also contain:
		// - manufacture code (https://gbdev.io/pandocs/The_Cartridge_Header.html#013f-0142--manufacturer-code)
		// - CGB flag (https://gbdev.io/pandocs/The_Cartridge_Header.html#0143--cgb-flag)
		TitleManufacturerCodeCGBFlag [4 * 4]byte

		// [0x0144:0x0146] https://gbdev.io/pandocs/The_Cartridge_Header.html#01440145--new-licensee-code
		NewLicenseeCode [2]byte

		// [0x0146] https://gbdev.io/pandocs/The_Cartridge_Header.html#0146--sgb-flag
		SGBFlag byte

		// [0x0147] https://gbdev.io/pandocs/The_Cartridge_Header.html#0147--cartridge-type
		CartridgeType byte

		// [0x0148] https://gbdev.io/pandocs/The_Cartridge_Header.html#0148--rom-size
		RomSize byte

		// [0x0149] https://gbdev.io/pandocs/The_Cartridge_Header.html#0149--ram-size
		RamSize byte

		// [0x014A] https://gbdev.io/pandocs/The_Cartridge_Header.html#014a--destination-code
		DestinationCode byte

		// [0x014B] https://gbdev.io/pandocs/The_Cartridge_Header.html#014b--old-licensee-code
		OldLicenseeCode byte

		// [0x014C] https://gbdev.io/pandocs/The_Cartridge_Header.html#014c--mask-rom-version-number
		MaskRomVersionNumber byte

		// [0x014D] https://gbdev.io/pandocs/The_Cartridge_Header.html#014d--header-checksum
		HeaderChecksum byte

		// [0x014E:0x0150] https://gbdev.io/pandocs/The_Cartridge_Header.html#014e-014f--global-checksum
		GlobalChecksum [2]byte
	}
}

func NewHeader(b []byte) (Header, error) {
	header := Header{}
	buf := bytes.NewReader(b)
	err := binary.Read(buf, binary.BigEndian, &header)
	return header, err
}

func (h *Header) CartridgeType() string {
	// https://gbdev.io/pandocs/The_Cartridge_Header.html#0147--cartridge-type
	cartridgeTypeMap := map[byte]string{
		0x00: "ROM ONLY",
		0x01: "MBC1",
		0x02: "MBC1+RAM",
		0x03: "MBC1+RAM+BATTERY",
		0x05: "MBC2",
		0x06: "MBC2+BATTERY",
		0x08: "ROM+RAM 9",
		0x09: "ROM+RAM+BATTERY 9",
		0x0B: "MMM01",
		0x0C: "MMM01+RAM",
		0x0D: "MMM01+RAM+BATTERY",
		0x0F: "MBC3+TIMER+BATTERY",
		0x10: "MBC3+TIMER+RAM+BATTERY 10",
		0x11: "MBC3",
		0x12: "MBC3+RAM 10",
		0x13: "MBC3+RAM+BATTERY 10",
		0x19: "MBC5",
		0x1A: "MBC5+RAM",
		0x1B: "MBC5+RAM+BATTERY",
		0x1C: "MBC5+RUMBLE",
		0x1D: "MBC5+RUMBLE+RAM",
		0x1E: "MBC5+RUMBLE+RAM+BATTERY",
		0x20: "MBC6",
		0x22: "MBC7+SENSOR+RUMBLE+RAM+BATTERY",
		0xFC: "POCKET CAMERA",
		0xFD: "BANDAI TAMA5",
		0xFE: "HuC3",
		0xFF: "HuC1+RAM+BATTERY",
	}

	return cartridgeTypeMap[h.Raw.CartridgeType]
}

type RomSizeInfo struct {
	RomSize          string
	NumberOfRomBanks int
	ExtraInfo        string
}

func (h *Header) RomSize() RomSizeInfo {
	// https://gbdev.io/pandocs/The_Cartridge_Header.html#0148--rom-size
	m := map[byte]RomSizeInfo{
		0x00: {"32 KiB", 2, "(no banking)"},
		0x01: {"64 KiB", 4, ""},
		0x02: {"128 KiB", 8, ""},
		0x03: {"256 KiB", 16, ""},
		0x04: {"512 KiB", 32, ""},
		0x05: {"1 MiB", 64, ""},
		0x06: {"2 MiB", 128, ""},
		0x07: {"4 MiB", 256, ""},
		0x08: {"8 MiB", 512, ""},
		0x52: {"1.1 MiB", 72, "unofficial, likely inaccurate"},
		0x53: {"1.2 MiB", 80, "unofficial, likely inaccurate"},
		0x54: {"1.5 MiB", 96, "unofficial, likely inaccurate"},
	}
	return m[h.Raw.RomSize]
}
