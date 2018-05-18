package mfe

import (
	"sort"

	fasthttp "github.com/valyala/fasthttp"
)

var rDir map[string]func(ctx *fasthttp.RequestCtx)
var rMap map[int]map[string]func(ctx *fasthttp.RequestCtx)
var lMap map[int]struct{}
var lS []int

var defRoute func(ctx *fasthttp.RequestCtx)

// AddDirectRoute -- Добавить роут точное имя
func AddDirectRoute(uri string, f func(ctx *fasthttp.RequestCtx)) {
	rDir[uri] = f
}

// AddRoute -- Добавить роут
func AddRoute(uri string, f func(ctx *fasthttp.RequestCtx)) {
	l := len(uri)
	if _, ok := lMap[l]; !ok {
		lMap[l] = struct{}{}
		rMap[l] = map[string]func(ctx *fasthttp.RequestCtx){}
		lS = append(lS, l)
		sort.Slice(lS, func(i, j int) bool { return lS[i] > lS[j] })
	}

	rMap[l][uri] = f
}

// AddDefaultRoute -- Добавить роут по умолчанию
func AddDefaultRoute(f func(ctx *fasthttp.RequestCtx)) {
	defRoute = f
}

func searchRoute(uri string) (fn func(ctx *fasthttp.RequestCtx)) {

	if f, ok := rDir[uri]; ok {
		return f
	}

	l := len(uri)

	for _, v := range lS {
		if v <= l {
			s := uri[:v]
			m := rMap[v]
			if f, ok := m[s]; ok {
				return f
			}
		}
	}

	return defRoute
}

// ListenAndServe Запуск fastHttp сервиса
func ListenAndServe(listenAddr string, defaultHandler func(ctx *fasthttp.RequestCtx)) {
	AddDefaultRoute(defaultHandler)
	fasthttp.ListenAndServe(listenAddr, func(ctx *fasthttp.RequestCtx) {

		f := searchRoute(string(ctx.Path()))

		f(ctx)
	})
}
