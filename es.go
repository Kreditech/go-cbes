package cbes

import (
    "gopkg.in/olivere/elastic.v2"
    "fmt"
)

// connect to elastic search and build the client
func connect (settings *Settings) (*elastic.Client, error) {
    client, err := elastic.NewClient(elastic.SetURL(settings.ElasticSearch.Urls...))
    if err != nil {
        return nil , err
    }

    return client, nil
}

// check if the index exists
func checkIndex(settings *Settings, client *elastic.Client) (bool, error) {
    exists, err := client.IndexExists(settings.ElasticSearch.Index).Do()
    if err != nil {
        return false, err
    }

    if !exists {
        return false, nil
    }

    return true, nil
}

// create
func createIndex(settings *Settings, client *elastic.Client) (bool, error) {
    builder := client.CreateIndex(settings.ElasticSearch.Index)
    fmt.Println("==========================")
    fmt.Println(builder)
    fmt.Println("==========================")

    return true, nil
}

// Open connection
func OpenEs (settings *Settings) (*elastic.Client, error) {
    client, err := connect(settings)
    if err != nil {
        return nil, err
    }

    indexExists, err := checkIndex(settings, client)
    if err != nil {
        return nil, err
    }

    if !indexExists {
        createIndex(settings, client)
    }

    return client, nil
}