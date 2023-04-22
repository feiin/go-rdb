package rdb

var db = NewTree()

func Get(key string) (interface{}, error) {
	return db.Search([]byte(key)), nil
}

func GetString(key string) (string, error) {
	return db.Search([]byte(key)).(string), nil
}

func Set(key string, value interface{}) bool {
	return db.Insert([]byte(key), value)
}

func SetString(key string, value string) bool {
	return db.Insert([]byte(key), value)
}
