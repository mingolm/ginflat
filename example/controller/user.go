package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/mingolm/ginflat"
	"github.com/mingolm/ginflat/example/middleware"
	"github.com/mingolm/ginflat/httperrors"
)

type User struct {
	Base
}

func (c *User) InitRouter(r ginflat.Router) {
	g := r.Group("/user").Use(middleware.UserMiddleware())
	g.Get("/detail", c.getDetail)
	g.Post("/delete", c.delete)
}

type (
	GetUserDetailRequest struct {
		Id   uint64 `form:"id" binding:"required"`
		Name string `form:"name" binding:"required"`
	}
	GetUserDetailResponse struct {
		Id   uint64 `json:"id"`
		Name string `json:"name"`
	}
)

func (c *User) getDetail(ctx *gin.Context, req *GetUserDetailRequest) (resp *GetUserDetailResponse, err error) {
	if req.Id > 10 {
		// return nil, httperrors.ErrInvalidArgument.Msg("user id")
		return nil, httperrors.ErrNotFound.Msg("user not found").ErrorType("USER_NOT_EXIST")
		// return nil, httperrors.NewWithCode(http.StatusBadGateway, "test error").ErrorType("AA")
		// return nil, httperrors.ErrNotFound.Localize("user.not_exist")
		// return nil, httperrors.ErrNotFound.ErrorType("USER_NOT_EXIST").Msg("user not exist")
		// return nil, httperrors.NewWithCode(http.StatusNotFound, "user not exist")
	}
	return &GetUserDetailResponse{
		Id:   req.Id,
		Name: req.Name,
	}, nil
}

type (
	DeleteRequest struct {
		Id uint64 `form:"id" binding:"required"`
	}
	DeleteResponse struct {
		IsDeleted bool `json:"is_deleted"`
	}
)

func (c *User) delete(ctx *gin.Context, req *DeleteRequest) (resp *DeleteResponse, err error) {
	return &DeleteResponse{
		IsDeleted: req.Id == 10,
	}, nil
}
