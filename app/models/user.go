package models

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/innovay-software/famapp-main/app/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	BaseDbModel
	FamilyID        uint64     `gorm:"column:family_id" json:"familyId"`
	UUID            uuid.UUID  `gorm:"type:uuid; uniqueIndex; default:uuid_generate_v3(); not null" json:"uuid"`
	Role            string     `gorm:"type:enum('superadmin', 'admin','manager','member')" json:"role"`
	Name            string     `gorm:"column:name" json:"name"`
	Email           string     `gorm:"uniqueIndex" json:"email"`
	EmailVerifiedAt *time.Time `gorm:"null" json:"emailVerifiedAt"`
	Mobile          string     `gorm:"uniqueIndex" json:"mobile"`
	Avatar          string     `gorm:"column:avatar" json:"avatar"`
	Password        string     `gorm:"column:password" json:"-"`
	LockerPasscode  string     `gorm:"column:locker_passcode" json:"lockerPasscode"`
	RefreshToken    *string    `gorm:"null" json:"-"`
	DeviceToken     *string    `gorm:"null" json:"-"`
	Notifications   JSONB      `gorm:"type:jsonb; default:'{\"total\":0}'" json:"notifications"`
	Status          bool       `gorm:"default: true" json:"status"`
	Family          Family     `gorm:"foreignKey:FamilyID"`
	Folders         []Folder   `gorm:"many2many:folder_invitees;foreignKey:ID;joinForeignKey:InviteeID;References:ID;joinReferences:FolderID" json:"folders"`
	// LockerNotes     []LockerNote `gorm:"foreignKey:OwnerID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.UUID = uuid.New()
	u.SetDefaultEmail()
	return nil
}

func (u *User) BeforeSave(tx *gorm.DB) error {
	// Check for super admin
	if u.ID == 1 {
		u.Role = "superadmin"
	}

	// Check if is valid
	if err := u.HasDataError(tx); err != nil {
		return err
	}

	// Set default email
	u.SetDefaultEmail()

	return nil
}

func (u *User) IsSuperAdmin() bool {
	if u.ID == 1 {
		return true
	}
	return strings.ToLower(u.Role) == "superadmin"
}

func (u *User) IsAdmin() bool {
	if u.IsSuperAdmin() {
		return true
	}
	return strings.ToLower(u.Role) == "admin"
}

func (u *User) IsManager() bool {
	return strings.ToLower(u.Role) == "manager"
}

func (u *User) HasDataError(db *gorm.DB) error {
	// Check for duplicate mobile
	var duplicateCount int64
	db.Model(&User{}).Where("mobile", u.Mobile).Where("id != ?", u.ID).Count(&duplicateCount)
	if duplicateCount > 0 {
		// found a user with same mobile
		return fmt.Errorf("duplicate mobile")
	}
	return nil
}

func (u *User) SetDefaultEmail() {
	if u.Email == "" {
		u.Email = u.Mobile + "@" + os.Getenv("EMAIL_DOMAIN")
	}
}

func (u *User) UpdateData(newData *map[string]any) error {
	allowedFields := []string{"Role", "Name", "Email", "Mobile", "Password", "LockerPasscode"}
	v := reflect.ValueOf(u)
	if v.Kind() != reflect.Ptr {
		return errors.New("model must be a pointer")
	}
	v = v.Elem()

	skippedFieldNames := []string{}
	for _, fieldName := range allowedFields {
		// Make sure fieldName is valid and can be set
		field := v.FieldByName(fieldName)
		if !field.IsValid() || !field.CanSet() {
			continue
		}

		// Check if field needs to be updated
		fieldValue, exists := (*newData)[utils.CamelToSnakeCase(fieldName)]
		if !exists {
			continue
		}

		switch field.Kind() {
		case reflect.String:
			field.SetString(fieldValue.(string))
		case reflect.Int:
			field.SetInt(int64(fieldValue.(int)))
		default:
			skippedFieldNames = append(skippedFieldNames, fieldName)
		}
	}

	if len(skippedFieldNames) > 0 {
		return fmt.Errorf("skipped field names: %s", strings.Join(skippedFieldNames, ","))
	}
	return nil
}

func (u *User) SetPassword(newPassword string) error {
	if newPassword == "" {
		return nil
	}
	pw := []byte(newPassword)
	result, err := bcrypt.GenerateFromPassword(pw, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	newPasswordHash := string(result)
	if u.Password != newPasswordHash {
		u.Password = newPasswordHash
	}
	return nil
}

func (u *User) SetAvatarUrl(newAvatarUrl string) error {
	u.Avatar = newAvatarUrl

	// if newAvatarUrl == "" || u.Avatar == newAvatarUrl {
	// 	return nil
	// }
	// oldAvatar := u.Avatar
	// if parts := strings.Split(newAvatarUrl, "/user-upload/"); len(parts) == 2 {
	// 	path := parts[1]
	// 	filename := filepath.Base(path)
	// 	filename = strconv.FormatInt(u.ID, 10) + "_" + filename
	// 	sourcePath, sourceExists := utils.GetStorageAbsPath("user-upload", path)
	// 	if !sourceExists {
	// 		// file not exist
	// 		return nil
	// 	}
	// 	targetPath, _ := utils.GetStorageAbsPath("avatars", filename)
	// 	os.Rename(sourcePath, targetPath)
	// 	u.Avatar = utils.GetUrlPath("avatars", filename)

	// 	if oldAvatarParts := strings.Split(oldAvatar, "/"); len(oldAvatarParts) > 0 {
	// 		oldAvatarFilename := oldAvatarParts[len(oldAvatarParts)-1]
	// 		oldAvatarFilePath, _ := utils.GetStorageAbsPath("avatars", oldAvatarFilename)
	// 		os.Remove(oldAvatarFilePath)
	// 	}
	// }

	return nil
}
