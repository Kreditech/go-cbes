package cbes_test
import (
    "testing"
    "fmt"
    "github.com/Kreditech/go-cbes"
    "reflect"
)

func TestColorLog (t *testing.T) {
    var err error = fmt.Errorf("test")
    cbes.ColorLog("[WARN] %s", err)
    cbes.ColorLog("[SUCC] %s", err)
    cbes.ColorLog("[ERRO] %s", err)
    cbes.ColorLog("[TRAC] %s", err)
    cbes.ColorLog("[INFO] %s", err)
}

func TestColorLogS (t *testing.T) {
    var err error = fmt.Errorf("test")
    log := cbes.ColorLogS("[TEST] %s", err)

    if reflect.TypeOf(log).Kind() != reflect.String {
        t.Fatalf("Type not matching")
    }
}