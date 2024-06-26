package helpers

import (
	"bufio"
	"debug/elf"
	"encoding/binary"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
)

// GetOffset takes an AppImage (either ELF or shappimage), returning the offset
// of its SquashFS archive
func GetOffset(src string) (int, error) {
	format, err := GetAppImageType(src)
	if err != nil {
		return -1, err
	}

	if format == -2 {
		return getShappImageSize(src)
	} else if format == 2 {
		return getElfSize(src)
	} else if format == 0 {
		return -1, errors.New("AppImage missing `AI\\0x02` magic at offset 0x08!")
	}

	return -1, errors.New("unsupported AppImage type")
}

// Takes a src file as argument, returning the size of the shImg header and
// an error if fail
func getShappImageSize(src string) (int, error) {
	f, err := os.Open(src)
	defer f.Close()
	if err != nil {
		return -1, err
	}

	_, err = f.Stat()
	if err != nil {
		return -1, err
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if len(scanner.Text()) > 10 && scanner.Text()[0:11] == "sfs_offset=" &&
			len(strings.Split(scanner.Text(), "=")) == 2 {

			offHex := strings.Split(scanner.Text(), "=")[1]
			offHex = strings.ReplaceAll(offHex, "'", "")
			offHex = strings.ReplaceAll(offHex, "\"", "")
			o, err := strconv.Atoi(offHex)

			return int(o), err
		}
	}

	return -1, errors.New("unable to find shappimage offset from `sfs_offset` variable")
}

// Function from <github.com/probonopd/go-appimage/internal/helpers/elfsize.go>
// credit goes to respective author; modified from original
// getElfSize takes a src file as argument, returning its size as an int
// and an error if unsuccessful
func getElfSize(src string) (int, error) {
	f, _ := os.Open(src)
	defer f.Close()
	e, err := elf.NewFile(f)
	if err != nil {
		return -1, err
	}

	// Find offsets based on arch
	sr := io.NewSectionReader(f, 0, 1<<63-1)
	var shoff, shentsize, shnum int

	switch e.Class {
	case elf.ELFCLASS64:
		hdr := new(elf.Header64)

		_, err = sr.Seek(0, 0)
		if err != nil {
			return -1, err
		}
		err = binary.Read(sr, e.ByteOrder, hdr)
		if err != nil {
			return -1, err
		}

		shoff = int(hdr.Shoff)
		shnum = int(hdr.Shnum)
		shentsize = int(hdr.Shentsize)
	case elf.ELFCLASS32:
		hdr := new(elf.Header32)

		_, err = sr.Seek(0, 0)
		if err != nil {
			return -1, err
		}
		err := binary.Read(sr, e.ByteOrder, hdr)
		if err != nil {
			return -1, err
		}

		shoff = int(hdr.Shoff)
		shnum = int(hdr.Shnum)
		shentsize = int(hdr.Shentsize)
	default:
		return 0, nil
	}

	return shoff + (shentsize * shnum), nil
}

// Find the type of AppImage
// Returns strings either `1` for ISO disk image AppImage, `2` for type 2
// SquashFS AppImage, `0` for unknown valid ELF or `-2` for shell script
// SquashFS AppImage (shappimage)
func GetAppImageType(src string) (int, error) {
	f, err := os.Open(src)
	defer f.Close()
	if err != nil {
		return -1, err
	}

	_, err = f.Stat()
	if err != nil {
		return -1, err
	}

	if HasMagic(f, "\x7fELF", 0) {
		if HasMagic(f, "AI\x01", 8) {
			// AppImage type is type 1 (standard)
			return 1, nil
		} else if HasMagic(f, "AI\x02", 8) {
			// AppImage type is type 2 (standard)
			return 2, nil
		}
		// Unknown AppImage, but valid ELF
		return 0, nil
	} else if HasMagic(f, "#!/bin/sh\n#.shImg.#", 0) {
		// AppImage is shappimage (shell script SquashFS implementation)
		return -2, nil
	}

	err = errors.New("unable to get AppImage type")
	return -1, err
}

// Checks the magic of a given file against the byte array provided
// if identical, return true
func HasMagic(r io.ReadSeeker, str string, offset int) bool {
	magic := make([]byte, len(str))

	r.Seek(int64(offset), io.SeekStart)

	_, err := io.ReadFull(r, magic[:])
	if err != nil {
		return false
	}

	for i := 0; i < len(str); i++ {
		if magic[i] != str[i] {
			return false
		}
	}

	return true
}
