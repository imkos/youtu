/*
* File Name:	net.go
* Description:
* Author:	Chapman Ou <ochapman.cn@gmail.com>
* Created:	2015-08-25
 */
package youtu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func (y *Youtu) interfaceURL(ifname string) string {
	return fmt.Sprintf("http://%s/youtu/api/%s", y.host, ifname)
}

func (y *Youtu) interfaceRequest(ifname string, req, rsp interface{}) (err error) {
	url := y.interfaceURL(ifname)
	if y.debug {
		fmt.Printf("req: %#v\n", req)
	}
	b_data, err := json.Marshal(req)
	if err != nil {
		return
	}
	b_body, err := y.get(url, b_data)
	if err != nil {
		return
	}
	err = json.Unmarshal(b_body, rsp)
	if err != nil {
		if y.debug {
			fmt.Fprintf(os.Stderr, "body:%s\n", string(b_body))
		}
		return fmt.Errorf("json.Unmarshal() rsp: %s failed: %s\n", rsp, err)
	}
	return
}

func (y *Youtu) get(addr string, req []byte) (rsp []byte, err error) {
	client := &http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
	httpreq, err := http.NewRequest("POST", addr, bytes.NewReader(req))
	if err != nil {
		return
	}
	auth := y.sign()
	if y.debug {
		fmt.Fprintf(os.Stderr, "Authorization: %s\n", auth)
	}
	httpreq.Header.Add("Authorization", auth)
	httpreq.Header.Add("Content-Type", "text/json")
	httpreq.Header.Add("User-Agent", "")
	httpreq.Header.Add("Accept", "*/*")
	httpreq.Header.Add("Expect", "100-continue")
	resp, err := client.Do(httpreq)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	rsp, err = ioutil.ReadAll(resp.Body)
	return
}
