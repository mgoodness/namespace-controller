package pimapi

const (
    linkBasePath = "link"
)

type LinkService struct {
    client *Client
}

type Link struct {
    Id              int64                   `json:"id,omitempty"`
    LinkRuleId      int64                   `json:"linkRuleId,omitempty"`
    ParentAssetId   int64                   `json:"parentAssetId,omitempty"`
    ChildAssetId    int64                   `json:"childAssetId,omitempty"`
    ParentAssetName int64                   `json:"parentAssetName,omitempty"`
    ChildAssetName  int64                   `json:"childAssetName,omitempty"`
    Created         string                  `json:"created,omitempty"`
    Updated         string                  `json:"updated,omitempty"`
    User            string                  `json:"user,omitempty"`
}

type LinkItems struct {
    Items []Link   `json:"items"`
    Page  *Page    `json:"page,omitempty"`
}

// Get Link by Id
func (srv LinkService) ID(id string) (Link, error) {
    var l Link

    err := srv.client.ID(linkBasePath, id, &l)

    return l, err
}

// Get Links
func (srv LinkService) Get(params map[string]string) (LinkItems, error) {
    var l LinkItems

    err := srv.client.Get(linkBasePath, params, &l)

    return l, err
}

// Create Link
func (srv LinkService) Create(createRequest LinkItems) (IdsResponse, error) {
    id, err := srv.client.Create(linkBasePath, createRequest)

    return id, err
}

// Update Link
func (srv LinkService) Update(updateRequest LinkItems) (IdsResponse, error) {
    id, err := srv.client.Update(linkBasePath, updateRequest)

    return id, err
}

// Delete Link
func (srv LinkService) Delete(deleteRequest IdsRequest) (IdsResponse, error) {
    id, err := srv.client.Delete(linkBasePath, deleteRequest)

    return id, err
}
