package models

import "errors"

type RenderRequest struct {
	Url string `json:"url,omitempty"`
	UserAgentType string `json:"user_agent_type,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
	Engine string `json:"engine,omitempty"`
	WaitTime float64 `json:"wait_time,omitempty"`
	ProxyUrl string `json:"proxy,omitempty"`
	Headless bool `json:"headless,omitempty"`
}

var supportedEngines = []string{"chromedp", "surf"}

func (req *RenderRequest) Validate() error {
	engineValid := false
	for _, eng := range supportedEngines {
		if eng == req.Engine {
			engineValid = true
		}
	}
	if !engineValid {
		return errors.New("invalid engine provided")
	}
	return nil
}