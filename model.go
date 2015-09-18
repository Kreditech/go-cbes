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

type Model struct {
    ID        int64  `type:"integer" analyzer:"standard"`
    TYPE      string `type:"string" analyzer:"keyword" index:"analyzed"`
    CreatedAt string `type:"date" format:"dateOptionalTime"`
    UpdatedAt string `type:"date" format:"dateOptionalTime"`
}

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
    name := getModelName(model)

    added := modelsCache.add(name, model)
    if !added {
        return fmt.Errorf("%s model allready registered", name)
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

        if field == "ttl" {
            ttl, err := strconv.Atoi(m.Type().Field(i).Tag.Get("ttl"))
            if err != nil {
                panic(err)
            }

            ttl = ttl * 1000
            prop := modelMapping[modelName].(map[string]interface{})
            prop["_ttl"] = map[string]interface{}{
                "enabled": true,
                "default": ttl,
            }

            continue
        }

        mapping := convertModelTags(strings.Split(string(m.Type().Field(i).Tag), " "))
        if mapping != nil {
            prop := modelMapping[modelName].(map[string]interface{})["properties"].(map[string]interface{})
            prop[field] = mapping
        }
    }

    d := reflect.ValueOf(new(Model)).Elem()
    for i := 0; i < d.NumField(); i++ {
        field := d.Type().Field(i).Name
        mapping := convertModelTags(strings.Split(string(d.Type().Field(i).Tag), " "))

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

// check the values of the given model and if they are not set, set the default values from model Tags
func setModelDefaults(model interface{}) interface{} {
    m := reflect.ValueOf(model).Elem()

    for i := 0; i < m.NumField(); i++ {
        field     := m.Type().Field(i).Name
        def       := m.Type().Field(i).Tag.Get("default")
        fieldVal  := m.FieldByName(field)
        fieldKind := fieldVal.Kind()

        switch fieldKind {
        case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
            if fieldVal.Int() == 0 && def != "" {
                i, err := strconv.ParseInt(def, 10, 64)
                if err == nil {
                    fieldVal.SetInt(i)
                }

            }
        case reflect.Float32, reflect.Float64:
            if fieldVal.Float() == 0 && def != "" {
                i, err := strconv.ParseFloat(def, 64)
                if err == nil {
                    fieldVal.SetFloat(i)
                }
            }
        case reflect.String:
            if fieldVal.String() == "" && def != "" {
                fieldVal.SetString(def)
            }
        case reflect.Bool:
            if fieldVal.Bool() == false  && def != "" {
                i, err := strconv.ParseBool(def)
                if err == nil {
                    fieldVal.SetBool(i)
                }
            }
        }
    }

    return model
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