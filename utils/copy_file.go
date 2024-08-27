package utils

import (
	"bufio"
	"io"
	"os"
)

// CopyFile copies a file from src to dst.
func CopyFile(src, dst string) error {
	fi, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fi.Close()
	r := bufio.NewReader(fi)

	fo, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer fo.Close()

	w := bufio.NewWriter(fo)

	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		if _, err := w.Write(buf[:n]); err != nil {
			return err
		}
	}

	if err = w.Flush(); err != nil {
		return err
	}
	return nil
}
