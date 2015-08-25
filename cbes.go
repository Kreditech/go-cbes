package cbes

import (
    _ "go-cbes/es"
    "gopkg.in/olivere/elastic.v2"
    "go-cbes/es"
)

type esSettings struct {
    name              string
    numberOfShards    int
    numberOfReplicas  int
}

type cbesSetting struct {
    esSettings esSettings
}

type cbesConnection struct {
    es elastic.Client
}

func (cbesConn *cbesConnection) Client(esOptions ...string) {
    cbesConn.es = es.Connect(esOptions)
}

