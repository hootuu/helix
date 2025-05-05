package helix

func Use(helix Helix) {
	doRegister(helix)
}

func Startup() {
	doStartup()
}

func OnStartupSuccess(call func()) {
	gHelixStartupSuccessOn.Add(call)
}
