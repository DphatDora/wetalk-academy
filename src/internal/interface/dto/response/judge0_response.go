package response

type SubmissionResponse struct {
	Token         string  `json:"token"`
	Stdout        *string `json:"stdout"`
	Stderr        *string `json:"stderr"`
	CompileOutput *string `json:"compileOutput"`
	Message       *string `json:"message"`
	Time          *string `json:"time"`
	Memory        *int    `json:"memory"`
	StatusID      int     `json:"statusID"`
	StatusDesc    string  `json:"statusDescription"`
}
