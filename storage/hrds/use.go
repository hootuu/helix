package hrds

func Register(code string) {
	doRegister(code)
}

func GetCache(code string) *Cache {
	return doGetCache(code)
}
