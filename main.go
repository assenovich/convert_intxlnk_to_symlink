package main

import (
	"bytes"
	"encoding/binary"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"unicode/utf16"
)

func DecodeUtf16(b []byte) (string, error) {
	ints := make([]uint16, len(b)/2)
	if err := binary.Read(bytes.NewReader(b), binary.BigEndian, &ints); err != nil {
		return "", err
	}
	return string(utf16.Decode(ints)), nil
}

func main() {
	walk := func(path string) error {
		return filepath.Walk(path,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				mode := info.Mode()
				if (mode&os.ModeSymlink == 0) && (mode&os.ModeDir == 0) && (info.Size() < 4096) {
					data, err := ioutil.ReadFile(path)
					if err != nil {
						return err
					}
					if len(data) > 8 {
						if string(data[0:7]) == "IntxLNK" {
							log.Println(path)
							data[7] = 0
							dst, err := DecodeUtf16(data[7:])
							if err != nil {
								return err
							}
							log.Println(dst)
							if err := os.Remove(path); err != nil {
								return err
							}
							if err := os.Symlink(dst, path); err != nil {
								return err
							}
							log.Println()
						}
					}
				}
				return nil
			})
	}
	if err := walk("/home/daw/9_data"); err != nil {
		log.Println(err)
		return
	}
}
