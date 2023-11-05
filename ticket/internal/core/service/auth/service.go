package auth

import (
	"context"
	"ticket/internal/core/domain"
	"time"

	"github.com/golang-jwt/jwt"
)

func (s Service) createAccessToken(claims domain.UserClaims) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedAccessToken, err := accessToken.SignedString([]byte(s.config.Secret))
	if err != nil {
		return "", err
	}
	return signedAccessToken, nil
}

func (s Service) createRefreshToken(claims jwt.StandardClaims) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedAccessToken, err := accessToken.SignedString([]byte(s.config.Secret))
	if err != nil {
		return "", err
	}
	return signedAccessToken, nil
}

func (s Service) parseAccessToken(accessToken string) *domain.UserClaims {
	parsedAccessToken, _ := jwt.ParseWithClaims(accessToken, &domain.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.Secret), nil
	})
	return parsedAccessToken.Claims.(*domain.UserClaims)

}

func (s Service) parseRefreshToken(refreshToken string) *jwt.StandardClaims {
	parsedRefreshToken, _ := jwt.ParseWithClaims(refreshToken, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.Secret), nil
	})
	return parsedRefreshToken.Claims.(*jwt.StandardClaims)
}

func (s Service) Login(ctx context.Context, data domain.LoginRequest) (*domain.TokenResponse, error) {
	user, err := s.userSrv.GetUserByEmail(ctx, data.Email)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.createAccessToken(domain.UserClaims{
		ID:   user.ID,
		Role: user.Role,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(s.config.AccessTokenExpire).Unix(),
		},
	})
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.createRefreshToken(jwt.StandardClaims{
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
func (s Service) Logout(ctx context.Context, data domain.LogoutRequest) error {
	userClaims := s.parseAccessToken(data.AccessToken)
	err := userClaims.Valid()
	if err != nil {
		return err
	}

	standardClaims := s.parseRefreshToken(data.RefreshToken)
	err = standardClaims.Valid()
	if err != nil {
		return err
	}

	return nil
}

func (s Service) RefreshAccessToken(ctx context.Context, data domain.RefreshTokenRequest) (*domain.TokenResponse, error) {
	userClaims := s.parseAccessToken(data.AccessToken)
	standardClaims := s.parseRefreshToken(data.RefreshToken)
	err := standardClaims.Valid()
	if err != nil {
		return nil, err
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, domain.UserClaims{
		ID:   userClaims.ID,
		Role: userClaims.Role,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(s.config.AccessTokenExpire).Unix(),
		},
	})
	signedAccessToken, err := accessToken.SignedString([]byte(s.config.Secret))
	if err != nil {
		return nil, err
	}

	return &domain.TokenResponse{
		AccessToken:  signedAccessToken,
		RefreshToken: data.RefreshToken,
	}, nil

}
func (s Service) Authentication(ctx context.Context, accessToken string) (*domain.AuthenticationResponse, error) {
	userClaims := s.parseAccessToken(accessToken)
	err := userClaims.Valid()
	if err != nil {
		return nil, err
	}

	user, err := s.userSrv.GetUserByID(ctx, userClaims.ID)
	if err != nil {
		return nil, err
	}

	return &domain.AuthenticationResponse{
		User: &domain.User{
			BaseUser: domain.BaseUser{
				ID:       user.ID,
				Name:     user.Name,
				Email:    user.Email,
				ImageUrl: user.ImageUrl,
			},
			Role: user.Role,
		},
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
