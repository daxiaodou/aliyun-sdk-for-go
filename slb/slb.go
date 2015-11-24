// SLB API package
package slb

import (
	"runtime"
	"strings"
	"time"

	"encoding/json"

	"github.com/jiangshengwu/aliyun-sdk-for-go/log"
	"github.com/jiangshengwu/aliyun-sdk-for-go/util"
)

const (
	// SLB API Host
	SLBHost string = "https://slb.aliyuncs.com/?"

	// All SLB APIs only support GET method
	SLBHttpMethod = "GET"

	// SDK only supports JSON format
	Format = "JSON"

	Version          = "2014-05-15"
	SignatureMethod  = "HMAC-SHA1"
	SignatureVersion = "1.0"
)

// struct for SLB client
type SlbClient struct {
	Common *CommonParam

	// Access to API call from this client
	LoadBalancer      LoadBalancerService
	ServerCertificate ServerCertificateService
	Listener          ListenerService
	BackendServer     BackendServerService
}

// Initialize an SLB client
func NewClient(accessKeyId string, accessKeySecret string, resourceOwnerAccount string) *SlbClient {
	client := &SlbClient{}
	client.Common = &CommonParam{}
	client.Common.AccessKeyId = accessKeyId
	client.Common.AccessKeySecret = accessKeySecret
	client.Common.ResourceOwnerAccount = resourceOwnerAccount
	ps := map[string]string{
		"Format":           Format,
		"Version":          Version,
		"AccessKeyId":      client.Common.AccessKeyId,
		"SignatureMethod":  SignatureMethod,
		"SignatureVersion": SignatureVersion,
	}
	client.Common.attr = ps

	client.LoadBalancer = &LoadBalancerOperator{client.Common}
	client.ServerCertificate = &ServerCertificateOperator{client.Common}
	client.Listener = &ListenerOperator{client.Common}
	client.BackendServer = &BackendServerOperator{client.Common}
	return client
}

func (client *SlbClient) GetClientName() string {
	return "SLB Client"
}

func (client *SlbClient) GetVersion() string {
	return client.Common.attr["Version"]
}

func (client *SlbClient) GetSignatureMethod() string {
	return client.Common.attr["SignatureMethod"]
}

func (client *SlbClient) GetSignatureVersion() string {
	return client.Common.attr["SignatureVersion"]
}

// struct for common parameters
type CommonParam struct {
	AccessKeyId          string
	AccessKeySecret      string
	ResourceOwnerAccount string
	attr                 map[string]string
}

func RequestAPI(params map[string]string) (string, error) {
	query := util.GetQueryFromMap(params)
	req := &util.AliyunRequest{}
	req.Url = SLBHost + query
	log.Debug(req.Url)
	result, err := req.DoGetRequest()
	return result, err
}

// Get function name by skip
// which means the differs between Caller and Callers
func GetFuncName(skip int) string {
	pc, _, _, _ := runtime.Caller(skip)
	name := runtime.FuncForPC(pc).Name()
	i := strings.LastIndex(name, ".")
	if i >= 0 {
		name = name[i+1:]
	}
	return name
}

// Generate all parameters include Signature
func (c *CommonParam) ResolveAllParams(action string, params map[string]string) map[string]string {
	if params == nil {
		params = make(map[string]string, len(c.attr))
	}
	// Process parameters
	for key, value := range c.attr {
		params[key] = value
	}
	params["Action"] = action
	if c.ResourceOwnerAccount != "" {
		params["ResourceOwnerAccount"] = c.ResourceOwnerAccount
	}
	params["TimeStamp"] = time.Now().UTC().Format("2006-01-02T15:04:05Z")
	params["SignatureNonce"] = util.GetGuid()
	sign := util.MapToSign(params, c.AccessKeySecret, SLBHttpMethod)
	params["Signature"] = sign
	return params
}

func (c *CommonParam) Request(action string, params map[string]string, response interface{}) error {
	p := c.ResolveAllParams(action, params)
	result, err := RequestAPI(p)
	if err != nil {
		return err
	}
	log.Debug(result)
	err = json.Unmarshal([]byte(result), response)
	if err != nil {
		log.Error(err)
	}
	return nil
}
