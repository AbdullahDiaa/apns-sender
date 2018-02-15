package apns

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"golang.org/x/net/http2"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"sync"
	"time"
)

const (
	//APNS Connection endpoints
	//Development server: api.development.push.apple.com:443
	//Production server: api.push.apple.com:443
	//You can alternatively use port 2197 when communicating with APNs. You might do this, for example, to allow APNs traffic through your firewall but to block other HTTPS traffic.
	//https://developer.apple.com/library/content/documentation/NetworkingInternet/Conceptual/RemoteNotificationsPG/CommunicatingwithAPNs.html

	APNSDevelopmentServer = "https://api.development.push.apple.com"
	APNSProductionServer  = "https://api.push.apple.com"
)

type Notification struct {
	Topic      string
	Expiration time.Time
	Payload    interface{}
}

type Server struct {
	host   string
	client *http.Client
}

var DialTLS = func(network, addr string, cfg *tls.Config) (net.Conn, error) {
	dialer := &net.Dialer{
		//Connection timeout with certificate
		Timeout: 20 * time.Second,
		//Duration of activity after not using the connection
		KeepAlive: 60 * time.Second,
	}
	return tls.DialWithDialer(dialer, network, addr, cfg)
}

func CreateServer(host, certificatestring string) *Server {
	s := new(Server)
	s.host = host

	cert, _ := tls.X509KeyPair([]byte(certificatestring), []byte(certificatestring))

	tlsConfig := new(tls.Config)
	tlsConfig.Certificates = []tls.Certificate{cert}

	if len(cert.Certificate) > 0 {
		tlsConfig.BuildNameToCertificate()
		fmt.Println(tlsConfig.NameToCertificate)
	}

	s.client = &http.Client{
		//HTTP2 transport layer
		Transport: &http2.Transport{
			TLSClientConfig: tlsConfig,
			DialTLS:         DialTLS,
		},
	}
	return s
}

func sendNotification(n Notification, s *Server, wg *sync.WaitGroup, deviceToken string) {
	defer wg.Done()

	payload, err := json.Marshal(n)
	if err != nil {
		fmt.Println(err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/3/device/%s", s.host, deviceToken), bytes.NewBuffer(payload))

	if err != nil {
		fmt.Println(err)
	}

	response, err := s.client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	if response != nil {
		defer response.Body.Close() // close response body to reuse connection
		fmt.Println(response)
		io.Copy(ioutil.Discard, response.Body) // read entire body to reuse connection
	}
}

func Push(n Notification, s *Server, d []string) error {
	fmt.Println(n)
	var wg sync.WaitGroup
	for i := 0; i < len(d); i++ {
		wg.Add(1)
		go sendNotification(n, s, &wg, d[i])
	}
	wg.Wait()
	return nil
}
