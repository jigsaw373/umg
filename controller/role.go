package controller

import (
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/boof/umg/auth"
	"github.com/boof/umg/rbac/domains"
	"github.com/boof/umg/rbac/policies"
	"github.com/boof/umg/rbac/products"
	"github.com/boof/umg/rbac/properties"
	"github.com/boof/umg/rbac/roles"
	"github.com/boof/umg/rbac/users"
	"github.com/boof/umg/services"
	"github.com/boof/umg/util/request"
	"github.com/boof/umg/util/response"
)

// GetProperties returns all properties
func GetProperties(c echo.Context) error {
	user, err := auth.GetUser(c)
	if err != nil {
		return echo.ErrUnauthorized
	}

	domain := c.QueryParam("domain")
	var res []properties.Property

	if domain == "" {
		res, err = services.GetProperties(user)
	} else {
		res, err = services.GetPropertiesForDomain(user, domain)
	}

	if err != nil {
		return response.InternalErr(c, err.Error())
	}

	return response.OK(c, res)
}

// GetDomains returns all domains
func GetDomains(c echo.Context) error {
	res, err := services.GetAllDomains()
	if err != nil {
		return response.InternalErr(c, "unable to retrieve domains from database")
	}
	if res == nil {
		res = make([]domains.Domain, 0)
	}

	return response.OK(c, res)
}

func GetRoles(c echo.Context) error {
	format := c.QueryParam("format")

	order := c.QueryParam("order")
	if order != roles.Desc {
		order = roles.Asc
	}

	sortBy := c.QueryParam("sort")
	if !roles.IsValidSort(sortBy) {
		sortBy = roles.ID
	}

	if format == "with-policy" {
		return GetRolesWithNamedPolicies(c, order, sortBy)
	} else {
		return GetRawRoles(c, order, sortBy)
	}
}

func GetRawRoles(c echo.Context, order, sortBy string) error {
	// pagination parameters
	count, page := request.GetPagination(c)

	res, getErr := services.GetAllRoles(count, page, order, sortBy)
	if getErr != nil {
		return getErr.Echo(c)
	}
	if res == nil {
		res = make([]roles.Role, 0)
	}

	// set pagination header
	response.SetPageCountHeader(&c, roles.Pages(count))

	return response.OK(c, res)
}

func GetRolesWithNamedPolicies(c echo.Context, order, sortBy string) error {
	res := make([]map[string]interface{}, 0)

	// pagination parameters
	count, page := request.GetPagination(c)

	all, getErr := services.GetAllRoles(count, page, order, sortBy)
	if getErr != nil {
		return getErr.Echo(c)
	}

	for _, role := range all {
		out := make(map[string]interface{})
		out["id"] = role.ID
		out["name"] = role.Name

		policies, err := services.GetNamedPolicies(role.ID)
		if err == nil {
			out["policies"] = policies
		}

		res = append(res, out)
	}

	// set pagination header
	response.SetPageCountHeader(&c, roles.Pages(count))

	return response.OK(c, res)
}

func GetPolicies(c echo.Context) error {
	format := c.QueryParam("format")
	if format == "raw" {
		return GetRawPolicies(c)
	} else if format == "named" {
		return GetNamedPolicies(c)
	} else {
		return response.BadReq(c, "invalid format")
	}
}

func GetRawPolicies(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return response.BadReq(c, "bad request")
	}

	pols, err := (&policies.Policy{RoleID: id}).GetRolePolicies()
	if err != nil {
		return response.InternalErr(c, "unable to get policies")
	}

	if pols == nil {
		pols = make([]policies.Policy, 0)
	}

	return response.OK(c, pols)
}

func GetNamedPolicies(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return response.BadReq(c, "bad request")
	}

	policies, err := services.GetNamedPolicies(id)
	if err != nil {
		return response.InternalErr(c, "unable to get policies")
	}

	return response.OK(c, policies)
}

func GetProducts(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return response.BadReq(c, "bad request")
	}

	prods, err := products.GetProductsByDomain(id)
	if err != nil {
		return response.InternalErr(c, "unable to get products")
	}

	if prods == nil {
		prods = make([]products.Product, 0)
	}

	return response.OK(c, prods)
}

func AddRole(c echo.Context) error {
	role := new(roles.Role)
	if err := c.Bind(role); err != nil {
		return response.BadReq(c, "bad request")
	}

	err := role.Save()
	if err != nil {
		return response.BadReq(c, err.Error())
	}

	return response.Created(c, role)
}

func AddRoleWithPolicy(c echo.Context) error {
	type Req struct {
		Name     string            `json:"name"`
		Policies []policies.Policy `json:"policies"`
	}

	req := new(Req)
	if err := c.Bind(req); err != nil {
		return response.BadReq(c, "bad request")
	}

	role := &roles.Role{Name: req.Name}
	if err := role.Save(); err != nil {
		return response.BadReq(c, err.Error())
	}

	saved := make([]policies.Policy, 0)
	for _, policy := range req.Policies {
		policy.RoleID = role.ID
		err := policy.Save()
		if err != nil {
			for _, s := range saved {
				s.RemoveByID()
			}

			services.RemoveRoleByID(role.ID)

			return response.BadReq(c, err.Error())
		}

		saved = append(saved, policy)
	}

	return response.Created(c, echo.Map{"id": role.ID})
}

func AssignRole(c echo.Context) error {
	type Req struct {
		UserID int64 `json:"user_id"`
		RoleID int64 `json:"role_id"`
	}

	req := new(Req)
	if err := c.Bind(req); err != nil {
		return response.BadReq(c, "bad request")
	}

	user, err := (&users.User{ID: req.UserID}).GetByID()
	if err != nil {
		return response.NotFound(c, "user not found")
	}

	err = user.AssignRole(req.RoleID)
	if err != nil {
		return response.BadReq(c, err.Error())
	}

	return response.Done(c)
}

func DisallowRole(c echo.Context) error {
	type Req struct {
		UserID int64 `json:"user_id"`
		RoleID int64 `json:"role_id"`
	}

	req := new(Req)
	if err := c.Bind(req); err != nil {
		return response.BadReq(c, "bad request")
	}

	user, err := (&users.User{ID: req.UserID}).GetByID()
	if err != nil {
		return response.NotFound(c, "user not found")
	}

	err = user.DisallowRole(req.RoleID)
	if err != nil {
		return response.BadReq(c, err.Error())
	}

	return response.Done(c)
}

func AddDomain(c echo.Context) error {
	domain := new(domains.Domain)
	if err := c.Bind(domain); err != nil {
		return response.BadReq(c, "bad request")
	}

	err := domain.Save()
	if err != nil {
		return response.BadReq(c, err.Error())
	}

	return response.Created(c, domain)
}

func AddProduct(c echo.Context) error {
	prod := new(products.Product)
	if err := c.Bind(prod); err != nil {
		return response.BadReq(c, "bad request")
	}

	err := prod.Save()
	if err != nil {
		return response.BadReq(c, err.Error())
	}

	return response.Created(c, prod)
}

func AddDomPolicy(c echo.Context) error {
	policy := new(policies.Policy)
	if err := c.Bind(policy); err != nil {
		return response.BadReq(c, "bad request")
	}

	policy.Type = policies.DomPolicy
	err := policy.Save()
	if err != nil {
		return response.BadReq(c, err.Error())
	}

	return response.Created(c, echo.Map{"id": policy.ID})
}

func AddProdPolicy(c echo.Context) error {
	policy := new(policies.Policy)
	if err := c.Bind(policy); err != nil {
		return response.BadReq(c, "bad request")
	}

	policy.Type = policies.ProdPolicy
	err := policy.Save()
	if err != nil {
		return response.BadReq(c, err.Error())
	}

	return response.Created(c, echo.Map{"id": policy.ID})
}

func AddAllProdPolicy(c echo.Context) error {
	policy := new(policies.Policy)
	if err := c.Bind(policy); err != nil {
		return response.BadReq(c, "bad request")
	}

	policy.Type = policies.AllProdPolicy
	err := policy.Save()
	if err != nil {
		return response.BadReq(c, err.Error())
	}

	return response.Created(c, echo.Map{"id": policy.ID})
}

func DelPolicy(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return response.BadReq(c, "bad request")
	}

	err = (&policies.Policy{ID: int64(id)}).RemoveByID()
	if err != nil {
		return response.InternalErr(c, err.Error())
	}

	return response.Done(c)
}

func DelRole(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return response.BadReq(c, "bad request")
	}

	err = services.RemoveRoleByID(id)
	if err != nil {
		return response.InternalErr(c, err.Error())
	}

	return response.Done(c)
}

func EditRole(c echo.Context) error {
	role := new(roles.Role)
	if err := c.Bind(role); err != nil {
		return response.BadReq(c, "bad request")
	}

	err := role.Update()
	if err != nil {
		return response.BadReq(c, err.Error())
	}

	return response.Done(c)
}
