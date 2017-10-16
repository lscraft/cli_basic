package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
)

var s, e, l int
var d string
var f bool

func arginit() {
	flag.IntVar(&s, "s", 1, "start page, default 1")
	flag.IntVar(&e, "e", 1, "end page, deafault 1")
	flag.IntVar(&l, "l", 72, "how many lines in single page, deafault 72")
	flag.StringVar(&d, "d", "", "pipe or not, default nil")
	flag.BoolVar(&f, "f", false, "search for the page break word")
	flag.Parse()
}

func judge() {
	if f == true && l != 72 {
		io.WriteString(os.Stderr, "-f and -l should not appear together")
		os.Exit(1)
	}
	if s < 1 || e < 1 || l < 1 {
		io.WriteString(os.Stderr, "input integer please")
		os.Exit(1)
	}
}
func pipeToLp(cmdinfo string) {
	cmdArg := "-d" + d
	cmd := exec.Command("lp", cmdArg)
	cmdin, cmderr := cmd.StdinPipe()
	if cmderr != nil {
		cmdin.Close()
		os.Exit(1)
	}
	cmd.Start()
	io.WriteString(cmdin, cmdinfo)
	cmdin.Close()
	cmd.Wait()
}

// function to deal with -f
func splitWithF(fpointer io.Reader) string {
	cmdinfo := ""
	begin := s - 1
	end := e
	count := 1
	reader := bufio.NewReader(fpointer)
	for {
		by, err := reader.ReadByte()
		if err == io.EOF {
			break
		}
		if string(by) == "\f" {
			count = count + 1
		} else if count > begin && count <= end {
			if d == "" {
				fmt.Printf("%s", string(by))
			} else {
				cmdinfo = cmdinfo + string(by)
			}
		}
	}
	if d == "" {
		fmt.Printf("\n")
	} else {
		cmdinfo = cmdinfo + "\n"
	}
	return cmdinfo
}

//normal scan function
func splitWithoutF(fpointer io.Reader) string {
	begin := (s - 1) * l
	end := e * l
	count := 1
	cmdinfo := ""
	scanner := bufio.NewScanner(fpointer)
	for scanner.Scan() {
		line := scanner.Text()
		if count > begin && count <= end {
			if d == "" {
				fmt.Println(line)
			} else {
				cmdinfo = cmdinfo + line + "\n"
			}
		}
		count = count + 1
	}
	return cmdinfo
}
func processInfo(fp io.Reader) string {

	if f == true {
		return splitWithF(fp)
	} else {
		return splitWithoutF(fp)
	}
}
func getFilePointer() io.Reader {
	var fp io.Reader
	if len(flag.Args()) == 0 {
		fp = os.Stdin
	} else {
		file, err := os.Open(flag.Args()[0])
		fp = bufio.NewReader(file)
		if err != nil {
			file.Close()
			os.Exit(1)
		}
	}
	return fp
}
func main() {
	arginit()

	judge()

	fp := getFilePointer()

	cmdinfo := processInfo(fp)

	//with -d option ,it will pipe info to command "lp -d***"
	if d != "" {
		pipeToLp(cmdinfo)
	}
}
