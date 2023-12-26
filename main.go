package main

import (
	"fmt"
	"net/http"
	"os"
)

const base_path = "/home/pedro/ubuntu-nas/media"

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/example", exampleHandler)

	fmt.Println("Levantando server")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
func homeHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Fprint(res, "Foi a home")
}

func exampleHandler(w http.ResponseWriter, r *http.Request) {
	testStr := base_path + "/Jovem-Nerd/NerdOffice/S01/[00001] - Piloto [XkG1F-wJ4Is].webm"

	fmt.Print("reading: " + testStr)

	fileBytes, err := os.ReadFile(base_path + "/movies/Top Gun Maverick (2022) [IMAX] [REPACK] [1080p] [WEBRip] [5.1] [YTS.MX]/Top.Gun.Maverick.2022.IMAX.REPACK.1080p.WEBRip.x264.AAC5.1-[YTS.MX].mp4")

	if err != nil {
		fmt.Printf("Invalid path received on readfile, %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(fileBytes)

}
