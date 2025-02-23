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

func (j *JWT) Create(data JWTData) (string, error) {
	//метод шифрования
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": data.Email,
		//данные
	})
	s, err := t.SignedString([]byte(j.Secret)) // подпись
	if err != nil {
		return "", err
	}
	return s, nil
}
func (j *JWT) Parse(token string) (bool, *JWTData) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.Secret), nil // передача секрета для парсинга токена
	})
	if err != nil {
		return false, nil
	}
	email := t.Claims.(jwt.MapClaims)["email"]
	return t.Valid, &JWTData{
		Email: email.(string),
	}

}
