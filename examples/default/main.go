package main

import (
	"fmt"
	cmdtk "github.com/Sinusoid42/commandlinetoolkit"
	"html/template"
	"net/http"
	"os"
	"strings"
)

var templates *template.Template

func init() {
	
	templates = template.Must(template.ParseGlob("./tmpl/*.tmpl"))
}

func Index(w http.ResponseWriter, r *http.Request) {
	
	m := make(map[string]interface{})
	fmt.Println("\n>>>")
	for x, value := range r.Header {
		fmt.Println(x, ": ", value[0])
		if strings.Contains(value[0], "Mobile") || strings.Contains(value[0], "mobile") {
			fmt.Println("MOBILE! : ", value[0])
			templates = template.Must(template.ParseGlob("./tmpl/*.tmpl"))
			templates.ExecuteTemplate(w, "__m.tmpl", m)
			return
		}
	}
	
	//if r.Header.Get("")
	templates = template.Must(template.ParseGlob("./tmpl/*.tmpl"))
	templates.ExecuteTemplate(w, "_.tmpl", m)
	
}

func main() {
	
	fmt.Println("The stdin input")
	
	fmt.Println(os.Args)
	
	cmdLine := cmdtk.NewCommandLine()
	
	cmdLine.ReadJSON("config.json")
	
	//cmdLine.Clear()
	//fmt.Println(cmdLine.Get())
	
	cmdLine.Set(cmdtk.SHELL|cmdtk.HISTORYFILE|cmdtk.HISTORY|cmdtk.PREDICTIONS|cmdtk.SUGGESTIONS, cmdtk.CLI_TRUE)
	
	fs := http.FileServer(http.Dir("./static/"))
	
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	
	http.HandleFunc("/", Index)
	
	server := http.Server{
		Addr: ":64766",
	}
	
	server.ListenAndServe()
	
	//set active shell
	//cmdline.Set(commandlinetoolkit.SHELL, true)
	
	cmdLine.Wait()
	
}
