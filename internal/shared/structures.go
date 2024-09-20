package shared

import (
	"net"
	"net/http"
	"os"
)

// CORE
type Params struct {
	Domain               string
	Subdomain            string
	FilePathSubdomains   string
	FileContentSubdoms   string
	FilePathIPv4Addrs    string
	FilePathIPv6Addrs    string
	FileContentIPv4Addrs string
	FileContentIPv6Addrs string
}

type Args struct {
	Verbose             bool
	Domain              string // target domain
	NewOutputDirPath    string // custom output dir path
	HttpCode            bool
	WordlistPath        string
	ExcHttpCodes        string // results to hide (specified by HTTP status code)
	FilHttpCodes        string // results to display (specified by HTTP status code)
	SubOnlyIp           bool
	AnalyzeHeader       bool
	PortScan            string // port range
	DbExtendPath        string // File path containing endpoints
	Timeout             int    // in seconds
	TorRoute            bool
	DnsLookup           bool
	DnsLookupCustom     string // Custom DNS server (args)
	DnsLookupTimeout    int
	HttpRequestDelay    int // in milliseconds
	DisableAllOutput    bool
	AnalyseHeaderSingle bool // HTTP
	Subdomain           string
	HttpRequestMethod   string
	ShowAllHeaders      bool
	DetectPurpose       bool // DNS lookups, HTTP header analysis (mail server, API etc.)
	MisconfTest         bool // CORS, header injection, request smuggling etc.
	AllowRedirects      bool
	IpAddress           string
	EnableVHostEnum     bool
	FilterHttpSize      string
}

type SetTestResults struct {
	TestName   string
	TestResult string
	Subdomain  string
}

type EnumerationMethod struct {
	MethodKey string
	Action    func(*Args, *http.Client, *FilePaths)
}

// OPTION-MANAGER
type SettingsHandler struct {
	Streams       *FileStreams
	Args          *Args
	Params        Params
	HttpClient    *http.Client
	ConsoleOutput chan<- string
	CodeFilterExc []string
	CodeFilter    []string
	IpAddrs       []string
	IpAddrsOut    string
	URL           string
}

// ENUM
type DnsLookupOptions struct {
	Subdomain string
	IpAddress net.IP
}

type HttpHeaders struct {
	Server string
	Hsts   string
	PowBy  string
	Csp    string
}

// OUTPUT
type FilePaths struct {
	FilePathSubdomain string
	FilePathIPv4      string
	FilePathIPv6      string
	FilePathJSON      string
}

type FileStreams struct {
	Ipv4AddrStream  *os.File
	Ipv6AddrStream  *os.File
	SubdomainStream *os.File
}

type ParamsSetupFilesBase struct {
	FileParams *Params
	CliArgs    *Args
	FilePaths  *FilePaths
	Subdomain  string
}

type SubdomainIpAddresses struct {
	IPv4 []net.IP `json:"ipv4Addresses"`
	IPv6 []net.IP `json:"ipv6Addresses"`
}

type SubdomainBase struct {
	Subdomain   []string `json:"subdomain"`
	OpenPorts   []int    `json:"openPorts"`
	IpAddresses SubdomainIpAddresses
}

type JsonResult struct {
	Subdomains []SubdomainBase `json:"subdomains"`
}
