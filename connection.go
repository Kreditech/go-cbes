package cbes

import (
    "gopkg.in/couchbaselabs/gocb.v0"
    "gopkg.in/olivere/elastic.v2"
    "sync"
    "fmt"
)

var dbCache = &dataBaseCache{cache: make(map[string]*alias)}

type DB struct {
    es *elastic.Client
    cb *gocb.Bucket
}

type alias struct {
    Name       string
    Connection *DB
}

type dataBaseCache struct {
    mux   sync.RWMutex
    cache map[string]*alias
}

// Add database connection
func (ch *dataBaseCache) add(name string, connection *alias) (added bool) {
    ch.mux.Lock()
    defer ch.mux.Unlock()

    if _, ok := ch.cache[name]; ok == false {
        ch.cache[name] = connection
        added = true
    }

    return
}

// get database alias if cached.
func (ch *dataBaseCache) get(name string) (al *alias, ok bool) {
    ch.mux.RLock()
    defer ch.mux.RUnlock()

    al, ok = ch.cache[name]
    return
}

// add database alias with original name.
func addAlias(aliasName string, db *DB) (*alias, error) {
    al := new(alias)

    al.Name = aliasName
    al.Connection = db

    if dbCache.add(aliasName, al) == false {
        return nil, fmt.Errorf("DataBase alias name `%s` already registered, cannot reuse", aliasName)
    }

    return al, nil
}

// Opens an DB specified by its aliasName
func Open(aliasName string, settings *Settings) error {
    var err error
    db := new(DB)

    db.es, err = OpenEs(settings)
    if err != nil {
        err = fmt.Errorf("register es `%s`, %s", aliasName, err.Error())
        goto end
    }

    fmt.Print(db)

    db.cb, err = OpenCb(settings)
    if err != nil {
        err = fmt.Errorf("register cb `%s`, %s", aliasName, err.Error())
        goto end
    }

    _, err = addAlias(aliasName, db)
    if err != nil {
        goto end
    }

end:
    return err
}

