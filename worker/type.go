package main

import "encoding/xml"

type Message struct {
	AccessKey string `json:"access_key"`
	XML       string `json:"xml"`
}

type NFEProc struct {
	XMLName xml.Name `xml:"nfeProc"`
	Amount  string   `xml:"NFe>infNFe>amount>ICMSTot>vNF"`
}

type NFE struct {
	XMLName xml.Name `xml:"NFe"`
	Amount  string   `xml:"infNFe>amount>ICMSTot>vNF"`
}
