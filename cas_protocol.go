package cas_connector

import (
	"encoding/xml"
	"strings"
)

// --- CAS protocol XML response parsing -------------------------------
//
// A successful /serviceValidate response looks like:
//
//	<cas:serviceResponse xmlns:cas="http://www.yale.edu/tp/cas">
//	  <cas:authenticationSuccess>
//	    <cas:user>jdoe</cas:user>
//	    <cas:attributes>
//	      <cas:mail>jdoe@example.com</cas:mail>
//	      <cas:displayName>Jane Doe</cas:displayName>
//	    </cas:attributes>
//	  </cas:authenticationSuccess>
//	</cas:serviceResponse>
//
// A failed one looks like:
//
//	<cas:serviceResponse xmlns:cas="http://www.yale.edu/tp/cas">
//	  <cas:authenticationFailure code="INVALID_TICKET">
//	    Ticket ST-xxx not recognized
//	  </cas:authenticationFailure>
//	</cas:serviceResponse>

type casServiceResponse struct {
	XMLName xml.Name                  `xml:"serviceResponse"`
	Success *casAuthenticationSuccess `xml:"authenticationSuccess"`
	Failure *casAuthenticationFailure `xml:"authenticationFailure"`
}

type casAuthenticationSuccess struct {
	User       string          `xml:"user"`
	Attributes casAttributeSet `xml:"attributes"`
}

// casAttributeSet uses encoding/xml's special ",any" tag to capture every
// child element of <cas:attributes> regardless of its name -- this is the
// correct way to do a wildcard match in Go's XML decoder (a bare
// "attributes>Any" path, which is NOT a wildcard, would silently match
// nothing).
type casAttributeSet struct {
	Items []casAttribute `xml:",any"`
}

// casAttribute captures a single <cas:attributes><cas:xxx>value</cas:xxx>
// element generically, since different CAS servers release different
// attribute sets under different names.
type casAttribute struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

// attribute returns the value of the first matching attribute name found
// (case-insensitive), trying each candidate name in order. Different CAS
// deployments release the same information under different attribute
// names (mail vs. email, displayName vs. cn, etc.), so callers pass a
// short list of names to try.
func (s *casAuthenticationSuccess) attribute(names ...string) string {
	for _, want := range names {
		for _, attr := range s.Attributes.Items {
			if strings.EqualFold(attr.XMLName.Local, want) {
				return strings.TrimSpace(attr.Value)
			}
		}
	}
	return ""
}

type casAuthenticationFailure struct {
	Code    string `xml:"code,attr"`
	Message string `xml:",chardata"`
}
