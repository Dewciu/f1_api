package common

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func Bind(c *gin.Context, obj interface{}) error {
	b := binding.Default(c.Request.Method, c.ContentType())
	return c.ShouldBindWith(obj, b)
}

func GetColumnFromUniqueErrorDetails(details string) string {
	var ret string
	flag := false
	for _, char := range details {
		if string(char) == ")" {
			flag = false
			break
		}

		if flag {
			ret += string(char)
		}

		if string(char) == "(" {
			flag = true
		}
	}
	return ret
}
