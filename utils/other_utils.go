package utils

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
)

type ResolverBehavior struct {
	Dnssec        bool  `json:"dnssec"`
	QnameEncode   bool  `json:"0x20"`
	QueryMerge    int   `json:"merge_dup"`
	MaxNsDepth    int   `json:"max_ns_depth"`
	MaxCnameDepth int   `json:"max_cname_depth"`
	RetryLimit    int   `json:"retry_limit"`
	FetchLimit    int   `json:"fetch_limit"`
	TimeoutStart  int64 `json:"timeout_start"`
	TimeoutEnd    int64 `json:"timeout_end"`
}

func IsValidIP(ip string) bool {
	res := net.ParseIP(ip)
	return res != nil
}

func IPv4ToInt(ip string) (uint32, error) {
	if !IsValidIP(ip) {
		return 0, fmt.Errorf("invalid IP address: %s", ip)
	}
	labels := strings.Split(ip, ".")
	var result uint32
	for _, label := range labels {
		label_val, _ := strconv.Atoi(label)
		result = (result << 8) | uint32(label_val)
	}
	return result, nil
}

func GetResult(addr string, token int) (*ResolverBehavior, error) {
	data := new(ResolverBehavior)
	url := "http://" + addr + "/results/" + strconv.Itoa(token)
	response, err := http.Get(url)
	if err != nil {
		return data, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return data, fmt.Errorf("wrong HTTP status code : %v", response.StatusCode)
	}

	err = json.NewDecoder(response.Body).Decode(data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func OutputJson(data interface{}) (string, error) {
	json_byte, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(json_byte), nil
}
