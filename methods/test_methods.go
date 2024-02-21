package methods

import (
	"errors"
	"fmt"
	"fpdns_client/utils"
	"strconv"
	"time"

	"github.com/miekg/dns"
)

func ProbeStart(addr string, sld string) (int, error) {
	time_now := strconv.FormatInt(time.Now().Unix(), 10)
	domain := dns.Fqdn("start-" + time_now + "." + sld)
	res, err := utils.DNSQuery(addr, domain, false, 1, 1, false)
	if err == nil {
		if len(res.Answer) == 1 {
			if a, ok := res.Answer[0].(*dns.A); ok {
				ip_str := a.A.String()
				ip_val, _ := utils.IPv4ToInt(ip_str)
				return int(ip_val), nil
			}
		}
	} else {
		return 0, err
	}
	return 0, errors.New("wrong token response")
}

func InitProbe(addr string, token int, sld string) error {
	time_now := strconv.FormatInt(time.Now().Unix(), 10)
	domain := dns.Fqdn("init-" + time_now + "-" + strconv.Itoa(token) + "." + sld)
	res, err := utils.DNSQuery(addr, domain, true, 1, 1, false)
	if err != nil {
		return err
	} else {
		if len(res.Answer) == 1 {
			if a, ok := res.Answer[0].(*dns.A); ok {
				ip_str := a.A.String()
				if ip_str == "200.200.200.200" {
					return nil
				}
			} else {
				return errors.New("invalid Reply Code")
			}
		}
	}
	return errors.New("init query failure")
}

func ProbeTerminate(addr string, token int, sld string) error {
	time_now := strconv.FormatInt(time.Now().Unix(), 10)
	domain := dns.Fqdn("end-" + time_now + "-" + strconv.Itoa(token) + "." + sld)
	res, err := utils.DNSQuery(addr, domain, false, 1, 1, false)
	if err != nil {
		return err
	} else {
		if len(res.Answer) == 1 {
			if a, ok := res.Answer[0].(*dns.A); ok {
				ip_str := a.A.String()
				if ip_str == "200.200.200.200" {
					return nil
				}
			} else {
				return errors.New("invalid Reply Code")
			}
		}
	}
	return errors.New("terminate query failure")
}

func BehaviorTest(target string, server string, api_port string, sld string) (string, error) {
	target_addr := target + ":53"
	server_addr := server + ":53"
	result_addr := server + ":" + api_port
	token, err := ProbeStart(server_addr, sld)
	if err != nil {
		return "", err
	}
	err = InitProbe(target_addr, token, sld)
	if err != nil {
		return "", err
	}
	err = ProbeTerminate(server_addr, token, sld)
	if err != nil {
		return "", err
	}
	data, err := utils.GetResult(result_addr, token)
	if err != nil {
		return "", err
	}
	json_str, _ := utils.OutputJson(data)
	fmt.Println(json_str)
	return "", nil
}
