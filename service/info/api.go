package info

import (
	"crypto/rsa"
	"github.com/gin-gonic/gin"
)

var Engine *gin.Engine
var JwtPublicKey *rsa.PublicKey
