package utils

import "github.com/google/uuid"

func Generate_uuid() string {
	return uuid.New().String()
}

func Parse_uuid(uid string) (bool, error) {
	_, err := uuid.Parse(uid)
	if err != nil {
		return false, err
	}
	return true, nil
}
