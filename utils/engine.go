package utils

import (
	"bytes"
	"net/http"

	"github.com/gin-gonic/gin"
)

var globalRouter *gin.Engine

// 添加一个函数用于获取全局 router 实例
func GetGlobalRouter() *gin.Engine {
	return globalRouter
}

func SetGlobalRouter(engine *gin.Engine) {
	globalRouter = engine
}

// 自定义结构体实现 http.ResponseWriter 接口，用于渲染模板到 buffer
type bufferResponseWriter struct {
	*bytes.Buffer
}

func (b *bufferResponseWriter) Header() http.Header {
	return make(http.Header)
}

func (b *bufferResponseWriter) WriteHeader(statusCode int) {
	// 忽略状态码写入
}

func GetRenderedTemplateContent(ctx *gin.Context, name string) string {
	h := ctx.MustGet("templateData").(map[string]any)

	return RenderedTemplateContent(ctx, name, h)
}

func RenderedTemplateContent(ctx *gin.Context, name string, data any) string {
	// 创建一个 buffer 来捕获渲染结果
	buf := new(bytes.Buffer)

	// 使用模板实例渲染到 buffer
	tmpl := globalRouter.HTMLRender.Instance(name, data)
	if tmpl == nil {
		return "template not found: " + name
	}

	// 手动构造一个 fake ResponseWriter 实现 Render 接口所需参数
	// 注意：Render 方法需要 http.ResponseWriter，但我们可以用 buf 伪装写入内容
	err := tmpl.Render(&bufferResponseWriter{buf})
	if err != nil {
		return err.Error()
	}

	return buf.String()
}
