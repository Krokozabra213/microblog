package service

import "github.com/google/uuid"

func GenerateRandomID() string {
	return uuid.NewString()
}
