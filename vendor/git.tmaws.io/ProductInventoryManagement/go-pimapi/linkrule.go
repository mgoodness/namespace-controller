package pimapi

const (
    linkRuleBasePath = "metadata/linkrule"
)

type LinkRuleService struct {
    client *Client
}

type LinkRule struct {
    Id              int64                   `json:"id,omitempty"`
    ParentAssetId   int64                   `json:"parentAssetId,omitempty"`
    ChildAssetId    int64                   `json:"childAssetId,omitempty"`
    Name            string                  `json:"name,omitempty"`
    Cardinality     string                  `json:"cardinality,omitempty"`
    MustLink        bool                    `json:"mustLink,omitempty"`
    StickyLink      bool                    `json:"stickyLink,omitempty"`
    Inherit         bool                    `json:"inherit,omitempty"`
}

type LinkRuleItems struct {
    Items []LinkRule    `json:"items"`
    Page  *Page         `json:"page,omitempty"`
}

// Get Asset by Id
func (srv LinkRuleService) ID(id string) (LinkRule, error) {
    var lr LinkRule

    err := srv.client.ID(linkRuleBasePath, id, &lr)

    return lr, err
}

// Get Assets
func (srv LinkRuleService) Get(params map[string]string) (LinkRuleItems, error) {
    var lr LinkRuleItems

    err := srv.client.Get(linkRuleBasePath, params, &lr)

    return lr, err
}

// Create Asset
func (srv LinkRuleService) Create(createRequest interface{}) (IdsResponse, error) {
    id, err := srv.client.Create(linkRuleBasePath, createRequest)

    return id, err
}

// Update Asset
func (srv LinkRuleService) Update(updateRequest interface{}) (IdsResponse, error) {
    id, err := srv.client.Update(linkRuleBasePath, updateRequest)

    return id, err
}

// Delete Asset
func (srv LinkRuleService) Delete(deleteRequest interface{}) (IdsResponse, error) {
    id, err := srv.client.Delete(linkRuleBasePath, deleteRequest)

    return id, err
}