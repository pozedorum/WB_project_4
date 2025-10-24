package server

import (
	"net/http"
	"net/http/pprof"

	"github.com/pozedorum/wbf/ginext"
)

func (ss *ShortURLServer) EnableProfiling(pprofGroup *ginext.RouterGroup) {
	pprofGroup.GET("/", func(c *ginext.Context) {
		http.Redirect(c.Writer, c.Request, "/debug/pprof/", http.StatusFound)
	})
	pprofGroup.GET("/:profile", func(c *ginext.Context) {
		profile := c.Param("profile")
		switch profile {
		case "profile":
			pprof.Profile(c.Writer, c.Request)
		case "trace":
			pprof.Trace(c.Writer, c.Request)
		case "heap":
			pprof.Handler("heap").ServeHTTP(c.Writer, c.Request)
		case "goroutine":
			pprof.Handler("goroutine").ServeHTTP(c.Writer, c.Request)
		case "allocs":
			pprof.Handler("allocs").ServeHTTP(c.Writer, c.Request)
		case "block":
			pprof.Handler("block").ServeHTTP(c.Writer, c.Request)
		case "mutex":
			pprof.Handler("mutex").ServeHTTP(c.Writer, c.Request)
		default:
			pprof.Index(c.Writer, c.Request)
		}
	})
}
