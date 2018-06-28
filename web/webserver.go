package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"sort"
)

var (
	port = flag.Int("port", 8080, "port to listen on")

	servingHost = flag.String("serving_host", 
		"localhost", "address of TF Serving host")
	modelPath = flag.String("model_path",
		"/v1/models/default:classify",
		"path to model in TF Serving")
	protocol = flag.String("protocol", "http",
		"protocol to use for TF Serving host")
	numResults = flag.Int("num_results", 5,
		"num_results of results to show the user")

	tmpl = template.Must(template.ParseFiles("templates/main.html"))
)

type B64Bytes struct {
	B64 string `json:"b64"`
}

type Example struct {
	Image B64Bytes `json:"image"`
}

type Request struct {
	Examples []Example `json:"examples"`
}

func requestFrom(im string) Request {
	return Request{
		[]Example{
			Example{
				B64Bytes{im},
			},
		},
	}

}

type Response struct {
	Results [][][2]interface{} `json:"results"`
	Error string `json:"error"`
}

type Score struct {
	Label string
	Score float64
}

type Result struct {
	scores []Score
}

func resultFrom(resp *Response) (*Result, error) {
	if len(resp.Results) != 1 {
		return nil, fmt.Errorf("got %d results, expected 1", len(resp.Results))
	}

	pairs := resp.Results[0]
	res := Result{make([]Score, len(pairs))}
	
	for i, pair := range pairs {
		label, ok1 := pair[0].(string)
		score, ok2 := pair[1].(float64)

		if !ok1 || !ok2 {
			return nil, fmt.Errorf("got pair of types %T/%T, want string/float64", pair[0], pair[1])
		}
		if label == "" {
			log.Print(pair)
		}
		res.scores[i] = Score{label, score}
	}

	return &res, nil
}

type TemplateParams struct {
	ImageB64 string
	ImageType string
	Scores []Score
}

func main() {
	flag.Parse()

	http.HandleFunc("/healthy", func (w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/classify", imageUploadHandler)

	listenAddress := fmt.Sprintf(":%d", *port)
	log.Printf("Listening at %s...", listenAddress)
	err := http.ListenAndServe(listenAddress, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.Execute(w, TemplateParams{})
}

func imageUploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, file)
	if err != nil {
		log.Fatal(err)
	}

	b64Image := base64.StdEncoding.EncodeToString(buf.Bytes())
	body := requestFrom(b64Image)

	buf.Reset()
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		log.Fatalf("error encoding JSON: %v", err)
	}

	p := TemplateParams{
		ImageB64: b64Image,
		ImageType: header.Header["Content-Type"][0],
	}
	defer tmpl.Execute(w, &p)

	servingAddress := fmt.Sprintf("%s://%s%s", *protocol, *servingHost, *modelPath)
	req, err := http.NewRequest("POST", servingAddress, &buf)
	if err != nil {
		log.Fatalf("error constructing POST request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("error executing POST request: %v", err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		log.Fatalf("error decoding response: %v", err)
	}

	res, err := resultFrom(&respBody)
	if err != nil {
		log.Printf("error parsing response: %v", err)
		log.Panic(respBody.Error)
		return
	}

	// Sort by descending score
	sort.Slice(res.scores, func(i, j int) bool {
		return res.scores[i].Score > res.scores[j].Score
	})

	p.Scores = res.scores[:*numResults]
}