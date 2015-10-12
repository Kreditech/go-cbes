package cbes_test
import (
    "testing"
    "go-cbes"
)

type TestModel struct {
    cbes.Model
    Name        string          `default:"TestName" type:"string" analyzer:"keyword" index:"analyzed"`
    LastName    string          `default:"Test Last Name" type:"string" analyzer:"standard" index:"analyzed"`
    Age         int64           `default:"25" type:"integer" analyzer:"standard" index:"not_analyzed"`
    Active      bool            `default:"true" type:"boolean"`
    Float       float64         `default:"12345.00" type:"float"`
    Int         int64           `default:"321" type:"long"`
    StringArray []string        `type:"string" analyzer:"keyword" index:"analyzed"`
    IntArray    []int64         `type:"integer" analyzer:"keyword" index:"analyzed"`
    FloatArray  []float64       `type:"float" analyzer:"keyword" index:"analyzed"`
    Interface   map[string]interface{} `type:"object" properties:"{'name':{'type':'object','enabled':false},'sid':{'type':'string','index':'not_analyzed'}}"`
    Nested      []interface{}   `type:"nested" properties:"{'first': {'type': 'string'}, 'last':{'type': 'string'}}"`
}

type TestModelTTL struct {
    cbes.Model
    Name        string          `default:"TestName" type:"string" analyzer:"keyword" index:"analyzed"`
    LastName    string          `default:"Test Last Name" type:"string" analyzer:"standard" index:"analyzed"`
    Age         int64           `default:"25" type:"integer" analyzer:"standard" index:"not_analyzed"`
    Active      bool            `default:"true" type:"boolean"`
    Float       float64         `default:"12345.00" type:"float"`
    Int         int64           `default:"321" type:"long"`
    StringArray []string        `type:"string" analyzer:"keyword" index:"analyzed"`
    IntArray    []int64         `type:"integer" analyzer:"keyword" index:"analyzed"`
    FloatArray  []float64       `type:"float" analyzer:"keyword" index:"analyzed"`
    Interface   map[string]interface{} `type:"object" properties:"{'name':{'type':'object','enabled':false},'sid':{'type':'string','index':'not_analyzed'}}"`
    Nested      []interface{}   `type:"nested" properties:"{'first': {'type': 'string'}, 'last':{'type': 'string'}}"`
    ttl         int64           `ttl:"60"` //ttl in seconds
}

func TestRegisterDataBase(t *testing.T) {
    settings := new(cbes.Settings)
    settings.ElasticSearch.Urls = []string{"http://192.168.33.10:9200"}
    settings.ElasticSearch.Index = "testindex"
    settings.ElasticSearch.NumberOfShards = 0
    settings.ElasticSearch.NumberOfReplicas = 0

    settings.CouchBase.Host = "192.168.33.10:8091"
    settings.CouchBase.UserName = "root"
    settings.CouchBase.Pass = "root123"

    bucket := new(cbes.Bucket)
    bucket.Name = "test"
    bucket.Pass = ""
    bucket.OperationTimeout = 5 // seconds

    settings.CouchBase.Bucket = bucket

    viewsOptions := new(cbes.ViewsOptions)
    viewsOptions.UpdateInterval = 5000
    viewsOptions.UpdateMinChanges = 5
    viewsOptions.ReplicaUpdateMinChanges = 5

    settings.CouchBase.ViewsOptions = viewsOptions

    err := cbes.RegisterDataBase(settings)
    if err != nil {
        t.Fatal(err)
    }
}

func TestRegisterModel(t *testing.T) {
    err := cbes.RegisterModel(new(TestModel), new(TestModelTTL))
    if err != nil {
        t.Fatal(err)
    }
}