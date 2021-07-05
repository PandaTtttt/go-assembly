package api

import (
	"github.com/PandaTtttt/go-assembly/errs"
	"github.com/PandaTtttt/go-assembly/util/m"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 封装并发送200结果信息
func ResultOK(c *gin.Context, data interface{}) {
	c.Writer.Header().Set("Cache-Control", "no-store")
	c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	if data == nil {
		c.String(http.StatusOK, "")
	} else {
		c.JSON(http.StatusOK, data)
	}

	c.Abort()
}

// 封装并发送400错误的信息
func Result400(c *gin.Context, err error) {
	myErr, ok := err.(*errs.Error)
	if !ok {
		c.String(http.StatusBadRequest, err.Error())
		c.Abort()
		return
	}
	if myErr.RetMsg == "" {
		myErr.RetMsg = myErr.Error()
	}

	c.Writer.Header().Set("Cache-Control", "no-store")
	c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")

	c.String(http.StatusBadRequest, myErr.Json())
	c.Abort()
}

// 封装并发送429错误的信息
func Result429(c *gin.Context, msg interface{}) {
	c.Writer.Header().Set("Cache-Control", "no-store")
	c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")

	c.JSON(http.StatusTooManyRequests, m.M{"message": msg})
	c.Abort()
}
