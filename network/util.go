package network

import (
	"chat_socket_server/types"
	"github.com/gin-gonic/gin"
)

func response(c *gin.Context, s int, res interface{}, data ...string) {
	c.JSON(s, types.NewRes(s, res, data...))
}
