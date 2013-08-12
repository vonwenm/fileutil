package fileutil

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"errors"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

var (
	ErrNotFile = errors.New("not a file")
	ErrNotDir  = errors.New("not a directory")
)

func MkRandFile(path string, size int) error {
	file, err := os.Create(path)

	if err != nil {
		return err
	}
	defer file.Close()

	for i := 0; i < size; i++ {
		if _, err := io.CopyN(file, rand.Reader, 1024); err != nil {
			return err
		}
	}

	return nil
}

func Sha1(filename string) (string, error) {
	return hashSum(filename, sha1.New())
}

func Sha256(filename string) (string, error) {
	return hashSum(filename, sha256.New())
}

func MD5(filename string) (string, error) {
	return hashSum(filename, md5.New())
}

func CopyFileN(src string, destinations ...string) {
	var wg sync.WaitGroup
	for _, dst := range destinations {
		wg.Add(1)
		go func(file string) {
			defer wg.Done()
			if err := CopyFile(src, file); err != nil {
				fmt.Println(err)
			}
		}(dst)
	}
	wg.Wait()
}

func CopyFile(src, dst string) error {
	fi, err := os.Stat(src)
	if err != nil {
		return err
	}
	if fi.IsDir() {
		return ErrNotFile
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	if Exist(dst) {
		return os.ErrExist
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if err := dstFile.Chmod(fi.Mode()); err != nil {
		return err
	}

	_, err = io.Copy(dstFile, srcFile)

	return err
}

func CopyDir(srcDir, dstDir string) error {
	fi, err := os.Stat(srcDir)
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		return ErrNotDir
	}

	if Exist(dstDir) {
		return os.ErrExist
	} else {
		if err := os.Mkdir(dstDir, fi.Mode()); err != nil {
			return err
		}
	}

	fis, err := ioutil.ReadDir(srcDir)
	if err != nil {
		return err
	}

	for _, fi := range fis {
		dstName := filepath.Join(dstDir, fi.Name())
		srcName := filepath.Join(srcDir, fi.Name())
		if fi.IsDir() {
			if err := CopyDir(srcName, dstName); err != nil {
				return err
			}
		} else {
			if err := CopyFile(srcName, dstName); err != nil {
				return err
			}
		}
	}

	return nil
}

func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func IsSameFile(src, dst string) (bool, error) {
	srcDigest, err := Sha1(src)
	if err != nil {
		return false, err
	}

	dstDigest, err := Sha1(dst)
	if err != nil {
		return false, err
	}

	return srcDigest == dstDigest, nil
}

func IsSameFileN(src string, destinations ...string) (bool, error) {
	srcDigest, err := Sha1(src)
	if err != nil {
		return false, err
	}

	for _, dst := range destinations {
		if dstDigest, err := Sha1(dst); err != nil || srcDigest != dstDigest {
			return false, err
		}
	}

	return true, nil
}

func IsSameDir(srcDir, dstDir string) (bool, error) {
	for _, dir := range [...]string{srcDir, dstDir} {
		if fi, err := os.Stat(dir); err != nil {
			return false, err
		} else if !fi.IsDir() {
			return false, ErrNotDir
		}
	}

	fis, err := ioutil.ReadDir(srcDir)
	if err != nil {
		return false, err
	}

	for _, fi := range fis {
		dstName := filepath.Join(dstDir, fi.Name())
		srcName := filepath.Join(srcDir, fi.Name())
		if fi.IsDir() {
			if same, err := IsSameDir(srcName, dstName); err != nil || !same {
				return same, err
			}
		} else {
			if same, err := IsSameFile(srcName, dstName); err != nil || !same {
				return same, err
			}
		}
	}

	return true, nil
}

func hashSum(filename string, digest hash.Hash) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if _, err := io.Copy(digest, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", digest.Sum(nil)), nil
}
