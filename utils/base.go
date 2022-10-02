package utils

type ITest interface {
	Add(num int) int
}

type Test struct {
	num int
}

func (t Test) Add(num int) int {
	return t.num + num
}

func Exec() {
	t := &Test{num: 10}
	println(t.Add(10))
}
