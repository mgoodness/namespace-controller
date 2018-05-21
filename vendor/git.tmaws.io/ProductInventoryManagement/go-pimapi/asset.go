package pimapi

const (
	assetBasePath = "asset"
)

type AssetService struct {
	client *Client
}

type AssetDoc struct {

	// Product type
	PCI    string `json:"PCI"`
	Goal   string `json:"Goal"`
	JIRA   string `json:"JIRA"`
	Name   string `json:"Name"`
	Custom []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"Custom"`
	Runbook      string  `json:"Runbook"`
	ThreeRdParty string  `json:"3rdParty"`
	AWSReady     float64 `json:"AWSReady"`
	Category     string  `json:"Category"`
	SuppTeam     string  `json:"SuppTeam"`
	AWSAssets    int     `json:"AWSAssets"`
	TechOwner    string  `json:"TechOwner"`
	AWSMigrate   string  `json:"AWSMigrate"`
	OnCallInfo   string  `json:"OnCallInfo"`
	Description  string  `json:"Description"`
	DisplayName  string  `json:"DisplayName"`
	// BusinessTier      int     `json:"BusinessTier"`
	OnpremAssets       int     `json:"OnpremAssets"`
	ProductOwner       string  `json:"ProductOwner"`
	SeniorProductOwner string  `json:"SeniorProductOwner"`
	ProgramOwner       string  `json:"ProgramOwner"`
	SlackChannel       string  `json:"SlackChannel"`
	SupportEmail       string  `json:"SupportEmail"`
	AWSMigrateEnd      string  `json:"AWSMigrateEnd"`
	Documentation      string  `json:"Documentation"`
	ProductQuality     float64 `json:"ProductQuality"`
	AWSMigrateStart    string  `json:"AWSMigrateStart"`
	OperatingSystem    string  `json:"OperatingSystem"`
	SeniorTechOwner    string  `json:"SeniorTechOwner"`
	TechMaturityURL    string  `json:"TechMaturityURL"`
	AWSMigrateComment  string  `json:"AWSMigrateComment"`

	// PhysicalServer type
	OS        string `json:"OS"`
	Site      string `json:"Site"`
	Type      string `json:"Type"`
	Model     string `json:"Model"`
	RtID      int    `json:"rt_id"`
	Uname     string `json:"uname"`
	CPUnum    int    `json:"CPUnum"`
	Creator   string `json:"Creator"`
	NetAddr   string `json:"NetAddr"`
	CPUmodel  string `json:"CPUmodel"`
	Location  string `json:"Location"`
	NlyteTag  string `json:"NlyteTag"`
	PoolName  string `json:"PoolName"`
	MemoryMiB int64  `json:"MemoryMiB"`
	// RunStatus              interface{} `json:"RunStatus"`
	OSRevision   string `json:"OSRevision"`
	ServiceTag   string `json:"ServiceTag"`
	Environment  string `json:"Environment"`
	Architecture string `json:"Architecture"`
	// BusinessUnit           []string `json:"BusinessUnit"`
	ComputeUnits     float32 `json:"ComputeUnits"`
	Manufacturer     string  `json:"Manufacturer"`
	SerialNumber     string  `json:"SerialNumber"`
	ServerStatus     string  `json:"ServerStatus"`
	OSDescription    string  `json:"OSDescription"`
	AvailabilityZone string  `json:"AvailabilityZone"`
	// KickstartMacAddress    string  `json:"KickstartMacAddress"`
	TechnologyOrganization string `json:"TechnologyOrganization"`

	// VirtualServer type
	Operator    string `json:"Operator"`
	TenantID    int    `json:"TenantID"`
	XstagSr     string `json:"xstag_sr"`
	TMrelease   int    `json:"TMrelease"`
	LastAction  string `json:"LastAction"`
	XstagDistro string `json:"xstag_distro"`
}

type Asset struct {
	Id              int64       `json:"id,omitempty"`
	AssetTypeId     int64       `json:"assetTypeId"`
	ExtId           string      `json:"extId,omitempty"`
	ExtSourceId     int64       `json:"extSourceId,omitempty"`
	Deleted         bool        `json:"deleted,omitempty"`
	Doc             AssetDoc    `json:"doc"`
	Links           *AssetLinks `json:"links,omitempty"`
	ProductCodeId   int64       `json:"productCodeId,omitempty"`
	ProductCodeName string      `json:"productCodeName,omitempty"`
	Created         string      `json:"created,omitempty"`
	Updated         string      `json:"updated,omitempty"`
	User            string      `json:"user,omitempty"`
}

type AssetLink struct {
	AssetId     int64  `json:"assetId,omitempty"`
	AssetTypeId int64  `json:"assetTypeId,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	LinkId      int64  `json:"linkId,omitempty"`
	LinkRuleId  int64  `json:"linkRuleId,omitempty"`
}

type AssetLinks struct {
	Children []*AssetLink `json:"children"`
	Parents  []*AssetLink `json:"parents"`
}

type AssetItems struct {
	Items      []Asset `json:"items"`
	Page       *Page   `json:"page,omitempty"`
	CountTotal int     `json:"countTotal"`
}

// Get Asset by Id
func (srv AssetService) ID(id string) (Asset, error) {
	var a Asset

	err := srv.client.ID(assetBasePath, id, &a)

	return a, err
}

// Get Assets
func (srv AssetService) Get(params map[string]string) (AssetItems, error) {
	var a AssetItems

	err := srv.client.Get(assetBasePath, params, &a)

	return a, err
}

// Create Asset
func (srv AssetService) Create(createRequest AssetItems) (IdsResponse, error) {
	id, err := srv.client.Create(assetBasePath, createRequest)

	return id, err
}

// Update Asset
func (srv AssetService) Update(updateRequest AssetItems) (IdsResponse, error) {
	id, err := srv.client.Update(assetBasePath, updateRequest)

	return id, err
}

// Delete Asset
func (srv AssetService) Delete(deleteRequest IdsRequest) (IdsResponse, error) {
	id, err := srv.client.Delete(assetBasePath, deleteRequest)

	return id, err
}

// Get All Assets
func (srv AssetService) All(params map[string]string) ([]Asset, error) {
	var a []Asset

	err := srv.client.Get(assetBasePath+"/all", params, &a)

	return a, err
}
