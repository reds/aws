package dynamodb

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/reds/aws/config"
)

type AttributeUpdateValue struct {
	Action string
	Value  AttributeValue
}
type UpdateItemRequest struct {
	Expected                    map[string]interface{}    `json:",omitempty"`
	Key                         map[string]AttributeValue `json:",omitempty"`
	AttributeUpdates            map[string]AttributeUpdateValue
	ReturnConsumedCapacity      string `json:",omitempty"`
	ReturnItemCollectionMetrics string `json:",omitempty"`
	ReturnValues                string `json:",omitempty"`
	TableName                   string `json:",omitempty"`
}

type UpdateItemResult struct {
	Attributes       map[string]AttributeValue `json:",omitempty"`
	ConsumedCapacity struct {
		CapacityUnits int
		TableName     string
	} `json:",omitempty"`
	ItemCollectionMetrics struct {
		ItemCollectionKey   AttributeValue
		SizeEstimateRangeGB int
	}
}

func UpdateItem(cfg *config.Config4, req *UpdateItemRequest) (*UpdateItemResult, error) {
	if req == nil {
		return nil, errors.New("UpdateItem: need a request")
	}
	r, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := doRequest(cfg, "UpdateItem", bytes.NewReader(r))
	if err != nil {
		return nil, err
	}
	var tr UpdateItemResult
	err = decodeResult(resp, &tr)
	if err != nil {
		return nil, err
	}
	return &tr, nil
}
