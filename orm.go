package cbes

import (
    "fmt"
    "time"
    "reflect"
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

// Set find query in ElasticSearch
func (o *orm) Find(model interface{}, query interface{}) *orm {
    tmpQuery["query"] = queryTemplate["query"]

    fmt.Println(queryTemplate)
    return o
}

// Set limit to Find() query in ElasticSearch
func (o *orm) Limit(limit int) *orm {
    if len(tmpQuery) == 0 {
        panic("You must declare Find() first!")
    }

    tmpQuery["size"] = limit
    fmt.Println(tmpQuery)
    return o
}

// Create new document in CouchBase and Elasticsearch
func (o *orm) Create(m interface{}) error {
    t             := time.Now()
    timeFormatted := t.Format(time.RFC3339)
    model         := setModelDefaults(m)

    reflect.ValueOf(model).Elem().FieldByName("CreatedAt").SetString(timeFormatted)
    reflect.ValueOf(model).Elem().FieldByName("UpdatedAt").SetString(timeFormatted)
    reflect.ValueOf(model).Elem().FieldByName("TYPE").SetString(getModelName(model))

    id, err := createCB(model)
    if err != nil {
        return fmt.Errorf("cbes.Create() CouchBase %s", err.Error())
    }

    err = createEs(id, model)
    if err != nil {
        return fmt.Errorf("cbes.Create() ElasticSearch %s", err.Error())
    }

    return nil
}

// Create a variadic of documents in CouchBase and ElasticSearch
func (o *orm) CreateEach(models ...interface{}) error {
    var err error

    for _, model := range models {
        err = o.Create(model)
        if err != nil {
            return fmt.Errorf("cbes.CreateEach() CouchBase %s", err.Error())
        }
    }

    return nil
}

// Destroy a document in CouchBase and ElasticSearch
func (o *orm) Destroy(filterQuery string) error {
    //TODO Insert Find function here and return the model id
    var modelId = "user:4" //TODO TO DELETE

    err := destroyCB(modelId)
    if err != nil {
        return fmt.Errorf("cbes.Destroy() CouchBase %s", err.Error())
    }
    return nil
}

// Update a document in CouchBase and ElasticSearch
func (o *orm) Update(filterQuery string, model interface{}) error {
    //TODO Insert Find function here and return the model
    var modelId = "user:10" //TODO TO DELETE

    err := updateCB(modelId, model)
    if err != nil {
        return fmt.Errorf("cbes.Update() CouchBase %s", err.Error())
    }
    return nil
}

// Create a new ORM object with
func NewOrm() *orm {
    return new(orm)
}