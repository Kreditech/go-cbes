package cbes

type model interface {
    Find()       string
    FindOne()
    Create()
    CreateEach()
    Update()
    Destroy()
    GetRaw()
    Aggregate()
    ReIndex()
}

type Model struct {
}

func (m *Model) Find() string {
    return "test test test"
}