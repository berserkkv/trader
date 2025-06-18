package ctxImpl

import "github.com/gin-gonic/gin"

type GinContext struct {
	C *gin.Context
}

func (g *GinContext) BindJSON(obj any) error {
	return g.C.BindJSON(obj)
}
func (g *GinContext) JSON(code int, obj any) {
	g.C.JSON(code, obj)
}
func (g *GinContext) Query(key string) string {
	return g.C.Query(key)
}
func (g *GinContext) Param(key string) string {
	return g.C.Param(key)
}
