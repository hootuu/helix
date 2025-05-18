package helix

func Use(helix Helix) {
	doRegister(helix)
}

func Startup() {
	doStartup()
}

func AfterStartup(call func()) {
	gHelixStartupSuccessOn.Add(call)
}
