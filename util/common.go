package util

import (
	"fmt"
	"runtime"
	"sort"
)

// CheckSign 签名验证
func CheckSign(data map[string]interface{}, privateKey string) bool {
	sign := data["sign"]
	delete(data, "sign")
	getSign := Sign(data, privateKey)
	return sign == getSign
}

// Sign 获取sign
func Sign(data map[string]interface{}, privateKey string) string {
	var keyList []string
	var sortString string
	for k := range data {
		keyList = append(keyList, k)
	}
	sort.Strings(keyList)
	for _, v := range keyList {
		sortString += fmt.Sprintf("%s=%v&", v, data[v])
	}
	sortString += fmt.Sprintf("key=%v", privateKey)
	return Md5String(sortString)
}

// RunFuncName 获取正在运行的函数名
func RunFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}
