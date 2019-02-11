package main

type Gear struct{}

func GetGear() interface{} {
	return &Gear{}
}

func (*Gear) GetName() string {
	return "hello"
}
