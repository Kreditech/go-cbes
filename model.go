package cbes

import (
    "reflect"
    "sync"
    "fmt"
    "strings"
    "encoding/json"
    "strconv"
    "regexp"
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
    reg, err := regexp.Compile("[^-A-Za-z0-9_:+{},' *-]+")

    if err != nil {
        panic(err)
    }

    for _, val := range _tags {
        clearString := reg.ReplaceAllString(val, "")
        tag := strings.Split(clearString, ":")

        if tag[0] == "default" || tag[0] == "json" {
            continue
        }

        if len(tag) < 2 {
            return nil
        }

        tagVal := tag[1]
        if len(tag) > 2 {
            tagVal = strings.Join(tag[1:], ":")
            tagVal = strings.Replace(tagVal, "'", "\"", -1)

            j := make(map[string]interface{})
            err := json.Unmarshal([]byte(tagVal), &j)
            if err != nil {
                fmt.Println(err)
                continue
            }

            data[string(tag[0])] = j
        } else {
            i, err := strconv.ParseInt(tagVal, 10, 64)
            if err == nil {
                data[string(tag[0])] = i
                continue
            }

            b, err := strconv.ParseBool(tagVal)
            if err == nil {
                data[string(tag[0])] = b
                continue
            }

            data[string(tag[0])] = tag[1]
        }
    }

    return data
}

// get the view name from model interface
func getModelName(model interface{}) (string){
    var name string
    _m := reflect.TypeOf(model)

    if _m.Kind() == reflect.Struct {
        name = _m.Name()
    } else {
        name = _m.Elem().Name()
    }
    return strings.ToLower(name)
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

        mapping := convertModelTags(strings.Split(string(m.Type().Field(i).Tag), "\" "))
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

// Receive ElasticSearch Hit and set a clean model with it
func setModel(_model, responseModel interface{}) interface{} {
    m := reflect.New(reflect.ValueOf(_model).Elem().Type()).Elem()
    r := responseModel.(map[string]interface{})

    for i := 0; i < m.NumField(); i++ {
        field := m.Type().Field(i).Name

        if field == "Model" {
            continue
        }

        set   := m.FieldByName(field)
        kind  := set.Kind()
        val   := r[field]

        switch kind {
        case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
            valType := reflect.TypeOf(val).Kind()

            if valType == reflect.Float64 {
                val = int64(val.(float64))
            }

            if valType == reflect.Float32 {
                val = int64(val.(float32))
            }

            set.SetInt(reflect.ValueOf(val).Int())
        case reflect.Float32, reflect.Float64:
            set.SetFloat(val.(float64))
        case reflect.String:
            set.SetString(val.(string))
        case reflect.Bool:
            set.SetBool(val.(bool))
        case reflect.Map:
            resMap := reflect.ValueOf(val)
            set.Set(resMap)
        case reflect.Slice:
            resSliceType := reflect.TypeOf(val)

            if resSliceType == reflect.TypeOf([]interface{}{}) {
                v := reflect.ValueOf(val)

                if v.Len() < 1 {
                    continue
                }

                modelSliceType := m.Field(i).Type()
                switch modelSliceType {
                case reflect.TypeOf([]string{}):
                    tmp := []string{}

                    for i := 0; i < v.Len(); i++ {
                        tmp = append(tmp, v.Index(i).Elem().String())
                    }

                    set.Set(reflect.ValueOf(tmp))
                case reflect.TypeOf([]int{}), reflect.TypeOf([]int8{}), reflect.TypeOf([]int16{}), reflect.TypeOf([]int32{}), reflect.TypeOf([]int64{}):
                    tmp := []int64{}

                    for i := 0; i < v.Len(); i++ {
                        tmp = append(tmp, int64(v.Index(i).Elem().Float()))
                    }

                    set.Set(reflect.ValueOf(tmp))
                case reflect.TypeOf([]float32{}), reflect.TypeOf([]float64{}):
                    tmp := []float64{}

                    for i := 0; i < v.Len(); i++ {
                        tmp = append(tmp, v.Index(i).Elem().Float())
                    }

                    set.Set(reflect.ValueOf(tmp))
                case reflect.TypeOf([]interface{}{}):
                    set.Set(v)
                }
            }
        }
    }

    id := r["ID"]
    idType := reflect.TypeOf(id).Kind()
    if idType == reflect.Float64 {
        id = int(id.(float64))
    }

    if idType == reflect.Float32 {
        id = int(id.(float32))
    }

    m.FieldByName("ID").SetInt(reflect.ValueOf(id).Int())
    m.FieldByName("TYPE").SetString(r["TYPE"].(string))
    m.FieldByName("CreatedAt").SetString(r["CreatedAt"].(string))
    m.FieldByName("UpdatedAt").SetString(r["UpdatedAt"].(string))

    return m.Interface()
}