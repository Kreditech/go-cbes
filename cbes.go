package cbes

import (
    "gopkg.in/olivere/elastic.v2"
    "os"
)

var (
    DebugLog = NewLog(os.Stderr)
)

// cbes configuration
type Setting struct {
    ElasticSearch struct {
                      Urls             string
                      Bucket           string
                      NumberOfShards   int
                      NumberOfReplicas int
                  }
    CouchBase     struct {
                      Host             string
                      Port             int
                  }
}

// connections
type cbesConnection struct {
    es elastic.Client
}

func RegisterDataBase(aliasName string, settings *Setting) {
    err := Open(aliasName, settings)

    if err != nil {

    }
}

func Client(settings Setting) (Setting) {
    return settings
    //    cbesConn.es = es.Connect(esOptions)
}



