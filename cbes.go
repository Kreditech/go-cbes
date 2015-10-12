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
                      Host         string
                      Bucket       *Bucket
                      UserName     string
                      Pass         string
                      ViewsOptions *ViewsOptions
                  }
}

var dbSettings = &Settings{}

// connections
type cbesConnection struct {
    es elastic.Client
}

// Register DataBase connection
func RegisterDataBase(settings *Settings) error {
    var err error

    err = open(settings)
    if err != nil {
        goto printError
    }

    dbSettings = settings

    err = importAllModels()
    if err != nil {
        goto printError
    }

printError:
    if err != nil {
        ColorLog("[ERRO] CBES: %s\n", err)
        return err
    }

    return nil
}

// Register a model or array of models
func RegisterModel(models ...interface{}) error{
    for _, model := range models {
        err := registerModel(model)
        if err != nil {
            ColorLog("[ERRO] CBES: register mode failed %s\n", err)
            return err
        }
    }

    return nil
}