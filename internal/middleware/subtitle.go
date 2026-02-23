package middleware

import (
	"MediaWarp/internal/config"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
)

func SubtitleCache(ttl time.Duration, reg *regexp.Regexp) gin.HandlerFunc {
	cachePool, err := config.CreateOptimizedCache(ttl)
	if err != nil {
		panic(fmt.Sprintf("create subtitle cache pool failed: %v", err))
	}
	cacheFunc := getCacheBaseFunc(cachePool, "字幕", reg.String())

	return func(ctx *gin.Context) {
		if ctx.Request.Method != http.MethodGet || !reg.MatchString(ctx.Request.URL.Path) {
			ctx.Next()
			return
		}
		cacheFunc(ctx)
	}
}
