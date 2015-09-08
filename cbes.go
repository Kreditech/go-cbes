package cbes

import (
    "gopkg.in/olivere/elastic.v2"
    "fmt"
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
                      Host     string
                      Bucket   *Bucket
                      UserName string
                      Pass     string
                  }
}

// connections
type cbesConnection struct {
    es elastic.Client
}

// Register DataBase connection
func RegisterDataBase(aliasName string, settings *Settings) {
    err := Open(aliasName, settings)

    if err != nil {
        ColorLog("[ERRO] CBES: %s\n", err)
    }
}

// Register a model or array of models
func RegisterModel(models ...interface{}) {
    for _, model := range models {
        fmt.Println(model)
    }
}



