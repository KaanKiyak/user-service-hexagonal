package mapper

import "user-service-hexagonal/internal/core/dto"

// ToRefreshTokenResponse - string tokenları DTO'ya çevirir
func ToRefreshTokenResponse(accessToken, refreshToken string) dto.RefreshTokenResponse {
	return dto.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}
