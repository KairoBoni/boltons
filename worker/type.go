package main

import "encoding/xml"

type Message struct {
	AccessKey string `json:"access_key"`
	XML       string `json:"xml"`
}

type TotalNFE struct {
	AccessKey string
	Total     string
}

type NFEProc struct {
	XMLName xml.Name `xml:"nfeProc"`
	Total   string   `xml:"NFe>infNFe>total>ICMSTot>vNF"`
}

type NFE struct {
	XMLName xml.Name `xml:"NFe"`
	Total   string   `xml:"infNFe>total>ICMSTot>vNF"`
}
