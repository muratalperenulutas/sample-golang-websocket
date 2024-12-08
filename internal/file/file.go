package file

import (
	"fmt"
	"os"
)

func SaveTextFile(data string) {

	file, err := os.OpenFile("iot.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(data)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
	fmt.Println("Text successfully saved")
}
