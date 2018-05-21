package pimapi

const (
    assetTypeBasePath = "metadata/assettype"
)

type AssetTypeService struct {
    client *Client
}

type AssetType struct {
    Id              int64                   `json:"id,omitempty"`
    Name            string                  `json:"name,omitempty"`
    DisplayName     string                  `json:"displayName,omitempty"`
    ReadOnly        bool                    `json:"readOnly,omitempty"`
    GeneratePRD     bool                    `json:"generatePrd,omitempty"`
    Created         string                  `json:"created,omitempty"`
    Updated         string                  `json:"updated,omitempty"`
    User            string                  `json:"user,omitempty"`
}

type AssetTypeItems struct {
    Items []AssetType   `json:"items"`
    Page  *Page         `json:"page,omitempty"`
}

// Get Asset by Id
func (srv AssetTypeService) ID(id string) (AssetType, error) {
    var a AssetType

    err := srv.client.ID(assetTypeBasePath, id, &a)

    return a, err
}

// Get Assets
func (srv AssetTypeService) Get(params map[string]string) (AssetTypeItems, error) {
    var a AssetTypeItems

    err := srv.client.Get(assetTypeBasePath, params, &a)

    return a, err
}

// Create Asset
func (srv AssetTypeService) Create(createRequest interface{}) (IdsResponse, error) {
    id, err := srv.client.Create(assetTypeBasePath, createRequest)

    return id, err
}

// Update Asset
func (srv AssetTypeService) Update(updateRequest interface{}) (IdsResponse, error) {
    id, err := srv.client.Update(assetTypeBasePath, updateRequest)

    return id, err
}

// Delete Asset
func (srv AssetTypeService) Delete(deleteRequest interface{}) (IdsResponse, error) {
    id, err := srv.client.Delete(assetTypeBasePath, deleteRequest)

    return id, err
}