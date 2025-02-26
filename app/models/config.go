package models

type Config struct {
	BaseDbModel
	ConfigKey    string `gorm:"column:config_key" json:"configKey"`
	ConfigValue  string `gorm:"column:config_value" json:"configValue"`
	ConfigType   string `gorm:"column:config_type; default:string" json:"configType"`
	ConfigParams JSONB  `gorm:"column:config_params; type:jsonb; default:'{}'" json:"configParams"`
	Group        string `gorm:"column:group; default:basic" json:"group"`
	Title        string `gorm:"column:title" json:"title"`
	Remark       string `gorm:"column:remark; null" json:"remark"`
	Status       bool   `gorm:"column:status" json:"status"`
}

func (Config) TableName() string {
	return "configs"
}
