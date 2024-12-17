package dto

type GetConfigResponse struct {
	ApiResponseBase `json:",squash"`
	ConfigValue     string `json:"configValue"`
}

// type CheckForUpdateRequestUri struct {
// 	ApiRequestUriBase
// 	Os             string `uri:"os" binding:"required"`
// 	CurrentVersion string `uri:"currentVersion" binding:"required"`
// }

type CheckForUpdateResponse struct {
	ApiResponseBase `json:",squash"`
	HasUpdate       bool   `json:"hasUpdate"`
	ForceUpdate     bool   `json:"forceUpdate"`
	Version			string `json:"version"`
	Title           string `json:"title"`
	Content         string `json:"content"`
	Url             string `json:"url"`
}

type PingResponse struct {
	ApiResponseBase `json:",squash"`
	Ping            string `json:"ping"`
}

// type UserAvatarRequestUri struct {
// 	ApiRequestUriBase
// 	UserId int64 `uri:"userId" binding:"number"`
// }
