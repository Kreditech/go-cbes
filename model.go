package cbes

import (
    "reflect"
    "sync"
    "fmt"
    "strings"
    "encoding/json"
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

func buildModelMapping(model interface{}) string {
    m := reflect.ValueOf(model).Elem()

    for i := 0; i < m.NumField(); i++ {
        field := m.Type().Field(i).Name
        val := m.Field(i).Interface()
        tags := strings.Split(string(m.Type().Field(i).Tag), " ")

        if field == "Mapping" {
            x, err := json.Marshal(val)
            if err != nil {
                fmt.Println(err)
            } else {
                fmt.Println(string(x))
            }
        }

        fmt.Println("-----------------------")
        fmt.Println(field)
        fmt.Println(val)
        fmt.Println(tags)
    }

    return ""
}

// import all models mapping and view into CouchBase and ElasticSearch
func importAllModels() error {
    //for _, model := range modelsCache.cache {
        //mapping := buildModelMapping(model)
        //createViewsCB(getViewName(model), getView(model))
        //fmt.Println(mapping)
//    }
    err := createViewCB(modelsCache.cache)
    if err != nil {
        fmt.Println(err)
        return err
    }

    return nil
}