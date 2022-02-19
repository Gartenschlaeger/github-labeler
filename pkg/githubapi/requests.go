package githubapi

type CreateLabelRequest struct {
	Name        string `json:"name"`
	Color       string `json:"color"`
	Description string `json:"description"`
}

type UpdateLabelRequest struct {
	Name        string `json:"new_name"`
	Color       string `json:"color"`
	Description string `json:"description"`
}
