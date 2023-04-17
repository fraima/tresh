package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
)

func main() {
	// указываем подсеть в формате "адрес IP/маска"
	subnet := "10.102.0.0/22"
	url := "http://51.250.78.167/v1/sync"
	ips := 1024
	// парсим подсеть
	ip, ipNet, err := net.ParseCIDR(subnet)
	if err != nil {
		panic(err)
	}

	// получаем первый IP адрес в подсети
	srcIP := ip.Mask(ipNet.Mask)
	// dstIP := ip.Mask(ipNet.Mask)

	for src := 0; src < ips; src++ {
		dstIP := ip.Mask(ipNet.Mask)
		incIP(srcIP)
		fmt.Println(src)
		for dst := 0; dst < ips; dst++ {

			incIP(dstIP)

			s := sha256.New()
			d := sha256.New()

			dstIPStr := dstIP.String()
			srcIPStr := srcIP.String()

			s.Write([]byte(srcIPStr))
			d.Write([]byte(dstIPStr))

			shaSrc := hex.EncodeToString(s.Sum(nil))[:10]
			shaDst := hex.EncodeToString(d.Sum(nil))[:10]

			// fmt.Println(shaSrc, ":", shaDst)
			fmt.Println(srcIPStr, ":", dstIPStr)

			// создаем JSON payload
			payload := map[string]interface{}{
				"sgRules": map[string]interface{}{
					"rules": []map[string]interface{}{
						{
							"portsTo": []map[string]interface{}{
								{"from": 443, "to": 443},
							},
							"sgFrom": map[string]interface{}{
								"name": shaSrc,
							},
							"sgTo": map[string]interface{}{
								"name": shaDst,
							},
							"transport": "TCP",
						},
					},
				},
				"syncOp": "Upsert",
			}

			// преобразуем payload в JSON
			jsonPayload, err := json.Marshal(payload)
			if err != nil {
				panic(err)
			}

			// создаем HTTP клиент и запрос
			client := &http.Client{}
			req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
			if err != nil {
				panic(err)
			}

			// устанавливаем заголовки
			req.Header.Set("Content-Type", "application/json")

			// отправляем запрос
			resp, err := client.Do(req)
			if err != nil {
				panic(err)
			}

			// читаем тело ответа
			defer resp.Body.Close()
			respBody, err := io.ReadAll(resp.Body)
			if err != nil {
				panic(err)
			}

			// выводим ответ
			fmt.Println(resp.Status)
			fmt.Println(string(respBody))
		}
	}
}

func incIP(ip net.IP) {
	for i := len(ip) - 1; i >= 0; i-- {
		ip[i]++
		if ip[i] > 0 {
			break
		}
	}
}
