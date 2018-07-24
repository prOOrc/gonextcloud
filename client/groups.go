package client

import (
	req "github.com/levigross/grequests"
	"github.com/partitio/gonextcloud/client/types"
	"net/http"
)

func (c *Client) GroupList() ([]string, error) {
	res, err := c.baseRequest(routes.groups, "", "", nil, http.MethodGet)
	if err != nil {
		return nil, err
	}
	var r types.GroupListResponse
	res.JSON(&r)
	if r.Ocs.Meta.Statuscode != 100 {
		e := types.ErrorFromMeta(r.Ocs.Meta)
		return nil, &e
	}
	return r.Ocs.Data.Groups, nil
}

func (c *Client) GroupUsers(name string) ([]string, error) {
	res, err := c.baseRequest(routes.groups, name, "", nil, http.MethodGet)
	if err != nil {
		return nil, err
	}
	var r types.UserListResponse
	res.JSON(&r)
	if r.Ocs.Meta.Statuscode != 100 {
		e := types.ErrorFromMeta(r.Ocs.Meta)
		return nil, &e
	}
	return r.Ocs.Data.Users, nil
}

func (c *Client) GroupSearch(search string) ([]string, error) {
	ro := &req.RequestOptions{
		Params: map[string]string{"search": search},
	}
	res, err := c.baseRequest(routes.groups, "", "", ro, http.MethodGet)
	if err != nil {
		return nil, err
	}
	var r types.GroupListResponse
	res.JSON(&r)
	if r.Ocs.Meta.Statuscode != 100 {
		e := types.ErrorFromMeta(r.Ocs.Meta)
		return nil, &e
	}
	return r.Ocs.Data.Groups, nil
}

func (c *Client) GroupCreate(name string) error {
	ro := &req.RequestOptions{
		Data: map[string]string{
			"groupid": name,
		},
	}
	if err := c.groupBaseRequest("", "", ro, http.MethodPost); err != nil {
		return err
	}
	return nil
}

func (c *Client) GroupDelete(name string) error {
	if err := c.groupBaseRequest(name, "", nil, http.MethodDelete); err != nil {
		return err
	}
	return nil
}

func (c *Client) GroupSubAdminList(name string) ([]string, error) {
	res, err := c.baseRequest(routes.groups, name, "subadmins", nil, http.MethodGet)
	if err != nil {
		return nil, err
	}
	var r types.UserListResponse
	res.JSON(&r)
	if r.Ocs.Meta.Statuscode != 100 {
		e := types.ErrorFromMeta(r.Ocs.Meta)
		return nil, &e
	}
	return r.Ocs.Data.Users, nil
}

func (c *Client) groupBaseRequest(name string, route string, ro *req.RequestOptions, method string) error {
	res, err := c.baseRequest(routes.groups, name, route, ro, method)
	if err != nil {
		return err
	}
	var r types.GroupListResponse
	res.JSON(&r)
	if r.Ocs.Meta.Statuscode != 100 {
		e := types.ErrorFromMeta(r.Ocs.Meta)
		return &e
	}
	return nil
}
