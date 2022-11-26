package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mingolm/ginflat"
	"github.com/mingolm/ginflat/example/middleware"
)

type User struct {
	Base
}

func (c *User) InitRouter(r ginflat.Router) {
	g := r.Group("/user").Use(middleware.UserMiddleware())
	g.Get("/detail", c.getDetail)
}

type (
	GetUserDetailRequest struct {
		Id uint64 `form:"id" binding:"required"`
	}
	GetUserDetailResponse struct {
		Id   uint64 `json:"id"`
		Name string `json:"name"`
	}
)

func (c *User) getDetail(ctx *gin.Context, req *GetUserDetailRequest) (resp *GetUserDetailResponse, err error) {
	if req.Id != 1 {
		fmt.Println("user return")
		return nil, errors.New("user not found")
	}
	return &GetUserDetailResponse{
		Id:   1,
		Name: "user",
	}, nil
}
