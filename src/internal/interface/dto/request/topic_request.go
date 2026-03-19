package request

type CreateTopicRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Author      Author `json:"author" binding:"required"`
}

type Author struct {
	Avatar string `json:"avatar"`
	Name   string `json:"name" binding:"required"`
}

type UpdateTopicRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
}
