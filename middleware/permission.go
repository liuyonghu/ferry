package middleware

import (
	"ferry/global/orm"
	"ferry/models/system"
	mycasbin "ferry/pkg/casbin"
	"ferry/pkg/jwtauth"
	_ "ferry/pkg/jwtauth"
	"ferry/pkg/logger"
	"ferry/tools"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

//权限检查中间件
func AuthCheckRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		var menuValue system.Menu

		data, _ := c.Get("JWT_PAYLOAD")
		v := data.(jwtauth.MapClaims)
		e, err := mycasbin.Casbin()
		tools.HasError(err, "", 500)
		//检查权限
		res, err := e.Enforce(v["rolekey"], c.Request.URL.Path, c.Request.Method)
		logger.Info(v["rolekey"], c.Request.URL.Path, c.Request.Method)

		tools.HasError(err, "", 500)

		err = orm.Eloquent.Model(&menuValue).
			Where("path = ? and action = ?", c.Request.URL.Path, c.Request.Method).
			Find(&menuValue).Error
		tools.HasError(err, "", 500)

		if res {
			c.Next()
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code": 403,
				"msg":  fmt.Sprintf("对不起，您没有 <%v> 访问权限，请联系管理员", menuValue.Title),
			})
			c.Abort()
			return
		}
	}
}
