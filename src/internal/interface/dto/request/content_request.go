package request

type ContentSectionRequest struct {
	Type     string `json:"type" binding:"required"`
	Content  string `json:"content,omitempty"`
	Language string `json:"language,omitempty"`
	URL      string `json:"url,omitempty"`
}

type CreateContentRequest struct {
	Sections []ContentSectionRequest `json:"sections" binding:"required,min=1,dive"`
}

type UpdateContentRequest struct {
	Sections []ContentSectionRequest `json:"sections" binding:"required,min=1,dive"`
}
