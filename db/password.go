package db

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
)

// GenResetPassToken generates reset password token
func GenResetPassToken(userID int64, duration int) (string, error) {
	id := uuid.New()

	set, err := redisClient.SetNX(id.String(), fmt.Sprintf("%d", userID), time.Duration(duration)*time.Minute).Result()
	if err != nil {
		return "", err
	}

	if !set {
		return "", fmt.Errorf("unable to set value")
	}

	return id.String(), nil
}

// RevokeResetPassToken revokes a reset password token
func RevokeResetPassToken(token string) error {
	_, err := redisClient.Del(token).Result()
	return err
}

// ValidateResetPassToken validates reset password token
func ValidateResetPassToken(token string) (int64, error) {
	value, err := redisClient.Get(token).Result()
	if err == redis.Nil {
		return -1, fmt.Errorf("invalid token")
	} else if err != nil {
		return -1, err
	}

	userID, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return -1, fmt.Errorf("unable to parse user id")
	}

	return userID, nil
}

// setOnline set user online
func SetOnline(userID int64) error {
	set, err := onlineUsers.SetNX(strconv.FormatInt(userID, 10), "online", 5*time.Minute).Result()
	if err != nil {
		return err
	}

	if !set {
		return fmt.Errorf("unable to set value")
	}

	return nil
}

func IsOnline(userID int64) bool {
	_, err := onlineUsers.Get(strconv.FormatInt(userID, 10)).Result()
	if err != nil {
		return false
	}

	return true
}
