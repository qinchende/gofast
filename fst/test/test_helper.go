package test

import "fmt"

func Cover(ver int) {
	switch ver {
	case 1:
		fmt.Println("GO 1")
	case 2:
		fmt.Println("GO 2")
	case 3:
		fmt.Println("GO 3")
	default:
		fmt.Println("GO def")
	}
}

func Cover2(ver int) {
	switch ver {
	case 1:
		fmt.Println("GoFast 1")
	default:
		fmt.Println("GoFast def")
	}
}
