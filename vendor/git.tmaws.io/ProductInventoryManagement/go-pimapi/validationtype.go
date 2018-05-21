package pimapi

const (
    validationTypeBasePath = "metadata/validationtype"
)

type ValidationTypeService struct {
    client *Client
}

type ValidationType struct {
    Id              int64                   `json:"id,omitempty"`
    Name            string                  `json:"name,omitempty"`
    DisplayName     string                  `json:"displayName,omitempty"`
    ValidRegex      string                  `json:"ValidRegex,omitempty"`
    ErrorMessage    string                 `json:"errorMessage,omitempty"`
    Created         string                  `json:"created,omitempty"`
    Updated         string                  `json:"updated,omitempty"`
    User            string                  `json:"user,omitempty"`
}

type ValidationTypeItems struct {
    Items []ValidationType  `json:"items"`
    Page  *Page             `json:"page,omitempty"`
}

// Get Asset by Id
func (srv ValidationTypeService) ID(id string) (ValidationType, error) {
    var v ValidationType

    err := srv.client.ID(validationTypeBasePath, id, &v)

    return v, err
}

// Get Assets
func (srv ValidationTypeService) Get(params map[string]string) (ValidationTypeItems, error) {
    var v ValidationTypeItems

    err := srv.client.Get(validationTypeBasePath, params, &v)

    return v, err
}

// Create Asset
func (srv ValidationTypeService) Create(createRequest interface{}) (IdsResponse, error) {
    id, err := srv.client.Create(validationTypeBasePath, createRequest)

    return id, err
}

// Update Asset
func (srv ValidationTypeService) Update(updateRequest interface{}) (IdsResponse, error) {
    id, err := srv.client.Update(validationTypeBasePath, updateRequest)

    return id, err
}

// Delete Asset
func (srv ValidationTypeService) Delete(deleteRequest interface{}) (IdsResponse, error) {
    id, err := srv.client.Delete(validationTypeBasePath, deleteRequest)

    return id, err
}