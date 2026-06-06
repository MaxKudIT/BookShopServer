package dto

type AskDTO struct {
	Question string `json:"question" binding:"required"`
}

type AskResponseDTO struct {
	Answer string `json:"answer"`
}
