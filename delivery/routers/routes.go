// Delivery/routers/router.go
package routers

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
    router := gin.Default()

    SetUpUser(router)
        
    SetUpAdmin(router)

    return router
}