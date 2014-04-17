package dns

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

func CreateHostedZone(domain, ref, comment string) (*CreateHostedZoneResult, *Error, error) {
	v := &CreateHostedZoneRequest{Name: domain, CallerReference: ref,
		Xmlns: "https://route53.amazonaws.com/doc/2012-02-29/"}
	v.HostedZoneConfig.Comment = comment
	b, err := xml.Marshal(v)
	if err != nil {
		return nil, nil, err
	}
	e := &Error{Message: string(b)}
	return nil, e, nil
}

/*
	<?xml version="1.0" encoding="UTF-8"?>
	<CreateHostedZoneRequest xmlns="https://route53.amazonaws.com/doc/2012-02-29/">
   <Name><DNS domain name></Name>
   <CallerReference><unique description></CallerReference>
   <HostedZoneConfig>
      <Comment><optional comment></Comment>
   </HostedZoneConfig>
</CreateHostedZoneRequest>
*/

type CreateHostedZoneRequest struct {
	Xmlns            string `xml:"xmlns,attr"`
	Name             string
	CallerReference  string
	HostedZoneConfig struct{ Comment string }
}

type CreateHostedZoneResult struct {
}

func dnsGetDate() (string, error) {
	resp, err := http.Get("https://route53.amazonaws.com/date")
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		if resp.Body != nil {
			resp.Body.Close()
		}
		return "", err
	}
	err = resp.Body.Close()
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", nil // os.NewError ( "bad response for date request" )
	}

	return resp.Header.Get("Date"), nil
}

func dnsSign(req *http.Request, id, sk string) error {
	// http://docs.amazonwebservices.com/Route53/latest/DeveloperGuide/RESTAuthentication.html
	date, err := dnsGetDate()
	if err != nil {
		return err
	}
	h := hmac.New(sha256.New, []byte(sk))
	h.Write([]byte(date))
	s := h.Sum(nil)
	e := make([]byte, base64.StdEncoding.EncodedLen(len(s)))
	base64.StdEncoding.Encode(e, s)
	authheader := "AWS3-HTTPS AWSAccessKeyId=" + id + ",Algorithm=HmacSHA256,Signature=" + string(e)

	req.Header.Set("X-Amzn-Authorization", authheader)
	req.Header.Set("x-amz-date", date)
	//	req.Header.Set ( "date", date )
	return nil
}

func dnsGet(path, id, sk, zoneid string) (string, *ErrorResponse, error) {
	req, err := http.NewRequest("GET", "https://route53.amazonaws.com/2012-02-29/hostedzone/"+
		zoneid+
		path, nil)
	if err != nil {
		return "", nil, err
	}
	err = dnsSign(req, id, sk)
	if err != nil {
		return "", nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", nil, err
	}
	fmt.Println(req)

	var buf bytes.Buffer
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return "", nil, err
	}
	err = resp.Body.Close()
	if err != nil {
		return "", nil, err
	}
	body := buf.Bytes()
	if resp.StatusCode != 200 {
		var x ErrorResponse
		err = xml.Unmarshal(body, &x)
		if err != nil {
			return "", nil, err
		}
		return "", &x, nil
	}
	return string(body), nil, nil
}

func GetHostedZone() {
	//	body, e, err := dnsGet ( "hostedzone/" + route53ZoneIdMinepass )

}

type ResourceRecord struct {
	Value string
}

type Change struct {
	Action            string
	ResourceRecordSet struct {
		Name            string
		Type            string
		TTL             int
		ResourceRecords []ResourceRecord `xml:"ResourceRecords>ResourceRecord,omitempty"`
	}
}

type ChangeResourceRecordSetsRequest struct {
	Xmlns       string `xml:"xmlns,attr"`
	ChangeBatch struct {
		Comment string
		Changes []Change `xml:"Changes>Change"`
	}
}

func ChangeResourceRecordSets(comment, id, sk, zoneid string, changes []Change) error {
	c := &ChangeResourceRecordSetsRequest{Xmlns: "https://route53.amazonaws.com/doc/2012-02-29/"}
	c.ChangeBatch.Comment = comment
	c.ChangeBatch.Changes = changes
	b, err := xml.Marshal(c)
	if err != nil {
		return err
	}
	post := bytes.NewBuffer([]byte(`<?xml version="1.0"?>`))
	post.Write(b)
	req, err := http.NewRequest("POST", "https://route53.amazonaws.com/2012-02-29/hostedzone/"+
		zoneid+
		"/rrset", post)
	if err != nil {
		return err
	}
	err = dnsSign(req, id, sk)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return err
	}
	err = resp.Body.Close()
	if err != nil {
		return err
	}
	body := buf.Bytes()
	if resp.StatusCode != 200 {
		var x ErrorResponse
		err = xml.Unmarshal(body, &x)
		if err != nil {
			return err
		}
		return nil
	}

	return nil
}

func AddDns(domain, ip, id, sk, zoneid string) error {
	cs := make([]Change, 1)
	cs[0].Action = "CREATE"
	cs[0].ResourceRecordSet.Name = domain
	cs[0].ResourceRecordSet.Type = "A"
	cs[0].ResourceRecordSet.TTL = 60
	cs[0].ResourceRecordSet.ResourceRecords = make([]ResourceRecord, 1)
	cs[0].ResourceRecordSet.ResourceRecords[0].Value = ip

	return ChangeResourceRecordSets("AddDns", id, sk, zoneid, cs)
}

type ListResourceRecordSetsResponse struct {
	ResourceRecordSets struct {
		ResourceRecordSet []struct {
			Name            string
			Type            string
			TTL             int
			ResourceRecords struct {
				ResourceRecord []struct {
					Value string
				}
			}
		}
	}
}

func GetIpForDomain(domain, id, sk, zoneid string) (string, int, error) {
	resp, eresp, err := dnsGet("/rrset?name="+url.QueryEscape(domain)+"&type=A&maxitems=1",
		id, sk, zoneid)
	if err != nil {
		return "", 0, err
	}
	if eresp != nil {
		return "", 0, fmt.Errorf("%+v", eresp)
	}
	v := ListResourceRecordSetsResponse{}
	err = xml.Unmarshal([]byte(resp), &v)
	if err != nil {
		return "", 0, err
	}
	if len(v.ResourceRecordSets.ResourceRecordSet) != 1 ||
		len(v.ResourceRecordSets.ResourceRecordSet[0].ResourceRecords.ResourceRecord) != 1 {
		return "", 0, errors.New("not resource found")
	}
	return v.ResourceRecordSets.ResourceRecordSet[0].ResourceRecords.ResourceRecord[0].Value,
		v.ResourceRecordSets.ResourceRecordSet[0].TTL, nil
}

func RemoveDns(domain, id, sk, zoneid string) error {
	ip, ttl, err := GetIpForDomain(domain, id, sk, zoneid)
	if err != nil {
		return err
	}
	cs := make([]Change, 1)
	cs[0].Action = "DELETE"
	cs[0].ResourceRecordSet.Name = domain
	cs[0].ResourceRecordSet.Type = "A"
	cs[0].ResourceRecordSet.TTL = ttl
	cs[0].ResourceRecordSet.ResourceRecords = make([]ResourceRecord, 1)
	cs[0].ResourceRecordSet.ResourceRecords[0].Value = ip

	return ChangeResourceRecordSets("RemoveDns", id, sk, zoneid, cs)
}

func RemoveIp(domain, ip, id, sk, zoneid string) error {
	cs := make([]Change, 1)
	cs[0].Action = "DELETE"
	cs[0].ResourceRecordSet.Name = domain
	cs[0].ResourceRecordSet.Type = "A"
	cs[0].ResourceRecordSet.TTL = 600
	cs[0].ResourceRecordSet.ResourceRecords = make([]ResourceRecord, 1)
	cs[0].ResourceRecordSet.ResourceRecords[0].Value = ip

	return ChangeResourceRecordSets("RemoveIp", id, sk, zoneid, cs)
}

// POST /2012-02-29/hostedzone/Z1PA6795UKMFR9/rrset HTTP/1.1
/*
 <?xml version="1.0"?>
	<ChangeResourceRecordSetsRequest xmlns="https://route53.amazonaws.com/
doc/2012-02-29/">
   <ChangeBatch>
      <Comment>
      This change batch creates a TXT record for www.example.com. and
      changes the A record for foo.example.com. from
      192.0.2.3 to 192.0.2.1.
      </Comment>
      <Changes>
         <Change>
            <Action>CREATE</Action>
            <ResourceRecordSet>
               <Name>www.example.com.</Name>
               <Type>TXT</Type>
               <TTL>600</TTL>
               <ResourceRecords>
                  <ResourceRecord>
                     <Value>"item 1" "item 2" "item 3"</Value>
                  </ResourceRecord>
               </ResourceRecords>
            </ResourceRecordSet>
         </Change>
         <Change>
            <Action>DELETE</Action>
            <ResourceRecordSet>
               <Name>foo.example.com.</Name>
               <Type>A</Type>
               <TTL>600</TTL>
               <ResourceRecords>
                  <ResourceRecord>
                     <Value>192.0.2.3</Value>
                  </ResourceRecord>
               </ResourceRecords>
            </ResourceRecordSet>
         </Change>
         <Change>
            <Action>CREATE</Action>
            <ResourceRecordSet>
               <Name>foo.example.com.</Name>
               <Type>A</Type>
               <TTL>600</TTL>
               <ResourceRecords>
                  <ResourceRecord>
                     <Value>192.0.2.1</Value>
                  </ResourceRecord>
               </ResourceRecords>
            </ResourceRecordSet>
         </Change>
      </Changes>
   </ChangeBatch>
</ChangeResourceRecordSetsRequest>
*/
