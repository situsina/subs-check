package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bestruirui/bestsub/internal/database/op"
	storageModel "github.com/bestruirui/bestsub/internal/models/storage"
	"github.com/bestruirui/bestsub/internal/modules/storage"
	"github.com/bestruirui/bestsub/internal/server/middleware"
	"github.com/bestruirui/bestsub/internal/server/resp"
	"github.com/bestruirui/bestsub/internal/server/router"
	"github.com/bestruirui/bestsub/internal/utils/log"
	"github.com/gin-gonic/gin"
)

func init() {
	router.NewGroupRouter("/api/v1/storage").
		Use(middleware.Auth()).
		AddRoute(
			router.NewRoute("/", router.POST).
				Handle(createStorage),
		).
		AddRoute(
			router.NewRoute("/", router.GET).
				Handle(getStorage),
		).
		AddRoute(
			router.NewRoute("/channel", router.GET).
				Handle(getStorageChannel),
		).
		AddRoute(
			router.NewRoute("/channel/config", router.GET).
				Handle(getStorageChannelConfig),
		).
		AddRoute(
			router.NewRoute("/", router.PUT).
				Handle(updateStorage),
		).
		AddRoute(
			router.NewRoute("/:id", router.DELETE).
				Handle(deleteStorage),
		)
}

// @Summary 创建存储
// @Description 创建存储
// @Tags 存储
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body storageModel.CreateRequest true "存储配置数据"
// @Success 200 {object} resp.SuccessStruct{data=storageModel.Response} "创建成功"
// @Failure 401 {object} resp.ErrorStruct "未授权"
// @Failure 500 {object} resp.ErrorStruct "服务器内部错误"
// @Router /api/v1/storage [post]
func createStorage(c *gin.Context) {
	var req storageModel.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Errorf("createStorage: %v", err)
		resp.ErrorBadRequest(c)
		return
	}
	configBytes, err := json.Marshal(req.Config)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	configStr := string(configBytes)
	data := storageModel.Data{
		Name:   req.Name,
		Type:   req.Type,
		Config: configStr,
	}
	if err := op.CreateStorage(c.Request.Context(), &data); err != nil {
		resp.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	resp.Success(c, storageModel.Response{
		ID:     data.ID,
		Name:   data.Name,
		Type:   data.Type,
		Config: req.Config,
	})
}

// @Summary 获取存储
// @Description 获取存储
// @Tags 存储
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} resp.SuccessStruct{data=[]storageModel.Response} "获取成功"
// @Failure 401 {object} resp.ErrorStruct "未授权"
// @Failure 500 {object} resp.ErrorStruct "服务器内部错误"
// @Router /api/v1/storage [get]
func getStorage(c *gin.Context) {
	storages, err := op.GetStorageList(c.Request.Context())
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	var result = make([]storageModel.Response, 0, len(storages))
	for _, v := range storages {
		var config any
		if err := json.Unmarshal([]byte(v.Config), &config); err != nil {
			resp.Error(c, http.StatusInternalServerError, err.Error())
			return
		}
		result = append(result, storageModel.Response{
			ID:     v.ID,
			Name:   v.Name,
			Type:   v.Type,
			Config: config,
		})
	}
	resp.Success(c, result)
}

// getStorageChannel 获取存储渠道
// @Summary 获取存储渠道
// @Tags 存储
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} resp.SuccessStruct{data=[]string} "获取成功"
// @Failure 400 {object} resp.ErrorStruct "请求参数错误"
// @Failure 401 {object} resp.ErrorStruct "未授权"
// @Failure 500 {object} resp.ErrorStruct "服务器内部错误"
// @Router /api/v1/storage/channel [get]
func getStorageChannel(c *gin.Context) {
	channels := make([]string, 0, len(storage.GetInfoMap()))
	for channel := range storage.GetInfoMap() {
		channels = append(channels, channel)
	}
	resp.Success(c, channels)
}

// getStorageChannelConfig 获取渠道配置
// @Summary 获取渠道配置
// @Tags 存储
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param channel query string false "渠道"
// @Success 200 {object} resp.SuccessStruct{data=map[string][]storage.Desc} "获取成功"
// @Failure 400 {object} resp.ErrorStruct "请求参数错误"
// @Failure 401 {object} resp.ErrorStruct "未授权"
// @Failure 500 {object} resp.ErrorStruct "服务器内部错误"
// @Router /api/v1/storage/channel/config [get]
func getStorageChannelConfig(c *gin.Context) {
	channel := c.Query("channel")
	if channel == "" {
		resp.Success(c, storage.GetInfoMap())
	} else {
		resp.Success(c, storage.GetInfoMap()[channel])
	}
}

// @Summary 更新存储
// @Description 更新存储
// @Tags 存储
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body storageModel.UpdateRequest true "存储配置数据"
// @Success 200 {object} resp.SuccessStruct{data=storageModel.Response} "更新成功"
// @Failure 401 {object} resp.ErrorStruct "未授权"
// @Failure 500 {object} resp.ErrorStruct "服务器内部错误"
// @Router /api/v1/storage [put]
func updateStorage(c *gin.Context) {
	var req storageModel.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.ErrorBadRequest(c)
		return
	}
	configBytes, err := json.Marshal(req.Config)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	configStr := string(configBytes)
	data := storageModel.Data{
		ID:     req.ID,
		Name:   req.Name,
		Type:   req.Type,
		Config: configStr,
	}
	if err := op.UpdateStorage(c.Request.Context(), &data); err != nil {
		resp.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	resp.Success(c, storageModel.Response{
		ID:     data.ID,
		Name:   data.Name,
		Type:   data.Type,
		Config: req.Config,
	})
}

// @Summary 删除存储
// @Description 删除存储
// @Tags 存储
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "存储ID"
// @Success 200 {object} resp.SuccessStruct "删除成功"
// @Failure 401 {object} resp.ErrorStruct "未授权"
// @Failure 500 {object} resp.ErrorStruct "服务器内部错误"
// @Router /api/v1/storage/{id} [delete]
func deleteStorage(c *gin.Context) {
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 16)
	if err != nil {
		resp.ErrorBadRequest(c)
		return
	}
	if err := op.DeleteStorage(c.Request.Context(), uint16(idUint)); err != nil {
		resp.ErrorBadRequest(c)
		return
	}
	resp.Success(c, nil)
}
