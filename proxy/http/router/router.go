package router

import (
	"net/http"

	"github.com/laohanlinux/riot/proxy/http/controller"
	"github.com/laohanlinux/riot/proxy/http/middleware"

	macaron "gopkg.in/macaron.v1"
)

func NewRouter(token ...string) http.Handler {
	macaronRouter := macaron.Classic()
	{
		macaronRouter.Group("/riot/bucket", func() {
			// bucket router
			macaronRouter.Post("/", controller.CreateBucket)
			macaronRouter.Get("/:bucket", controller.BucketInfo)
			macaronRouter.Delete("/:bucket", controller.DelBucket)
			// kv router
			macaronRouter.Group("/:bucket", func() {
				macaronRouter.Get("/key/:key", controller.GetValue)
				macaronRouter.Post("/key/:key", controller.SetValue)
				macaronRouter.Delete("/key/:key", controller.DelValue)
			})
		}, middleware.ContextMiddleware, middleware.AuthorMiddleware, middleware.OutputMiddleware)
	}

	// adm router
	{
		macaronRouter.Group("/riot/admin", func() {
			macaronRouter.Get("/leader", controller.Leader)
			macaronRouter.Get("/states", controller.States)
			macaronRouter.Get("/peers", controller.Peers)
			macaronRouter.Get("/snapshot", controller.SnapshotInfo)
			macaronRouter.Post("/remove", controller.RemovePeer)
			macaronRouter.Get("/router-test", controller.RouterTest)
		}, middleware.ContextMiddleware, middleware.AuthorMiddleware, middleware.OutputMiddleware)
	}
	macaronRouter.NotFound(controller.Contr404)
	macaronRouter.Handlers(middleware.ElaspRequest)

	if len(token) > 0 {
		macaronRouter.Map(token[0])
	} else {
		macaronRouter.Map("")
	}
	return macaronRouter
}
