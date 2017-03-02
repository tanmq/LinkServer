package main

import (
	"linkServer/config"
	"fmt"
)

func main() {
	cf := config.NewConfiger()
	fmt.Println(cf.Int("serverport"))
}


