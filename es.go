package cbes

import (
    "gopkg.in/olivere/elastic.v2"
)

// connect to elastic search
func connect (settings *Setting) (elastic.Client, error) {
    client, err := elastic.NewClient(elastic.SetURL(settings.ElasticSearch.Urls))

    if err != nil  {
        return nil, err
    }

    return client, nil
}

// Open connection
func OpenEs (setting *Setting) (elastic.Client, error) {
    client, err := connect(setting)

    if err != nil {
        return nil, err
    }

    return client, nil
}