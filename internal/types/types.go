package types

type Chromium struct {
	Id    int
	Url   string
	Title string
}

type RequestData struct {
	Hwid     string `json:"hwid"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type GoInfoObject struct {
	GoOS     string
	Kernel   string
	Core     string
	Platform string
	OS       string
	Hostname string
	CPUs     int
}
