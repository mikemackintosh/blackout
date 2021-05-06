package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/jszwec/csvutil"
	"github.com/mikemackintosh/blackout"
	"github.com/mitchellh/go-homedir"
)

var (
	flagProfile string
	flagDomain  string
	flagOutfile string
)

const cookieProfilePath = "%s/Library/Application Support/Google/Chrome/%s/Cookies"

func init() {
	flag.StringVar(&flagProfile, "p", "Default", "Chrome Profile")
	flag.StringVar(&flagDomain, "d", "", "Domain to filter on")
	flag.StringVar(&flagOutfile, "o", "", "Outfile to write to")
}

func main() {
	flag.Parse()

	// Expand homedir for chrome config
	home, err := homedir.Dir()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	// Build path
	cookieDb := fmt.Sprintf(cookieProfilePath,
		home,
		flagProfile,
	)

	// Get cookies
	cookies, err := blackout.GetCookies(cookieDb)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	// If filtering enabled, filter
	if len(flagDomain) > 0 {
		fmt.Println("filtering on domains")
		var outCookies blackout.Cookies
		for _, cookie := range cookies {
			if strings.Contains(cookie.Domain, flagDomain) {
				fmt.Printf("%s %s add\n", cookie.Domain, flagDomain)
				outCookies = append(outCookies, cookie)
			}
		}
		cookies = outCookies
	}

	// Build CSV
	var b []byte
	if len(flagOutfile) > 0 {
		if strings.Contains(flagOutfile, "json") {
			b, err = json.Marshal(cookies)
			if err != nil {
				fmt.Println("error:", err)
			}
		} else {
			b, err = csvutil.Marshal(cookies)
			if err != nil {
				fmt.Println("error:", err)
			}
		}

		ioutil.WriteFile(flagOutfile, b, 0644)
		os.Exit(0)
	}

	os.Stdout.Write(b)
}
