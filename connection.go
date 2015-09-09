package cbes

import (
    "gopkg.in/couchbaselabs/gocb.v0"
    "gopkg.in/olivere/elastic.v2"
    "fmt"
)

type db struct {
    es *elastic.Client
    cb *gocb.Bucket
}

// Opens DB connection
func Open(settings *Settings) error {
    var err error
    db := new(db)

    db.es, err = OpenEs(settings)
    if err != nil {
        err = fmt.Errorf("register ElasticSearch %s", err.Error())
        goto end
    }

    db.cb, err = OpenCb(settings)
    if err != nil {
        err = fmt.Errorf("register CouchBase %s", err.Error())
        goto end
    }

end:
    return err
}

