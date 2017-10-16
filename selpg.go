package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func main() {
	var s, e, l int
	var d string
	flag.IntVar(&s, "s", 1, "start page, default 1")
	flag.IntVar(&e, "e", 1, "end page, deafault 1")
	flag.IntVar(&l, "l", 72, "how many lines in single page, deafault 72")
	flag.StringVar(&d, "d", "", "pipe or not, default nil")
	flag.Parse()
	begin := (s - 1) * l
	end := e * l
	count := 1
	cmdinfo := ""
	//without filename ,read info from standard input
	if len(flag.Args()) == 0 {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			if count > begin && count <= end {
				if d == "" {
					fmt.Println(line)
				} else {
					cmdinfo = cmdinfo + line + "\n"
				}
			}
		}
		//with filename, read info from file
	} else {
		fileName := flag.Args()[0]
		f, err := os.Open(fileName)
		if err != nil {
			f.Close()
			os.Exit(1)
		}
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			if count > begin && count <= end {
				if d == "" {
					fmt.Println(line)
				} else {
					cmdinfo = cmdinfo + line + "\n"
				}
			}
		}
	}
	//with -d option ,it will pipe info to command "lp -d***"
	if d != "" {
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
}
