# systemstat
System statistics library for Go 
# Example
	package main

	import (
		"fmt"
		"github.com/permannur/systemstat"
		"time"
	)

	func main() {

		cpu, memory, err := systemstat.FedoraTop()
		if err != nil {
			fmt.Printf("err = %s", err)
			return
		}

		for {
			fmt.Printf("%10v", "cpu total")
			for _, detail := range cpu.DetailList {
				fmt.Printf("%10v", detail.Name)
			}
			fmt.Printf("%10v", "memory")
			fmt.Printf("\n")
			fmt.Printf("%10v", cpu.Percent)
			for _, detail := range cpu.DetailList {
				fmt.Printf("%10v", detail.Percent)
			}
			fmt.Printf("%10v", memory.Percent)
			fmt.Printf("\n")
			time.Sleep(time.Millisecond * 500)
		}

	}
