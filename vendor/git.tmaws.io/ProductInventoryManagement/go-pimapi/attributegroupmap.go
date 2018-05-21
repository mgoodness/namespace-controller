package pimapi

const (
    attributeGroupMapBasePath = "metadata/attributegroupmap"
)

type AttributeGroupMapService struct {
    client *Client
}

type AttributeGroupMap struct {
    Id                  int64   `json:"id,omitempty"`
    AssetTypeId         int64   `json:"assetTypeId,omitempty"`
    AttributeGroupId    int64   `json:"attributeGroupId,omitempty"`
    DisplayOrder        int64   `json:"displayOrder,omitempty"`
}

type AttributeGroupMapItems struct {
    Items []AttributeGroupMap   `json:"items"`
    Page  *Page                 `json:"page,omitempty"`
}

// Get Asset by Id
func (srv AttributeGroupMapService) ID(id string) (AttributeGroupMap, error) {
    var a AttributeGroupMap

    err := srv.client.ID(attributeGroupMapBasePath, id, &a)

    return a, err
}

// Get Assets
func (srv AttributeGroupMapService) Get(params map[string]string) (AttributeGroupMapItems, error) {
    var a AttributeGroupMapItems

    err := srv.client.Get(attributeGroupMapBasePath, params, &a)

    return a, err
}

// Create Asset
func (srv AttributeGroupMapService) Create(createRequest interface{}) (IdsResponse, error) {
    id, err := srv.client.Create(attributeGroupMapBasePath, createRequest)

    return id, err
}

// Update Asset
func (srv AttributeGroupMapService) Update(updateRequest interface{}) (IdsResponse, error) {
    id, err := srv.client.Update(attributeGroupMapBasePath, updateRequest)

    return id, err
}

// Delete Asset
func (srv AttributeGroupMapService) Delete(deleteRequest interface{}) (IdsResponse, error) {
    id, err := srv.client.Delete(attributeGroupMapBasePath, deleteRequest)

    return id, err
}