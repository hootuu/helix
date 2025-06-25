package hdb

func Register(code string) {
	doRegister(code)
}

func GetDatabase(code string) *Database {
	return doGetDb(code)
}
