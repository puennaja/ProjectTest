package auth

import (
	"context"
	"ticket/internal/core/domain"
	"ticket/pkg/errors"
	"ticket/pkg/utils"
	"time"

	"github.com/golang-jwt/jwt"
)

func (s Service) createToken(claims jwt.StandardClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.config.Secret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func (s Service) parseToken(token string) (*jwt.StandardClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	return parsedToken.Claims.(*jwt.StandardClaims), nil
}

func (s Service) Login(ctx context.Context, data *domain.LoginRequest) (*domain.TokenResponse, error) {
	user, err := s.userSrv.GetUserByEmail(ctx, data.Email)
	if err != nil {
		return nil, err
	}

	hash := utils.SHA3Hash(data.Pwd)
	if hash != user.Password {
		return nil, errors.ErrUnauthorized.SetError(errors.New("wrong password"))
	}

	accessToken, err := s.createToken(jwt.StandardClaims{
		Id:        user.ID,
		Subject:   user.Role,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(s.config.AccessTokenExpire).Unix(),
	})
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.createToken(jwt.StandardClaims{
		Id:        user.ID,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(s.config.RefreshTokenExpire).Unix(),
	})
	if err != nil {
		return nil, err
	}
	return &domain.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}
func (s Service) Logout(ctx context.Context, data *domain.LogoutRequest) error {
	accessClaims, err := s.parseToken(data.AccessToken)
	if err != nil {
		return err
	}

	err = accessClaims.Valid()
	if err != nil {
		return err
	}

	return nil
}

func (s Service) RefreshAccessToken(ctx context.Context, data *domain.RefreshTokenRequest) (*domain.TokenResponse, error) {
	claims, err := s.parseToken(data.RefreshToken)
	if err != nil {
		return nil, err
	}

	err = claims.Valid()
	if err != nil {
		return nil, err
	}

	user, err := s.userSrv.GetUserByID(ctx, claims.Id)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.createToken(jwt.StandardClaims{
		Id:        user.ID,
		Subject:   user.Role,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(s.config.AccessTokenExpire).Unix(),
	})
	if err != nil {
		return nil, err
	}

	return &domain.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: data.RefreshToken,
	}, nil

}
func (s Service) Authentication(ctx context.Context, accessToken string) (*domain.AuthenticationResponse, error) {
	claims, err := s.parseToken(accessToken)
	if err != nil {
		return nil, err
	}

	err = claims.Valid()
	if err != nil {
		return nil, err
	}

	user, err := s.userSrv.GetUserByID(ctx, claims.Id)
	if err != nil {
		return nil, err
	}

	return &domain.AuthenticationResponse{
		User: *user,
	}, nil

}
func (s Service) Authorization(ctx context.Context, role, method, path string) (*domain.AuthorizationResponse, error) {
	ok, err := s.casbinEnforcer.Enforce(role, path, method)
	if err != nil {
		return nil, err
	}

	out := &domain.AuthorizationResponse{
		GrantAccess: ok,
	}
	return out, nil
}
