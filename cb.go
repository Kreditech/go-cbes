package cbes

import (
    "gopkg.in/couchbase/gocb.v1"
    "time"
    "fmt"
    "strconv"
    "reflect"
)

type Bucket struct {
    Name             string
    Pass             string
    OperationTimeout int //seconds
}

// connect to CouchBase
func connectCb(settings *Settings) (*gocb.Cluster, error) {
    cluster, err := gocb.Connect(settings.CouchBase.Host)

    if err != nil {
        return nil, err
    }

    return cluster, err;
}

// open CouchBase bucket
func openBucket(settings *Settings, cluster *gocb.Cluster) (*gocb.Bucket, error) {
    bucket := settings.CouchBase.Bucket

    b, err := cluster.OpenBucket(bucket.Name, bucket.Pass)
    b.SetOperationTimeout(time.Duration(bucket.OperationTimeout) * time.Second)

    if err != nil {
        return nil, err
    }

    return b, err
}

// open CouchBase cluster and bucket
func openCb(settings *Settings) (*gocb.Bucket, error) {
    cluster, err := connectCb(settings)
    if err != nil {
        return nil, err
    }

    bucket, err := openBucket(settings, cluster)
    if err != nil {
        return nil, err
    }

    return bucket, nil
}

// Generate CouchBase view script for given model
func generateModelViewScript (model interface{}) string {
    viewName := getModelName(model)
    return "function (doc, meta) {if(doc.TYPE && doc.TYPE == '" + viewName + "') {emit(meta.id, {doc: doc, meta: meta});}}"
}

// Create and upsert the view in CouchBase
func createModelViewsCB(models map[string]interface{}) error {
    manager := *connection.cb.Manager(dbSettings.CouchBase.UserName, dbSettings.CouchBase.Pass)
    views := map[string]gocb.View{}

    for _, model := range models {
        newView := gocb.View{}
        newView.Map = generateModelViewScript(model)

        views[getModelName(model)] = newView
    }

    dDocument := gocb.DesignDocument{}
    dDocument.Name = dbSettings.CouchBase.Bucket.Name
    dDocument.Views = views

    err := manager.UpsertDesignDocument(&dDocument)
    if err != nil {
        fmt.Println("InsertDesignDocument Error: ")
        fmt.Println(err)
        return err
    }

    return nil
}

// Insert on CouchBase and return the Id
func createCB(model interface{}) (int64, error) {
    var err error
    var count string

    modelName := getModelName(model)
    cb := *connection.cb

    num, _, err := cb.Counter(modelName + ":count", 1, 1, 0)
    if err != nil {
        return 0, err
    }

    count = strconv.FormatUint(num, 10)
    key := modelName + ":" + count

    id, err := strconv.ParseInt(count, 10, 64)
    if err != nil {
        return 0, err
    }

    m := reflect.ValueOf(model).Elem()
    m.FieldByName("ID").SetInt(id)

    ttl := 0
    for i := 0; i < m.NumField(); i++ {
        field := m.Type().Field(i).Name

        if field == "ttl" {
            _ttl, err := strconv.Atoi(m.Type().Field(i).Tag.Get("ttl"))
            if err != nil {
                panic(err)
            }

            ttl = _ttl
        }
    }

    _, err = cb.Insert(key, model, uint32(ttl))
    if err != nil {
        return 0, err
    }

    return id, err
}

// Remove on CouchBase
func destroyCB(modelId string) error {
    _, err := connection.cb.Remove(modelId, 0)
    if err != nil {
        return err
    }

    return nil
}

// Update on CouchBase
func updateCB(modelId string, replaceModel interface{}) error {
    _, err := connection.cb.Replace(modelId, &replaceModel, 0, 0)
    if err != nil {
        return err
    }

    return nil
}

// Get all the data for a specific model
func getByView(model interface{}) ([]interface{}, error) {
    modelName := getModelName(model)
    data := []interface{}{}
    cb := *connection.cb

    query := gocb.NewViewQuery(dbSettings.CouchBase.Bucket.Name, modelName)
    rows, err := cb.ExecuteViewQuery(query)
    if err != nil {
        return data, err
    }

    var row interface{}
    for rows.Next(&row) {
        data = append(data, row)
    }

    rows.Close()
    return data, nil
}
