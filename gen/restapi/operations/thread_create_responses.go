// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "study/tp/db-forum-server/gen/models"
)

// ThreadCreateCreatedCode is the HTTP code returned for type ThreadCreateCreated
const ThreadCreateCreatedCode int = 201

/*ThreadCreateCreated Ветка обсуждения успешно создана.
Возвращает данные созданной ветки обсуждения.


swagger:response threadCreateCreated
*/
type ThreadCreateCreated struct {

	/*
	  In: Body
	*/
	Payload *models.Thread `json:"body,omitempty"`
}

// NewThreadCreateCreated creates ThreadCreateCreated with default headers values
func NewThreadCreateCreated() *ThreadCreateCreated {

	return &ThreadCreateCreated{}
}

// WithPayload adds the payload to the thread create created response
func (o *ThreadCreateCreated) WithPayload(payload *models.Thread) *ThreadCreateCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the thread create created response
func (o *ThreadCreateCreated) SetPayload(payload *models.Thread) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ThreadCreateCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// ThreadCreateNotFoundCode is the HTTP code returned for type ThreadCreateNotFound
const ThreadCreateNotFoundCode int = 404

/*ThreadCreateNotFound Автор ветки или форум не найдены.


swagger:response threadCreateNotFound
*/
type ThreadCreateNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewThreadCreateNotFound creates ThreadCreateNotFound with default headers values
func NewThreadCreateNotFound() *ThreadCreateNotFound {

	return &ThreadCreateNotFound{}
}

// WithPayload adds the payload to the thread create not found response
func (o *ThreadCreateNotFound) WithPayload(payload *models.Error) *ThreadCreateNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the thread create not found response
func (o *ThreadCreateNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ThreadCreateNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// ThreadCreateConflictCode is the HTTP code returned for type ThreadCreateConflict
const ThreadCreateConflictCode int = 409

/*ThreadCreateConflict Ветка обсуждения уже присутсвует в базе данных.
Возвращает данные ранее созданной ветки обсуждения.


swagger:response threadCreateConflict
*/
type ThreadCreateConflict struct {

	/*
	  In: Body
	*/
	Payload *models.Thread `json:"body,omitempty"`
}

// NewThreadCreateConflict creates ThreadCreateConflict with default headers values
func NewThreadCreateConflict() *ThreadCreateConflict {

	return &ThreadCreateConflict{}
}

// WithPayload adds the payload to the thread create conflict response
func (o *ThreadCreateConflict) WithPayload(payload *models.Thread) *ThreadCreateConflict {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the thread create conflict response
func (o *ThreadCreateConflict) SetPayload(payload *models.Thread) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ThreadCreateConflict) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(409)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
