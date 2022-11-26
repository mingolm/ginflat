package ginflat

import "github.com/gin-gonic/gin"

type Middleware = gin.HandlerFunc

type MiddlewareChain = gin.HandlersChain
