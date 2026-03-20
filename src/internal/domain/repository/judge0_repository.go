package repository

import (
	"context"
	"wetalk-academy/internal/infrastructure/judge0"
)

type Judge0Repository interface {
	Submit(ctx context.Context, req *judge0.SubmissionRequest) (*judge0.StatusResponse, error)
}
