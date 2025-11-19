package service

import "errors"

type ServiceErr struct {
	Err error
	Msg string
}

func NewServiceError(message string, err error) *ServiceErr {
	return &ServiceErr{Msg: message, Err: err}
}

func (r *ServiceErr) Error() string {
	return r.Msg
}

func (r *ServiceErr) Unwrap() error {
	return r.Err
}

var (
	ErrOTPSent = &ServiceErr{
		Msg: "OTP already sent",
		Err: errors.New("OTP already sent"),
	}
	ErrCacheRetrieval = &ServiceErr{
		Msg: "failed to retrieve cached OTP code",
		Err: errors.New("cache retrieval error"),
	}
	ErrOTPGeneration = &ServiceErr{
		Msg: "failed to create OTP",
		Err: errors.New("otp generation error"),
	}
	ErrCacheSet = &ServiceErr{
		Msg: "failed to cache the OTP code with TTL",
		Err: errors.New("cache setting error"),
	}
	ErrInvalidOTP = &ServiceErr{
		Msg: "code is not valid",
		Err: errors.New("invalid code"),
	}
	ErrRefreshTokenCreation = &ServiceErr{
		Msg: "failed to create refresh token",
		Err: errors.New("refresh token creation error"),
	}
	ErrInvalidRoomName = &ServiceErr{
		Msg: "invalid room name",
		Err: errors.New("invalid room name"),
	}
	ErrInvalidRoomCapacity = &ServiceErr{
		Msg: "invalid room cap",
		Err: errors.New("invalid room cap"),
	}
	ErrRoomNameExists = &ServiceErr{
		Msg: "room name exists",
		Err: errors.New("rome already exists"),
	}
)
