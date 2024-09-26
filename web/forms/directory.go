package forms

import (
	"bialekredki/atik/models"
)

type BooleanLikeValues string

const (
	BooleanLikeValue_On    BooleanLikeValues = "on"
	BooleanLikeValue_Off   BooleanLikeValues = "off"
	BooleanLikeValue_True  BooleanLikeValues = "true"
	BooleanLikeValue_False BooleanLikeValues = "false"
)

func (blv BooleanLikeValues) Bool() bool {
	switch blv {
	case BooleanLikeValue_True, BooleanLikeValue_On:
		return true
	default:
		return false
	}
}

type CreateDirectory struct {
	Name               string            `form:"name" binding:"required"`
	ParentId           *uint             `form:"parentId"`
	IsColdStorageClass BooleanLikeValues `form:"isColdStorageClass"`
}

func (form CreateDirectory) StorageClass() models.StorageTypeEnum {
	var storageClass models.StorageTypeEnum = models.StorageType_Hot
	if form.IsColdStorageClass.Bool() {
		storageClass = models.StorageType_Cold
	}
	return storageClass
}
