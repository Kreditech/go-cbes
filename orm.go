package cbes
import "fmt"

type orm struct {
    alias *alias
}

func (o *orm) Find(model interface{}) string {
    return "test test test !"
}

// Switch to another registered database driver by given name.
func (o *orm) Using(name string) error {
    if al, ok := dbCache.get(name); ok {
        o.alias = al
    } else {
        return fmt.Errorf("<Orm.Using> unknown db alias name `%s`", name)
    }

    return nil
}

// Create a new ORM object with
func NewOrm() *orm {
    o := new(orm)
    err := o.Using("default")

    if err != nil {
        panic(err)
    }

    return o
}