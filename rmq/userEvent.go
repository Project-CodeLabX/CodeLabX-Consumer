package rmq

type UserEvent struct {
	UserName string `json:"username"`
	Language string `json:"language,omitempty"`
	Code     string `json:"code,omitempty"`
	FileName string `json:"filename,omitempty"`
}
