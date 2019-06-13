package test

import (
	"fmt"
	"github.com/fanliao/go-promise"
	"testing"
)

func TestTask (t *testing.T){
	task1 := func() (r interface{}, err error) {
		return 10, nil
	}
	task2 := func(v interface{}) (r interface{}, err error) {
		return v.(int) * 2, nil
	}

	f, _ := promise.Start(task1,true).Pipe(task2)
	fmt.Println(f.OnSuccess(func(v interface{}) {
		fmt.Println(v)
	}))
}

