package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logc"
	"watchAlert/pkg/response"
)

func Service(ctx *gin.Context, fu func() (interface{}, interface{})) {
	data, err := fu()
	if err != nil {
		logc.Error(context.Background(), err)
		response.Fail(ctx, err.(error).Error(), "failed")
		ctx.Abort()
		return
	} else {
		response.Success(ctx, data, "success")
	}
}

func BindJson(ctx *gin.Context, req interface{}) error {
	err := ctx.ShouldBindJSON(req)
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		ctx.Abort()
		return err
	}
	return nil
}

func BindQuery(ctx *gin.Context, req interface{}) error {
	err := ctx.ShouldBindQuery(req)
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		ctx.Abort()
		return err
	}
	return nil
}

func ValidateAndRun(ctx *gin.Context, req interface{}, fu func() (interface{}, interface{})) {
	if err := BindJson(ctx, req); err != nil {
		return
	}
	Service(ctx, fu)
}

func ValidateAndRunQuery(ctx *gin.Context, req interface{}, fu func() (interface{}, interface{})) {
	if err := BindQuery(ctx, req); err != nil {
		return
	}
	Service(ctx, fu)
}

func HandleError(ctx *gin.Context, err error) {
	if err != nil {
		logc.Error(context.Background(), err)
		response.Fail(ctx, err.Error(), "failed")
		ctx.Abort()
	}
}
