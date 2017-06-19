package utils

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	//   "time"
)

var ()

func Sysexec(cmd string, args ...string) (result string, err error) {
	arg := append([]string{cmd}, args...)
	arg_str := fmt.Sprintf("%s", strings.Join(arg, " "))

	ori_output, err := exec.Command("/bin/bash", "-l", "-c", arg_str).CombinedOutput()
	return strings.TrimSpace(string(ori_output)), err
}

func IsDirExist(path string) (exist bool) {
	fi, err := os.Stat(path)

	if err != nil {
		exist = os.IsExist(err)
	} else {
		exist = fi.IsDir()
	}
	return
}

func IsFileExist(path string) (exist bool) {
	fi, err := os.Stat(path)

	if err != nil {
		exist = os.IsExist(err)
	} else {
		exist = !fi.IsDir()
	}
	return
}

func ListDir(dirPth string, suffix string) (files []string, err error) {
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}

	files = make([]string, 0, len(dir))

	for _, fi := range dir {
		// 若该文件为目录，或结尾不符合，则忽略
		if fi.IsDir() || !strings.HasSuffix(fi.Name(), suffix) {
			continue
		}

		files = append(files, fi.Name())
	}

	return
}

func Include(str string, strings []string) bool {
	for _, element := range strings {
		if element == str {
			return true
		}
	}

	return false
}

func AtoUint64(s string) (i uint64, err error) {
	return strconv.ParseUint(s, 10, 64)
}

func IsDiskFS(fs string) bool {
	switch fs {
	case "ext4", "ext3", "ext2":
		return true
	default:
		return false
	}
}

func CloneMap(src map[string]string) (dst map[string]string) {
	dst = make(map[string]string)
	for key, value := range src {
		dst[key] = value
	}

	return
}

func UnitConvert(i float64, oldunit, newunit string) (j float64) {
	m1 := GetMagnificationFromUnit(oldunit)
	m2 := GetMagnificationFromUnit(newunit)
	j = i * m1 / m2

	return
}

func GetMagnificationFromUnit(unit string) (magnification float64) {
	switch unit {
	case "bit":
		magnification = 1
	case "byte":
		magnification = 8
	case "sector":
		magnification = 512
	case "kb":
		magnification = 1000
	case "kib":
		magnification = 1024
	case "mb":
		magnification = 1000000
	case "mib":
		magnification = 1048576
	case "gb":
		magnification = 1000000000
	case "gib":
		magnification = 1073741824
	}

	return
}

func GobEncode(data interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func GobDecode(data []byte, to interface{}) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(to)
}
