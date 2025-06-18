package httpctx

type Context interface {
	BindJSON(obj any) error
	JSON(code int, obj any)
	Query(key string) string
	Param(key string) string
}
