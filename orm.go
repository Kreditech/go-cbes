package cbes

import (
    "fmt"
)

var tmpQuery = make(map[string]interface{})
var queryTemplate = map[string]interface{}{
    "query": map[string]interface{}{
        "filtered": map[string]interface{}{
            "query": map[string]interface{}{
                "bool": map[string]interface{}{
                    "must": []interface{}{
                        map[string]interface{}{
                            "term": map[string]interface{}{
                                "_type": map[string]string{
                                    "value": "",
                                },
                            },
                        },
                    },
                },
            },
            "filter": make(map[string]interface{}),
        },
    },
}

type functions interface  {
    Find()       *orm
    FindOne()    *orm
    Create()
    CreateEach()
    Update()
    Destroy()
    Reindex()
    Aggregate()  *orm
    Limit()      *orm
    Order()      *orm
    Skip()       *orm
    Do()         interface{}
}


type orm struct {
    db *db
}

// set find query in ElasticSearch
func (o *orm) Find(model interface{}, query interface{}) *orm {
    tmpQuery["query"] = queryTemplate["query"]

    fmt.Println(queryTemplate)
    return o
}

// set limit to Find() query in ElasticSearch
func (o *orm) Limit(limit int) *orm {
    if len(tmpQuery) == 0 {
        panic("You must declare Find() first!")
    }

    tmpQuery["size"] = limit
    fmt.Println(tmpQuery)
    return o
}

// create new document in CouchBase and Elasticsearch
func (o *orm) Create(model interface{}) error {
    return nil
}

// Create a new ORM object with
func NewOrm() *orm {
    return new(orm)
}