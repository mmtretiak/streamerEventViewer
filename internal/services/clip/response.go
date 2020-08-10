package clip

type TotalViewsResp struct {
	Count int `json:"count"`
}

type TotalViewsByStreamerResp struct {
	StreamerID string `json:"streamer_id"`
	Count      int    `json:"count"`
}

type CreateClipResp struct {
	EditURL string `json:"edit_url"`
}
