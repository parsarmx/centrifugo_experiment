package dto

import rpc_service "golang_template/proto"

// PublishRequestDTO is the domain-level representation of a PublishRequest.
type PublishRequestDTO struct {
	Channel        string
	Data           []byte
	B64Data        string
	SkipHistory    bool
	Tags           map[string]string
	IdempotencyKey string
	Delta          bool
	Version        uint64
	VersionEpoch   string
}

// PublishResponseDTO represents the domain-level response of a publish request.
type PublishResponseDTO struct {
	ErrorCode    string
	ErrorMessage string
	PublishID    string
	Offset       uint64
}

// ToPublishRequestDTO converts a gRPC PublishRequest to the domain model.
func ToPublishRequestDTO(req *rpc_service.PublishRequest) *PublishRequestDTO {
	if req == nil {
		return nil
	}

	return &PublishRequestDTO{
		Channel:        req.Channel,
		Data:           req.Data,
		B64Data:        req.B64Data,
		SkipHistory:    req.SkipHistory,
		Tags:           req.Tags,
		IdempotencyKey: req.IdempotencyKey,
		Delta:          req.Delta,
		Version:        req.Version,
		VersionEpoch:   req.VersionEpoch,
	}
}
