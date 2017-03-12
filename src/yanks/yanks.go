package yanks

import (
	"bytes"
	"encoding/binary"
	"flag"
	"io"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

var (
	apiKey = flag.String("yank-key", "", "yanks API key")
	name   = flag.String("yank-name", "", "yanks name")
)

func tick() {
	b := bytes.Buffer{}
	runtime.GC() // get up-to-date statistics
	if err := pprof.WriteHeapProfile(&b); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}
	if err := sendData(b.Bytes()); err != nil {
		log.Printf("Error: %v\n", err)
	}
}

func sendData(p []byte) error {

	// connect to this socket
	conn, err := net.Dial("tcp", "127.0.0.1:7000")
	if err != nil {
		return err
	}
	defer conn.Close()
	size := make([]byte, 4)
	binary.LittleEndian.PutUint32(size, uint32(len(p)))
	if _, err = conn.Write(size); err != nil {
		return err
	}
	if _, err = conn.Write(p); err != nil {
		return err
	}
	log.Printf("DONE")
	return nil
}

func ticker(c <-chan time.Time) {
	for _ = range c {
		log.Printf("Sending stats ...")
		tick()
	}
}

func getProfile(url string) {
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	} else {
		defer response.Body.Close()
		_, err := io.Copy(os.Stdout, response.Body)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func init() {

	flag.Parse()
	c := time.Tick(10 * time.Second)
	go ticker(c)
	log.Printf("yanks: using %s for %s ...", *apiKey, *name)
	log.Fatal(http.ListenAndServe(":6000", nil))
}
