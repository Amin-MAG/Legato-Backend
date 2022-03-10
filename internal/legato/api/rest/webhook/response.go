package webhook

type BriefWebhookResponse struct {
	Id        uint   `json:"id"`
	Token     string `json:"token"`
	IsEnable  bool   `json:"isEnable"`
	ServiceID uint   `json:"serviceID"`
}
