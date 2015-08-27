package cbes

import (
    "gopkg.in/olivere/elastic.v2"
    "gopkg.in/couchbaselabs/gocb.v0"

    "sync"
    "fmt"
)

var dbCache = &dataBaseCache{cache: make(map[string]*alias)}

type DB struct {
    es elastic.Client
    cb gocb.Cluster
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
func Open(aliasName string, settings *Setting) error {
    var (
        err error
//        db *DB
    )

//    db.cb, err = OpenCb(settings)
//    if err != nil {
//        err = fmt.Errorf("register cb `%s`, %s", aliasName, err.Error())
//        goto end
//    }

//    db.es, err = OpenEs(settings)
//    if err != nil {
//        err = fmt.Errorf("register es `%s`, %s", aliasName, err.Error())
//        goto end
//    }
//
//    _, err = addAlias(aliasName, db)
//    if err != nil {
//        goto end
//    }

//end:
    return err
}

