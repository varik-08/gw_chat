package jwt

import (
    "crypto/rand"
    "crypto/rsa"
    "crypto/x509"
    "encoding/pem"
    "errors"
    "time"

    gjwt "github.com/golang-jwt/jwt/v5"
)

// Claims — минимальные клеймы нашего токена.
type Claims struct {
    UserID   int64  `json:"userId"`
    Username string `json:"username"`
    gjwt.RegisteredClaims
}

// RSAKeyPair хранит приватный ключ и kid.
type RSAKeyPair struct {
    PrivateKey *rsa.PrivateKey
    KeyID      string
}

// GenerateRSAKey генерирует новый RSA ключ.
func GenerateRSAKey(bits int, kid string) (*RSAKeyPair, error) {
    if bits == 0 {
        bits = 2048
    }
    pk, err := rsa.GenerateKey(rand.Reader, bits)
    if err != nil {
        return nil, err
    }
    return &RSAKeyPair{PrivateKey: pk, KeyID: kid}, nil
}

// MarshalRSAPrivateKeyPEM кодирует приватный ключ в PEM.
func MarshalRSAPrivateKeyPEM(priv *rsa.PrivateKey) ([]byte, error) {
    if priv == nil {
        return nil, errors.New("nil private key")
    }
    b := x509.MarshalPKCS1PrivateKey(priv)
    blk := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: b}
    return pem.EncodeToMemory(blk), nil
}

// ParseRSAPrivateKeyPEM парсит PEM приватного ключа.
func ParseRSAPrivateKeyPEM(pemBytes []byte) (*rsa.PrivateKey, error) {
    blk, _ := pem.Decode(pemBytes)
    if blk == nil {
        return nil, errors.New("invalid pem")
    }
    return x509.ParsePKCS1PrivateKey(blk.Bytes)
}

// SignTokenRS256 подписывает токен RS256, проставляя kid в заголовок.
func SignTokenRS256(priv *rsa.PrivateKey, kid string, userID int64, username string, ttl time.Duration) (string, time.Time, error) {
    if ttl <= 0 {
        ttl = 15 * time.Minute
    }
    exp := time.Now().Add(ttl)
    claims := &Claims{
        UserID:   userID,
        Username: username,
        RegisteredClaims: gjwt.RegisteredClaims{
            ExpiresAt: gjwt.NewNumericDate(exp),
        },
    }
    token := gjwt.NewWithClaims(gjwt.SigningMethodRS256, claims)
    if kid != "" {
        token.Header["kid"] = kid
    }
    s, err := token.SignedString(priv)
    if err != nil {
        return "", time.Time{}, err
    }
    return s, exp, nil
}

