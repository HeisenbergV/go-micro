package micro

import "testing"

func Test_Init(t *testing.T) {
	serv := NewService()
	serv.Init()
	serv.Server().Handle(nil)
	serv.Run()
}
