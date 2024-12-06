package main

import (
	"bytes"
	"encoding/json"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
)

type Interface struct {
	Name         string   `json:"name"`
	Index        int      `json:"index"`
	Flag         string   `json:"flag"`
	Mac          string   `json:"mac"`
	Mtu          int      `json:"mtu"`
	IP           []string `json:"ip"`
	Connectivity string   `json:"Connectivity"`
}

func HttpPost(url, header, method string, body []byte) error {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	keyValue := strings.Split(header, ":")
	if len(keyValue) == 2 {
		req.Header.Set(keyValue[0], keyValue[1])
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", runtime.GOOS+"/"+runtime.GOARCH)

	client := &http.Client{Timeout: time.Second * 30}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	logs.Info("restful response status:", resp.Status)

	return nil
}

func queryIpInfo(connectivity string) ([]byte, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	filter := ConfigGet().FilterInterface
	lookup := make([]Interface, 0)

	for _, iface := range ifaces {
		var info Interface

		if filter != "" && !strings.EqualFold(iface.Name, filter) {
			continue
		}

		info.Name = iface.Name
		info.Index = iface.Index
		info.Mac = iface.HardwareAddr.String()
		info.Mtu = iface.MTU
		info.Flag = iface.Flags.String()
		info.Connectivity = connectivity

		address, err := iface.Addrs()
		if err == nil {
			for _, ip := range address {
				info.IP = append(info.IP, ip.String())
			}
		}
		lookup = append(lookup, info)
	}

	body, err := json.MarshalIndent(lookup, "", "\t")
	if err != nil {
		return nil, err
	}

	return body, nil
}

func MonitorInit() error {
	go func() {
		for {
			url := ConfigGet().ConnectivityURL
			if url != "" {
				connectivityIP, err := HttpRequest(url)
				if err != nil {
					logs.Warning("http request failed, %s", err.Error())
				} else {
					logs.Info("http request %s success, body: [%s]", url, string(connectivityIP))
					StatusUpdate(string(connectivityIP))
				}
			}

			ipinfo, err := queryIpInfo(statusConnectivity)
			if err != nil {
				logs.Warning("query ip info failed, %s", err.Error())

			} else {
				outputFile := filepath.Join(ConfigGet().OutputDirectory, "ip_info.json")
				latest, err := os.ReadFile(outputFile)
				if err != nil || !bytes.Equal(ipinfo, latest) {
					err = os.WriteFile(outputFile, ipinfo, 0644)
					if err != nil {
						logs.Warning("write file %s failed, %s", outputFile, err.Error())
					}
				}

				if ConfigGet().RestfulURL != "" {
					err := HttpPost(ConfigGet().RestfulURL, ConfigGet().RestfulHeader, ConfigGet().RestfulMethod, ipinfo)
					if err != nil {
						logs.Warning("http post ip info failed, %s", err.Error())
					}
				}
			}

			time.Sleep(time.Duration(ConfigGet().Interval) * time.Second)
		}
	}()
	return nil
}
