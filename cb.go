package cbes

import (
    "gopkg.in/couchbaselabs/gocb.v0"
    "time"
)

type Bucket struct {
    Name string
    Pass string
    OperationTimeout int //seconds
}


func OpenCb(setting *Setting) (gocb.Cluster, error){
    cluster, err := gocb.Connect(setting.CouchBase.Host + ":" + setting.CouchBase.Port)

    if err != nil {
        return nil, err
    }

    return cluster, err;
}

func RegisteBucket (bucket *Bucket) (gocb.Bucket, error) {
    cluster := gocb.Cluster
    b, err := cluster.OpenBucket(bucket.Name, bucket.Pass)
    b.SetOperationTimeout(time.Duration(bucket.OperationTimeout)* time.Second)

    if err != nil {
        return nil, err
    }
     return b, err
}
