package models

import "github.com/google/uuid"

type StorageTypeEnum string
type MetadataObjectType uint8

const (
	StorageType_Hot  StorageTypeEnum = "HOT"
	StorageType_Cold StorageTypeEnum = "COLD"
)

const (
	FileObjectType      MetadataObjectType = iota
	DirectoryObjectType MetadataObjectType = iota
)

func (mot MetadataObjectType) String() string {
	switch mot {
	case FileObjectType:
		return "file"
	case DirectoryObjectType:
		return "directory"
	}
	return ""
}

type MetadataOwner struct {
	ID      uint   `gorm:"primaryKey"`
	Name    string `gorm:"type:varchar(128);not null;uniqueIndex"`
	Created int64  `gorm:"autoCreateTime"`
}

type MetadataObject struct {
	ID           uint               `gorm:"primaryKey"`
	Name         string             `gorm:"type:varchar(128);not null;uniqueIndex:idx_metadataobject_unique"`
	Created      int64              `gorm:"autoCreateTime"`
	OwnerId      uint               `gorm:"not null;uniqueIndex:idx_metadataobject_unique"`
	ObjectType   MetadataObjectType `gorm:"not null;index"`
	StorageClass StorageTypeEnum    `gorm:"not null;type:storage_type_enum"`
	Owner        MetadataOwner
	ParentId     *uint           `gorm:"uniqueIndex:idx_metaobject_unique"`
	Parent       *MetadataObject `gorm:"foreignKey:ParentId"`
}

type MetadataDirectory struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name     string    `gorm:"type:varchar(64);not null;index;uniqueIndex:idx_metadirectories_unique"`
	Created  int64     `gorm:"autoCreateTime"`
	OwnerId  uint      `gorm:"type:uuid;not null;index;uniqueIndex:idx_metadirectories_unique"`
	Owner    MetadataOwner
	ParentID *uuid.UUID         `gorm:"type:uuid;uniqueIndex:idx_metadirectories_unique"`
	Parent   *MetadataDirectory `gorm:"foreignKey:ParentID"` // Define the relationship
}

type MetadataFile struct {
	ID          uuid.UUID       `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name        string          `gorm:"type:varchar(64);not null"`
	StorageType StorageTypeEnum `gorm:"type:storage_type_enum"`
	Created     int64           `gorm:"autoCreateTime"`
	Updated     int64           `gorm:"autoDeletedTime"`
	DirectoryID uuid.UUID
	Directory   *MetadataDirectory
}
