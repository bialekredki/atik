package handlers

import (
	"bialekredki/atik/web/forms"
	"bialekredki/atik/web/templates"
	"bialekredki/atik/web/templates/components"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *Handler) CreateDirectory(ctx *gin.Context) {
	ownerId := h.getOwnerId(ctx)
	form := &forms.CreateDirectory{}
	if err := ctx.ShouldBind(form); err != nil {
		ctx.String(http.StatusUnprocessableEntity, err.Error())
		return
	}
	object, err := h.metadataRepository.CreateNewDirectory(form.Name, ownerId, form.ParentId, form.StorageClass())
	if err != nil {
		h.logger.Error("internal server error", zap.Error(err))
		ctx.String(http.StatusInternalServerError, "Unknown error")
		return
	}
	ctx.HTML(http.StatusCreated, "", components.ObjectTableRow(*object))
}

func (h *Handler) GetDirectoryRow(ctx *gin.Context) {
	ownerId := h.getOwnerId(ctx)
	objectId, ok := ctx.Params.Get("objectId")
	if !ok {
		ctx.AbortWithError(http.StatusUnprocessableEntity, errors.New("missing objectId in params"))
		return
	}
	intObjectId, err := strconv.Atoi(objectId)
	if err != nil {
		ctx.AbortWithError(http.StatusUnprocessableEntity, errors.New("objectId is not an integer value"))
		return
	}
	uintObjectId := uint(intObjectId)
	object, err := h.metadataRepository.ById(uintObjectId)
	if err != nil {
		ctx.AbortWithError(http.StatusNotFound, errors.New("object with ObjectId can't be found"))
		return
	}
	if object.OwnerId != ownerId {
		ctx.AbortWithError(http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}
	subobjects, err := h.metadataRepository.ListMetadataObjects(ownerId, &object.ID, 20)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	session := sessions.Default(ctx)
	profile := getProfileFromSession(session)
	parents := h.metadataRepository.ListObjectParentsById(object.ID)
	h.logger.Sugar().Debugf("%+v", parents)
	ctx.HTML(http.StatusOK, "", templates.Home(templates.HomeTemplateParams{
		Profile:          profile,
		Directory:        object,
		Objects:          subobjects,
		DirectoryParents: parents,
	}))
}
