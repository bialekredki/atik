package models

import "github.com/google/uuid"

type MetadataOwner struct {
	ID      uint   `gorm:"primaryKey"`
	Name    string `gorm:"type:varchar(128);not null"`
	Created int64  `gorm:"autoCreateTime"`
}

type MetadataDirectory struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name     string    `gorm:"type:varchar(64);not null"`
	ParentID uuid.UUID `gorm:"type:uuid"`
	Parent   *MetadataDirectory
}

type MetadataFile struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name        string    `gorm:"type:varchar(64);not null"`
	Created     int64     `gorm:"autoCreateTime"`
	Updated     int64     `gorm:"autoDeletedTime"`
	DirectoryID uuid.UUID
	Directory   *MetadataDirectory
}
