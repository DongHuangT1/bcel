package bcel

import (
	"bytes"
	"compress/gzip"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strings"
)

var (
	CHAR_MAP    = make([]int, 48)
	MAP_CHAR    = make([]int, 256)
	BCEL_TOKEN  = "$$BCEL$$"
	ESCAPE_CHAR = byte('$')
)

func init() {
	var j, i int

	for i = 'A'; i <= 'Z'; i++ {
		CHAR_MAP[j] = i
		MAP_CHAR[i] = j
		j++
	}

	for i = 'g'; i <= 'z'; i++ {
		CHAR_MAP[j] = i
		MAP_CHAR[i] = j
		j++
	}

	CHAR_MAP[j] = '$'
	MAP_CHAR['$'] = j
	j++
	CHAR_MAP[j] = '_'
	MAP_CHAR['_'] = j
}

func Encode(b []byte, compress bool) (string, error) {
	if compress {
		rw := new(bytes.Buffer)
		gz := gzip.NewWriter(rw)

		_, err := gz.Write(b)
		if err != nil {
			return "", err
		}

		err = gz.Close()
		if err != nil {
			return "", err
		}

		b = rw.Bytes()
	}

	var buf bytes.Buffer

	for _, v := range b {
		n := byte(v & 0x000000ff)
		if isJavaIdentifierPart(n) && n != ESCAPE_CHAR {
			buf.WriteByte(n)
		} else {
			buf.WriteByte(ESCAPE_CHAR)
			if n >= 48 {
				fmt.Fprintf(&buf, "%02x", n)
			} else {
				buf.WriteByte(byte(CHAR_MAP[n]))
			}
		}
	}

	return BCEL_TOKEN + buf.String(), nil
}

func Decode(s string, uncompress bool) ([]byte, error) {
	if !strings.HasPrefix(s, BCEL_TOKEN) {
		return nil, errors.New("invalid bcel")
	}

	var buf bytes.Buffer
	s = s[len(BCEL_TOKEN):]

	for r := strings.NewReader(s); ; {
		b, err := r.ReadByte()
		if err != nil {
			break
		}

		if b != ESCAPE_CHAR {
			buf.WriteByte(b)
		} else {
			i, err := r.ReadByte()
			if err != nil {
				break
			}

			if i >= '0' && i <= '9' || i >= 'a' && i <= 'f' {
				j, err := r.ReadByte()
				if err != nil {
					break
				}

				tmp := []byte{i, j}

				n, err := hex.Decode(tmp, tmp)
				if err != nil {
					return nil, err
				}

				buf.Write(tmp[:n])
			} else {
				buf.WriteByte(byte(MAP_CHAR[i]))
			}
		}
	}

	if uncompress {
		gz, err := gzip.NewReader(&buf)
		if err != nil {
			return nil, err
		}
		defer gz.Close()

		return io.ReadAll(gz)
	}

	return buf.Bytes(), nil
}

func isJavaIdentifierPart(ch byte) bool {
	return ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch >= '0' && ch <= '9' || ch == '_'
}
