aws
===

Amazon Web Services tools


dynamodb
========

A thin veneer over the aws dynamodb service.

DescribeTable

ListTables

GetItem

PutItem

```go
    	cfg, err := config.LoadConfig("/tmp/aws.cfg", "dynamodb")
    	if err != nil {
    		t.Fatal(cfg, err)
    	}
    	tr, err := PutItem(cfg, &PutItemRequest{TableName: "mprtest", Item: map[string]AttributeValue{"id": {S: "dyntest"}, "favs": {SS: []string{"f1", "f2", "f4"}}, "ts":{N: "90"}}})
    	if err != nil {
    		t.Fatal(err)
    	}
```

UpdateItem