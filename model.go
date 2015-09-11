package cbes

import (
    "reflect"
    "sync"
    "fmt"
    "strings"
    "regexp"
    "encoding/json"
    "strconv"
)

type _models struct  {
    mux sync.RWMutex
    cache map[string]interface{}
}

var modelsCache = &_models{cache: make(map[string]interface{})}

// check and add a model to cache
func (m *_models) add(name string, model interface{}) (added bool) {
    m.mux.Lock()
    defer m.mux.Unlock()

    if _, ok := m.cache[name]; ok == false {
        m.cache[name] = model
        added = true
    }

    return
}

// get the model from cache
func (m *_models) get(name string) (model interface{}, ok bool) {
    m.mux.RLock()
    defer m.mux.RUnlock()

    model, ok = m.cache[name]
    return
}

// register a model
func registerModel (model interface{}) error {
    _model := reflect.TypeOf(model).Elem()
    modelName := strings.ToLower(_model.Name())

    added := modelsCache.add(modelName, model)
    if !added {
        return fmt.Errorf("%s model allready registered", _model.Name())
    }

    return nil
}

// convert model tags into objects and prepare them for ElasticSearch mapping builder
func convertModelTags(_tags []string) interface{} {
    data := make(map[string]interface{})

    reg, err := regexp.Compile("[^-A-Za-z0-9_:]+")
    if err != nil {
        panic(err)
    }

    for _, val := range _tags {
        clearString := reg.ReplaceAllString(val, "")
        tags := strings.Split(clearString, ":")

        if (tags[0] == "default") {
            continue
        }

        if len(tags) < 2 {
            return nil
        }

        if len(tags) == 2 {
            i, err := strconv.ParseInt(tags[1], 10, 64)
            if err == nil {
                data[string(tags[0])] = i
                continue
            }

            b, err := strconv.ParseBool(tags[1])
            if err == nil {
                data[string(tags[0])] = b
                continue
            }

            data[string(tags[0])] = tags[1]
        } else if len(tags) > 2 {
            obj := []string{strings.Join(tags[1:], ":")}
            data[tags[0]] = convertModelTags(obj)
        }
    }

    return data
}

// get the view name from model interface
func getModelName(model interface{}) (string){
    name := strings.ToLower(reflect.TypeOf(model).Elem().Name())
    return name
}

// build ElasticSearch mapping from model struct tags
func buildModelMapping(model interface{}) string {
    m := reflect.ValueOf(model).Elem()
    modelName := getModelName(model)

    modelMapping := make(map[string]interface{})
    modelMapping[modelName] = map[string]interface{}{
        "properties": make(map[string]interface{}),
    }

    for i := 0; i < m.NumField(); i++ {
        field := m.Type().Field(i).Name
        mapping := convertModelTags(strings.Split(string(m.Type().Field(i).Tag), " "))

        if mapping != nil {
            prop := modelMapping[modelName].(map[string]interface{})["properties"].(map[string]interface{})
            prop[field] = mapping
        }
    }

    mappingJson, err := json.Marshal(modelMapping)
    if err != nil {
        fmt.Println(err)
    }

    return string(mappingJson)
}

// import all models mapping and view into CouchBase and ElasticSearch
func importAllModels() error {
    err := createModelViewsCB(modelsCache.cache)
    if err != nil {
        fmt.Println(err)
        return err
    }

    for _, model := range modelsCache.cache {
        modelMapping := buildModelMapping(model)

        err := addMapping(modelMapping, getModelName(model))
        if err != nil  {
            return err
        }
    }

    return nil
}