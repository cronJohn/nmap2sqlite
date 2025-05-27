package nmap

type NmapRunAttributes struct {
	Scanner          string   `xml:"scanner,attr"          json:"scanner"`
	Args             string   `xml:"args,attr"             json:"args"`
	Start            UnixTime `xml:"start,attr"            json:"start"`
	StartStr         string   `xml:"startstr,attr"         json:"startstr"`
	Version          string   `xml:"version,attr"          json:"version"`
	XMLOutputVersion string   `xml:"xmloutputversion,attr" json:"xmloutputversion"`
}

type UnixTime int64

type ScanInfo struct {
	Type        string `xml:"type,attr"        json:"type"`
	Protocol    string `xml:"protocol,attr"    json:"protocol"`
	NumServices int64  `xml:"numservices,attr" json:"numservices"`
	Services    string `xml:"services,attr"    json:"services"`
	// Not sure why some packages include ScanFlags when that isn't even outputted
}

type Verbose struct {
	Level int64 `xml:"level,attr" json:"level"`
}

type Debugging struct {
	Level int64 `xml:"level,attr" json:"level"`
}
type Host struct {
	StartTime     UnixTime      `xml:"starttime,attr"     json:"starttime"`
	EndTime       UnixTime      `xml:"endtime,attr"       json:"endtime"`
	Status        Status        `xml:"status"             json:"status"`
	Addresses     []Address     `xml:"address"            json:"addresses"`
	Hostnames     []Hostname    `xml:"hostnames>hostname" json:"hostnames"`
	Ports         []Port        `xml:"ports>port"         json:"ports"`
	ExtraPorts    []ExtraPorts  `xml:"ports>extraports"   json:"extraports"`
	HostScripts   []Script      `xml:"hostscript>script"  json:"hostscripts"`
	Os            Os            `xml:"os"                 json:"os"`
	Distance      Distance      `xml:"distance"           json:"distance"`
	Uptime        Uptime        `xml:"uptime"             json:"uptime"`
	TcpSequence   TcpSequence   `xml:"tcpsequence"        json:"tcpsequence"`
	IpIdSequence  IpIdSequence  `xml:"ipidsequence"       json:"ipidsequence"`
	TcpTsSequence TcpTsSequence `xml:"tcptssequence"      json:"tcptssequence"`
	Trace         Trace         `xml:"trace"              json:"trace"`
	Times         Times         `xml:"times"              json:"times"`
}

type Address struct {
	Addr     string `xml:"addr,attr"     json:"addr"`
	AddrType string `xml:"addrtype,attr" json:"addrtype"`
	Vendor   string `xml:"vendor,attr"   json:"vendor"`
}

type Status struct {
	State     string `xml:"state,attr"      json:"state"`
	Reason    string `xml:"reason,attr"     json:"reason"`
	ReasonTTL int64  `xml:"reason_ttl,attr" json:"reason_ttl"`
}

type Hostname struct {
	Name string `xml:"name,attr" json:"name"`
	Type string `xml:"type,attr" json:"type"`
}

type Port struct {
	Protocol string   `xml:"protocol,attr" json:"protocol"`
	PortId   int64    `xml:"portid,attr"   json:"id"`
	State    State    `xml:"state"         json:"state"`
	Service  Service  `xml:"service"       json:"service"`
	Scripts  []Script `xml:"script"        json:"scripts"`
}

type State struct {
	State     string `xml:"state,attr"      json:"state"`
	Reason    string `xml:"reason,attr"     json:"reason"`
	ReasonTTL int64  `xml:"reason_ttl,attr" json:"reason_ttl"`
}

type Service struct {
	Name      string `xml:"name,attr"      json:"name"`
	Conf      int64  `xml:"conf,attr"      json:"conf"`
	Method    string `xml:"method,attr"    json:"method"`
	Version   string `xml:"version,attr"   json:"version"`
	Product   string `xml:"product,attr"   json:"product"`
	ExtraInfo string `xml:"extrainfo,attr" json:"extrainfo"`
	Tunnel    string `xml:"tunnel,attr"    json:"tunnel"`
	Proto     string `xml:"proto,attr"     json:"proto"`
	Rpcnum    string `xml:"rpcnum,attr"    json:"rpcnum"`
	Lowver    string `xml:"lowver,attr"    json:"lowver"`
	Highver   string `xml:"highver,attr"   json:"highver"`
	OsType    string `xml:"ostype,attr"    json:"ostype"`
	ServiceFp string `xml:"servicefp,attr" json:"servicefp"`
	CPEs      []CPE  `xml:"cpe"            json:"cpes"`
}

type CPE string

type ExtraPorts struct {
	State        string        `xml:"state,attr"   json:"state"`
	Count        int64         `xml:"count,attr"   json:"count"`
	ExtraReasons []ExtraReason `xml:"extrareasons" json:"reasons"`
}

type ExtraReason struct {
	Reason string `xml:"reason,attr" json:"reason"`
	Count  int64  `xml:"count,attr"  json:"count"`
	Proto  string `xml:"proto,attr"  json:"proto"`
	Ports  string `xml:"ports,attr"  json:"ports"`
}

type Os struct {
	OsMatches []OsMatch  `xml:"osmatch"  json:"osmatches"`
	PortsUsed []PortUsed `xml:"portused" json:"portsused"`
}

type PortUsed struct {
	State  string `xml:"state,attr"  json:"state"`
	Proto  string `xml:"proto,attr"  json:"proto"`
	PortId int64  `xml:"portid,attr" json:"portid"`
}

type OsMatch struct {
	Name      string    `xml:"name,attr"     json:"name"`
	Accuracy  string    `xml:"accuracy,attr" json:"accuracy"`
	Line      string    `xml:"line,attr"     json:"line"`
	OsClasses []OsClass `xml:"osclass"       json:"osclasses"`
}

type OsClass struct {
	Type     string `xml:"type,attr"     json:"type"`
	Vendor   string `xml:"vendor,attr"   json:"vendor"`
	OsFamily string `xml:"osfamily,attr" json:"osfamily"`
	OsGen    string `xml:"osgen,attr"`
	Accuracy string `xml:"accuracy,attr" json:"accuracy"`
	CPEs     []CPE  `xml:"cpe"           json:"cpes"`
}

type Distance struct {
	Value int64 `xml:"value,attr" json:"value"`
}

type Uptime struct {
	Seconds  int64  `xml:"seconds,attr"  json:"seconds"`
	Lastboot string `xml:"lastboot,attr" json:"lastboot"`
}

type TcpSequence struct {
	Index      int64  `xml:"index,attr"      json:"index"`
	Difficulty string `xml:"difficulty,attr" json:"difficulty"`
	Values     string `xml:"values,attr"     json:"values"`
}

type (
	IpIdSequence  Sequence
	TcpTsSequence Sequence
)

type Sequence struct {
	Class  string `xml:"class,attr"  json:"class"`
	Values string `xml:"values,attr" json:"values"`
}

type Script struct {
	Id       string    `xml:"id,attr"     json:"id"`
	Output   string    `xml:"output,attr" json:"output"`
	Elements []Element `xml:"elem"        json:"elements"`
	Tables   []Table   `xml:"table"       json:"tables"`
}

type Table struct {
	Key      string    `xml:"key,attr" json:"key"`
	Elements []Element `xml:"elem"     json:"elements"`
	Table    []Table   `xml:"table"    json:"tables"`
}

type Element struct {
	Key   string `xml:"key,attr"  json:"key"`
	Value string `xml:",chardata" json:"value"`
}

type Trace struct {
	Port  int64  `xml:"port,attr"  json:"port"`
	Proto string `xml:"proto,attr" json:"proto"`
	Hops  []Hop  `xml:"hop"        json:"hops"`
}

type Hop struct {
	TTL    int64   `xml:"ttl,attr"    json:"ttl"`
	IPAddr string  `xml:"ipaddr,attr" json:"ipaddr"`
	RTT    float64 `xml:"rtt,attr"    json:"rtt"`
	Host   string  `xml:"host,attr"   json:"host"`
}

type Times struct {
	SRTT   int64 `xml:"srtt,attr"   json:"srtt"`
	RTTVAR int64 `xml:"rttvar,attr" json:"rttvar"`
	To     int64 `xml:"to,attr"     json:"to"`
}

type RunStats struct {
	Finished Finished  `xml:"finished" json:"finished"`
	Hosts    HostStats `xml:"hosts"    json:"hosts"`
}

type Finished struct {
	Time     UnixTime `xml:"time,attr"     json:"time"`
	TimeStr  string   `xml:"timestr,attr"  json:"timestr"`
	Elapsed  float64  `xml:"elapsed,attr"  json:"elapsed"`
	Summary  string   `xml:"summary,attr"  json:"summary"`
	Exit     string   `xml:"exit,attr"     json:"exit"`
	ErrorMsg string   `xml:"errormsg,attr" json:"errormsg"`
}

type HostStats struct {
	Up    int64 `xml:"up,attr"    json:"up"`
	Down  int64 `xml:"down,attr"  json:"down"`
	Total int64 `xml:"total,attr" json:"total"`
}
