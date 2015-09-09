package cbes

import "reflect"

func registerModel (model interface{}) error {

    view := reflect.TypeOf(model).Elem()

    err := InsertDesignDocument(view.Name())
    if err != nil {
        return err
    }
    return nil
}