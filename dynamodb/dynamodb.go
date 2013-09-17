package dynamodb

import (
	"bytes"
	"encoding/json"
	"github.com/reds/aws/config"
	"github.com/reds/aws/sign"
	"io"
	"io/ioutil"
	"net/http"
)

type AttributeValue struct {
	B  string   `json:",omitempty"`
	BS []string `json:",omitempty"`
	N  string   `json:",omitempty"`
	NS []string `json:",omitempty"`
	S  string   `json:",omitempty"`
	SS []string `json:",omitempty"`
}

type GetItemRequest struct {
	TableName              string
	Key                    map[string]AttributeValue
	AttributesToGet        []string `json:",omitempty"`
	ConsistentRead         bool     `json:",omitempty"`
	ReturnConsumedCapacity string   `json:",omitempty"`
}

type GetItemResult struct {
	ConsumedCapacity struct {
		CapacityUnits int
		TableName     string
	} `json:",omitempty"`
	Item map[string]AttributeValue `json:",omitempty"`
}

func GetItem(cfg *config.Config4, req *GetItemRequest) (*GetItemResult, error) {
	r := []byte("{}")
	var err error
	if req != nil {
		r, err = json.Marshal(req)
		if err != nil {
			return nil, err
		}
	}
	body := bytes.NewReader(r)
	resp, err := doRequest(cfg, "GetItem", body)
	if err != nil {
		return nil, err
	}
	var tr GetItemResult
	err = decodeResult(resp, &tr)
	if err != nil {
		return nil, err
	}
	return &tr, nil
}

type PutItemRequest struct {
	Expected                    map[string]interface{}    `json:",omitempty"`
	Item                        map[string]AttributeValue `json:",omitempty"`
	ReturnConsumedCapacity      string                    `json:",omitempty"`
	ReturnItemCollectionMetrics string                    `json:",omitempty"`
	ReturnValues                string                    `json:",omitempty"`
	TableName                   string                    `json:",omitempty"`
}

type PutItemResult struct {
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

func PutItem(cfg *config.Config4, req *PutItemRequest) (*PutItemResult, error) {
	r := []byte("{}")
	var err error
	if req != nil {
		r, err = json.Marshal(req)
		if err != nil {
			return nil, err
		}
	}
	body := bytes.NewReader(r)
	resp, err := doRequest(cfg, "PutItem", body)
	if err != nil {
		return nil, err
	}
	var tr PutItemResult
	err = decodeResult(resp, &tr)
	if err != nil {
		return nil, err
	}
	return &tr, nil
}

type DescribeTableRequest struct {
	TableName string
}

type DescribeTableResult struct {
	Table struct {
		AttributeDefinitions []struct {
			AttributeName string
			AttributeType string
		}
		CreationDateTime float32
		ItemCount        int
		KeySchema        []struct {
			AttributeName string
			KeyType       string
		}
		LocalSecondaryIndexes []struct {
			IndexName      string
			IndexSizeBytes int
			ItemCount      int
			KeySchema      []struct {
				AttributeName string
				KeyType       string
			}
			Projection struct {
				NonKeyAttributes []string
				ProjectionType   string
			}
		}
		ProvisionedThroughput struct {
			LastDecreaseDateTime   int
			LastIncreaseDateTime   int
			NumberOfDecreasesToday int
			ReadCapacityUnits      int
			WriteCapacityUnits     int
		}
		TableName      string
		TableSizeBytes int
		TableStatus    string
	}
}

func DescribeTable(cfg *config.Config4, req *DescribeTableRequest) (*DescribeTableResult, error) {
	r := []byte("{}")
	var err error
	if req != nil {
		r, err = json.Marshal(req)
		if err != nil {
			return nil, err
		}
	}
	body := bytes.NewReader(r)
	resp, err := doRequest(cfg, "DescribeTable", body)
	if err != nil {
		return nil, err
	}
	var tr DescribeTableResult
	err = decodeResult(resp, &tr)
	if err != nil {
		return nil, err
	}
	return &tr, nil
}

type ListTablesRequest struct {
	ExclusiveStartTableName string
	Limit                   int
}

type ListTablesResult struct {
	LastEvaluatedTableName string
	TableNames             []string
}

func ListTables(cfg *config.Config4, req *ListTablesRequest) (*ListTablesResult, error) {
	r := []byte("{}")
	var err error
	if req != nil {
		r, err = json.Marshal(req)
		if err != nil {
			return nil, err
		}
	}
	body := bytes.NewReader(r)
	resp, err := doRequest(cfg, "ListTables", body)
	if err != nil {
		return nil, err
	}
	var ltr ListTablesResult
	err = decodeResult(resp, &ltr)
	if err != nil {
		return nil, err
	}
	return &ltr, nil
}

func decodeResult(resp *http.Response, r interface{}) error {
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	resp.Body.Close()
	err = json.Unmarshal(b, r)
	if err != nil {
		return err
	}
	return nil
}

func doRequest(cfg *config.Config4, target string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("POST", "https://"+cfg.Host+"/", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Host", cfg.Host)
	req.Header.Set("Content-Type", "application/x-amz-json-1.0")
	req.Header.Set("x-amz-target", "DynamoDB_20120810."+target)
	sign.SignV4(req, "dynamodb", cfg.Region, cfg.AccessKeyId, cfg.Secret)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
