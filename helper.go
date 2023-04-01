package goksei

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func getExpireTime(rawToken string) (*time.Time, error) {
	token, _, err := jwt.NewParser().ParseUnverified(rawToken, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid jwt map claims")
	}

	exp, ok := claims["exp"]
	if !ok {
		return nil, fmt.Errorf("cannot find exp claim in the token")
	}

	expUnix, ok := exp.(float64)
	if !ok {
		return nil, fmt.Errorf("exp claim invalid")
	}

	t := time.Unix(int64(expUnix), 0)

	return &t, nil
}
