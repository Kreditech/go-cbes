package cbes

import (
    "fmt"
    "time"
    "reflect"
    "encoding/json"
    "strconv"
)

var tmpModel interface{}
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
    "sort": []interface{}{},
}

type functions interface  {
    Find()       *Orm
    Where()      *Orm
    Create()
    CreateEach()
    Update()
    Destroy()
    Reindex()
    Aggregate()  *Orm
    Limit()      *Orm
    Order()      *Orm
    From()       *Orm
    Do()         interface{}
}


type Orm struct {
    db *db
}

// Execute builded query
func (o *Orm) Do() []interface{} {
    jsonQuery, err := json.Marshal(tmpQuery)
    if err != nil {
        panic(err)
    }

    res := searchEs(string(jsonQuery))
    data := []interface{}{}

    if (res.TotalHits() == 0) {
        return data
    }

    for _, hit := range res.Hits.Hits {
        item := make(map[string]interface{})
        err := json.Unmarshal(*hit.Source, &item)
        if err != nil {
            panic(err)
        }

        data = append(data, setModel(tmpModel, item))
    }

    return data
}

// Set the model witch you want to find
// If you don't set Limit Find ElasticSearch will use the default Limit (10)
func (o *Orm) Find(model interface{}) *Orm {
    typeName := getModelName(model)
    _model, ok := modelsCache.get(typeName); if ok {
        tmpModel = _model
    }

    // clone queryTemplate into tmpQuery
    _copy, err := json.Marshal(queryTemplate)
    if err != nil {
        panic(err)
    }

    err = json.Unmarshal(_copy, &tmpQuery)
    if err != nil {
        panic(err)
    }

    // set tmpQuery model type
    typeVal := tmpQuery["query"].(map[string]interface{})["filtered"]
    typeVal = typeVal.(map[string]interface{})["query"]
    typeVal = typeVal.(map[string]interface{})["bool"]
    typeVal = typeVal.(map[string]interface{})["must"].([]interface{})[0]
    typeVal = typeVal.(map[string]interface{})["term"]
    typeVal = typeVal.(map[string]interface{})["_type"]
    typeVal.(map[string]interface{})["value"] = typeName

    return o
}

// Set query for ElasticSearch
func (o *Orm) Where(query string) *Orm {
    if len(tmpQuery) == 0 {
        panic("You must declare Find() first!")
    }

    var q map[string]interface{}

    err := json.Unmarshal([]byte(query), &q)
    if err != nil {
        return o
    }

    tmpQuery["query"].(map[string]interface{})["filtered"].(map[string]interface{})["filter"] = q
    return o
}

// Pagination of results.  The from parameter defines the offset from the first result you want to fetch
// The Limit() function allows you to configure the maximum amount of hits to be returned. From() default is 0, and
// Limit() default is 10
func (o *Orm) From(from int) *Orm {
    if len(tmpQuery) == 0 {
        panic("You must declare Find() first!")
    }

    tmpQuery["from"] = from
    return o
}

// Set limit to Find() query in ElasticSearch. Max limit = 999999999
func (o *Orm) Limit(limit int) *Orm {
    if len(tmpQuery) == 0 {
        panic("You must declare Find() first!")
    }

    tmpQuery["size"] = limit
    return o
}

// Set ElasticSearch sort. direction: true = 'asc', false = 'desc'
func (o *Orm) Order(field string, direction bool) *Orm {
    if len(tmpQuery) == 0 {
        panic("You must declare Find() first!")
    }

    dir := "asc"
    if direction == false {
        dir = "desc"
    }

    sort := map[string]interface{}{
        field: map[string]interface{}{
            "order": dir,
        },
    }

    typeVal := tmpQuery["sort"].([]interface{})
    typeVal = append(typeVal, sort)

    tmpQuery["sort"] = typeVal
    return o
}

// Create new document in CouchBase and Elasticsearch
func (o *Orm) Create(m interface{}) error {
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
func (o *Orm) CreateEach(models ...interface{}) error {
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
func (o *Orm) Destroy(model interface{}, query string) error {
    var err error

    models := o.Find(model).Where(query).Limit(999999999).Do()
    for _, m := range models {
        _m := reflect.ValueOf(m)
        id := _m.FieldByName("TYPE").String() + ":" + strconv.FormatInt(_m.FieldByName("ID").Int(), 10)

        err = destroyCB(id)
        if err != nil {
            return fmt.Errorf("cbes.Destroy() CouchBase %s", err.Error())
        }

        err = destroyES(id, model)
        if err != nil {
            return fmt.Errorf("cbes.Destroy() ElastiSearch %s", err.Error())
        }
    }


    return nil
}

// Update a document in CouchBase and ElasticSearch
func (o *Orm) Update(model interface{}, query string) error {
    var err error
    models := o.Find(model).Where(query).Limit(999999999).Order("ID", true).Do()

    for _, m := range models {
        _m            := reflect.ValueOf(m)
        modelType     := _m.FieldByName("TYPE").String()
        modelID       := _m.FieldByName("ID").Int()
        createdAt     := _m.FieldByName("CreatedAt").String()
        id            := modelType + ":" + strconv.FormatInt(modelID, 10)
        timeFormatted := time.Now().Format(time.RFC3339)
        _model        := reflect.ValueOf(model).Elem()
        modelClone    := reflect.New(_model.Type()).Elem()

        for i := 0; i < _model.NumField(); i++ {
            field      := _model.Field(i)
            name       := _model.Type().Field(i).Name
            _mField    := _m.FieldByName(name)
            cloneField := modelClone.FieldByName(name)

            switch field.Kind() {
            case reflect.Int64:
                var empty, val int64

                if field.Int() != empty {
                    val = field.Int()
                } else {
                    val = _mField.Int()
                }

                cloneField.SetInt(val)
            case reflect.Float64:
                var empty, val float64

                if field.Float() != empty {
                    val = field.Float()
                } else {
                    val = _mField.Float()
                }

                cloneField.SetFloat(val)
            case reflect.String:
                var empty, val string

                if field.String() != empty {
                    val = field.String()
                } else {
                    val = _mField.String()
                }

                cloneField.SetString(val)
            case reflect.Bool:
                var empty, val bool

                if field.Bool() != empty {
                    val = field.Bool()
                } else {
                    val = _mField.Bool()
                }

                cloneField.SetBool(val)
            case reflect.Slice, reflect.Map:
                val := reflect.ValueOf(field.Interface())

                if val.Len() == 0 {
                    val = reflect.ValueOf(_mField.Interface())
                }

                cloneField.Set(val)
            }
        }

        modelClone.FieldByName("UpdatedAt").SetString(timeFormatted)
        modelClone.FieldByName("CreatedAt").SetString(createdAt)
        modelClone.FieldByName("TYPE").SetString(modelType)
        modelClone.FieldByName("ID").SetInt(modelID)

        setModel := modelClone.Interface()
        err = updateCB(id, setModel)
        if err != nil {
            return fmt.Errorf("cbes.Update() CouchBase %s", err.Error())
        }

        err = updateES(id, setModel)
        if err != nil {
            return fmt.Errorf("cbes.Update() ElasticSearch %s", err.Error())
        }
    }

    return nil
}

// Create a new ORM object with
func NewOrm() *Orm {
    return new(Orm)
}