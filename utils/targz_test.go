package utils

import "testing"

func TestCompress(t *testing.T) {
	err := Compress("../cli", "utils1.tar.gz")
	if err != nil {
		return
	}
	err = Extract("utils.tar.gz", "/Users/maybemia/Documents/github/maygit/cli")
	if err != nil {
		return
	}
}
