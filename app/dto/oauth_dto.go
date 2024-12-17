package dto

import "github.com/innovay-software/famapp-main/app/models"

// type LoginRequest struct {
// 	ApiRequestBase
// 	Mobile      string `json:"mobile" binding:"required"`
// 	Password    string `json:"password" binding:"required"`
// 	DeviceToken string `json:"deviceToken" binding:"omitempty"`
// }

// type AccessTokenLoginRequest struct {
// 	ApiRequestBase
// 	DeviceToken string `json:"deviceToken" binding:"omitempty"`
// }

type LoginResponse struct {
	ApiResponseBase                   `json:",squash"`
	ClientOssAccessKeyId              string               `json:"clientOssAccessKeyId"`
	ClientOssAccessKeySecret          string               `json:"clientOssAccessKeySecret"`
	OssDomain                         string               `json:"ossDomain"`
	OssBucketName                     string               `json:"ossBucketName"`
	AliyunOssAccessKeyId              string               `json:"aliyunOssAccessKeyId"`
	AliyunOssAccessKeySecret          string               `json:"aliyunOssAccessKeySecret"`
	AliyunOssDomain                   string               `json:"aliyunOssDomain"`
	AliyunOssBucketName               string               `json:"aliyunOssBucketName"`
	GoogleCloudStorageAccessId        string               `json:"googleCloudStorageAccessId"`
	GoogleCloudStorageAccessSecretKey string               `json:"googleCloudStorageAccessSecretKey"`
	GoogleCloudStorageDomain          string               `json:"googleCloudStorageDomain"`
	GoogleCloudStorageBucketName      string               `json:"googleCloudStorageBucketName"`
	HwObsAccessId                     string               `json:"hwObsAccessId"`
	HwObsAccessSecretKey              string               `json:"hwObsAccessSecretKey"`
	HwObsDomain                       string               `json:"hwObsDomain"`
	HwObsBucketName                   string               `json:"hwObsBucketName"`
	ImCenterData                      string               `json:"imCenterData"`
	User                              *models.User         `json:"userData"`
	Members                           []*models.UserMember `json:"memberList"`
	Folders                           []*models.Folder     `json:"folders"`
}
