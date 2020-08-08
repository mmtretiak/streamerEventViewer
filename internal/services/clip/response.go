package clip

type TotalViewsResp struct {
	Count int `json:"count"`
}

type TotalViewsByStreamer struct {
	StreamerID string `json:"streamer_id"`
	Count      int    `json:"count"`
}
