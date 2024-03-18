package main

import "fmt"

func printNumberTriangle(n int) {
    for i := 1; i <= n; i++ {
        for j := 1; j <= i; j++ {
            fmt.Printf("%s ", "*")
        }
        fmt.Println()
    }
    for i := n - 1; i >= 1; i-- {
        for j := 1; j <= i; j++ {
            fmt.Printf("%s ", "*")
        }
        fmt.Println()
    }
}

func main() {
    var x int
	fmt.Scanf("%d", &x)
    printNumberTriangle(x)
}