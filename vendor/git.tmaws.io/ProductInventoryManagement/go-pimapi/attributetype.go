package pimapi

const (
    attributeTypeBasePath = "metadata/attributetype"
)

type AttributeTypeService struct {
    client *Client
}

type AttributeType struct {
    Id                  int64                   `json:"id,omitempty"`
    Name                string                  `json:"name,omitempty"`
    DisplayName         string                  `json:"displayName,omitempty"`
    Type                string                  `json:"type,omitempty"`
    Suggest             bool                    `json:"suggest,omitempty"`
    InputType           string                  `json:"inputType,omitempty"`
    AllowedValues       string                  `json:"allowedValues,omitempty"`
    ValidationTypeId    int64                   `json:"validationTypeId,omitempty"`
    Hint                string                  `json:"hint,omitempty"`
    Placeholder         string                  `json:"placeholder,omitempty"`
    Boost               int64                   `json:"boost,omitempty"`
    DefaultValue        int64                   `json:"defaultValue,omitempty"`
    Created             string                  `json:"created,omitempty"`
    Updated             string                  `json:"updated,omitempty"`
    User                string                  `json:"user,omitempty"`
}

type AttributeTypeItems struct {
    Items []AttributeType   `json:"items"`
    Page  *Page             `json:"page,omitempty"`
}

// Get Asset by Id
func (srv AttributeTypeService) ID(id string) (AttributeType, error) {
    var a AttributeType

    err := srv.client.ID(attributeTypeBasePath, id, &a)

    return a, err
}

// Get Assets
func (srv AttributeTypeService) Get(params map[string]string) (AttributeTypeItems, error) {
    var a AttributeTypeItems

    err := srv.client.Get(attributeTypeBasePath, params, &a)

    return a, err
}

// Create Asset
func (srv AttributeTypeService) Create(createRequest interface{}) (IdsResponse, error) {
    id, err := srv.client.Create(attributeTypeBasePath, createRequest)

    return id, err
}

// Update Asset
func (srv AttributeTypeService) Update(updateRequest interface{}) (IdsResponse, error) {
    id, err := srv.client.Update(attributeTypeBasePath, updateRequest)

    return id, err
}

// Delete Asset
func (srv AttributeTypeService) Delete(deleteRequest interface{}) (IdsResponse, error) {
    id, err := srv.client.Delete(attributeTypeBasePath, deleteRequest)

    return id, err
}