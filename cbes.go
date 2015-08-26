package cbes

import (
    _ "go-cbes/es"
    "gopkg.in/olivere/elastic.v2"
//    "go-cbes/es"
    "os"
    "go-cbes/connection"
)

var (
    DebugLog = NewLog(os.Stderr)
)

// cbes configuration
type Setting struct {
    ElasticSearch struct {
        Urls              string
        Bucket            string
        NumberOfShards    int
        NumberOfReplicas  int
    }
    CouchBase string
}

// connections
type cbesConnection struct {
    es elastic.Client
}

func RegisterDataBase (aliasName string, settings Setting) {
    conn, err := connection.Open(aliasName, settings)

    if err != nil {

    }
}

func Client (settings Setting) (Setting) {
    return settings
//    cbesConn.es = es.Connect(esOptions)
}



