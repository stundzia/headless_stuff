package models

type HeadlessResponse struct {
	TargetUrl string `json:"target_url"`
	ResUrl string `json:"res_url"`
	Body string `json:"content"`
	Title string `json:"title"`
	StatusCode int `json:"status_code"`
	ContentBytes []byte
	ContentString string `json:"content_string"`
}

func (sr *HeadlessResponse) Write(p []byte) (n int, err error){
	sr.ContentBytes = p
	sr.ContentString = string(p)
	return len(p), nil
}