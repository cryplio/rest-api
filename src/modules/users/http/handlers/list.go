package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Nivl/go-rest-tools/types/apierror"
	"github.com/cryplio/rest-api/src/modules/users"

	"github.com/Nivl/go-rest-tools/paginator"
	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/guard"
	"github.com/Nivl/go-rest-tools/security/auth"
)

var listEndpoint = &router.Endpoint{
	Verb:    http.MethodGet,
	Path:    "/users",
	Handler: List,
	Guard: &guard.Guard{
		Auth:        guard.AdminAccess,
		ParamStruct: &ListParams{},
	},
}

// ListParams represents the params accepted by the Add endpoint
type ListParams struct {
	paginator.HandlerParams
	Sort string `from:"query" json:"sort" default:"created_at"`
}

// List is an endpoint used to list all Organization
func List(req router.HTTPRequest, deps *router.Dependencies) error {
	params := req.Params().(*ListParams)
	paginator := params.Paginator()

	sort, err := listParamsGetSort(params.Sort)
	if err != nil {
		return err
	}

	stmt := `
	SELECT profile.*, ` + auth.JoinUserSQL("users") + `
	FROM user_profiles profile
	JOIN users
	  ON users.id = profile.user_id
	WHERE users.deleted_at IS NULL
	ORDER BY %s
	OFFSET $1 LIMIT $2`
	stmt = fmt.Sprintf(stmt, sort)

	profiles := users.Profiles{}
	err = deps.DB.Select(&profiles, stmt, paginator.Offset(), paginator.Limit())
	if err != nil {
		return err
	}
	return req.Response().Ok(profiles.ExportPrivate())
}

func listParamsGetSort(sortFields string) (string, error) {
	sortableField := map[string]bool{
		"is_featured": true,
		"created_at":  true,
		"name":        true,
	}

	// We remove all empty sort (ex: "name,,created_at" will be "name,created_at")
	sortFields = strings.Replace(sortFields, ",,", ",", -1)
	// we also remove the leading and trailing coma
	sortFields = strings.Trim(sortFields, ",")

	// Set the default ordering if none is provided
	if sortFields == "" || sortFields == "," {
		sortFields = "is_featured,created_at"
	}

	sort := ""
	fields := strings.Split(sortFields, ",")
	for _, f := range fields {
		// Check the ordering
		order := "ASC"
		if f[0] == '-' {
			order = "DESC"
			f = f[1:]
		}
		// Make sure the field is valid
		if _, found := sortableField[f]; !found {
			return "", apierror.NewBadRequest("sort", "field %s is not valid", f)
		}
		sort += fmt.Sprintf("%s %s,", f, order)
	}
	return strings.TrimSuffix(sort, ","), nil
}
