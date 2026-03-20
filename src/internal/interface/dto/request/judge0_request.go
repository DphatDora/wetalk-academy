package request

type SubmitCodeRequest struct {
	SourceCode     string `json:"sourceCode" binding:"required"`
	LanguageID     int    `json:"languageID" binding:"required"`
	Stdin          string `json:"stdin"`
	ExpectedOutput string `json:"expectedOutput"`
}
