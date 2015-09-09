package cbes

import (
    "gopkg.in/couchbaselabs/gocb.v0"
    "time"
    "fmt"
    "encoding/json"
)

type Bucket struct {
    Name string
    Pass string
    OperationTimeout int //seconds
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

func InsertDesignDocument (name string) error {
    aux := dbSettings
    fmt.Println(name)
    json1, err := json.Marshal(aux)
    if err != nil {
                fmt.Println("Error: ")
                fmt.Println(err)
                return nil
            }
    fmt.Println(string(json1))
//    fmt.Println(aux.CouchBase.UserName)
//    fmt.Println(aux.CouchBase.Pass)
//    bManager := Connection.cb.Manager(dbSettings.CouchBase.UserName,dbSettings.CouchBase.Pass)
//
//    fmt.Printf("% +v\n", bManager)
//
//    dDocuments, err := bManager.GetDesignDocuments()
//    if err != nil {
//        fmt.Println("Error: ")
//        fmt.Println(err)
//        return err
//    }
//
//    for i := range dDocuments {
//        json1, err := json.Marshal(dDocuments[i])
//        if err != nil {
//            fmt.Println("Error: ")
//            fmt.Println(err)
//        }
//        fmt.Println(string(json1))
//    }
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
