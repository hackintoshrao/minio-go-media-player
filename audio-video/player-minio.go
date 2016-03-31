package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/minio/minio-go"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type Entry struct {
	Name  string // name of the object
	IsDir bool
	Mode  os.FileMode
}

const (
	filePrefix = "/f/"
)

var (
	addr = flag.String("http", ":8080", "http listen address")
	root = flag.String("root", "/Users/hackintoshrao/Music/", "music root")
)

func main() {
	flag.Parse()
	http.HandleFunc("/jplayer/", func(w http.ResponseWriter, r *http.Request) {
		log.Print("serve static : " + r.URL.Path)
		http.ServeFile(w, r, r.URL.Path[1:])
	})
	http.HandleFunc("/songs/", func(w http.ResponseWriter, r *http.Request) {
		log.Print("serve static : " + r.URL.Path)
		http.ServeFile(w, r, r.URL.Path[1:])
	})
	http.HandleFunc("/", Index)
	http.HandleFunc("/list", ListObjects)
	http.HandleFunc(filePrefix, File)
	http.ListenAndServe(*addr, nil)
}

func Index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./index.html")
	log.Print("index called")
}
func ListObjects(w http.ResponseWriter, r *http.Request) {

	// Note: YOUR-ACCESSKEYID, YOUR-SECRETACCESSKEY, my-bucketname and my-prefixname
	// are dummy values, please replace them with original values.

	// Requests are always secure (HTTPS) by default. Set insecure=true to enable insecure (HTTP) access.
	// This boolean value is the last argument for New().

	// New returns an Amazon S3 compatible client object. API copatibality (v2 or v4) is automatically
	// determined based on the Endpoint value.
	s3Client, err := minio.New("s3.amazonaws.com", "AKIAIBVF3NRPLX5ZEYZA", "jAX9g55X4p+FnGW5iyBrW1p9/+D8jZ6BrMgaWMFY", false)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create a done channel to control 'ListObjects' go routine.
	doneCh := make(chan struct{})

	// Indicate to our routine to exit cleanly upon return.
	defer close(doneCh)
	type responseJson struct {
		Key  string
		User string
	}
	var objectInfos []responseJson
	// List all objects from a bucket-name with a matching prefix.
	for objectInfo := range s3Client.ListObjects("karthicminio", "", true, doneCh) {
		if objectInfo.Err != nil {
			fmt.Println(objectInfo.Err)
			return
		}
		res := responseJson{Key: objectInfo.Key, User: "Karthic"}
		objectInfos = append(objectInfos, res)
		fmt.Printf("%+v", objectInfo)

	}
	json.NewEncoder(w).Encode(objectInfos)
	//	fmt.Fprintf(w, "%v", objectInfos)
}

func File(w http.ResponseWriter, r *http.Request) {

	log.Print("Path: " + r.URL.Path)
	fn := filepath.Join(*root, r.URL.Path[len(filePrefix):])
	fi, err := os.Stat(fn)
	log.Print("File called: ", fn)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if fi.IsDir() {
		serveDirectory(fn, w, r)
		return
	}
	http.ServeFile(w, r, fn)
}

func serveDirectory(fn string, w http.ResponseWriter,
	r *http.Request) {
	defer func() {
		if err, ok := recover().(error); ok {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}()
	d, err := os.Open(fn)
	if err != nil {
		panic(err)
	}
	log.Print("serverDirectory called: ", fn)

	files, err := d.Readdir(-1)
	if err != nil {
		panic(err)
	}

	// Json Encode isn't working with the FileInfo interface,
	// therefore populate an Array of Entry and add the Name method
	entries := make([]Entry, len(files), len(files))

	for k := range files {
		//log.Print(files[k].Name())
		entries[k].Name = files[k].Name()
		entries[k].IsDir = files[k].IsDir()
		entries[k].Mode = files[k].Mode()
	}

	j := json.NewEncoder(w)

	if err := j.Encode(&entries); err != nil {
		panic(err)
	}
}
func getPresigndURL(w http.ResponseWriter, r *http.Request) {
	s3Client, err := minio.New("s3.amazonaws.com", "AKIAIBVF3NRPLX5ZEYZA", "jAX9g55X4p+FnGW5iyBrW1p9/+D8jZ6BrMgaWMFY", false)
	if err != nil {
		log.Fatalln(err)
	}

	// Set request parameters
	reqParams := make(url.Values)
	//reqParams.Set("response-content-disposition", "attachment; filename=\"your-filename.txt\"")

	// Gernerate presigned get object url.
	presignedURL, err := s3Client.PresignedGetObject("karthicminio", object, time.Duration(1000)*time.Second, reqParams)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(presignedURL)

}
