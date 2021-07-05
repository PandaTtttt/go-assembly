package main

import (
	"errors"
	"fmt"
	"github.com/PandaTtttt/go-assembly/errs"
)

const (
	// 自定义的业务错误类型，在项目中自行添加。
	ErrWrongName errs.RetCode = 100001
)

func main() {
	err := getStandardError()
	fmt.Println(errs.NotFound.Is(err)) // false
	fmt.Println(err.Error())           // this is a standard error

	err = getCustomError()
	fmt.Println(errs.NotFound.Is(err)) // true

	fmt.Println(err.Error())              // NotFound[this is a custom error]
	fmt.Println(err.(*errs.Error).Json()) // {"retCode":1,"retMsg":"this is a custom error"}

	err = getBusinessError()
	fmt.Println(ErrWrongName.Is(err)) // true

	fmt.Println(err.Error())              // Code(100001)[this is a business error]
	fmt.Println(err.(*errs.Error).Json()) // {"retCode":100001,"retMsg":"this is a business error"}
}

func getStandardError() error {
	return errors.New("this is a standard error")
}

func getCustomError() error {
	return errs.NotFound.New("this is a custom error")
}

func getBusinessError() error {
	return ErrWrongName.New("this is a business error")
}
