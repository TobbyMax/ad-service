package grpc

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"homework10/internal/app"
	"net/mail"
)

func (s *AdService) CreateAd(ctx context.Context, request *CreateAdRequest) (*AdResponse, error) {
	if request.UserId == nil {
		return nil, status.Error(codes.InvalidArgument, ErrMissingArgument.Error())
	}
	ad, err := s.app.CreateAd(ctx, request.GetTitle(), request.GetText(), request.GetUserId())

	if err != nil {
		return nil, status.Error(GetErrorCode(err), err.Error())
	}

	return AdSuccessResponse(ad), nil
}

func (s *AdService) ChangeAdStatus(ctx context.Context, request *ChangeAdStatusRequest) (*AdResponse, error) {
	if request.AdId == nil || request.UserId == nil {
		return nil, status.Error(codes.InvalidArgument, ErrMissingArgument.Error())
	}
	ad, err := s.app.ChangeAdStatus(ctx, request.GetAdId(), request.GetUserId(), request.GetPublished())

	if err != nil {
		return nil, status.Error(GetErrorCode(err), err.Error())
	}
	return AdSuccessResponse(ad), nil
}

func (s *AdService) UpdateAd(ctx context.Context, request *UpdateAdRequest) (*AdResponse, error) {
	if request.AdId == nil || request.UserId == nil {
		return nil, status.Error(codes.InvalidArgument, ErrMissingArgument.Error())
	}
	ad, err := s.app.UpdateAd(ctx, request.GetAdId(), request.GetUserId(), request.GetTitle(), request.GetText())

	if err != nil {
		return nil, status.Error(GetErrorCode(err), err.Error())
	}
	return AdSuccessResponse(ad), nil
}

func (s *AdService) GetAd(ctx context.Context, request *GetAdRequest) (*AdResponse, error) {
	if request.AdId == nil {
		return nil, status.Error(codes.InvalidArgument, ErrMissingArgument.Error())
	}
	ad, err := s.app.GetAd(ctx, request.GetAdId())

	if err != nil {
		return nil, status.Error(GetErrorCode(err), err.Error())
	}
	return AdSuccessResponse(ad), nil
}

func (s *AdService) ListAds(ctx context.Context, request *ListAdRequest) (*ListAdResponse, error) {
	date, err := app.ParseDate(request.Date)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	al, err := s.app.ListAds(ctx, app.ListAdsParams{
		Published: request.Published,
		Uid:       request.UserId,
		Date:      date,
		Title:     request.Title,
	})

	if err != nil {
		return nil, status.Error(GetErrorCode(err), err.Error())
	}
	return AdListSuccessResponse(al), nil
}

func (s *AdService) CreateUser(ctx context.Context, request *CreateUserRequest) (*UserResponse, error) {
	_, err := mail.ParseAddress(request.GetEmail())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	u, err := s.app.CreateUser(ctx, request.GetName(), request.GetEmail())

	if err != nil {
		return nil, status.Error(GetErrorCode(err), err.Error())
	}
	return UserSuccessResponse(u), nil
}

func (s *AdService) UpdateUser(ctx context.Context, request *UpdateUserRequest) (*UserResponse, error) {
	if request.Id == nil {
		return nil, status.Error(codes.InvalidArgument, ErrMissingArgument.Error())
	}
	_, err := mail.ParseAddress(request.GetEmail())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	u, err := s.app.UpdateUser(ctx, request.GetId(), request.GetName(), request.GetEmail())

	if err != nil {
		return nil, status.Error(GetErrorCode(err), err.Error())
	}
	return UserSuccessResponse(u), nil
}

func (s *AdService) GetUser(ctx context.Context, request *GetUserRequest) (*UserResponse, error) {
	if request.Id == nil {
		return nil, status.Error(codes.InvalidArgument, ErrMissingArgument.Error())
	}
	u, err := s.app.GetUser(ctx, request.GetId())

	if err != nil {
		return nil, status.Error(GetErrorCode(err), err.Error())
	}
	return UserSuccessResponse(u), nil
}

func (s *AdService) DeleteUser(ctx context.Context, request *DeleteUserRequest) (*emptypb.Empty, error) {
	if request.Id == nil {
		return nil, status.Error(codes.InvalidArgument, ErrMissingArgument.Error())
	}
	err := s.app.DeleteUser(ctx, request.GetId())
	if err != nil {
		return nil, status.Error(GetErrorCode(err), err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *AdService) DeleteAd(ctx context.Context, request *DeleteAdRequest) (*emptypb.Empty, error) {
	if request.AdId == nil || request.AuthorId == nil {
		return nil, status.Error(codes.InvalidArgument, ErrMissingArgument.Error())
	}
	err := s.app.DeleteAd(ctx, request.GetAdId(), request.GetAuthorId())
	if err != nil {
		return nil, status.Error(GetErrorCode(err), err.Error())
	}
	return &emptypb.Empty{}, nil
}
