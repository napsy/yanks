package main

type memEntry struct {
	flat  int
	flatP float64
	sum   int
	cum   int
	cumP  float64
	fn    string
}

type yanks struct {
	db   Db
	apps map[string]*app
}

func main() {
	yanks := &yanks{
		apps: make(map[string]*app),
	}
	yanks.collector()
}
