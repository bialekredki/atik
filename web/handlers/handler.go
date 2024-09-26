package handlers

import (
	"bialekredki/atik/lib/aws"
	"bialekredki/atik/lib/metadata"
	"bialekredki/atik/pkg/auth"
	"errors"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Handler struct {
	auth               *auth.Authenticator
	s3repository       *aws.S3Repository
	metadataRepository *metadata.MetadataRepository
	logger             *zap.Logger
}

type HandlerParams struct {
	fx.In
	Auth               *auth.Authenticator
	S3Repository       *aws.S3Repository
	MetadataRepository *metadata.MetadataRepository
	Logger             *zap.Logger
}

func NewHandler(p HandlerParams) *Handler {
	return &Handler{
		s3repository:       p.S3Repository,
		auth:               p.Auth,
		metadataRepository: p.MetadataRepository,
		logger:             p.Logger,
	}
}

func (h *Handler) getOwnerIdOrError(ctx *gin.Context) (uint, error) {
	session := sessions.Default(ctx)
	owner_id, ok := session.Get("owner_id").(uint)
	if !ok {
		return 0, errors.New("not able to get owner_id from session")
	}
	return owner_id, nil
}

func (h *Handler) getOwnerId(ctx *gin.Context) uint {
	ownerId, _ := h.getOwnerIdOrError(ctx)
	return ownerId
}
