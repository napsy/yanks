package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

func parse(buf io.Reader) ([]memEntry, error) {
	b := bufio.NewReader(buf)
	entries := []memEntry{}
	n := 0
	for ; ; n++ {
		l, _, err := b.ReadLine()
		if err != nil {
			break
		}
		if n < 2 {
			continue
		}
		line := strings.Trim(string(l), " \t")
		cols := strings.Fields(line)
		if len(cols) != 6 {
			continue
		}
		for i := 0; i < len(cols); i++ {
			p := strings.Index(cols[i], "%")
			if p > -1 {
				cols[i] = cols[i][:p]
			}
		}
		e := memEntry{}
		if e.flat, err = strconv.Atoi(cols[0]); err != nil {
			return nil, fmt.Errorf("parsing 'flat': %v", err)
		}
		if e.flatP, err = strconv.ParseFloat(cols[1], 64); err != nil {
			return nil, fmt.Errorf("parsing 'flatP' (%v): %v", cols[1], err)
		}
		if e.sum, err = strconv.Atoi(cols[2]); err != nil {
			return nil, fmt.Errorf("parsing 'sum': %v", err)
		}
		if e.cum, err = strconv.Atoi(cols[3]); err != nil {
			return nil, fmt.Errorf("parsing 'cum': %v", err)
		}
		if e.cumP, err = strconv.ParseFloat(cols[4], 64); err != nil {
			return nil, fmt.Errorf("parsing 'cumP': %v", err)
		}
		e.fn = cols[5]
		entries = append(entries, e)
	}
	return entries, nil
}

func (yanks *yanks) handleRequest(c net.Conn) {
	l := make([]byte, 4)
	c.Read(l)
	size := binary.LittleEndian.Uint32(l)
	b := make([]byte, size)
	c.Read(b)
	if err := syscall.Unlink("/tmp/bla"); err != nil {
	}
	if err := syscall.Mkfifo("/tmp/bla", 0644); err != nil {
		fmt.Printf("mkfifo err: %v\n", err)
		return
	}
	cmd := exec.Command("go", "tool", "pprof", "-top", "/tmp/bla")
	out, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("out err: %v\n", err)
		return
	}
	err = cmd.Start()
	if err != nil {
		fmt.Printf("ERROR: %v\n")
		return
	}
	go func() {
		//f.Write(b)
		if err := ioutil.WriteFile("/tmp/bla", b, 0644); err != nil {
			fmt.Printf("write rr: %v\n", err)
		}
		syscall.Unlink("/tmp/bla")
	}()

	entries, err := parse(out)
	if err != nil {
		fmt.Printf("Parse error: %v\n", err)
	}
	cmd.Wait()
	fmt.Printf("%+v\n", entries)
}

func (yanks *yanks) collector() {
	// Listen for incoming connections.
	l, err := net.Listen("tcp", ":7000")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go yanks.handleRequest(conn)
	}
}
