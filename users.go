package gonextcloud

import (
	"encoding/json"
	"github.com/fatih/structs"
	req "github.com/levigross/grequests"
	"github.com/partitio/gonextcloud/types"
	"net/http"
	"path"
	"strconv"
	"strings"
	"sync"
)

// UserList return the Nextcloud'user list
func (c *Client) UserList() ([]string, error) {
	res, err := c.baseRequest(routes.users, "", "", nil, http.MethodGet)
	//res, err := c.session.Get(u.String(), nil)
	if err != nil {
		return nil, err
	}
	var r types.UserListResponse
	res.JSON(&r)
	return r.Ocs.Data.Users, nil
}

// User return the details about the specified user
func (c *Client) User(name string) (*types.User, error) {
	if name == "" {
		return nil, &types.APIError{Message: "name cannot be empty"}
	}
	res, err := c.baseRequest(routes.users, name, "", nil, http.MethodGet)
	if err != nil {
		return nil, err
	}
	var r types.UserResponse
	js := res.String()
	// Nextcloud does not encode JSON properly
	js = reformatJSON(js)
	if err := json.Unmarshal([]byte(js), &r); err != nil {
		return nil, err
	}
	return &r.Ocs.Data, nil
}

// UserSearch returns the users whose name match the search string
func (c *Client) UserSearch(search string) ([]string, error) {
	ro := &req.RequestOptions{
		Params: map[string]string{"search": search},
	}
	res, err := c.baseRequest(routes.users, "", "", ro, http.MethodGet)
	if err != nil {
		return nil, err
	}
	var r types.UserListResponse
	res.JSON(&r)
	return r.Ocs.Data.Users, nil
}

// UserCreate create a new user
func (c *Client) UserCreate(username string, password string, user *types.User) error {
	ro := &req.RequestOptions{
		Data: map[string]string{
			"userid":   username,
			"password": password,
		},
	}
	if err := c.userBaseRequest("", "", ro, http.MethodPost); err != nil {
		return err
	}
	if user == nil {
		return nil
	}
	return c.UserUpdate(user)
}

//UserDelete delete the user
func (c *Client) UserDelete(name string) error {
	return c.userBaseRequest(name, "", nil, http.MethodDelete)
}

//UserEnable enables the user
func (c *Client) UserEnable(name string) error {
	ro := &req.RequestOptions{
		Data: map[string]string{},
	}
	return c.userBaseRequest(name, "enable", ro, http.MethodPut)
}

//UserDisable disables the user
func (c *Client) UserDisable(name string) error {
	ro := &req.RequestOptions{
		Data: map[string]string{},
	}
	return c.userBaseRequest(name, "disable", ro, http.MethodPut)
}

//UserSendWelcomeEmail (re)send the welcome mail to the user (return an error if the user has not configured his email)
func (c *Client) UserSendWelcomeEmail(name string) error {
	return c.userBaseRequest(name, "welcome", nil, http.MethodPost)
}

//UserUpdate takes a *types.User struct to update the user's information
func (c *Client) UserUpdate(user *types.User) error {
	m := structs.Map(user)
	errs := make(chan types.UpdateError)
	var wg sync.WaitGroup
	for k := range m {
		if !ignoredUserField(k) && m[k].(string) != "" {
			wg.Add(1)
			go func(key string, value string) {
				defer wg.Done()
				if err := c.userUpdateAttribute(user.ID, strings.ToLower(key), value); err != nil {
					errs <- types.UpdateError{
						Field: key,
						Error: err,
					}
				}
			}(k, m[k].(string))
		}
	}
	go func() {
		wg.Wait()
		close(errs)
	}()
	return types.NewUpdateError(errs)
}

//UserUpdateEmail update the user's email
func (c *Client) UserUpdateEmail(name string, email string) error {
	return c.userUpdateAttribute(name, "email", email)
}

//UserUpdateDisplayName update the user's display name
func (c *Client) UserUpdateDisplayName(name string, displayName string) error {
	return c.userUpdateAttribute(name, "displayname", displayName)
}

//UserUpdatePhone update the user's phone
func (c *Client) UserUpdatePhone(name string, phone string) error {
	return c.userUpdateAttribute(name, "phone", phone)
}

//UserUpdateAddress update the user's address
func (c *Client) UserUpdateAddress(name string, address string) error {
	return c.userUpdateAttribute(name, "address", address)
}

//UserUpdateWebSite update the user's website
func (c *Client) UserUpdateWebSite(name string, website string) error {
	return c.userUpdateAttribute(name, "website", website)
}

//UserUpdateTwitter update the user's twitter
func (c *Client) UserUpdateTwitter(name string, twitter string) error {
	return c.userUpdateAttribute(name, "twitter", twitter)
}

//UserUpdatePassword update the user's password
func (c *Client) UserUpdatePassword(name string, password string) error {
	return c.userUpdateAttribute(name, "password", password)
}

//UserUpdateQuota update the user's quota (bytes)
func (c *Client) UserUpdateQuota(name string, quota int) error {
	return c.userUpdateAttribute(name, "quota", strconv.Itoa(quota))
}

//UserGroupList lists the user's groups
func (c *Client) UserGroupList(name string) ([]string, error) {
	res, err := c.baseRequest(routes.users, name, "groups", nil, http.MethodGet)
	if err != nil {
		return nil, err
	}
	var r types.GroupListResponse
	res.JSON(&r)
	return r.Ocs.Data.Groups, nil
}

//UserGroupAdd adds a the user to the group
func (c *Client) UserGroupAdd(name string, group string) error {
	ro := &req.RequestOptions{
		Data: map[string]string{
			"groupid": group,
		},
	}
	return c.userBaseRequest(name, "groups", ro, http.MethodPost)
}

//UserGroupRemove removes the user from the group
func (c *Client) UserGroupRemove(name string, group string) error {
	ro := &req.RequestOptions{
		Data: map[string]string{
			"groupid": group,
		},
	}
	return c.userBaseRequest(name, "groups", ro, http.MethodDelete)
}

//UserGroupPromote promotes the user as group admin
func (c *Client) UserGroupPromote(name string, group string) error {
	ro := &req.RequestOptions{
		Data: map[string]string{
			"groupid": group,
		},
	}
	return c.userBaseRequest(name, "subadmins", ro, http.MethodPost)
}

//UserGroupDemote demotes the user
func (c *Client) UserGroupDemote(name string, group string) error {
	ro := &req.RequestOptions{
		Data: map[string]string{
			"groupid": group,
		},
	}
	return c.userBaseRequest(name, "subadmins", ro, http.MethodDelete)
}

//UserGroupSubAdminList lists the groups where he is subadmin
func (c *Client) UserGroupSubAdminList(name string) ([]string, error) {
	if !c.loggedIn() {
		return nil, unauthorized
	}
	u := c.baseURL.ResolveReference(routes.users)
	u.Path = path.Join(u.Path, name, "subadmins")
	res, err := c.session.Get(u.String(), nil)
	if err != nil {
		return nil, err
	}
	var r types.BaseResponse
	res.JSON(&r)
	return r.Ocs.Data, nil
}

func (c *Client) userUpdateAttribute(name string, key string, value string) error {
	ro := &req.RequestOptions{
		Data: map[string]string{
			"key":   key,
			"value": value,
		},
	}
	return c.userBaseRequest(name, "", ro, http.MethodPut)
}

func (c *Client) userBaseRequest(name string, route string, ro *req.RequestOptions, method string) error {
	_, err := c.baseRequest(routes.users, name, route, ro, method)
	return err
}

func ignoredUserField(key string) bool {
	keys := []string{"ID", "Quota", "Enabled", "Groups", "Language"}
	for _, k := range keys {
		if key == k {
			return true
		}
	}
	return false
}
