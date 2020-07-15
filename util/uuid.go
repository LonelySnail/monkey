package utils

import uuid "github.com/satori/go.uuid"

func UUid() string {
	return uuid.NewV4().String()
}
