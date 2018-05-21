package pimapi

const (
    attributeTypeMapBasePath = "metadata/attributetypemap"
)

type AttributeTypeMapService struct {
    client *Client
}

type AttributeTypeMap struct {
    Id                          int64       `json:"id,omitempty"`
    AssetTypeId                 int64       `json:"assetTypeId,omitempty"`
    AttributeTypeId             int64       `json:"attributeGroupId,omitempty"`
    Required                    bool        `json:"required,omitempty"`
    GroupName                   string      `json:"groupName,omitempty"`
    AttributeGroupId            int64       `json:"attributeGroupId,omitempty"`
    ValidationTypeOverrideId    int64       `json:"validationTypeOverrideId,omitempty"`
    DisplayOrder                int64       `json:"displayOrder,omitempty"`
    Created                     string      `json:"created,omitempty"`
    Updated                     string      `json:"updated,omitempty"`
    User                        string      `json:"user,omitempty"`
}

type AttributeTypeMapItems struct {
    Items []AttributeTypeMap    `json:"items"`
    Page  *Page                 `json:"page,omitempty"`
}

// Get Asset by Id
func (srv AttributeTypeMapService) ID(id string) (AttributeTypeMap, error) {
    var a AttributeTypeMap

    err := srv.client.ID(attributeTypeMapBasePath, id, &a)

    return a, err
}

// Get Assets
func (srv AttributeTypeMapService) Get(params map[string]string) (AttributeTypeMapItems, error) {
    var a AttributeTypeMapItems

    err := srv.client.Get(attributeTypeMapBasePath, params, &a)

    return a, err
}

// Create Asset
func (srv AttributeTypeMapService) Create(createRequest interface{}) (IdsResponse, error) {
    id, err := srv.client.Create(attributeTypeMapBasePath, createRequest)

    return id, err
}

// Update Asset
func (srv AttributeTypeMapService) Update(updateRequest interface{}) (IdsResponse, error) {
    id, err := srv.client.Update(attributeTypeMapBasePath, updateRequest)

    return id, err
}

// Delete Asset
func (srv AttributeTypeMapService) Delete(deleteRequest interface{}) (IdsResponse, error) {
    id, err := srv.client.Delete(attributeTypeMapBasePath, deleteRequest)

    return id, err
}