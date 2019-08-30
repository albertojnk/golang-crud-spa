package endpoint

import (
	"golang-crud-spa/backend/datasource"
	"golang-crud-spa/backend/search"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

// CreateRepositoryRequest is the struct we will use to unmarshal the request
type CreateRepositoryRequest struct {
	Username string `json:"username" bson:"username"`
}

// CreateRepository is the handler that will create our repository in the database
func CreateRepository(c echo.Context) error {
	r := c.Request()

	defer r.Body.Close()

	reqData, err := Decode(r.Body, "create")
	if err != nil {
		log.Printf("something went wrong decoding body, err: %s", err)
		status, err := HandleErrors(err)
		c.JSON(status, err)
		return err
	}

	data := reqData.(CreateRepositoryRequest)

	respData, err := GetStarredRepositories(data.Username)

	err = datasource.CreateUserRepositories(data.Username, respData)
	if err != nil {
		log.Printf("error while creating, err: %s", err)
		status, err := HandleErrors(err)
		c.JSON(status, err)
		return err
	}

	user, err := datasource.ListUserRepositories(data.Username)
	if err != nil {
		log.Printf("error while listing, err: %s", err)
		status, err := HandleErrors(err)
		c.JSON(status, err)
		return err
	}

	if len(user.Repositories) > 0 {
		err = search.CreateIndex(indexName, user)
		if err != nil {
			log.Printf("Error creating index: %s", err)
			status, err := HandleErrors(err)
			c.JSON(status, err)
			return err
		}
	}

	// Encode response into json
	c.JSON(http.StatusCreated, user)

	return nil
}
