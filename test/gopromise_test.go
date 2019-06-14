package test

import (
	"fmt"
	"testing"
	"time"
)

func TestTask(tt *testing.T) {

	t := time.NewTicker(3 * time.Second)
	defer t.Stop()

	t2 := time.NewTicker(4 * time.Second)
	defer t2.Stop()

	for {

		select {

		case <-t.C:

			fmt.Println(fmt.Sprintf("t.C ,%s",time.Now()))

		case <-t2.C:

			fmt.Println(fmt.Sprintf("t2.C ,%s",time.Now()))

		}

	}

}
