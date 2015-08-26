package connection

import (
    "gopkg.in/olivere/elastic.v2"
    "go-cbes"
    "sync"
    "go-cbes/es"
    "fmt"
)

type DB struct {
    es elastic.Client
    cb string
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

    if dataBaseCache.add(aliasName, al) == false {
        return nil, fmt.Errorf("DataBase alias name `%s` already registered, cannot reuse", aliasName)
    }

    return al, nil
}

// Opens an DB specified by its aliasName
func Open (aliasName string, settings *cbes.Setting) error {
    var (
        err error
        db *DB
    )

    db.es, err = es.Open(settings)
    if err != nil {
        err = fmt.Errorf("register es `%s`, %s", aliasName, err.Error())
        goto end
    }

    _, err = addAlias(aliasName, db)
    if err != nil {
        goto end
    }

end:
    return err
}

