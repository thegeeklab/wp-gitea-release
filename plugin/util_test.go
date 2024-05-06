package plugin

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

//nolint:lll
func TestChecksum(t *testing.T) {
	tests := []struct {
		name     string
		method   string
		input    []byte
		expected string
		wantErr  bool
	}{
		{
			name:     "md5",
			method:   "md5",
			input:    []byte("hello"),
			expected: "5d41402abc4b2a76b9719d911017c592",
			wantErr:  false,
		},
		{
			name:     "sha1",
			method:   "sha1",
			input:    []byte("hello"),
			expected: "aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d",
			wantErr:  false,
		},
		{
			name:     "sha256",
			method:   "sha256",
			input:    []byte("hello"),
			expected: "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
			wantErr:  false,
		},
		{
			name:     "sha512",
			method:   "sha512",
			input:    []byte("hello"),
			expected: "9b71d224bd62f3785d96d46ad3ea3d73319bfbc2890caadae2dff72519673ca72323c3d99ba5c11d7c7acc6e14b8c5da0c4663475c2e5c3adef46f73bcdec043",
			wantErr:  false,
		},
		{
			name:     "adler32",
			method:   "adler32",
			input:    []byte("hello"),
			expected: "062c0215",
			wantErr:  false,
		},
		{
			name:     "crc32",
			method:   "crc32",
			input:    []byte("hello"),
			expected: "3610a686",
			wantErr:  false,
		},
		{
			name:     "blake2b",
			method:   "blake2b",
			input:    []byte("hello"),
			expected: "324dcf027dd4a30a932c441f365a25e86b173defa4b8e58948253471b81b72cf",
			wantErr:  false,
		},
		{
			name:     "blake2s",
			method:   "blake2s",
			input:    []byte("hello"),
			expected: "19213bacc58dee6dbde3ceb9a47cbb330b3d86f8cca8997eb00be456f140ca25",
			wantErr:  false,
		},
		{
			name:     "unsupported_method",
			method:   "unsupported",
			input:    []byte("hello"),
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bytes.NewReader(tt.input)

			sum, err := Checksum(r, tt.method)
			if tt.wantErr {
				assert.Error(t, err)

				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, sum)
		})
	}
}

func BenchmarkChecksum(b *testing.B) {
	input := []byte("hello")
	methods := []string{"md5", "sha1", "sha256", "sha512", "adler32", "crc32", "blake2b", "blake2s"}

	for _, method := range methods {
		b.Run(method, func(b *testing.B) {
			r := bytes.NewReader(input)

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_, _ = Checksum(r, method)
				_, _ = r.Seek(0, io.SeekStart)
			}
		})
	}
}

func TestWriteChecksums(t *testing.T) {
	tempDir := t.TempDir()
	files := []string{
		filepath.Join(tempDir, "file1.txt"),
		filepath.Join(tempDir, "file2.txt"),
	}

	for _, file := range files {
		err := os.WriteFile(file, []byte("hello"), 0o600)
		if err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}
	}

	methods := []string{"md5", "sha256"}

	tests := []struct {
		name     string
		files    []string
		methods  []string
		expected []string
		wantErr  bool
	}{
		{
			name:    "valid_input",
			files:   files,
			methods: methods,
			expected: []string{
				files[0],
				files[1],
				filepath.Join(tempDir, "md5sum.txt"),
				filepath.Join(tempDir, "sha256sum.txt"),
			},
			wantErr: false,
		},
		{
			name:    "empty_files",
			files:   []string{},
			methods: methods,
			expected: []string{
				filepath.Join(tempDir, "md5sum.txt"),
				filepath.Join(tempDir, "sha256sum.txt"),
			},
			wantErr: false,
		},
		{
			name:    "empty_methods",
			files:   files,
			methods: []string{},
			expected: []string{
				files[0],
				files[1],
			},
			wantErr: false,
		},
		{
			name:     "non_existent_file",
			files:    append(files, filepath.Join(tempDir, "non_existent.txt")),
			methods:  methods,
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := WriteChecksums(tt.files, tt.methods)
			if tt.wantErr {
				assert.Error(t, err)

				return
			}

			assert.NoError(t, err)
			sort.Strings(result)
			sort.Strings(tt.expected)
			assert.Equal(t, tt.expected, result)
		})
	}
}
