package helix

func Use(helix Helix) {
	doRegister(helix)
}

func Startup() {
	doStartup()
}
