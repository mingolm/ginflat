package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mingolm/ginflat"
	"github.com/mingolm/ginflat/httperrors"
)

type User struct {
	Base
}

func (c *User) InitRouter(r ginflat.Router) {
	g := r.Group("/user")
	g.Get("/get", c.get)
	g.Post("/add", c.add)
	g.Put("/update", c.update)
	g.Delete("/delete", c.delete)
}

type (
	GetRequest struct {
		Id uint64 `form:"id" binding:"required"`
	}
	GetResponse struct {
		Id   uint64 `json:"id"`
		Name string `json:"name"`
	}
)

func (c *User) get(ctx *gin.Context, req *GetRequest) (resp *GetResponse, err error) {
	key := fmt.Sprintf("user:%d", req.Id)
	userVal, ok := c.dbs.Load(key)
	if !ok {
		return nil, httperrors.ErrNotFound.Msg("user %d not found", req.Id)
	}

	userRow := userVal.(*UserModel)
	return &GetResponse{
		Id:   userRow.Id,
		Name: userRow.Name,
	}, nil
}

type (
	AddRequest struct {
		Id   uint64 `form:"id" binding:"required"`
		Name string `form:"name" binding:"required"`
	}
	AddResponse struct{}
)

func (c *User) add(ctx *gin.Context, req *AddRequest) (resp *AddResponse, err error) {
	key := fmt.Sprintf("user:%d", req.Id)

	if _, ok := c.dbs.Load(key); ok {
		return nil, httperrors.ErrAlreadyExists.Msg("user %d already exist", req.Id)
	}

	c.dbs.Store(fmt.Sprintf("user:%d", req.Id), &UserModel{
		Id:   req.Id,
		Name: req.Name,
	})

	return &AddResponse{}, nil
}

type (
	UpdateRequest struct {
		Id   uint64 `form:"id" binding:"required"`
		Name string `form:"name" binding:"required"`
	}
	UpdateResponse struct{}
)

func (c *User) update(ctx *gin.Context, req *UpdateRequest) (resp *UpdateResponse, err error) {
	key := fmt.Sprintf("user:%d", req.Id)

	if _, ok := c.dbs.Load(key); !ok {
		return nil, httperrors.ErrNotFound.Msg("user %d not found", req.Id)
	}
	c.dbs.Store(key, &UserModel{
		Id:   req.Id,
		Name: req.Name,
	})

	return &UpdateResponse{}, nil
}

type (
	DeleteRequest struct {
		Id uint64 `form:"id" binding:"required"`
	}
	DeleteResponse struct{}
)

func (c *User) delete(ctx *gin.Context, req *DeleteRequest) (resp *DeleteResponse, err error) {
	key := fmt.Sprintf("user:%d", req.Id)

	if _, ok := c.dbs.Load(key); !ok {
		return nil, httperrors.ErrNotFound.Msg("user %d not found", req.Id)
	}
	c.dbs.Delete(key)

	return &DeleteResponse{}, nil
}
