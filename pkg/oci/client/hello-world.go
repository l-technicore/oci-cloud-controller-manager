package client

//// You can edit this code!
//// Click here and start typing.
//package main
//
//import (
//	"fmt"
//	"net/url"
//	"os"
//)
//
//func main() {
//
//	url, err := url.Parse("http://www.google.com/search?q=hi")
//	if err != nil {
//		fmt.Println("Error: ", err.Error())
//	}
//	fmt.Println("Url: ", url.Host)
//
//	os.Setenv("OCI_PROXY", "http://www.bing.com")
//
//	ociProxy := os.Getenv("OCI_PROXY")
//	if ociProxy != "" {
//		proxyURL, err := url.Parse(ociProxy)
//		if err != nil {
//			fmt.Println(err.Error(), "failed to parse OCI proxy url: %s", ociProxy)
//		}
//		if proxyURL.Scheme == "https" {
//			os.Setenv("HTTPS_PROXY", ociProxy)
//		} else {
//			os.Setenv("HTTP_PROXY", ociProxy)
//		}
//	}
//
//	fmt.Println(os.Getenv("HTTP_PROXY"))
//
//	//transport := common.DefaultBaseClientWithSigner(common.HTTPRequestSigner)
//	//tls, err := transport.TLSConfigProvider.NewOrDefault()
//	//if err != nil {
//	//	fmt.Println("Error: ", err.Error())
//	//}
//	//proxy, err := transport.TransportTemplate.NewOrDefault(tls)
//	//if err != nil {
//	//	fmt.Println("Error: ", err.Error())
//	//}
//	//
//	//res, err := proxy.(*http.Transport).Proxy(&http.Request{URL: url})
//	//if err != nil {
//	//	fmt.Println("Error: ", err.Error())
//	//}
//	//if res != nil {
//	//	fmt.Println("Proxy: ", res.Host)
//	//}
//}
