package controllers

import (
	"fmt"
	"guide/helper"
	"net/http"
)

func TestDisplay(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Key query : %s\n",r.FormValue("query"))
	fmt.Printf("Key select : %s\n",r.FormValue("select"))
	fmt.Printf("Key check : %s\n",r.FormValue("check"))
	fmt.Printf("Key check_multip : %s\n",r.Form["check_multip"])
	fmt.Printf("Key radio : %s\n\n",r.FormValue("radio"))
	fmt.Println(r.Form["check_multip"])
	helper.RenderTemplate(w,r,"exemple_formulaire",nil)
}