package input

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ReadFile_ExistingFile_ReturnsContent(t *testing.T) {
	// given
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.py")
	content := "# comment\nprint('hello')"
	err := os.WriteFile(filePath, []byte(content), 0644)
	assert.NoError(t, err)

	// when
	result := ReadFile(filePath)

	// then
	assert.Contains(t, result, "# comment")
	assert.Contains(t, result, "print('hello')")
}

func Test_ReadFile_NonexistentFile_ReturnsEmptyString(t *testing.T) {
	// given
	filePath := "/nonexistent/file.py"

	// when
	result := ReadFile(filePath)

	// then
	assert.Equal(t, "", result)
}

func Test_ReadFile_UTF8EncodedFile_ReturnsContent(t *testing.T) {
	// given
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "utf8_test.py")
	content := "# 한글 주석\nprint('hello')"
	err := os.WriteFile(filePath, []byte(content), 0644)
	assert.NoError(t, err)

	// when
	result := ReadFile(filePath)

	// then
	assert.Contains(t, result, "한글 주석")
	assert.Contains(t, result, "print('hello')")
}

func Test_ReadFile_Latin1EncodedFile_FallbackToLatin1(t *testing.T) {
	// given
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "latin1_test.py")

	// Write with Latin-1 encoding (will fail UTF-8 validation)
	// "café" in Latin-1: 0x63 0x61 0x66 0xe9 (0xe9 is invalid standalone UTF-8)
	latin1Content := []byte{
		0x23, 0x20, 0x63, 0x61, 0x66, 0xe9, // "# café" in Latin-1
		0x0a,                                     // newline
		0x70, 0x72, 0x69, 0x6e, 0x74, 0x28, 0x27, // "print('"
		0x74, 0x65, 0x73, 0x74, 0x27, 0x29, // "test')"
	}
	err := os.WriteFile(filePath, latin1Content, 0644)
	assert.NoError(t, err)

	// when
	result := ReadFile(filePath)

	// then
	// Should fallback to latin-1 and successfully read
	assert.Contains(t, result, "café")
	assert.Contains(t, result, "print('test')")
}

func Test_ReadString_GivenContent_ReturnsSameContent(t *testing.T) {
	// given
	content := "# comment\ncode"

	// when
	result := ReadString(content)

	// then
	assert.Equal(t, content, result)
}

func Test_ReadString_EmptyString_ReturnsEmptyString(t *testing.T) {
	// given
	content := ""

	// when
	result := ReadString(content)

	// then
	assert.Equal(t, "", result)
}
