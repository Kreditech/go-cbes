package cbes

type orm struct {
    db *db
}

func (o *orm) Find(model interface{}) string {
    return "test test test !"
}

// Create a new ORM object with
func NewOrm() *orm {
    return new(orm)
}