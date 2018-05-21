package pimapi

const (
    attributeGroupBasePath = "metadata/attributegroup"
)

type AttributeGroupService struct {
    client *Client
}

type AttributeGroup struct {
    Id              int64                   `json:"id,omitempty"`
    Name            string                  `json:"name,omitempty"`
}

type AttributeGroupItems struct {
    Items []AttributeGroup  `json:"items"`
    Page  *Page             `json:"page,omitempty"`
}

// Get Asset by Id
func (srv AttributeGroupService) ID(id string) (AttributeGroup, error) {
    var a AttributeGroup

    err := srv.client.ID(attributeGroupBasePath, id, &a)

    return a, err
}

// Get Assets
func (srv AttributeGroupService) Get(params map[string]string) (AttributeGroupItems, error) {
    var a AttributeGroupItems

    err := srv.client.Get(attributeGroupBasePath, params, &a)

    return a, err
}

// Create Asset
func (srv AttributeGroupService) Create(createRequest interface{}) (IdsResponse, error) {
    id, err := srv.client.Create(attributeGroupBasePath, createRequest)

    return id, err
}

// Update Asset
func (srv AttributeGroupService) Update(updateRequest interface{}) (IdsResponse, error) {
    id, err := srv.client.Update(attributeGroupBasePath, updateRequest)

    return id, err
}

// Delete Asset
func (srv AttributeGroupService) Delete(deleteRequest interface{}) (IdsResponse, error) {
    id, err := srv.client.Delete(attributeGroupBasePath, deleteRequest)

    return id, err
}