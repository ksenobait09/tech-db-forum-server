// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"io"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	strfmt "github.com/go-openapi/strfmt"

	models "study/tp/db-forum-server/gen/models"
)

// NewPostsCreateParams creates a new PostsCreateParams object
// no default values defined in spec.
func NewPostsCreateParams() PostsCreateParams {

	return PostsCreateParams{}
}

// PostsCreateParams contains all the bound params for the posts create operation
// typically these are obtained from a http.Request
//
// swagger:parameters postsCreate
type PostsCreateParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*Список создаваемых постов.
	  Required: true
	  In: body
	*/
	Posts models.Posts
	/*Идентификатор ветки обсуждения.
	  Required: true
	  In: path
	*/
	SlugOrID string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewPostsCreateParams() beforehand.
func (o *PostsCreateParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	if runtime.HasBody(r) {
		defer r.Body.Close()
		var body models.Posts
		if err := route.Consumer.Consume(r.Body, &body); err != nil {
			if err == io.EOF {
				res = append(res, errors.Required("posts", "body"))
			} else {
				res = append(res, errors.NewParseError("posts", "body", "", err))
			}
		} else {
			// validate body object
			if err := body.Validate(route.Formats); err != nil {
				res = append(res, err)
			}

			if len(res) == 0 {
				o.Posts = body
			}
		}
	} else {
		res = append(res, errors.Required("posts", "body"))
	}
	rSlugOrID, rhkSlugOrID, _ := route.Params.GetOK("slug_or_id")
	if err := o.bindSlugOrID(rSlugOrID, rhkSlugOrID, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindSlugOrID binds and validates parameter SlugOrID from path.
func (o *PostsCreateParams) bindSlugOrID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.SlugOrID = raw

	return nil
}
