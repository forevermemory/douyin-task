package router

import (
	"douyin/web/db"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func route(f func(ctx *gin.Context) interface{}) gin.HandlerFunc {
	return func(context *gin.Context) {
		context.JSON(200, f(context))
	}
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method               // 请求方法
		origin := c.Request.Header.Get("Origin") // 请求头部
		var headerKeys []string                  // 声明请求头keys
		for k := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		// fmt.Println("origin: ", origin)
		if origin != "" {
			// c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Origin", origin) // 这是允许访问所有域

			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE") // 服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
			// header的类型
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			// 允许跨域设置                                                                                                      可以返回其他子段
			c.Writer.Header().Del("Access-Control-Expose-Headers")                                                                                                                                                 // 跨域关键设置 让浏览器可以解析
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar") // 跨域关键设置 让浏览器可以解析
			// c.Header("Access-Control-Max-Age", "0")                                                                                                                                                        // 缓存请求信息 单位为秒
			c.Header("Access-Control-Allow-Credentials", "true") //  跨域请求是否需要带cookie信息 默认设置为true

		}
		/*c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")*/
		// c.Header("Content-type", "application/json") // 设置返回格式是json
		c.Header("Accept", "video/*,text/html,application/xhtml+xml,application/xml,application/json;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
		// 放行所有OPTIONS方法
		//c.Header("Content-Encoding","gzip")
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}

		c.Next() //  处理请求
	}
}

// JWTAuth 中间件，检查token
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var param struct {
			Token string `header:"token"`
		}
		if err := c.ShouldBindHeader(&param); err != nil {
			c.Abort()
			return
		}

		token := param.Token
		// 1. 是否携带token
		if token == "" {
			c.JSON(http.StatusForbidden, gin.H{
				"code": -1,
				"msg":  -1,
			})
			c.Abort()
			return
		}

		j := db.NewJWT()
		// parseToken 解析token包含的信息
		// 2.解析token失败
		claims, err := j.ParseToken(token)
		if claims == nil {
			if err != nil {
				c.JSON(http.StatusForbidden, gin.H{
					"code": -1,
					"msg":  -1,
				})
				c.Abort()
				return
			}

			c.JSON(http.StatusForbidden, gin.H{
				"code": -1,
				"msg":  -1,
			})
			c.Abort()
			return
		}
		// 3.token 是否有效 是否过期 TODO
		// if ok, err := isTokenExist(token, claims.UserID); err != nil || !ok {
		// 	c.JSON(http.StatusForbidden, gin.H{
		// 		"code": -1,
		// 		"msg":  "token已过期，请重新登录",
		// 	})
		// 	c.Abort()
		// 	return
		// }
		//
		// 继续交由下一个路由处理,并将解析出的信息传递下去

		c.Set("token", claims.Token)
		// c.Set("username", claims.Username)
		// c.Set("is_vip", claims.IsVip)

	}
}
