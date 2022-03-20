package asciiweb

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
)

func OpenBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		fmt.Println("Launching browser...")
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		fmt.Println("Launching browser...")
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		fmt.Println("Launching browser...")
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}
