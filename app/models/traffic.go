package models

type Traffic struct {
	BaseDbModel
	IP           string `gorm:"column:ip" json:"ip"`
	UserID       int64  `gorm:"column:user_id" json:"userId"`
	Requester    string `gorm:"column:requester" json:"requester"`
	RequestURI   string `gorm:"column:request_uri" json:"requestUri"`
	RequestBody  string `gorm:"column:request_body" json:"requestBody"`
	ResponseCode string `gorm:"column:response_code" json:"responseCode"`
	ErrorMessage string `gorm:"column:error_message" json:"errorMessage"`
}

func (Traffic) TableName() string {
	return "traffics"
}
