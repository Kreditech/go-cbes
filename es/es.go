package es

import (
    "gopkg.in/olivere/elastic.v2"
    "go-cbes"
)

// connect to elastic search
func connect (settings *cbes.Setting) (elastic.Client, error) {
    client, err := elastic.NewClient(elastic.SetURL(settings.ElasticSearch.Urls))

    if err != nil  {
        return nil, err
    }

    return client, nil
}

// Open connection
func Open (setting *cbes.Setting) (elastic.Client, error) {
    client, err := connect(setting)

    if err != nil {
        return nil, err
    }

    return client, nil
}