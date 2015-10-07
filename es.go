package cbes

import (
    "gopkg.in/olivere/elastic.v2"
    "fmt"
    "strconv"
    "reflect"
)

// connect to elastic search and build the client
func connectEs (settings *Settings) (*elastic.Client, error) {
    client, err := elastic.NewClient(elastic.SetURL(settings.ElasticSearch.Urls...))
    if err != nil {
        return nil , err
    }

    return client, nil
}

// Check if the index exists
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

// Create Index
func createIndex(settings *Settings, client *elastic.Client) (bool, error) {
    builder, err := client.CreateIndex(settings.ElasticSearch.Index).Do()
    if err != nil {
        return false, err
    }

    return builder.Acknowledged, nil
}

// Open connection
func openEs (settings *Settings) (*elastic.Client, error) {
    client, err := connectEs(settings)
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

// put model mapping if it doesn't exist
func addMapping(mapping string, modelName string) error {
    var err error
    index := dbSettings.ElasticSearch.Index
    es := *connection.es

    currentMapping, err := es.GetMapping().Index(index).Type(modelName).Do()
    if err != nil {
        return fmt.Errorf("expected get mapping response; got: %v", err.Error())
    }

    if len(currentMapping) > 0 {
        return nil
    }

    res, err := es.PutMapping().
        IgnoreConflicts(true).
        Index(index).
        Type(modelName).
        BodyString(mapping).Do()

    if err != nil {
        return err
    }
    if res == nil {
        return fmt.Errorf("expected put mapping response; got: %v", res)
    }
    if !res.Acknowledged {
        return fmt.Errorf("expected put mapping ack; got: %v", res.Acknowledged)
    }

    return nil
}

// delete model mapping
func deleteMapping(model interface{}) error {
    var err error
    index := dbSettings.ElasticSearch.Index
    es := *connection.es
    modelName := getModelName(model)

    _, err = es.GetMapping().Index(index).Type(modelName).Do()
    if err != nil {
        return nil
    }

    _, err = es.DeleteMapping().Index(index).Type(modelName).Do()
    if err != nil {
        return err
    }

    return nil
}

// create ElasticSearch document based on model
func createEs(id int64, model interface{}) error {
    modelName := getModelName(model)
    es := *connection.es
    index := dbSettings.ElasticSearch.Index
    key := modelName + ":" + strconv.FormatInt(id, 10)

    reflect.ValueOf(model).Elem().FieldByName("ID").SetInt(id)

    _, err := es.Index().
        Index(index).
        Type(modelName).
        Id(key).
        BodyJson(model).Do()

    if err != nil {
        return err
    }

    return nil
}

// search in ElasticSearch
func searchEs(query string) *elastic.SearchResult {
    es := *connection.es
    index := dbSettings.ElasticSearch.Index

    res, err := es.Search().Index(index).Source(query).Do()
    if err != nil {
        panic(err)
    }

    return res
}

// update on ElasticSearch
func updateES(id string, model interface{}) error {
    modelName := getModelName(model)
    es := *connection.es
    index := dbSettings.ElasticSearch.Index

    _, err := es.Update().
        Index(index).
        Type(modelName).
        Id(id).
        Doc(model).Do()

    if err != nil {
        return err
    }

    return nil
}

// delete document from ElasticSearch
func destroyES(id string, model interface{}) error {
    modelName := getModelName(model)
    es := *connection.es
    index := dbSettings.ElasticSearch.Index

    _, err := es.Delete().
    Index(index).
    Type(modelName).
    Id(id).Do()

    if err != nil {
        return err
    }

    return nil
}

// Insert the given models array to ElasticSearch
// This method is used by orm.Reindex() method
func importModelsToEs (models []interface{}) error {
    for _, model := range models {
        id := reflect.ValueOf(model).FieldByName("ID").Int()
        modelName := getModelName(model)
        es := *connection.es
        index := dbSettings.ElasticSearch.Index
        key := modelName + ":" + strconv.FormatInt(id, 10)

        _, err := es.Index().
            Index(index).
            Type(modelName).
            Id(key).
            BodyJson(model).Do()

        if err != nil {
            return err
        }
    }

    return nil
}