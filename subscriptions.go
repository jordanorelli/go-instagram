package instagram

type Notification struct {
	ChangedAspect  string `json:"changed_aspect"`
	SubscriptionId int    `json:"subscription_id"`
	Object         string `json:"object"`
	ObjectId       string `json:"object_id"`
	Time           int    `json:"time"`
}
