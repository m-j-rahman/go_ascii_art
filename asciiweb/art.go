package asciiweb

import (
	"bufio"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Elements struct {
	Input  string
	Output string
}

func MainPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		Err("404 not found", http.StatusNotFound, w)
		return
	}

	tpl, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		Err("500 Internal Server Error", http.StatusInternalServerError, w)
		return
	}
	tpl.Execute(w, Elements{"", ""})
}

func ExportFile(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/export" {
		Err("404 not found", http.StatusNotFound, w)
		return
	}

	err := r.ParseForm()
	if err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	if err != nil {
		Err("500 Internal Server Error", http.StatusInternalServerError, w)
		return
	}

	banner := r.FormValue("banner")
	input := r.FormValue("userText")
	fName := "export.txt"

	if banner == "" || input == "" || input == "\"\"" {
		Err("400 Bad Request", http.StatusNotFound, w)
		return
	}

	input = strings.Replace(string(input), "\r\n", "\\n", -1)

	art := ReadFile(banner)
	output := AsciiString(RemoveQuotes(input), art)

	if strings.HasSuffix(output, " is not in Ascii range") {
		tpl, err := template.ParseFiles("./templates/index.html")
		if err != nil {
			Err("500 Internal Server Error", http.StatusInternalServerError, w)
			return
		}
		tpl.Execute(w, Elements{"hi", output})
	} else {
		file, err := os.Open("export.txt")
		if err != nil {
			Err("lol", http.StatusInternalServerError, w)
			return
		}

		writeFile(fName, output)

		temp := make([]byte, 512)
		file.Read(temp)

		FileContentType := http.DetectContentType(temp)
		FileStat, _ := file.Stat()
		FileSize := strconv.FormatInt(FileStat.Size(), 10)

		FileName := "export"

		w.Header().Set("Content-Type", FileContentType+";"+FileName)
		w.Header().Set("Content-Length", FileSize)
		w.Header().Set("Content-Disposition", "attachment; filename=export.txt")

		file.Seek(0, 0)

		io.Copy(w, file)

		file.Close()
	}
}

func writeFile(outputFile string, ascii string) {
	err := ioutil.WriteFile(outputFile, []byte(ascii), 0777)
	CheckError(err)
}

func AsciiArt(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/ascii-art" {
		Err("404 not found", http.StatusNotFound, w)
		return
	}

	err := r.ParseForm()
	if err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	tpl, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		Err("500 Internal Server Error", http.StatusInternalServerError, w)
		return
	}

	banner := r.FormValue("banner")
	input := r.FormValue("userText")

	if banner == "" || input == "" || input == "\"\"" {
		Err("400 Bad Request", http.StatusNotFound, w)
		return
	}

	input = strings.Replace(string(input), "\r\n", "\\n", -1)

	art := ReadFile(banner)
	output := AsciiString(RemoveQuotes(input), art)

	fmt.Printf("Your input:\n%s\n", input)
	fmt.Println("---------------------------")
	fmt.Printf("%s\n", output)

	tpl.Execute(w, Elements{"hi", output})
}

func ReadFile(banner string) []string {
	fileName := banner + ".txt"
	file, err := os.Open(fileName)
	CheckError(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	sS := make([]string, 0)

	for scanner.Scan() {
		sS = append(sS, scanner.Text())
	}

	CheckError(scanner.Err())

	return sS
}

func CheckError(e error) {
	if e != nil {
		panic(e)
	}
}

func AsciiString(text []string, art []string) (final string) {
	for _, str := range text {
		for i := 0; i < 8; i++ {
			for _, v := range str {
				if GetCharStartingIndex(int(v), i) == -1 {
					errString := string(v) + " is not in Ascii range"
					return errString
				}
				final += art[GetCharStartingIndex(int(v), i)]
			}
			final += "\n"
		}
	}

	return
}

func GetCharStartingIndex(char int, counter int) int {
	if char > 126 {
		return -1
	}
	return (1 + (char-32)*8) + (char - 32) + counter
}

func RemoveQuotes(input string) (newInput []string) {
	slicedInput := []rune(input)
	if slicedInput[0] == 34 {
		if len(slicedInput) != 1 {
			slicedInput = slicedInput[1:]
		} else {
			return []string{""}
		}
	}
	if slicedInput[len(slicedInput)-1] == 34 {
		slicedInput = slicedInput[:len(slicedInput)-1]
	}

	return strings.Split(string(slicedInput), "\\n")
}
