package io

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	uuid2 "github.com/satori/go.uuid"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
)

func ReadFile(file string) ([]byte, error) {
	var out = []byte{}
	var err error
	var f *os.File
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
		if f != nil {
			f.Close()
		}
	}()
	if _, err := os.Stat(file); err != nil {
		return []byte{}, err
	}
	f, _ = os.Open(file)
	out, err = ioutil.ReadAll(f)
	return out, err
}

func WriteFile(file string, data []byte, override bool) error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
	}()
	if _, err = os.Stat(file); err == nil && override{
		err = os.Remove(file)
		if err != nil {
			return err
		}
	}
	err = ioutil.WriteFile(file, data, 0664)
	return err
}

func CopyFile(source string, dest string) error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
	}()
	var fs os.FileInfo
	if fs, err = os.Stat(source); err != nil {
		return errors.New(fmt.Sprintf("CopyFile - Source File %s doesn't exist!!", source))
	}
	if _, err = os.Stat(dest); err == nil {
		_ = os.Remove(dest)
	}
	var src, dst *os.File
	dst, err = os.Create(dest)
	if err != nil {
		return err
	}
	defer dst.Close()
	src, err = os.Open(source)
	if err != nil {
		return err
	}
	defer src.Close()
	var sz int64 = 0
	sz, err = io.Copy(dst, src)
	if sz != fs.Size() {
		return errors.New(fmt.Sprintf("CopyFile - Witten bytes: %v B are not same of source bytes %v B", sz, fs.Size()))
	}
	return err
}

func GetUniqueId() string {
	uuid, err := uuid.NewUUID()
	if err == nil {
		return uuid.String()
	}
	uuid2, err := uuid2.NewV1()
	if err == nil {
		return uuid2.String()
	}
	t1 := int64(7 * rand.Int()) + int64(11 * rand.Int())
	t2 := int64(7 * rand.Int()) + int64(11 * rand.Int())
	t3 := int64(7 * rand.Int()) + int64(11 * rand.Int())
	t4 := int64(7 * rand.Int()) + int64(11 * rand.Int())
	t5 := int64(7 * rand.Int()) + int64(11 * rand.Int())
	return fmt.Sprintf("%v-%v-%v-%v-%v", t1, t2, t3, t4, t5)
}