package cbes

import (
    "gopkg.in/couchbase/gocb.v1"
    "gopkg.in/olivere/elastic.v2"
    "fmt"
)

var connection = new(db)

type db struct {
    es *elastic.Client
    cb *gocb.Bucket
}

// Opens DB connection
func open(settings *Settings) error {
    var err error

    connection.es, err = openEs(settings)
    if err != nil {
        err = fmt.Errorf("register ElasticSearch %s", err.Error())
        goto end
    }

    connection.cb, err = openCb(settings)
    if err != nil {
        err = fmt.Errorf("register CouchBase %s", err.Error())
        goto end
    }

end:
    return err
}

