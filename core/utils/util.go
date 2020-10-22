package utils

import (
	"encoding/json"
	"log"

	"github.com/sulenn/go-http/core/types"
)

func ParseHttpResponse(bytes []byte) (*types.ResponseJSON, error) {
	tempResponseJSON := &types.TempResponseJSON{}
	//responseJson := &types.ResponseJSON{}
	err := json.Unmarshal(bytes, tempResponseJSON)
	if err != nil {
		log.Printf("bytes unmarshal failed: %+v\n", err)
		return nil, err
	}
	err = json.Unmarshal([]byte(tempResponseJSON.Data), &tempResponseJSON.ResponseJSON.Data)
	if err != nil {
		log.Printf("[]byte(tempResponseJSON.Data) unmarshal failed: %+v\n", err)
		return nil, err
	}
	return tempResponseJSON.ResponseJSON, nil
}
