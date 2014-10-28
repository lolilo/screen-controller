package main

import(
	"flag"
	"fmt"
	"net/http"
	// "os"
	"os/exec"
	"strconv"
	"strings"
	// "path/filepath"
)

const (
	DEFAULT_LOCALHOST_PORT = 4000
	// DEFAULT_LOG_FILE = "/log/slave.log" 

	LINUX_DEFAULT_BROWSER_CMD = "chromium"
	LINUX_DEFAULT_BROWSER_ARGS = "--kiosk"

	OSX_DEFAULT_BROWSER_CMD = "open"
	OSX_DEFAULT_BROWSER_ARGS = "-a 'Google Chrome' --args --kiosk"
)

var port int
var browser_cmd string
var browser_args string
var current_dir string
var err error

func main() {
	setUp()
	http.HandleFunc("/", handleRequest)

	fmt.Printf("Listening on port: %v\n", port)
	fmt.Printf("You can send HTTP POST requests with a 'url' parameter to open it in a browser.\n")
	fmt.Printf("e.g.: curl localhost:%v -X POST -d \"url=http://www.google.com\"\n", port)


	// fmt.Println("CURRENT DIRECTORY IS: ", current_dir)
	
	// start HTTP server with given address and handler
	// handler=nil will default handler to DefaultServeMux
	err := http.ListenAndServe(":" + strconv.Itoa(port), nil)
	if err != nil {
		fmt.Printf("Error starting HTTP server: %v\n", err)
		fmt.Println("Aborting process.")
	}
}

func setUp() {
	OS := getOs()
	if (OS=="unknown") {
		fmt.Printf("Failed to detect operating system.\n")
	} else {
		fmt.Printf("Detected operating system: %v\n", OS)
	}

	switch OS {
	case "Linux":
		browser_cmd = LINUX_DEFAULT_BROWSER_CMD
		browser_args = LINUX_DEFAULT_BROWSER_ARGS
	case "OSX":
		browser_cmd = OSX_DEFAULT_BROWSER_CMD
		browser_args = OSX_DEFAULT_BROWSER_ARGS
	default:
		print("ERROR: Unknown operating system. \n")
	}

	flag.IntVar(&port, "port", DEFAULT_LOCALHOST_PORT, "the port to listen on for commands")
	// can pass flag argument: $ ./slave -port=8080
	// if flag not specified, will set DEFAULT_LOCALHOST_PORT
	flag.Parse()

	// fmt.Println(os.Args)
	// fmt.Println(filepath.Dir(os.Args[0]))
	// current_dir, err = filepath.Abs(filepath.Dir(os.Args[0]))
 //    if err != nil {
 //        fmt.Printf("Error getting the current directory %v\n", err)
 //    }
}

func getOs() string {
	operatingSystemName := exec.Command( "uname", "-a") // display operating system name...why do we need the -a?
	var kernel string
	kernalName, err := operatingSystemName.Output()
	if( err != nil ) {
		fmt.Printf("Error encountered while reading kernal: %v\n", err)
		kernel = "unknown"
	} else {
		kernel = strings.Split( string(kernalName), " " )[0]
	}
	var OS string
	switch kernel {
	case "Linux":
		OS = "Linux"
	case "Darwin":
		OS = "OSX"
	default:
		OS = "unknown"
	}
	return OS
}

func handleRequest(writer http.ResponseWriter, request *http.Request) {
	url := request.PostFormValue("url")
	// fmt.Printf("Executing: %v %v %v\n", browser_cmd, browser_args, url)
	// fmt.Println("CURRENT DIRECTORY IS: ", current_dir)
	// command := "open " + url
	// fmt.Printf("%T", command)

	browser_cmd = "open"
	fmt.Printf("Executing command: %v %v", browser_cmd, url)
	// for some reason the following doesn't work if I pass in "command" ... should be the same string
	err :=exec.Command(browser_cmd, url).Run()
	// err := exec.Command(current_dir+"/../scripts/OS_X_open_browser.sh", url).Run()
	// err := exec.Command(browser_cmd, url).Run()
	if err != nil {
		fmt.Printf("Error opening URL: %v\n", err)
	}
}
