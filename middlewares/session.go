package middlewares

import (
	"math/rand"
	"net/http"
	"webserver/models"
	"webserver/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next() // Step1: Process the request first.

		// Step2: Check if any errors were added to the context
		if len(ctx.Errors) > 0 {
			errors := make([]string, 0)
			for _, v := range ctx.Errors {
				errors = append(errors, v.Error())
			}

			if ctx.Request.Header.Get("X-Requested-With") == "XMLHttpRequest" {
				ctx.JSON(http.StatusOK, gin.H{
					"errors": errors,
				})
			} else {
				ctx.HTML(http.StatusOK, "errors.html", gin.H{
					"errors": errors,
				})
			}
		}

		// Any other steps if no errors are found
	}
}

func login(ctx *gin.Context) {
	//判断出是否是ajax请求
	if ctx.Request.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		ctx.JSON(http.StatusOK, gin.H{
			"redirect_url": "/",
		})
	} else {
		rCode := rand.Intn(900000) + 100000
		ctx.HTML(http.StatusOK, "login.html", gin.H{
			"feishu": utils.AppConfig.Feishu,
			"rCode":  rCode,
		})
	}

	ctx.Abort()
}

func AuthSessionMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := getSession(ctx)
		if userId == 0 {
			login(ctx)
			return
		}

		authUser, _ := models.GetUserByID(userId)
		if authUser == nil {
			login(ctx)
			return
		}

		ctx.Set("authUser", authUser)

		templateData := make(map[string]any)
		templateData["authUser"] = authUser
		templateData["appConfig"] = utils.AppConfig
		ctx.Set("templateData", templateData)

		ctx.Next()
	}
}

func EnableCookieSession() gin.HandlerFunc {
	// 创建基于 cookie 的存储引擎
	store := cookie.NewStore([]byte("worktime"))

	store.Options(sessions.Options{
		MaxAge:   3600 * 365,           // 有效期
		Path:     "/",                  // Cookie生效路径：根路径表示全站所有接口都能携带该Cookie
		HttpOnly: true,                 // 禁止客户端JS访问Cookie（防XSS攻击，避免Cookie被脚本窃取）
		Secure:   false,                // 仅HTTPS环境下生效（开发环境设false，生产环境必须设true）
		SameSite: http.SameSiteLaxMode, // 限制跨域请求携带Cookie（防CSRF攻击）
	})

	// 注册Session中间件（name为"session_id"的 Cookie，来存储sessionID）
	// 表示所有请求都会经过这个中间件处理，自动检查请求是否携带名为“session_iD”的Cookie，若有从store中加载对应session数据
	// 若没有则生成新的SessionID并创建空的Session实例，请求处理完自动将Session数据同步到存储引擎，发送SessionID给客户端

	return sessions.Sessions("session_id", store)
}

// register and login will save session
func SaveSession(ctx *gin.Context, userId int) {
	session := sessions.Default(ctx)
	session.Set("UserId", userId)
	session.Save()
}

// logout will clear session
func ClearSession(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Clear()
	session.Save()
}

func HasSession(ctx *gin.Context) bool {
	session := sessions.Default(ctx)
	if val := session.Get("UserId"); val == nil {
		return false
	}
	return true
}

func getSession(ctx *gin.Context) int {
	session := sessions.Default(ctx)
	val := session.Get("UserId")
	if val == nil {
		return 0
	}

	if val.(int) == 0 {
		return 0
	}
	return val.(int)
}
