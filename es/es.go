package es

import (
    "gopkg.in/olivere/elastic.v2"
    "fmt"
)

// connect to elastic search
func Connect (urls ...string) (elastic.Client) {
    client, err := elastic.NewClient(elastic.SetURL(urls))

    if err != nil  {
        fmt.Errorf(err)
    }

    return client
}
//
//func CheckIndex (name string, nr) {
//
//}