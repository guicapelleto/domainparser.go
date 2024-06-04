package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"slices"
	"strings"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	White  = "\033[37m"
	
)

var (
	So = runtime.GOOS
)

func clearScreen() {
	var cls *exec.Cmd
	if So == "windows" {
		cls = exec.Command("cmd.exe","/c","cls")
	}else{
		cls = exec.Command("clear")
	}
	cls.Stdout = os.Stdout
	cls.Run()
}

func colorPrint(color string, text string) {
	if So == "windows"{
		fmt.Println(text)
	}else{
		fmt.Println(color, text, Reset)
	}
}

func printBanner() {
	banner := "X19fICBfX19fIF8gIF8gX19fXyBfIF8gIF8gICAgX19fICBfX19fIF9fX18gX19fXyBfX19fIF9fX18gCnwgIFwgfCAgfCB8XC98IHxfX3wgfCB8XCB8ICAgIHxfX10gfF9ffCB8X18vIFtfXyAgfF9fXyB8X18vIAp8X18vIHxfX3wgfCAgfCB8ICB8IHwgfCBcfCAgICB8ICAgIHwgIHwgfCAgXCBfX19dIHxfX18gfCAgXCAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAK"
	decodedBytes, _ := base64.StdEncoding.DecodeString(banner)
	colorPrint(Blue, string(decodedBytes))
	colorPrint(Red, "\nby: guicapelleto\n\n")
}

func getDomain() (domain, site string) {
	if len(os.Args) < 3 {
		fmt.Println("Usage:")
		fmt.Println(os.Args[0], "domain.com domain_site.com\n")
		os.Exit(0)
	}
	arg := os.Args[2]
	if !strings.Contains(arg, ".") {
		log.Fatal("Not a valid site url!")
	}
	domain = os.Args[1]
	if !strings.Contains(domain, ".") {
		log.Fatal("Not a valid domain!")
	}
	if !(strings.Contains(arg, "https://") || strings.Contains(arg, "http://")) {
		site = "https://" + arg
	} else {
		site = arg
	}
	fmt.Println("Checking:", domain, "on", site,"\n\n")
	return
}

func printResult(domain, host string) {
	var result string
	var space int = 50
	if len(result) >= 40{
		space = len(result) + 30
	}
	totalBytes := strings.Repeat(" ", space - len(domain))
	if So == "windows"{
		result = domain + totalBytes + host
	}else{
		result = Green + domain + totalBytes + Purple + host + Reset
	}
	fmt.Println(result)
}

func Grep(texto, padrao string, delimitador ...string) (encontrados []string) {
	var separador string
	if len(delimitador) > 0 {
		separador = delimitador[0]
	} else {
		separador = "\n"
	}
	for _, parte := range strings.Split(strings.ReplaceAll(texto, "\n\r", "\n"), separador) {
		if strings.Contains(parte, padrao) {
			encontrados = append(encontrados, parte)
		}
	}
	return
}

func getIPV4(domain string) (ip string) {
	ips, _ := net.LookupIP(domain)
	ip = ""
	if len(ips) != 0 {
		ip = ips[0].To4().String()
	}
	return
}

func parseDomain(site, domain string) {
	var subdomains []string
	rep := strings.NewReplacer("\"", " ","//", " ","&" , " ", "/", " ", ";", " ", ",", " ")
	resposta, err := http.Get(site)
	if err != nil {
		log.Fatal("Error on loading page source code:", err.Error())
	}
	defer resposta.Body.Close()
	body, _ := io.ReadAll(resposta.Body)
	for _, line := range Grep(string(body), "." + domain) {
		part := rep.Replace(line)
		for _, subdomain := range Grep(part, "." + domain, " "){
			if strings.Contains(subdomain, "\\"){
				continue
			}
			if !slices.Contains(subdomains, subdomain){
				subdomains = append(subdomains, subdomain)
				current_ip := getIPV4(subdomain)
				if len(current_ip) != 0 && current_ip != "<nil>"{
					printResult(subdomain, current_ip)
				
				}
			}
		}
	}
}

func main() {
	clearScreen()
	printBanner()
	domain, site := getDomain()
	parseDomain(site, domain)
}
