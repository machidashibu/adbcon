package model

type Version struct {
	App string
	Api string
	Adb string
	Sdk string
}

func (v Version) APP() string {
	return v.App
}

func (v Version) API() string {
	return v.Api
}

func (v Version) ADB() string {
	return v.Adb
}

func (v Version) SDK() string {
	return v.Sdk
}
