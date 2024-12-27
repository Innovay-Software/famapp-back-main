package repositories

import (
	"fmt"

	"github.com/innovay-software/famapp-main/app/models"
	"gorm.io/gorm"
)

type userRepo struct {
	mainDBCon *gorm.DB
	readDBCon *gorm.DB
	rd        *redisRepo
}

func (u userRepo) FindMemberList(limit int, afterId uint64) (
	[]*models.UserMember, error,
) {
	db := u.readDBCon
	users := []*models.UserMember{}
	if err := db.Limit(limit).Where("id > ?", afterId).
		Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (u userRepo) FindAllUser() (*[]models.User, error) {
	db := u.readDBCon
	users := []models.User{}
	if err := db.Limit(1000).Find(&users).Error; err != nil {
		return nil, err
	}
	return &users, nil
}

func (u userRepo) FindUserByField(fieldName string, value string) (
	*models.User, error,
) {
	user := models.User{}
	db := u.readDBCon
	if err := db.Where(fieldName+" = ?", value).
		First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u userRepo) DeleteUser(fieldName string, value string) error {
	user := models.User{}
	db := u.mainDBCon
	if err := db.Where(fieldName+" = ?", value).
		First(&user).Error; err != nil {
		return err
	}
	MongoRepoIns.DeleteUserInMongo(&user)
	return DeleteDbModel(&user)
}

func (u userRepo) CreateUser(user *models.User) error {
	// check for duplicate mobile
	db := u.mainDBCon
	duplicateCount := int64(0)
	db.Model(&models.User{}).
		Where("mobile", user.Mobile).
		Where("id != ?", user.ID).
		Count(&duplicateCount)
	if duplicateCount > 0 {
		// found a user with same mobile
		return fmt.Errorf("duplicate mobile")
	}

	err := CreateDbModel(user)
	if err != nil {
		return err
	}

	MongoRepoIns.CreateUserInMongo(user)
	return nil
}

func (u userRepo) SaveUser(user *models.User) error {
	return SaveDbModel(user)
}

func (u userRepo) FindFolders(user *models.User) []*models.Folder {
	db := u.readDBCon
	var folders []*models.Folder
	db.Model(&user).Association("Folders").Find(&folders)
	return folders
}
