package metadata

import (
	"bialekredki/atik/models"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type MetadataRepository struct {
	connection *gorm.DB
	logger     *zap.Logger
}

func NewMetadataRepository(connection *gorm.DB) (*MetadataRepository, error) {
	err := connection.AutoMigrate(
		models.MetadataDirectory{},
		models.MetadataFile{},
		models.MetadataOwner{},
		models.MetadataObject{},
	)

	if err != nil {
		return nil, err
	}

	return &MetadataRepository{
		connection: connection,
	}, nil
}

func (r *MetadataRepository) createMetadataObject(name string, ownerId uint, parentId *uint, objectType models.MetadataObjectType, storageClass models.StorageTypeEnum) (*models.MetadataObject, error) {
	object := models.MetadataObject{
		Name:         name,
		OwnerId:      ownerId,
		ParentId:     parentId,
		ObjectType:   objectType,
		StorageClass: storageClass,
	}
	result := r.connection.Create(&object)
	if result.Error != nil {
		return nil, result.Error
	}
	return &object, nil
}

func (r *MetadataRepository) CreateNewDirectory(name string, ownerId uint, parentId *uint, storageClass models.StorageTypeEnum) (*models.MetadataObject, error) {
	return r.createMetadataObject(name, ownerId, parentId, models.DirectoryObjectType, storageClass)
}

func (r *MetadataRepository) ById(id uint) (*models.MetadataObject, error) {
	var object models.MetadataObject
	tx := r.connection.Where("id = ?", id).First(&object)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &object, nil
}

func (r *MetadataRepository) ListMetadataObjects(ownerId uint, parentId *uint, limit int) ([]models.MetadataObject, error) {
	objects := make([]models.MetadataObject, 0, limit)
	tx := r.connection.Where("owner_id = ?", ownerId)
	if parentId != nil {
		tx = tx.Where("parent_id = ?", *parentId)
	} else {
		tx = tx.Where("parent_id is NULL")
	}
	tx = tx.Order("name DESC").Limit(limit).Find(&objects)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return objects, nil
}

func (r *MetadataRepository) CreateNewOwner(name string) uint {
	owner := models.MetadataOwner{
		Name: name,
	}
	result := r.connection.Where(&owner).FirstOrCreate(&owner)
	if result.Error != nil {
		panic(result.Error)
	}
	return owner.ID
}

func (r *MetadataRepository) ListContentsOfDirectory(parentId *uuid.UUID, ownerId uint, limit int) ([]models.MetadataDirectory, []models.MetadataFile) {
	directories := []models.MetadataDirectory{}
	files := []models.MetadataFile{}
	return directories, files
	var directories_count int64
	r.connection.Where("owner_id = ?", ownerId).Order("name DESC").Limit(limit).Find(&directories).Count(&directories_count)
	if directories_count < int64(limit) {
		limit -= int(directories_count)
		r.connection.Where("directory_id = ?", parentId).Order("name DESC").Limit(limit).Find(&files)
	}
	fmt.Printf("ListContents %v\n%v", directories, files)

	return directories, files
}

func (r *MetadataRepository) ListObjectParentsById(objectId uint) []*models.MetadataObject {
	parents := make([]*models.MetadataObject, 0)
	var object models.MetadataObject
	tx := r.connection.Joins("Parent").First(&object, objectId)
	if tx.Error != nil {
		r.logger.Sugar().Errorln(tx.Error.Error())
		return parents
	}
	objectPtr := &object
	fmt.Printf("%v\n", objectPtr)
	for objectPtr.Parent != nil {
		fmt.Printf("%v\n", objectPtr)
		parents = append(parents, objectPtr.Parent)
		objectPtr = objectPtr.Parent
	}
	return parents
}
