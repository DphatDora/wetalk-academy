package service

import (
	"context"
	"wetalk-academy/internal/domain/repository"
	"wetalk-academy/internal/infrastructure/judge0"
	"wetalk-academy/internal/interface/dto/request"
	"wetalk-academy/internal/interface/dto/response"
)

type Judge0Service struct {
	judge0Repo repository.Judge0Repository
}

func NewJudge0Service(judge0Repo repository.Judge0Repository) *Judge0Service {
	return &Judge0Service{
		judge0Repo: judge0Repo,
	}
}

func (s *Judge0Service) SubmitCode(ctx context.Context, req *request.SubmitCodeRequest) (*response.SubmissionResponse, error) {
	submissionReq := &judge0.SubmissionRequest{
		SourceCode:     req.SourceCode,
		LanguageID:     req.LanguageID,
		Stdin:          req.Stdin,
		ExpectedOutput: req.ExpectedOutput,
	}

	result, err := s.judge0Repo.Submit(ctx, submissionReq)
	if err != nil {
		return nil, err
	}

	return &response.SubmissionResponse{
		Token:         "", // not returned when wait=true
		Stdout:        result.Stdout,
		Stderr:        result.Stderr,
		CompileOutput: result.CompileOutput,
		Message:       result.Message,
		Time:          result.Time,
		Memory:        result.Memory,
		StatusID:      result.Status.ID,
		StatusDesc:    result.Status.Description,
	}, nil
}
