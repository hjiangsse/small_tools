package myarch

import (
	"fmt"
	"log"
)

type MYARCH struct {
}

func (m *MYARCH) Main() error {
	log.Println("This is the Main function of MYARCH")
	return nil
}

func (m *MYARCH) Exit() {
	fmt.Println("This is the Exit function of MYARCH")
}
