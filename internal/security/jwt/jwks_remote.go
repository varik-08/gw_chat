package jwt

import (
    "context"
    "errors"
    "time"

    gjwt "github.com/golang-jwt/jwt/v5"
    "github.com/lestrrat-go/jwx/v2/jwk"
)

// NewRemoteJWKSProvider скачивает JWKS по URL один раз при старте.
func NewRemoteJWKSProvider(ctx context.Context, url string) (*JWKSProvider, error) {
    set, err := jwk.Fetch(ctx, url, jwk.WithHTTPClient(nil))
    if err != nil {
        return nil, err
    }
    return &JWKSProvider{set: set}, nil
}

// ValidateRS256 проверяет подпись RS256 токена используя JWKS.
func ValidateRS256(tokenString string, p *JWKSProvider) (*Claims, error) {
    if p == nil {
        return nil, errors.New("jwks provider is nil")
    }
    token, err := gjwt.ParseWithClaims(tokenString, &Claims{}, p.Keyfunc, gjwt.WithLeeway(2*time.Second))
    if err != nil {
        return nil, err
    }
    claims, ok := token.Claims.(*Claims)
    if !ok || !token.Valid {
        return nil, errors.New("invalid token")
    }
    return claims, nil
}

