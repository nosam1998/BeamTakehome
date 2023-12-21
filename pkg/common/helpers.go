package common

import "encoding/base64"

func EncodeToBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func DecodeBase64(data []byte) ([]byte, error) {
	decodedData, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return nil, err
	}
	return decodedData, nil
}
