package request

type CreateLessonRequest struct {
	TopicSlug  string `json:"topicSlug" binding:"required"`
	Title      string `json:"title" binding:"required"`
	OrderIndex int    `json:"orderIndex" binding:"required,min=1"`
}

type UpdateLessonRequest struct {
	Title      string `json:"title" binding:"required"`
	OrderIndex int    `json:"orderIndex" binding:"required,min=1"`
}
