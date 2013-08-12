package fileutil

import (
	"os"
	"testing"
)

func Test_CopyFile(t *testing.T) {
	src := "test/test.txt"
	dst := "test/copied.txt"

	defer func() {
		os.Remove(src)
		os.Remove(dst)
	}()

	MkRandFile(src, 1024)

	if !Exist(src) {
		t.Fatal("cannot create file")
	}

	if err := CopyFile(src, dst); err != nil {
		t.Fatal(err)
	}

	if same, _ := IsSameFile(src, dst); !same {
		t.Error("not same")
	}
}

func Test_CopyFileN(t *testing.T) {
	src := "test/test.txt"
	destinations := []string{"test/copied1.txt", "test/copied2.txt", "test/copied3.txt"}

	defer func() {
		os.Remove(src)
		for _, dst := range destinations {
			os.Remove(dst)
		}
	}()

	MkRandFile(src, 1024)

	if !Exist(src) {
		t.Fatal("cannot create file")
	}

	CopyFileN(src, destinations...)

	if same, err := IsSameFileN(src, destinations...); !same {
		t.Errorf("not same %s", err)
	}
}

func Test_CopyDir(t *testing.T) {
	src := "test/sub1"
	dst := "test/copied1"

	defer func() { os.RemoveAll(dst) }()

	if err := CopyDir(src, dst); err != nil {
		t.Fatal(err)
	}

	if same, _ := IsSameDir(src, dst); !same {
		t.Error("not same")
	}
}
