package cbes

import (
    "gopkg.in/olivere/elastic.v2"
)

// cbes configuration
type Settings struct {
    ElasticSearch struct {
                      Urls             []string
                      Index            string
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

func RegisterDataBase(aliasName string, settings *Settings) {
    err := Open(aliasName, settings)

    if err != nil {
        ColorLog("[ERRO] CBES: %s\n", err)
    }
}

func Client(settings Settings) (Settings) {
    return settings
    //    cbesConn.es = es.Connect(esOptions)
}



