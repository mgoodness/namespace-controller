package pimapi

const (
    externalSourceBasePath = "externalsource"
)

type ExternalSourceService struct {
    client *Client
}

type ExternalSource struct {
    Id              int64                   `json:"id,omitempty"`
    Name            string                  `json:"name,omitempty"`
    Created         string                  `json:"created,omitempty"`
    Updated         string                  `json:"updated,omitempty"`
    User            string                  `json:"user,omitempty"`
}

type ExternalSourceItems struct {
    Items []ExternalSource  `json:"items"`
    Page  *Page             `json:"page,omitempty"`
}

// Get ExternalSource by Id
func (srv ExternalSourceService) ID(id string) (ExternalSource, error) {
    var a ExternalSource

    err := srv.client.ID(externalSourceBasePath, id, &a)

    return a, err
}

// Get ExternalSources
func (srv ExternalSourceService) Get(params map[string]string) (ExternalSourceItems, error) {
    var a ExternalSourceItems

    err := srv.client.Get(externalSourceBasePath, params, &a)

    return a, err
}

// Create ExternalSource
func (srv ExternalSourceService) Create(createRequest ExternalSourceItems) (IdsResponse, error) {
    id, err := srv.client.Create(externalSourceBasePath, createRequest)

    return id, err
}

// Update ExternalSource
func (srv ExternalSourceService) Update(updateRequest ExternalSourceItems) (IdsResponse, error) {
    id, err := srv.client.Update(externalSourceBasePath, updateRequest)

    return id, err
}

// Delete ExternalSource
func (srv ExternalSourceService) Delete(deleteRequest IdsRequest) (IdsResponse, error) {
    id, err := srv.client.Delete(externalSourceBasePath, deleteRequest)

    return id, err
}
