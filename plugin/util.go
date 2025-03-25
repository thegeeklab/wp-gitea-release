package plugin

import (
	"crypto/md5"  //nolint:gosec
	"crypto/sha1" //nolint:gosec
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"fmt"
	"hash/adler32"
	"hash/crc32"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/blake2s"
)

var ErrHashMethodNotSupported = errors.New("hash method not supported")

// Checksum calculates the checksum of the given io.Reader using the specified hash method.
// Supported hash methods are: "md5", "sha1", "sha256", "sha512", "adler32", "crc32", "blake2b", "blake2s".
func Checksum(r io.Reader, method string) (string, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}

	switch method {
	case "md5":
		//nolint:gosec
		return fmt.Sprintf("%x", md5.Sum(b)), nil
	case "sha1":
		//nolint:gosec
		return fmt.Sprintf("%x", sha1.Sum(b)), nil
	case "sha256":
		return fmt.Sprintf("%x", sha256.Sum256(b)), nil
	case "sha512":
		return fmt.Sprintf("%x", sha512.Sum512(b)), nil
	case "adler32":
		return fmt.Sprintf("%08x", adler32.Checksum(b)), nil
	case "crc32":
		return fmt.Sprintf("%08x", crc32.ChecksumIEEE(b)), nil
	case "blake2b":
		return fmt.Sprintf("%x", blake2b.Sum256(b)), nil
	case "blake2s":
		return fmt.Sprintf("%x", blake2s.Sum256(b)), nil
	}

	return "", fmt.Errorf("%w: %q", ErrHashMethodNotSupported, method)
}

// WriteChecksums calculates the checksums for the given files using the specified hash methods,
// and writes the checksums to files named after the hash methods (e.g. "md5sum.txt", "sha256sum.txt").
func WriteChecksums(files, methods []string, outDir string) ([]string, error) {
	if len(files) == 0 {
		return files, nil
	}

	checksumFiles := make([]string, 0)

	for _, method := range methods {
		checksumFile := filepath.Join(outDir, method+"sum.txt")

		f, err := os.Create(checksumFile)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		for _, file := range files {
			handle, err := os.Open(file)
			if err != nil {
				return nil, fmt.Errorf("failed to read %q artifact: %w", file, err)
			}
			defer handle.Close()

			hash, err := Checksum(handle, method)
			if err != nil {
				return nil, fmt.Errorf("could not checksum %q file: %w", file, err)
			}

			_, err = fmt.Fprintf(f, "%s  %s\n", hash, file)
			if err != nil {
				return nil, err
			}
		}

		checksumFiles = append(checksumFiles, checksumFile)
	}

	return append(files, checksumFiles...), nil
}
