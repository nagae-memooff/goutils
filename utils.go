package utils

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
	"time"

	"crypto/sha1"
)

const (
	letterAlphabets           = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterNumbers             = "0123456789"
	letterAlphabetsAndNumbers = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	letterIdxBits             = 6                    // 6 bits to represent a letter index
	letterIdxMask             = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax              = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var ()

// 以登录式shell的形式执行命令。由于是登陆式shell，因此能加载环境变量
// 由于是把命令拼装传给bash，因此支持管道、重定向等，使用比较方便
func Sysexec(cmd string, args ...string) (result string, err error) {
	arg := append([]string{cmd}, args...)
	arg_str := fmt.Sprintf("%s", strings.Join(arg, " "))

	ori_output, err := exec.Command("/bin/bash", "-l", "-c", arg_str).CombinedOutput()
	return strings.TrimSpace(string(ori_output)), err
}

// 判断目录是否存在
func IsDirExist(path string) (exist bool) {
	fi, err := os.Stat(path)

	if err != nil {
		exist = os.IsExist(err)
	} else {
		exist = fi.IsDir()
	}
	return
}

// 判断文件是否存在
func IsFileExist(path string) (exist bool) {
	fi, err := os.Stat(path)

	if err != nil {
		exist = os.IsExist(err)
	} else {
		exist = !fi.IsDir()
	}
	return
}

func Mkdirp(dir string) {
	os.MkdirAll(dir, os.ModePerm)
}

// 列出目录中指定结尾的文件的文件名(不包含目录)
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

// 字符串数组中是否存在指定字符
func Include(str string, strings []string) bool {
	if strings == nil {
		return false
	}

	for _, element := range strings {
		if element == str {
			return true
		}
	}

	return false
}

// string转换成uint64
func AtoUint64(s string) (i uint64, err error) {
	return strconv.ParseUint(s, 10, 64)
}

// 判断指定的文件系统类型是否是磁盘文件系统，diamond里用
func IsDiskFS(fs string) bool {
	switch fs {
	case "btrfs", "ext2", "ext3", "ext4", "jfs", "reiser", "xfs", "ffs", "ufs", "jfs2", "vxfs", "hfs", "ntfs", "fat32", "zfs", "fuse.mfs":
		return true
	default:
		return false
	}
}

// 拷贝一个字符串类型的map。字符串类型，无所谓深拷贝浅拷贝。
func CloneMap(src map[string]string) (dst map[string]string) {
	dst = make(map[string]string)
	for key, value := range src {
		dst[key] = value
	}

	return
}

// 单位转换用
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

// 把一个data序列化成golang自己的gob格式，方便存储
func GobEncode(data interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// 把一个gob格式的序列化后的byte数组，反序列化到给定的结构中。
func GobDecode(data []byte, to interface{}) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(to)
}

// 获取当前进程内存占用
func GetMemUsage() (int, error) {
	pageSize := 4096

	pid := os.Getpid()

	f, err := os.Open(fmt.Sprintf("/proc/%d/stat", pid))
	if err != nil {
		return 0, fmt.Errorf("failed to get memory usage by pid: %d", pid)
	}

	defer f.Close()

	buff := bufio.NewReader(f)

	line, err := buff.ReadString('\n')
	cpu_metrics := strings.Fields(line)

	rss, _ := strconv.Atoi(cpu_metrics[23])

	return rss * pageSize, nil

}

func RandString(n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	rand.Seed(time.Now().UnixNano())
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterAlphabetsAndNumbers) {
			b[i] = letterAlphabetsAndNumbers[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

func RandNumberString(n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	rand.Seed(time.Now().UnixNano())
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterNumbers) {
			b[i] = letterNumbers[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

func RandAlphabetString(n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	rand.Seed(time.Now().UnixNano())
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterAlphabets) {
			b[i] = letterAlphabets[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

func Sha1Sum(data string) (sum string) {
	//产生一个散列值得方式是 sha1.New()，sha1.Write(bytes)，然后 sha1.Sum([]byte{})。这里我们从一个新的散列开始。
	h := sha1.New()
	//写入要处理的字节。如果是一个字符串，需要使用[]byte(s) 来强制转换成字节数组。
	h.Write([]byte(data))
	//这个用来得到最终的散列值的字符切片。Sum 的参数可以用来都现有的字符切片追加额外的字节切片：一般不需要要。
	bs := h.Sum(nil)
	//SHA1 值经常以 16 进制输出，例如在 git commit 中。使用%x 来将散列结果格式化为 16 进制字符串。
	sum = fmt.Sprintf("%x", bs)
	return
}

func SplitToIntSlice(str string) (slice []int) {
	str = strings.TrimSpace(str)
	string_slice := strings.Split(str, ",")
	slice = make([]int, 0, len(string_slice))

	for _, item := range string_slice {
		i, err := strconv.Atoi(strings.TrimSpace(item))
		if err == nil {
			slice = append(slice, i)
		}
	}

	return
}

func UniqDup(data interface{}) (interface{}, error) {
	slice := reflect.ValueOf(data)
	if slice.Kind() != reflect.Slice && slice.Kind() != reflect.Array {
		return data, errors.New("data is not a slice or array")
	}

	_map := make(map[interface{}]struct{})
	uniq_slice := reflect.MakeSlice(slice.Type(), 0, slice.Len())
	in := make([]reflect.Value, 0)

	for i := 0; i < slice.Len(); i++ {
		item := slice.Index(i)
		_, ok := reflect.TypeOf(item.Interface()).MethodByName("UniqId")

		if ok {
			retval := item.MethodByName("UniqId").Call(in)
			if len(retval) == 1 && retval[0].Kind() == reflect.Int {
				id := retval[0].Int()

				if _, ok := _map[id]; !ok {
					uniq_slice = reflect.Append(uniq_slice, slice.Index(i))
					_map[id] = struct{}{}
				}
			} else {
				return struct{}{}, errors.New("func (UniqId) return value not int")
			}
		} else {
			if _, ok := _map[item.Interface()]; !ok {
				uniq_slice = reflect.Append(uniq_slice, slice.Index(i))
				_map[item.Interface()] = struct{}{}
			}
		}
	}

	return uniq_slice.Interface(), nil
}

func Uniq(data interface{}) (err error) {

	value := reflect.ValueOf(data)
	if value.Kind() != reflect.Ptr {
		err = errors.New("not a pointer")
		return
	}

	uniq_slice, err := UniqDup(value.Elem().Interface())
	if err != nil {
		return
	}
	ptr := reflect.ValueOf(uniq_slice)

	value.Elem().Set(ptr)
	return nil
}

// 调试时避免烦人的未使用变量提示
func Unused(args ...interface{}) {
	for _, arg := range args {
		fmt.Printf("unused var: %v(%t)\n", arg, arg)
	}
}
