package qa

type CreateQaPairRequest struct {
	SelectedQas []struct {
		Question string `json:"question"`
		Answer   string `json:"answer"`
	} `json:"selectedQas"`
}
