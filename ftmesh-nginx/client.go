package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net/http"

	// "strings"
	"time"

	"github.com/google/uuid"
)

func generateDummyData(size uint32) []byte {
	data := make([]byte, size)
	for i := uint32(0); i < size; i++ {
		data[i] = 'A' + byte(i%26)
	}
	return data
}

func generateDummyString(size uint32) string {
	data := make([]byte, size)
	for i := uint32(0); i < size; i++ {
		data[i] = 'A' + byte(i%26)
	}
	return string(data)
}

func generateStateObj(size uint32) []byte {
	mtString := "10730\n0\n6000\n/svc-a/recover\nd9175a23-65b0-4e78-9802-1d29d0a019d6\n" +
		"svc-a-107\nsvc-a-110\nfib\n10.214.96.110\n10729\n"
	mtLen := uint32(len(mtString))

	var buf [4]byte
	fmt.Printf("mtLen :%d\n", mtLen)
	binary.LittleEndian.PutUint32(buf[:], mtLen)
	slice := append(buf[:], []byte(mtString)...)

	if mtLen+4 < size {
		dummy := generateDummyData(size - mtLen - 4)
		slice = append(slice, dummy...)
	}

	return slice
}

func readAllBody(closer io.ReadCloser) ([]byte, error) {
	defer func(closer io.ReadCloser) {
		err := closer.Close()
		if err != nil {
			panic(err)
		}
	}(closer)
	body, err := io.ReadAll(closer)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func RequestGenerator(dummy bool) *http.Request {
	plot := generateDummyString(10240)
	bts := []byte(plot)
	oriLen := len(bts)
	fmt.Printf("orilen: %d\n", oriLen)

	if dummy {
		bts = append(bts, generateStateObj(10245)...)
	}

	buf := bytes.NewBuffer(bts)
	req, _ := http.NewRequest("GET", "http://127.0.0.1:10729/db", buf)

	req.Header.Set("User-Agent", "Go-http-client/1.1")
	req.Header.Set("Accept", "text/plain")
	req.Header.Set("Connection", "close")

	if dummy {
		req.Header.Add("x-ftmesh-length", fmt.Sprintf("%d", oriLen))
		req.Header.Add("x-ftmesh-mode", "0")
	}

	return req
}

func LocalRequestGenerator(timeString string, dummy bool) *http.Request {
	plot := generateDummyString(9955)
	bts := []byte("{\"ExpressionAttributeNames\":{\"#0\":\"info\",\"#1\":\"rating\",\"#2\":\"plot\"},\"ExpressionAttributeValues\":{\":0\":{\"S\":\"5.0\"},\":1\":{\"S\":\"" + plot + "\"}},\"Key\":{\"title\":{\"S\":\"FooBar\"},\"year\":{\"N\":\"2024\"}},\"ReturnValues\":\"UPDATED_NEW\",\"TableName\":\"movieTable\",\"UpdateExpression\":\"SET #0.#1 = :0, #0.#2 = :1\\n\"}\n")
	oriLen := len(bts)
	//fmt.Printf("orilen: %d\n", oriLen)

	if dummy {
		bts = append(bts, generateStateObj(1024)...)
	}

	buf := bytes.NewBuffer(bts)
	req, _ := http.NewRequest("POST", "http://127.0.0.1:8000", buf)

	req.Header.Add("Amz-Sdk-Invocation-Id", uuid.NewString())
	req.Header.Add("Amz-Sdk-Request", "attempt=1; max=3")
	req.Header.Set("Authorization", "AWS4-HMAC-SHA256 Credential=AaBbCcDc1234/20240407/us_west_2/dynamodb/aws4_request, SignedHeaders=accept-encoding;amz-sdk-invocation-id;amz-sdk-request;content-length;content-type;host;x-amz-date;x-amz-target, Signature=0bf62c65363bbb400cda94abba6e60194d4a9e356f20581fd9e799184a1d43ac")
	req.Header.Set("Content-Type", "application/x-amz-json-1.0")
	req.Header.Add("X-Amz-Date", "20240502T"+timeString+"Z")
	req.Header.Add("X-Amz-Target", "DynamoDB_20120810.UpdateItem")
	req.Header.Set("Accept-Encoding", "identity")

	if dummy {
		req.Header.Add("x-ftmesh-length", fmt.Sprintf("%d", oriLen))
		req.Header.Add("x-ftmesh-mode", "0")
	}

	return req
}

func main() {
	var sum float64 = 0.0
	for i := 0; i < 50; i++ {
		req := RequestGenerator(false)
		localreq := LocalRequestGenerator("20590"+fmt.Sprintf("%d", i), false)
		req.Header.Add("x-ftmesh-cluster", "cluster_0")
		localreq.Header.Add("x-ftmesh-cluster", "cluster_0")
		start := time.Now()

		client := &http.Client{}
		fmt.Printf("Before do req ts: %v\n", time.Now().UnixMicro())
		//_, err := client.Do(req)
		resp, err := client.Do(localreq)
		if err != nil {
			panic(err)
		}

		resp, err = client.Do(req)
		if err != nil {
			panic(err)
		}

		duration := time.Since(start)
		fmt.Println(duration.Microseconds())
		sum += float64(duration.Microseconds())
		body, err := readAllBody(resp.Body)
		bodyString := string(body)
		fmt.Println(bodyString)
		time.Sleep(1000 * time.Millisecond)
	}
	fmt.Printf("Average time: %.2f\n", sum/50)
}
