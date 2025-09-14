package jwt

import (
    "crypto/rsa"
    "encoding/json"
    "errors"
    "net/http"

    "github.com/lestrrat-go/jwx/v2/jwk"
    gjwt "github.com/golang-jwt/jwt/v5"
)

// JWKSProvider держит набор публичных ключей для валидации токенов.
type JWKSProvider struct {
    set jwk.Set
}

// NewJWKSProvider из публичного ключа и kid создаёт набор из одного ключа.
func NewJWKSProvider(pub *rsa.PublicKey, kid string) (*JWKSProvider, error) {
    if pub == nil {
        return nil, errors.New("nil public key")
    }
    key, err := jwk.FromRaw(pub)
    if err != nil {
        return nil, err
    }
    _ = key.Set(jwk.KeyIDKey, kid)
    _ = key.Set(jwk.AlgorithmKey, "RS256")
    set := jwk.NewSet()
    set.AddKey(key)
    return &JWKSProvider{set: set}, nil
}

// Keyfunc для github.com/golang-jwt/jwt — ищет ключ по kid.
func (p *JWKSProvider) Keyfunc(token *gjwt.Token) (interface{}, error) {
    if p == nil || p.set == nil {
        return nil, errors.New("jwks not initialized")
    }
    kid, _ := token.Header["kid"].(string)
    if kid == "" {
        return nil, errors.New("missing kid")
    }
    key, ok := p.set.LookupKeyID(kid)
    if !ok {
        return nil, errors.New("key not found")
    }
    var pub rsa.PublicKey
    if err := key.Raw(&pub); err != nil {
        return nil, err
    }
    return &pub, nil
}

// Handler возвращает стандартный http.Handler, который отдаёт JWKS JSON.
func (p *JWKSProvider) Handler() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if p == nil || p.set == nil {
            w.WriteHeader(http.StatusServiceUnavailable)
            return
        }
        w.Header().Set("Content-Type", "application/json")
        _ = json.NewEncoder(w).Encode(p.set)
    })
}

