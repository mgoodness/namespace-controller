package pimapi

const (
    productCodeBasePath = "metadata/productcode"
)

type ProductCodeService struct {
    client *Client
}

type ProductCode struct {
    Id          int64   `json:"id,omitempty"`
    Name        string  `json:"name,omitempty"`
    Deleted     bool    `json:"deleted,omitempty"`
    Created     string  `json:"created,omitempty"`
    Updated     string  `json:"updated,omitempty"`
    User        string  `json:"user,omitempty"`
}

type ProductCodeItems struct {
    Items []ProductCode     `json:"items"`
    Page  *Page             `json:"page,omitempty"`
}

// Get Asset by Id
func (srv ProductCodeService) ID(id string) (ProductCode, error) {
    var p ProductCode

    err := srv.client.ID(productCodeBasePath, id, &p)

    return p, err
}

// Get Assets
func (srv ProductCodeService) Get(params map[string]string) (ProductCodeItems, error) {
    var p ProductCodeItems

    err := srv.client.Get(productCodeBasePath, params, &p)

    return p, err
}

// Create Product Code
// This endpoint is unique in its response, so will not use the base client Create
func (srv ProductCodeService) Create() (ProductCode, error) {
    var p ProductCode

    req, err := srv.client.NewRequest("POST", productCodeBasePath, nil)
    if err != nil {
        return p, err
    }

    err = srv.client.Do(req, &p)

    return p, err
}

// Update Asset
func (srv ProductCodeService) Update(updateRequest interface{}) (IdsResponse, error) {
    id, err := srv.client.Update(productCodeBasePath, updateRequest)

    return id, err
}

// Delete Asset
func (srv ProductCodeService) Delete(deleteRequest interface{}) (IdsResponse, error) {
    id, err := srv.client.Delete(productCodeBasePath, deleteRequest)

    return id, err
}