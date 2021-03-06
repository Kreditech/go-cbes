package cbes_test
import (
    "github.com/Kreditech/go-cbes"
    "testing"
    "reflect"
)

func TestNewOrm(t *testing.T) {
    o := cbes.NewOrm()

    if reflect.TypeOf(o) != reflect.TypeOf(new(cbes.Orm)) {
        t.Fatalf("NewOrm() different than Orm")
    }
}

func TestCreate(t *testing.T) {
    var err error
    o := cbes.NewOrm()

    resModel, err := o.Create(&testModel)
    if err != nil {
        t.Fatal(err)
    }

    _ = resModel.(TestModel)

    resModel, err = o.Create(&testModelTtl)
    if err != nil {
        t.Fatal(err)
    }

    _ = resModel.(TestModelTTL)
}

func TestCreateEach(t *testing.T) {
    var err error
    var models []interface{}
    var modelsTtl []interface{}
    o := cbes.NewOrm()

    for i := 0; i < 10; i++ {
        models = append(models, &testModel)
        modelsTtl = append(modelsTtl, &testModelTtl)
    }

    createdModels, err := o.CreateEach(models...)
    if err != nil {
        t.Fatal(err)
    }

    for _, m := range createdModels {
        _ = m.(TestModel)
    }

    createdModels, err = o.CreateEach(modelsTtl...)
    if err != nil {
        t.Fatal(err)
    }

    for _, m := range createdModels {
        _ = m.(TestModelTTL)
    }
}

func TestCount(t *testing.T) {
    o := cbes.NewOrm()
    q := `{
        "query": {
            "bool": {
                "must": [
                    {
                        "term": {
                            "Name": "` + testModel.Name + `"
                        }
                    }
                ]
            }
        }
    }`

    count := o.Find(&testModel).Where(q).Count()
    if count != 11 {
        t.Fatalf("Wrong Count")
    }
}

func TestUpdate(t *testing.T) {
    o := cbes.NewOrm()
    q := `{
        "query": {
            "bool": {
                "must": [
                    {
                        "term": {
                            "Name": "` + testModel.Name + `"
                        }
                    }
                ]
            }
        }
    }`

    res := o.Find(&testModel).Where(q).Do()
    if len(res) < 1 {
        t.Fatalf("No results found")
    }

    for i := 0; i < len(res); i++ {
        m := res[i].(TestModel)

        newAge := i * 100
        m.Age = int64(newAge)
        m.StringArray = []string{}

        err := o.Update(&m)
        if err != nil {
             t.Fatal(err)
        }
    }

    q = `{
        "query": {
            "bool": {
                "must": [
                    {
                        "term": {
                            "Age": 200
                        }
                    }
                ]
            }
        }
    }`

    res = o.Find(&testModel).Where(q).Do()
    if len(res) != 1 {
        t.Fatalf("Wrong update count!")
    }
}

func TestFind(t *testing.T) {
    o := cbes.NewOrm()
    q := `{
        "query": {
            "bool": {
                "must": [
                    {
                        "term": {
                            "Name": "` + testModel.Name + `"
                        }
                    }
                ]
            }
        }
    }`

    res := o.Find(&testModel).Where(q).Do()
    if len(res) < 1 {
        t.Fatalf("No results found")
    }

    if len(res) != 10 {
        t.Fatalf("Return wrong count")
    }

    m := res[0].(TestModel)
    if reflect.TypeOf(&m) != reflect.TypeOf(new(TestModel)) {
        t.Fatalf("Return type not matching")
    }
}

func errorWhere(t *testing.T) {
    defer func() {
        if e := recover(); e == nil {
            t.Fatal("Expecting error!")
        }
    }()

    o := cbes.NewOrm()
    _ = o.Where("").Find(&testModel)
}

func TestWhere(t *testing.T) {
    errorWhere(t)

    o := cbes.NewOrm()
    q := `{
        "query": {
            "bool": {
                "must": [
                    {
                        "term": {
                            "Name": "` + testModel.Name + `"
                        }
                    }
                ]
            }
        }
    }`

    res := o.Find(&testModel).Where(q)
    if reflect.TypeOf(res) != reflect.TypeOf(new(cbes.Orm)) {
        t.Fatalf("Where() different than Orm")
    }
}

func TestDo(t *testing.T) {
    o := cbes.NewOrm()
    q := `{
        "query": {
            "bool": {
                "must": [
                    {
                        "term": {
                            "Name": "` + testModel.Name + `"
                        }
                    }
                ]
            }
        }
    }`

    res := o.Find(&testModel).Where(q).Do()
    if reflect.TypeOf(res) != reflect.TypeOf([]interface{}{}) {
        t.Fatalf("Do() wrong type")
    }
}

func errorOrder(t *testing.T) {
    defer func() {
        if e := recover(); e == nil {
            t.Fatal("Expecting error!")
        }
    }()

    o := cbes.NewOrm()
    _ = o.Order("ID", true)
}

func TestOrder(t *testing.T) {
    errorOrder(t)

    o := cbes.NewOrm()
    q := `{
        "query": {
            "bool": {
                "must": [
                    {
                        "term": {
                            "Name": "` + testModel.Name + `"
                        }
                    }
                ]
            }
        }
    }`

    res := o.Find(&testModel).Where(q).Order("ID", true).Do()
    if len(res) < 2 {
        t.Fatalf("No results found")
    }

    m1 := res[0].(TestModel)
    m2 := res[1].(TestModel)

    if m1.ID > m2.ID {
        t.Fatalf("Order results is wrong")
    }

    res = o.Find(&testModel).Where(q).Order("ID", false).Do()
    if len(res) < 2 {
        t.Fatalf("No results found")
    }
}

func errorLimit(t *testing.T) {
    defer func() {
        if e := recover(); e == nil {
            t.Fatal("Expecting error!")
        }
    }()

    o := cbes.NewOrm()
    _ = o.Limit(1)
}

func TestLimit(t *testing.T) {
    errorLimit(t)

    o := cbes.NewOrm()
    q := `{
        "query": {
            "bool": {
                "must": [
                    {
                        "term": {
                            "Name": "` + testModel.Name + `"
                        }
                    }
                ]
            }
        }
    }`

    res := o.Find(&testModel).Where(q).Limit(1).Do()
    if len(res) != 1 {
        t.Fatalf("Limit() Returns wrong")
    }
}

func errorFrom(t *testing.T) {
    defer func() {
        if e := recover(); e == nil {
            t.Fatal("Expecting error!")
        }
    }()

    o := cbes.NewOrm()
    _ = o.From(1)
}

func TestFrom(t *testing.T) {
    errorFrom(t)

    o := cbes.NewOrm()
    q := `{
        "query": {
            "bool": {
                "must": [
                    {
                        "term": {
                            "Name": "` + testModel.Name + `"
                        }
                    }
                ]
            }
        }
    }`

    res := o.Find(&testModel).Where(q).Limit(2).From(3).Do()
    if len(res) != 2 {
        t.Fatalf("From() Returns wrong")
    }
}

func errorAggregate(t *testing.T) {
    defer func() {
        if e := recover(); e == nil {
            t.Fatal("Expecting error!")
        }
    }()

    o := cbes.NewOrm()
    _ = o.Aggregate("")
}

func TestAggregate(t *testing.T) {
    errorAggregate(t)

    o := cbes.NewOrm()
    q := `{
        "query": {
            "bool": {
                "must": [
                    {
                        "term": {
                            "Name": "` + testModel.Name + `"
                        }
                    }
                ]
            }
        }
    }`

    aggQuery := `{
          "test_agg": {
              "terms": {
                  "field": "Age"
              },
              "aggs":{
                  "IDS_max":{
                      "max": {
                          "field": "Age"
                      }
                  }
              }
          }
      }
    `

    res := o.Find(&testModel).Where(q).Aggregate(aggQuery).Do()
    if len(res) < 1 {
        t.Fatalf("Aggregate() Returns wrong")
    }
}

func TestGetCollection (t *testing.T) {
    o := cbes.NewOrm()
    collection, err := o.GetCollection(&testModel)
    if err != nil {
        t.Fatal(err)
    }

    if len(collection) < 1 {
        t.Fatalf("No results found")
    }

    m := collection[0].(TestModel)
    if reflect.TypeOf(&m) != reflect.TypeOf(new(TestModel)) {
        t.Fatalf("Return type not matching")
    }
}

func TestGetRawCollection (t *testing.T) {
    o := cbes.NewOrm()
    collection, err := o.GetRawCollection(&testModel)
    if err != nil {
        t.Fatal(err)
    }

    if len(collection) < 1 {
        t.Fatalf("No results found")
    }

    if reflect.TypeOf(collection) != reflect.TypeOf([]interface{}{}) {
        t.Fatalf("GetRawCollection() wrong type")
    }
}

func TestReindex(t *testing.T) {
    o := cbes.NewOrm()

    err := o.Reindex(&testModel)
    if err != nil {
        t.Fatal(err)
    }

    q := `{
        "query": {
            "bool": {
                "must": [
                    {
                        "term": {
                            "Name": "` + testModel.Name + `"
                        }
                    }
                ]
            }
        }
    }`

    res := o.Find(&testModel).Where(q).Do()
    if len(res) < 1 {
        t.Fatalf("No results found")
    }

    if len(res) != 10 {
        t.Fatalf("Return wrong count")
    }

    m := res[0].(TestModel)
    if reflect.TypeOf(&m) != reflect.TypeOf(new(TestModel)) {
        t.Fatalf("Return type not matching")
    }
}

func TestDestroy (t *testing.T) {
    o := cbes.NewOrm()
    q := `{
        "query": {
            "bool": {
                "must": [
                    {
                        "term": {
                            "Age": 300
                        }
                    }
                ]
            }
        }
    }`

    affected, err := o.Destroy(&testModel, q)
    if err != nil {
        t.Fatal(err)
    }

    if len(affected) == 0 {
        t.Fatalf("Objects not destroyed")
    }

    for _, deletedModel := range affected {
        _ = deletedModel.(TestModel)
    }

    affected, err = o.Destroy(&testModel, "")
    if err != nil {
        t.Fatal(err)
    }

    if len(affected) == 0 {
        t.Fatalf("Objects not destroyed")
    }

    for _, deletedModel := range affected {
        _ = deletedModel.(TestModel)
    }
}