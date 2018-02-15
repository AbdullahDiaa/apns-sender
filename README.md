# apns-sender
HTTP2 APNS sender using golang

##Usgae 

```golang
	n := apns.Notification{
		Topic:      "TOPIC",
		Expiration: time.Now(),
		Payload:    []byte(`{"aps":{"alert":"Hello!"}}`),
	}
	
	cert := "CERT_PEM_STRING"
	s := apns.CreateServer(apns.APNSDevelopmentServer, cert)
	d := []string{"TOKEN"}
	apns.Push(n, s, d)
```