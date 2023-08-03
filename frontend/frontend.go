package frontend

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var Templates = []string{
	"frontend/templates/index.tmpl",
}

func Index(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index.tmpl", nil)
}
