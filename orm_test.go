package cbes_test
import (
    "go-cbes"
    "testing"
    "reflect"
    "time"
    "strconv"
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

    err = o.Create(&testModel)
    if err != nil {
        t.Fatal(err)
    }

    err = o.Create(&testModelTtl)
    if err != nil {
        t.Fatal(err)
    }
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

    err = o.CreateEach(models...)
    if err != nil {
        t.Fatal(err)
    }

    err = o.CreateEach(modelsTtl...)
    if err != nil {
        t.Fatal(err)
    }

    time.Sleep(2500 * time.Millisecond)
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

    m := res[0].(TestModel)
    if reflect.TypeOf(&m) != reflect.TypeOf(new(TestModel)) {
        t.Fatalf("Return type not matching")
    }

    qUpdate := `{
        "query": {
            "bool": {
                "must": [
                    {
                        "term": {
                            "ID": ` + strconv.FormatInt(testModel.ID, 10) + `
                        }
                    }
                ]
            }
        }
    }`

    m.Age = 300
    err := o.Update(m, qUpdate)
    if err != nil {
        t.Fatal(err)
    }

    time.Sleep(2500 * time.Millisecond)
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

func TestWhere(t *testing.T) {
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

func TestOrder(t *testing.T) {
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

func TestLimit(t *testing.T) {
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

func TestFrom(t *testing.T) {
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

func TestAggregate(t *testing.T) {
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

    time.Sleep(2500 * time.Millisecond)

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

    err := o.Destroy(&testModel, q)
    if err != nil {
        t.Fatal(err)
    }

    time.Sleep(2500 * time.Millisecond)
    res := o.Find(&testModel).Where(q).Do()
    if len(res) > 0 {
        t.Fatalf("Objects not destroied")
    }

    err = o.Destroy(&testModel, "")
    if err != nil {
        t.Fatal(err)
    }

    time.Sleep(2500 * time.Millisecond)
    res = o.Find(&testModel).Where("").Do()
    if len(res) > 0 {
        t.Fatalf("Objects not destroied")
    }
}