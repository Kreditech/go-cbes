package cbes

import (
    "gopkg.in/couchbaselabs/gocb.v0"
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
func OpenCb(settings *Settings) (*gocb.Bucket, error) {
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

// get the view from the model interface
func getView (model interface{}) string {
    viewName := getModelName(model)
    return "function (doc, meta) {if(doc._TYPE && doc._TYPE == '" + viewName + "') {emit(meta.id, {doc: doc, meta: meta});}}"
}

// create and upsert the view in CouchBase
func createModelViewsCB(models map[string]interface{}) error {
    manager := Connection.cb.Manager(dbSettings.CouchBase.UserName, dbSettings.CouchBase.Pass)
    views := map[string]gocb.View{}

    for _, model := range models {
        newView := gocb.View{}
        newView.Map = getView(model)

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

// insert on CouchBase and return the Id
func createCB(model interface{}) (int64, error) {
    var err error
    var count string

    modelName := getModelName(model)
    cb := *Connection.cb

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

    reflect.ValueOf(model).Elem().FieldByName("ID").SetInt(id)

    _, err = cb.Insert(key, model, 0)
    if err != nil {
        return 0, err
    }

    return id, err
}

// insert each on CouchBase
func createEachCB(models []interface{}) error {
    for i := range models {
        _, err := createCB(models[i])
        if err != nil {
            return err
        }
    }

    return nil
}

// remove on CouchBase
func destroyCB(modelId string) error {
    _, err := Connection.cb.Remove(modelId, 0)
    if err != nil {
        return err
    }

    return nil
}

// update on CouchBase
func updateCB(modelId string, replaceModel interface{}) error {
    _, err := Connection.cb.Replace(modelId, &replaceModel, 0, 0)
    if err != nil {
        return err
    }

    return nil
}