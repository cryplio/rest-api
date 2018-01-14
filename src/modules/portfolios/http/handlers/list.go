package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Nivl/go-rest-tools/types/apierror"
	"github.com/cryplio/rest-api/src/modules/portfolios"

	"github.com/Nivl/go-rest-tools/paginator"
	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/guard"
)

var listEndpoint = &router.Endpoint{
	Verb:    http.MethodGet,
	Path:    "/portfolios",
	Handler: List,
	Guard: &guard.Guard{
		Auth:        guard.LoggedUserAccess,
		ParamStruct: &ListParams{},
	},
}

// ListParams represents the params accepted by the List endpoint
type ListParams struct {
	paginator.HandlerParams
	Sort string `from:"query" json:"sort" default:"name"`
}

// List is an endpoint used to list all portfolio of the requester
func List(req router.HTTPRequest, deps *router.Dependencies) error {
	params := req.Params().(*ListParams)
	paginator := params.Paginator()
	user := req.User()

	sort, err := listParamsGetSort(params.Sort)
	if err != nil {
		return err
	}

	stmt := `
	SELECT *
	FROM user_portfolios
	WHERE deleted_at IS NULL
		AND user_id=$1
	ORDER BY %s
	OFFSET $2 LIMIT $3`
	stmt = fmt.Sprintf(stmt, sort)

	profiles := portfolios.Portfolios{}
	err = deps.DB.Select(&profiles, stmt, user.ID, paginator.Offset(), paginator.Limit())
	if err != nil {
		return err
	}

	return req.Response().Ok(profiles.ExportPublic())
}

func listParamsGetSort(sortFields string) (string, error) {
	sortableField := map[string]bool{
		"created_at": true,
		"name":       true,
	}

	// We remove all empty sort (ex: "name,,created_at" will be "name,created_at")
	sortFields = strings.Replace(sortFields, ",,", ",", -1)
	// we also remove the leading and trailing coma
	sortFields = strings.Trim(sortFields, ",")

	// Set the default ordering if none is provided
	if sortFields == "" || sortFields == "," {
		sortFields = "name,created_at"
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
