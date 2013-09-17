package dynamodb

import (
	"encoding/json"
	"github.com/reds/aws/config"
	"testing"
)

//	tr, err := UpdateItem ( cfg, &UpdateItemRequest{TableName:"wcnfavs",Key:map[string]AttributeValue{"broadcaster":{S:"dyntest"}}, AttributeUpdates:map[string]AttributeUpdateValue{"favs":{Action:"ADD", Value:AttributeValue{"favs":{SS:[]string{"f1","f2","f4"}}}}}})

func TestJson(t *testing.T) {
	//	au := map[string]AttributeUpdateValue{"favs":{Action:"ADD",Value:AttributeValue{SS:[]string{"a","b"}}}}
	b, err := json.Marshal(&UpdateItemRequest{TableName: "wcnfavs", Key: map[string]AttributeValue{"broadcaster": {S: "dyntest"}}, AttributeUpdates: map[string]AttributeUpdateValue{"favs": {Action: "ADD", Value: AttributeValue{SS: []string{"a", "b"}}}}})
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
	r, err := UpdateItem(cfg, &UpdateItemRequest{TableName: "wcnfavs", Key: map[string]AttributeValue{"broadcaster": {S: "dyntest"}}, AttributeUpdates: map[string]AttributeUpdateValue{"hides": {Action: "DELETE", Value: AttributeValue{SS: []string{"a", "d"}}}}})
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
	tr, err := GetItem(cfg, &GetItemRequest{TableName: "wcnfavs", Key: map[string]AttributeValue{"broadcaster": {S: "dyntest"}}, AttributesToGet: []string{"favs", "hides"}})
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
	tr, err := PutItem(cfg, &PutItemRequest{TableName: "wcnfavs", Item: map[string]AttributeValue{"broadcaster": {S: "dyntest"}, "favs": {SS: []string{"f1", "f2", "f4"}}}})
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
	ltr, err := DescribeTable(cfg, &DescribeTableRequest{"wcnfavs"})
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
	ltr, err = ListTables(cfg, &ListTablesRequest{"wcn", 10})
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
