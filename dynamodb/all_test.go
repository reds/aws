package dynamodb

import (
	"encoding/json"
	"github.com/reds/aws/config"
	"testing"
)

//	tr, err := UpdateItem ( cfg, &UpdateItemRequest{TableName:"mprtest",Key:map[string]AttributeValue{"id":{S:"dyntest"}}, AttributeUpdates:map[string]AttributeUpdateValue{"favs":{Action:"ADD", Value:AttributeValue{"favs":{SS:[]string{"f1","f2","f4"}}}}}})

func TestJson(t *testing.T) {
	//	au := map[string]AttributeUpdateValue{"favs":{Action:"ADD",Value:AttributeValue{SS:[]string{"a","b"}}}}
	b, err := json.Marshal(&UpdateItemRequest{TableName: "mprtest", Key: map[string]AttributeValue{"id": {S: "dyntest"}}, AttributeUpdates: map[string]AttributeUpdateValue{"favs": {Action: "ADD", Value: AttributeValue{SS: []string{"a", "b"}}}}})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(b))
}

func TestUpdateItem(t *testing.T) {
	cfg, err := config.LoadConfig("/tmp/aws.cfg", "dynamodb")
	if err != nil {
		t.Fatal(cfg, err)
	}
	r, err := UpdateItem(cfg, &UpdateItemRequest{TableName: "mprtest", Key: map[string]AttributeValue{"id": {S: "dyntest"}}, AttributeUpdates: map[string]AttributeUpdateValue{"hides": {Action: "DELETE", Value: AttributeValue{SS: []string{"a", "d"}}}}})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(r)
}

func TestGetItem(t *testing.T) {
	cfg, err := config.LoadConfig("/tmp/aws.cfg", "dynamodb")
	if err != nil {
		t.Fatal(cfg, err)
	}
	tr, err := GetItem(cfg, &GetItemRequest{TableName: "mprtest", Key: map[string]AttributeValue{"id": {S: "dyntest"}, "ts": {N: "80"}}, AttributesToGet: []string{"favs", "hides"}})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tr)
}

func TestPutItem(t *testing.T) {
	cfg, err := config.LoadConfig("/tmp/aws.cfg", "dynamodb")
	if err != nil {
		t.Fatal(cfg, err)
	}
	tr, err := PutItem(cfg, &PutItemRequest{TableName: "mprtest", Item: map[string]AttributeValue{"id": {S: "dyntest"}, "favs": {SS: []string{"f1", "f2", "f4"}}, "ts": {N: "90"}}})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tr)
}

func TestDescribeTables(t *testing.T) {
	cfg, err := config.LoadConfig("/tmp/aws.cfg", "dynamodb")
	if err != nil {
		t.Fatal(err)
	}
	ltr, err := DescribeTable(cfg, &DescribeTableRequest{"mprtest"})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ltr)
}

func TestListTables(t *testing.T) {
	cfg, err := config.LoadConfig("/tmp/aws.cfg", "dynamodb")
	if err != nil {
		t.Fatal(err)
	}
	ltr, err := ListTables(cfg, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ltr)
	ltr, err = ListTables(cfg, &ListTablesRequest{"mprtest", 10})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ltr)
	ltr, err = ListTables(cfg, &ListTablesRequest{"a", 10})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ltr)
}
