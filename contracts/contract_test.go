package contracts

import (
	"testing"
)

func TestWriteFile(t *testing.T) {
	if err := WriteFile("/tmp/test_dir/", "code.bin", "test"); err != nil {
		t.Fatal(err)
	}
}
