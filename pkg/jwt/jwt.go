package jwt

import "github.com/golang-jwt/jwt/v5"

type JWTData struct {
	Email string
}

type JWT struct {
	Secret string
}

func NewJWT(secret string) *JWT {
	return &JWT{
		Secret: secret,
	}
}

func (j *JWT) CreateToken(data JWTData) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": data.Email,
	})
	secret, err := token.SignedString([]byte(j.Secret))
	if err != nil {
		return "", err
	}
	return secret, nil
}

func (j *JWT) ParseToken(token string) (bool, *JWTData) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.Secret), nil
	})
	if err != nil {
		return false, nil
	}
	email := t.Claims.(jwt.MapClaims)["email"].(string)
	return t.Valid, &JWTData{
		Email: email,
	}
}
