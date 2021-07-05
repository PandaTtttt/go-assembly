package util

import (
	"fmt"
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
