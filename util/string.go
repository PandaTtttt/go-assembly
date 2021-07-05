package util

import (
	"bufio"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"github.com/PandaTtttt/go-assembly/util/must"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unsafe"
)

func FormatConfig(path string) map[string]string {
	config := make(map[string]string)

	f, err := os.Open(path)
	defer must.Close(f)
	must.Must(err)

	r := bufio.NewReader(f)
	for {
		b, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			must.Must(err)
		}
		s := strings.TrimSpace(string(b))
		index := strings.Index(s, "=")
		if index < 0 {
			continue
		}
		key := strings.TrimSpace(s[:index])
		if len(key) == 0 {
			continue
		}
		value := strings.TrimSpace(s[index+1:])
		if len(value) == 0 {
			continue
		}
		config[key] = value
	}
	return config
}

func NameWithoutExt(path string) string {
	name := filepath.Base(path)
	i := strings.LastIndex(name, ".")
	if i == 0 {
		return name
	}
	return name[:i]
}

func RandomString(l int, t int) string {
	str1 := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	str2 := "0123456789abcdefghijklmnopqrstuvwxyz"
	str3 := "abcdefghijklmnopqrstuvwxyz"
	str4 := "0123456789"
	var bytes []byte
	switch t {
	case 1:
		bytes = []byte(str1)
	case 2:
		bytes = []byte(str2)
	case 3:
		bytes = []byte(str3)
	case 4:
		bytes = []byte(str4)
	default:
		bytes = []byte(str1)
	}

	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func Md5String(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func Sha256String(data []byte, key []byte) string {
	m := hmac.New(sha256.New, key)
	m.Write(data)
	return hex.EncodeToString(m.Sum(nil))
}

func GetSuffixBySep(s, sep string) string {
	i := strings.LastIndex(s, sep)
	return s[i+1:]
}

func ParseBool(s string) bool {
	s = strings.ToLower(s)
	if s == "true" {
		return true
	}
	return false
}

// BytesToString converts byte slice to string.
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// StringToBytes converts string to byte slice.
func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}
