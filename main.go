package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

func decode(s string) ([]byte, error) {
	compressed, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("base64: %v", err)
	}

	decompressor, err := gzip.NewReader(bytes.NewReader(compressed))
	if err != nil {
		return nil, fmt.Errorf("gzip: %v", err)
	}

	data, err := io.ReadAll(decompressor)
	if err != nil {
		return nil, fmt.Errorf("io: %v", err)
	}

	return data, nil
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v %v\n", r.Method, r.URL.RequestURI())
	encoded := r.FormValue("data")
	data, err := decode(encoded)
	if err != nil {
		msg := fmt.Sprintf("decode: %v\n", err)
		log.Print(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	_, err = w.Write(data)
	if err != nil {
		msg := fmt.Sprintf("Write: %v\n", err)
		log.Print(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
}

func serve() {
	http.HandleFunc("/", handleRoot)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func encode(r io.Reader) (string, error) {
	var buffer bytes.Buffer
	writer := gzip.NewWriter(&buffer)
	_, err := io.Copy(writer, r)
	if err != nil {
		return "", fmt.Errorf("io: %v", err)
	}

	writer.Close()
	encoded := base64.StdEncoding.EncodeToString(buffer.Bytes())

	v := url.Values{}
	v.Add("data", encoded)

	return v.Encode(), nil
}

func encodeStdin() {
	encoded, err := encode(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(encoded)
}

func encodeFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	encoded, err := encode(file)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(encoded)
}

func main() {
	argc := len(os.Args)
	switch {
	case argc == 1:
		serve()
	case argc == 2 && os.Args[1] == "url":
		encodeStdin()
	case argc == 3 && os.Args[1] == "url":
		encodeFile(os.Args[2])
	default:
		log.Panicf("Invalid args: %v", os.Args[1:])
	}
}
