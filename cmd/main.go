package main

import (
	"fmt"
	"math/big"
)

func main() {
	// ts := time.Now().Nanosecond()
	// hrbj5p7EGPdK
	// 11pP6uMXqo
	// 11pP6uMXe3
	ts := 1130063
	// us := uuid.NewString()
	// 2ad1da09
	// us := "2ad1da09-384e-47f1-b5c4-b635400023c87273673s"
	// us := "hello-world"
	var bs big.Int
	dd := bs.SetBytes([]byte(fmt.Sprintf("%d", ts)))
	fmt.Printf("\n\nUUID: %d\tBase62: %s\tNum: %d\n\n", ts, bs.Text(62), dd.Int64())
}
