package app

import (
	"bufio"
	"fmt"
	"github.com/Vakaram/testovoeMahazineSklad/internal/storage"
	"os"
)

type app struct {
	Store *storage.Store
}

func Start() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter text: ")
		text, _ := reader.ReadString('\n')
		// Обработка введенного текста
		fmt.Printf("Вот введеный вами текст : %s", text)
	}
}
