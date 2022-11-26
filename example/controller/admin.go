package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/mingolm/ginflat"
)

type Admin struct {
	Base
}

func (c *Admin) InitRouter(r ginflat.Router) {
	g := r.Group("/admin")
	g.Get("/detail", c.getDetail)
}

type (
	GetAdminDetailRequest struct {
	}
	GetAdminDetailResponse struct {
		Id   uint64 `json:"id"`
		Name string `json:"name"`
	}
)

func (c *Admin) getDetail(ctx *gin.Context, req *GetAdminDetailRequest) (resp *GetAdminDetailResponse, err error) {
	return &GetAdminDetailResponse{
		Id:   1,
		Name: "admin",
	}, nil
}
