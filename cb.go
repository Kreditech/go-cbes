package cbes

import (
    "gopkg.in/couchbaselabs/gocb.v0"
    "time"
    "fmt"
    "strings"
    "reflect"
)

type Bucket struct {
    Name string
    Pass string
    OperationTimeout int //seconds
}

func getViewName (model interface{}) (string){
    viewName := strings.ToLower(reflect.TypeOf(model).Elem().Name())
    return viewName
}

func getView (model interface{}) (string){
    viewName := strings.ToLower(reflect.TypeOf(model).Elem().Name())
    return "function (doc, meta) {if(doc._TYPE && doc._TYPE == '" + viewName + "') {emit(meta.id, {doc: doc, meta: meta});}}"
}

func connectCb(settings *Settings) (*gocb.Cluster, error) {
    cluster, err := gocb.Connect(settings.CouchBase.Host)

    if err != nil {
        return nil, err
    }

    return cluster, err;
}

func openBucket (settings *Settings, cluster *gocb.Cluster) (*gocb.Bucket, error) {
    bucket := settings.CouchBase.Bucket

    b, err := cluster.OpenBucket(bucket.Name, bucket.Pass)
    b.SetOperationTimeout(time.Duration(bucket.OperationTimeout)* time.Second)

    if err != nil {
        return nil, err
    }

    return b, err
}

func createViewCB(models map[string]interface{}) error{
    manager := Connection.cb.Manager(dbSettings.CouchBase.UserName,dbSettings.CouchBase.Pass)
    views := map[string]gocb.View{}

    for _, model := range models {
        newView := gocb.View{}
        newView.Map = getView(model)

        views[getViewName(model)] = newView
    }

    dDocument := gocb.DesignDocument{}
    dDocument.Name = "udk"
    dDocument.Views = views

    err := manager.UpsertDesignDocument(&dDocument)
    if err != nil {
        fmt.Println("InsertDesignDocument Error: ")
        fmt.Println(err)
        return err
    }

    return nil
}

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
