package sia

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/modeneis/uminer/src/model"
)

// NewClient creates a new SiadClient given a '[stratum+tcp://]host:port' connectionstring
func NewClient(connectionstring, pooluser string) (sc model.Client) {
	if strings.HasPrefix(connectionstring, "stratum+tcp://") {
		sc = &StratumClient{connectionstring: strings.TrimPrefix(connectionstring, "stratum+tcp://"), User: pooluser}
	} else {
		s := SiadClient{}
		s.siadurl = "http://" + connectionstring + "/miner/header"
		sc = &s
	}
	return sc
}

// SiadClient is a simple client to a siad
type SiadClient struct {
	siadurl string
}

func decodeMessage(resp *http.Response) (msg string, err error) {
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return msg, err
	}
	var data struct {
		Message string `json:"message"`
	}
	if err = json.Unmarshal(buf, &data); err == nil {
		msg = data.Message
	}
	return msg, err
}

//Start does nothing
func (sc *SiadClient) Start() {}

//SetDeprecatedJobCall does nothing
func (sc *SiadClient) SetDeprecatedJobCall(call model.DeprecatedJobCall) {}

//GetHeaderForWork fetches new work from the SIA daemon
func (sc *SiadClient) GetHeaderForWork() (target []byte, header []byte, deprecationChannel chan bool, job interface{}, err error) {
	//the deprecationChannel is not used but return a valid channel anyway
	deprecationChannel = make(chan bool)

	client := &http.Client{}

	req, err := http.NewRequest("GET", sc.siadurl, nil)
	if err != nil {
		return target, header, deprecationChannel, job, err
	}

	req.Header.Add("User-Agent", "Sia-Agent")
	resp, err := client.Do(req)
	if err != nil {
		return target, header, deprecationChannel, job, err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.WithError(err).Error("Failed to resp.Body.Close()")
		}
	}()
	switch resp.StatusCode {
	case 200:
	case 400:
		msg, errd := decodeMessage(resp)
		if errd != nil {
			err = fmt.Errorf("Status code %d", resp.StatusCode)
		} else {
			err = fmt.Errorf("Status code %d, message: %s", resp.StatusCode, msg)
		}
		return target, header, deprecationChannel, job, err
	default:
		err = fmt.Errorf("Status code %d", resp.StatusCode)
		return target, header, deprecationChannel, job, err
	}
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return target, header, deprecationChannel, job, err
	}

	if len(buf) < 112 {
		err = fmt.Errorf("Invalid response, only received %d bytes", len(buf))
		return target, header, deprecationChannel, job, err
	}

	target = buf[:32]
	header = buf[32:112]

	return target, header, deprecationChannel, job, err
}

//SubmitHeader reports a solved header to the SIA daemon
func (sc *SiadClient) SubmitHeader(header []byte, job interface{}) (err error) {
	req, err := http.NewRequest("POST", sc.siadurl, bytes.NewReader(header))
	if err != nil {
		return err
	}

	req.Header.Add("User-Agent", "Sia-Agent")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	switch resp.StatusCode {
	case 204:
	default:
		msg, errd := decodeMessage(resp)
		if errd != nil {
			err = fmt.Errorf("Status code %d", resp.StatusCode)
		} else {
			err = fmt.Errorf("%s", msg)
		}
		return err
	}
	return err
}
