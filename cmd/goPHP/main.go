package main

import (
	"GoPHP/cmd/goPHP/common"
	"GoPHP/cmd/goPHP/config"
	"GoPHP/cmd/goPHP/interpreter"
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path"
	"time"
)

var serverAddr string
var documentRoot string

func main() {
	file := flag.String("f", "", "Parse and execute <file>.")
	// Web server
	addr := flag.String("S", "", "Run with built-in web server. <addr>:<port>")
	docRoot := flag.String("t", "", "Specify document root <docroot> for built-in web server.")

	flag.Parse()

	// Serve with built-in web server
	if *addr != "" {
		serverAddr = *addr
		documentRoot = *docRoot
		webServer()
		os.Exit(0)
	}

	// Parse given file
	if *file != "" {
		absFilePath, err := common.GetAbsPath(*file)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		processFile(absFilePath)
	}

	// Read stdin or wait for it
	processStdin()
}

// ------------------- MARK: stdin -------------------

func processStdin() {
	content := ""
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if content != "" || scanner.Text() == "" {
			content += "\n"
		}
		content += scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error:", err)
	}

	output, exitCode := processContent(string(content), "main.php")
	fmt.Print(output)
	os.Exit(exitCode)
}

// ------------------- MARK: given file -------------------

func processFile(filename string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	output, exitCode := processContent(string(content), filename)
	fmt.Print(output)
	if exitCode == 500 {
		exitCode = 1
	}
	os.Exit(exitCode)
}

// ------------------- MARK: web server -------------------

func webServer() {
	if documentRoot == "" {
		documentRoot, _ = os.Getwd()
	} else {
		absPath, err := common.GetAbsPath(documentRoot)
		if err != nil {
			fmt.Println("Error: Could not find directory " + documentRoot)
			os.Exit(1)
		}
		documentRoot = absPath
	}

	fmt.Printf("[%s] GoPHP %s Development Server (%s) started\n",
		time.Now().Format("Mon Jan 02 15:04:05 2006"),
		config.Version,
		serverAddr,
	)
	fmt.Println("Document root is " + documentRoot)
	fmt.Println("Press Ctrl-C to quit")

	http.HandleFunc("/", requestHandler)
	if err := http.ListenAndServe(serverAddr, nil); err != nil {
		fmt.Println("Failed to start web server")
		fmt.Println(err)
	}
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	absFilePath := path.Join(documentRoot, r.URL.Path)
	_, err := os.Stat(absFilePath)
	if err != nil {
		fmt.Println("404", absFilePath)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, http.StatusText(http.StatusNotFound))
		return
	}

	content, err := os.ReadFile(absFilePath)
	if err != nil {
		fmt.Println("Error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	output, exitCode := processContent(string(content), absFilePath)
	if exitCode == 500 {
		w.WriteHeader(http.StatusInternalServerError)
	}
	fmt.Fprint(w, output)
}

// ------------------- MARK: core logic -------------------

func processContent(content string, filename string) (output string, exitCode int) {
	interpreter := interpreter.NewInterpreter(interpreter.NewDevConfig(), &interpreter.Request{}, filename)
	result, err := interpreter.Process(content)
	if err != nil {
		result += interpreter.ErrorToString(err)
		return result, 500
	}
	return result, interpreter.GetExitCode()
}
