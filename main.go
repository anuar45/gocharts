package main

func main() {

	grsDB := NewGithubRepoDB()

	gisDB := NewGoImportDB()

	srv := NewWebServer(gisDB, grsDB)

	srv.Run()

	// f, _ := os.Create("data.out")
	// bw := bufio.NewWriter(f)
	// for i, v := range gis {
	// 	s := strconv.Itoa(i) + "  " + v.URL + "  " + strconv.Itoa(v.Count) + "\n"
	// 	bw.WriteString(s)
	// }
	// bw.Flush()
}
