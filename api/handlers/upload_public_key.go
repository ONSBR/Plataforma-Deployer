package handlers

import (
	"io/ioutil"

	"github.com/ONSBR/Plataforma-Deployer/models/exceptions"

	"github.com/ONSBR/Plataforma-Deployer/actions"

	"github.com/gin-gonic/gin"
)

//UploadPublicKey uploads users's public key to git-server
func UploadPublicKey(c *gin.Context) {
	buf := c.Request.Body
	data, err := ioutil.ReadAll(buf)
	if err != nil {
		c.JSON(400, exceptions.NewInvalidArgumentException(err))
		return
	}
	info, ex := actions.InstallPublicKey(data, c.Param("solution"), c.Param("filename"))
	if ex != nil {
		c.JSON(ex.Status(), ex)
		return
	}
	c.JSON(200, info)
}
