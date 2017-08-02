package rrstorage

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/cheggaaa/pb"
	"github.com/songtianyi/rrframework/logs"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type UfileStorage struct {
	PublicKey  string
	PrivateKey string
	BucketName string

	usema chan struct{} // uploading concurrency limit
}

const (
	EXPIRE       = 3600
	SUFFIX       = ".ufile.ucloud.cn"
	MAX_PUT_SIZE = 50 * (1 << 20)
	MAX_GET_SIZE = 50 * (1 << 20)
	PARTIAL_SIZE = 4 * (1 << 20)
)

func CreateUfileStorage(pub, pri, bun string, ucl int) StorageWrapper {
	s := &UfileStorage{
		PublicKey:  pub,
		PrivateKey: pri,
		BucketName: bun,
		usema:      make(chan struct{}, ucl),
	}
	return s
}

func (s *UfileStorage) signheader(method, ctype, bucket, filename string) string {
	data := method + "\n"
	data += "\n"         //Content-MD5 empty
	data += ctype + "\n" //Content-Type
	data += "\n"         //Date empty
	data += "/" + bucket + "/" + filename

	h := hmac.New(sha1.New, []byte(s.PrivateKey))
	h.Write([]byte(data))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

type initResponse struct {
	UploadId string
	BlkSize  int
	Bucket   string
	Key      string
}

func (s *UfileStorage) initiateMultipartUpload(filename string) (*initResponse, error) {
	sign := s.signheader("POST", "application/octet-stream", s.BucketName, filename)

	auth := "UCloud" + " " + s.PublicKey + ":" + sign
	client := &http.Client{}
	url := "http://" + s.BucketName + SUFFIX + "/" + filename + "?uploads"
	req, err := http.NewRequest("POST", url, nil)

	req.Header.Add("Authorization", auth)
	req.Header.Add("Content-Type", "application/octet-stream")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("initiateMultipartUpload failed, %s", string(body))
	}
	var res initResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

type uploadResponse struct {
	PartNumber int
}

func (s *UfileStorage) uploadPart(content []byte, info *initResponse, partNum int) (*uploadResponse, string, error) {
	sign := s.signheader("PUT", "application/octet-stream", info.Bucket, info.Key)

	auth := "UCloud" + " " + s.PublicKey + ":" + sign
	client := &http.Client{}
	url := "http://" + info.Bucket + SUFFIX + "/" + info.Key + "?uploadId=" + info.UploadId + "&partNumber=" + strconv.Itoa(partNum)
	req, err := http.NewRequest("PUT", url, bytes.NewReader(content))

	req.Header.Add("Authorization", auth)
	req.Header.Add("Content-Type", "application/octet-stream")
	req.Header.Add("Content-Length", strconv.Itoa(info.BlkSize))

	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, "", err
	}
	if resp.StatusCode != 200 {
		return nil, "", fmt.Errorf("uploadPart failed, %s", string(body))
	}
	var res uploadResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, "", err
	}
	return &res, resp.Header.Get("ETag"), nil
}

type finishResponse struct {
	Bucket   string
	Key      string
	FileSize int
}

func (s *UfileStorage) finishMultipartUpload(info *initResponse, etags string) (*finishResponse, error) {
	sign := s.signheader("POST", "text/plain", info.Bucket, info.Key)

	auth := "UCloud" + " " + s.PublicKey + ":" + sign
	client := &http.Client{}
	url := "http://" + info.Bucket + SUFFIX + "/" + info.Key + "?uploadId=" + info.UploadId + "&newKey=" + info.Key
	req, err := http.NewRequest("POST", url, strings.NewReader(etags))

	req.Header.Add("Authorization", auth)
	req.Header.Add("Content-Length", strconv.Itoa(len(etags)))
	req.Header.Add("Content-Type", "text/plain")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("finishMultipartUpload failed, %s", string(body))
	}
	var res finishResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *UfileStorage) put(content []byte, filename string) error {
	// sign
	sign := s.signheader("PUT", "application/octet-stream", s.BucketName, filename)
	auth := "UCloud" + " " + s.PublicKey + ":" + sign
	client := &http.Client{}
	url := "http://" + s.BucketName + SUFFIX + "/" + filename
	req, err := http.NewRequest("PUT", url, bytes.NewReader(content))

	req.Header.Add("Authorization", auth)
	req.Header.Add("Content-Type", "application/octet-stream")
	req.Header.Add("Content-Length", strconv.Itoa(len(content)))

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		body, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			return err
		}
		return fmt.Errorf("put file failed, %s", string(body))
	}
	return nil
}

func (s *UfileStorage) Save(content []byte, filename string) error {

	size := len(content)
	if size > MAX_PUT_SIZE {
		// > 50M
		initRes, err := s.initiateMultipartUpload(filename)
		if err != nil {
			return err
		}
		num := size / initRes.BlkSize
		bar := pb.StartNew(num + 1)
		etags := make([]string, 0)
		var (
			wg sync.WaitGroup
			em sync.Mutex
		)
		for i := 0; i < num; i++ {
			s.usema <- struct{}{}
			wg.Add(1)
			go func(j int) {
				defer func() {
					wg.Done()
					<-s.usema
				}()
				part := content[j*initRes.BlkSize : (j+1)*initRes.BlkSize]
				_, etag, err := s.uploadPart(part, initRes, j)
				if err != nil {
					logs.Error(err)
					return
				}
				em.Lock()
				etags = append(etags, etag)
				bar.Increment()
				em.Unlock()
			}(i)
		}
		// TODO capture error
		wg.Wait()
		if num*initRes.BlkSize < size {
			// remaining part
			part := content[num*initRes.BlkSize:]
			_, etag, err := s.uploadPart(part, initRes, num)
			if err != nil {
				return err
			}
			etags = append(etags, etag)
			bar.Increment()
		}
		_, err = s.finishMultipartUpload(initRes, strings.Join(etags, ","))
		if err != nil {
			return err
		}
		bar.Finish()

	} else {
		return s.put(content, filename)
	}
	return nil
}

type fileItem struct {
	BucketName string
	FileName   string
	Hash       string
	MimeType   string
	Size       int
	CreateTime int
	ModifyTime int
}

type fileList struct {
	BucketName string
	BucketId   string
	NextMarker string
	DataSet    []fileItem
}

func (s *UfileStorage) PrefixFileList(prefix string) (*fileList, error) {
	// sign
	sign := s.signheader("GET", "", s.BucketName, "")
	auth := "UCloud" + " " + s.PublicKey + ":" + sign
	client := &http.Client{}
	url := "http://" + s.BucketName + SUFFIX + "/?list&prefix=" + prefix
	req, err := http.NewRequest("GET", url, nil)

	req.Header.Add("Authorization", auth)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("PrefixFileList failed, %s", string(body))
	}
	var res fileList
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *UfileStorage) getFile(filename, brange string) ([]byte, int, error) {
	// sign
	sign := s.signheader("GET", "", s.BucketName, filename)
	auth := "UCloud" + " " + s.PublicKey + ":" + sign
	client := &http.Client{}
	url := "http://" + s.BucketName + SUFFIX + "/" + filename
	req, err := http.NewRequest("GET", url, nil)

	req.Header.Add("Authorization", auth)
	req.Header.Add("Range", brange)

	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, 0, err
	}
	if resp.StatusCode != 206 && resp.StatusCode != 200 {
		return nil, 0, fmt.Errorf("getFile failed, %s", string(body))
	}
	size := 0
	if resp.StatusCode == 200 {
		// complete
		size, _ = strconv.Atoi(resp.Header.Get("Content-Length"))
	} else if resp.StatusCode == 206 {
		// partial
		size, _ = strconv.Atoi(strings.Split(resp.Header.Get("Content-Range"), "/")[1])
	}
	return body, size, nil
}

func (s *UfileStorage) Fetch(filename string) ([]byte, error) {
	b, size, err := s.getFile(filename, "bytes=0-"+strconv.Itoa(MAX_GET_SIZE-1))
	if err != nil {
		return b, err
	}
	lb := len(b)
	if lb == size {
		// downloaded
		return b, nil
	}
	// partial
	size -= lb
	num := size / PARTIAL_SIZE
	bar := pb.StartNew(num + 1)
	// TODO concurrency
	for i := 0; i <= num; i++ {
		brange := "bytes="
		brange += strconv.Itoa(i*PARTIAL_SIZE+lb) + "-"
		brange += strconv.Itoa((i+1)*PARTIAL_SIZE + lb - 1)
		bp, _, err := s.getFile(filename, brange)
		if err != nil {
			logs.Error(err)
			continue
		}
		b = append(b, bp...)
		bar.Increment()
	}
	bar.Finish()
	return b, nil
}
