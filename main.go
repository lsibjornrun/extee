// main.go

/*
extee

Copyright (c) 2013 Bjorn Runaker

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

/* Changes:
1.0 Initial version
*/


package main

import (
	"flag"
	"bufio"
	"fmt"
	"log"
    "github.com/codeskyblue/go-sh"
	"os"
	"io"
	"regexp"
	"strings"
)

var expression string
var command string
var logfile string
var bQuiet bool
var bVerbose bool
var bDryrun bool
var bDeletelogFirst bool


var Usage = func() {
	fmt.Fprintf(os.Stderr, "Usage of %s\n", os.Args[0])
	flag.PrintDefaults()
//	fmt.Fprintf(os.Stderr, "\nConfig file:\nportStart = <first port to be used on localhost>\nportEnd = <last port to use on localhost\n[proxy]\nport = <SOCKS proxy to create on localhost. OPTIONAL (used with -s parameter)>\naddress = \"<IP address to proxy. MANDATORY>\"\n")
//	fmt.Fprintf(os.Stderr, "user=\"<proxy username. MANDATORY>\"\n")
//	fmt.Fprintf(os.Stderr, "ssh=\"<ssh client with full path. Recommended if not using default ssh>\"\n")
}

func readLoop() {
	var re = regexp.MustCompile(expression)
	
	
	r := bufio.NewReader(os.Stdin)
	for {
		str, err := r.ReadString('\n')

		if len(str) > 0 {
			if !bQuiet { fmt.Print(str) }

			if (re.MatchString(str)) {
				if bVerbose { fmt.Print("*match*") }
			} else
			{
				continue
			}
			
			cmd := command
			
			match := re.FindStringSubmatch(str)
  			if match == nil {
				if bVerbose { fmt.Println("no match!!") }
				continue
  			}

			for i := 0; i < re.NumSubexp(); i++ {
				subre := regexp.MustCompile("<"+re.SubexpNames()[i+1]+">")
				cmd = subre.ReplaceAllString(cmd, match[i+1])
								
			}

			if bVerbose || bDryrun { fmt.Println("cmd = " + cmd) }
			
			if (!bDryrun) {
				cmdArray := strings.Split(cmd, " ")
				
				params := make([]interface{}, 0)
				cmd0 := ""
				for index,element := range cmdArray {
					if (index == 0) { 
						cmd0 = element 
					} else
					{ 
						params = append(params, element) 
					}
				} 
		
	
				
				if logfile == "" {
					sh.Command(cmd0, params...).Run()
				} else
				{
					file, err := os.OpenFile(logfile, os.O_RDWR|os.O_APPEND, 0666)
					if err != nil {
						file, err = os.Create(logfile)
						if (err != nil) {
							fmt.Println("Can't write to " + logfile)
							os.Exit(1)
						}
					}
					defer file.Close()
					w := bufio.NewWriter(file)
					
					c1 := sh.Command(cmd0, params...)
					c1.Stdout = w
					c1.Start()
					c1.Wait()
					wc, ok := c1.Stdout.(io.WriteCloser)
					if ok {
						wc.Close()
					}
					
					
					w.Flush()
				}
			}

		}
		if err == io.EOF {
			os.Exit(0)
		}
		if err != nil {
			log.Println("Read Line Error:", err)
			continue
		}


	}
}


func main() {
	flag.StringVar(&expression, "e", "", "regular Expression")
	flag.StringVar(&command, "x", "", "eXecute command")
	flag.StringVar(&logfile, "l", "", "logfile")
	flag.BoolVar(&bQuiet, "q", false, "Quiet. Used in scripts")
	flag.BoolVar(&bVerbose, "v", false, "Verbose. Used when testing regular expression")
	flag.BoolVar(&bDryrun, "n", false, "Dry run. Show but not execute command")
	flag.BoolVar(&bDeletelogFirst, "d", false, "Delete logfile before execution")

	flag.Usage = Usage
	flag.Parse()	
	
	
	if bDeletelogFirst {
		if logfile == "" {
			fmt.Println("Delete log require log file")
			os.Exit(1)
		}
		os.Remove(logfile)
	}
	
	readLoop()

}
