package main

import "encoding/xml"

//NFEProc to unmashal the XML to get the value total of XML
//when the first field of the XML ir nfeProc
type NFEProc struct {
	XMLName xml.Name `xml:"nfeProc"`
	Amount  string   `xml:"NFe>infNFe>total>ICMSTot>vNF"`
}

//NFE to unmashal the XML to get the value total of XML
//when the first field of the XML ir NFe
type NFE struct {
	XMLName xml.Name `xml:"NFe"`
	Amount  string   `xml:"infNFe>total>ICMSTot>vNF"`
}
