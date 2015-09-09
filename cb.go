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

func connectCb(settings *Settings) (*gocb.Cluster, error) {
    cluster, err := gocb.Connect(settings.CouchBase.Host)

    if err != nil {
        return nil, err
    }

    return cluster, err;
}

func openBucket (bucket *Bucket, cluster *gocb.Cluster) (*gocb.Bucket, error) {
    b, err := cluster.OpenBucket(bucket.Name, bucket.Pass)
    b.SetOperationTimeout(time.Duration(bucket.OperationTimeout)* time.Second)

    if err != nil {
        return nil, err
    }

    return b, err
}

func OpenCb(settings *Settings) (*gocb.Bucket, error) {
    cluster, err := connectCb(settings)
    if err != nil {
        return nil, err
    }

    bucket, err := openBucket(settings.CouchBase.Bucket, cluster)
    if err != nil {
        return nil, err
    }

    return bucket, nil
}

func createViewCb (name string) bool {
    return true
}