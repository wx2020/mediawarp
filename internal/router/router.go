package router

import (
	"MediaWarp/constants"
	"MediaWarp/internal/config"
	"MediaWarp/internal/handler"
	"MediaWarp/internal/logging"
	"MediaWarp/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	ginR := gin.New()
	ginR.Use(
		middleware.Logger(),
		middleware.Recovery(),
		middleware.SetRefererPolicy(constants.SameOrigin),
	)

	if config.ClientFilter.Enable {
		ginR.Use(middleware.ClientFilter())
		logging.Info("客户端过滤中间件已启用")
	} else {
		logging.Info("客户端过滤中间件未启用")
	}

	mediawarpRouter := ginR.Group("/MediaWarp")
	{
		mediawarpRouter.Any("/version", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, config.Version())
		})
		if config.Web.Enable { // 启用 Web 页面修改相关设置
			if config.Web.Custom { // 用户自定义静态资源目录
				mediawarpRouter.Static("/custom", config.CostomDir())
				logging.Info("使用自定义静态资源目录: ", config.CostomDir())
			} else {
				// 当没有自定义目录且没有内嵌资源时，提示用户
				logging.Info("Web功能已启用，但未配置自定义静态资源目录。")
				logging.Info("如需使用Web美化功能（如actor-plus、emby-swiper等），请在config.yaml中设置 web.custom: true 并提供自定义资源目录。")
			}
			if config.Web.Robots != "" { // 自定义 robots.txt
				ginR.GET(
					"/robots.txt",
					func(ctx *gin.Context) {
						ctx.String(http.StatusOK, config.Web.Robots)
					},
				)
			}
		}
	}

	handlers := make(gin.HandlersChain, 0, 3)

	handlers = append(handlers, getRegexpRouterHandler())
	ginR.NoRoute(handlers...)
	return ginR
}

// 正则表达式路由处理器
//
// 从媒体服务器处理结构体中获取正则路由规则
// 依次匹配请求, 找到对应的处理器
func getRegexpRouterHandler() gin.HandlerFunc {
	mediaServerHandler := handler.GetMediaServer()
	middlewareChain := NewMiddlewareChain().
		Add(QueryKeyCaseInsensitive).
		Add(DisableCompression)

	return func(ctx *gin.Context) {
		for _, rule := range mediaServerHandler.GetRegexpRouteRules() {
			if rule.Regexp.MatchString(ctx.Request.URL.Path) { // 不带查询参数的字符串：/emby/Items/54/Images/Primary
				logging.AccessDebugf(ctx, "匹配成功正则表达式: %s", rule.Regexp.String())

				middlewareChain.Execute(rule.Handler)(ctx)
				return
			}
		}

		// 未匹配路由
		mediaServerHandler.ReverseProxy(ctx.Writer, ctx.Request)
	}
}
