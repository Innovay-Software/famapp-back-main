package dto

type ApiRequest interface {
	DummyFunc()
}

type ApiRequestBase struct{}

func (req *ApiRequestBase) DummyFunc() {}

type ApiRequestUri interface {
	GetRequestUriCode() int
}

type ApiRequestUriBase struct{}

func (req ApiRequestUriBase) GetRequestUriCode() int {
	return 0
}

// ApiResponse interface to be used to interact with end points
type ApiResponse interface {
	SetAccessToken(string)
	GetAccessToken() string
	SetRefreshToken(string)
	GetRefreshToken() string
}

// ApiResponseBase type that implements ApiResponse Interface
// So it can be use as a base struct for all ApiResponse structs
type ApiResponseBase struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func (res *ApiResponseBase) SetAccessToken(token string) {
	if token != "" {
		res.AccessToken = token
	}
}

func (res *ApiResponseBase) SetRefreshToken(token string) {
	if token != "" {
		res.RefreshToken = token
	}
}

func (res *ApiResponseBase) GetAccessToken() string {
	return (*res).AccessToken
}

func (res *ApiResponseBase) GetRefreshToken() string {
	return (*res).RefreshToken
}
