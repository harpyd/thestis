// Package v1 provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.9.1 DO NOT EDIT.
package v1

import (
	"fmt"
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/go-chi/chi/v5"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Returns pipeline with such ID.
	// (GET /pipelines/{pipelineId})
	GetPipeline(w http.ResponseWriter, r *http.Request, pipelineId string)
	// Restart pipeline with such ID.
	// (PUT /pipelines/{pipelineId})
	RestartPipeline(w http.ResponseWriter, r *http.Request, pipelineId string)
	// Cancels pipeline with such ID.
	// (PUT /pipelines/{pipelineId}/canceled)
	CancelPipeline(w http.ResponseWriter, r *http.Request, pipelineId string)
	// Returns specification with such ID.
	// (GET /specifications/{specificationId})
	GetSpecification(w http.ResponseWriter, r *http.Request, specificationId string)
	// Returns test campaigns.
	// (GET /test-campaigns)
	GetTestCampaigns(w http.ResponseWriter, r *http.Request)
	// Creates test campaign for testing services logic using BDD specification style.
	// (POST /test-campaigns)
	CreateTestCampaign(w http.ResponseWriter, r *http.Request)
	// Removes test campaign with such ID.
	// (DELETE /test-campaigns/{testCampaignId})
	RemoveTestCampaign(w http.ResponseWriter, r *http.Request, testCampaignId string)
	// Returns test campaign with such ID.
	// (GET /test-campaigns/{testCampaignId})
	GetTestCampaign(w http.ResponseWriter, r *http.Request, testCampaignId string)
	// Asynchronously starts pipeline of test campaign's active specification.
	// (POST /test-campaigns/{testCampaignId}/pipeline)
	StartPipeline(w http.ResponseWriter, r *http.Request, testCampaignId string)
	// Returns pipeline history.
	// (GET /test-campaigns/{testCampaignId}/pipelines)
	GetPipelineHistory(w http.ResponseWriter, r *http.Request, testCampaignId string)
	// Loads specification to test campaign.
	// (POST /test-campaigns/{testCampaignId}/specification)
	LoadSpecification(w http.ResponseWriter, r *http.Request, testCampaignId string)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.HandlerFunc) http.HandlerFunc

// GetPipeline operation middleware
func (siw *ServerInterfaceWrapper) GetPipeline(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "pipelineId" -------------
	var pipelineId string

	err = runtime.BindStyledParameter("simple", false, "pipelineId", chi.URLParam(r, "pipelineId"), &pipelineId)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "pipelineId", Err: err})
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetPipeline(w, r, pipelineId)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// RestartPipeline operation middleware
func (siw *ServerInterfaceWrapper) RestartPipeline(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "pipelineId" -------------
	var pipelineId string

	err = runtime.BindStyledParameter("simple", false, "pipelineId", chi.URLParam(r, "pipelineId"), &pipelineId)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "pipelineId", Err: err})
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.RestartPipeline(w, r, pipelineId)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// CancelPipeline operation middleware
func (siw *ServerInterfaceWrapper) CancelPipeline(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "pipelineId" -------------
	var pipelineId string

	err = runtime.BindStyledParameter("simple", false, "pipelineId", chi.URLParam(r, "pipelineId"), &pipelineId)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "pipelineId", Err: err})
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CancelPipeline(w, r, pipelineId)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// GetSpecification operation middleware
func (siw *ServerInterfaceWrapper) GetSpecification(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "specificationId" -------------
	var specificationId string

	err = runtime.BindStyledParameter("simple", false, "specificationId", chi.URLParam(r, "specificationId"), &specificationId)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "specificationId", Err: err})
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetSpecification(w, r, specificationId)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// GetTestCampaigns operation middleware
func (siw *ServerInterfaceWrapper) GetTestCampaigns(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetTestCampaigns(w, r)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// CreateTestCampaign operation middleware
func (siw *ServerInterfaceWrapper) CreateTestCampaign(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CreateTestCampaign(w, r)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// RemoveTestCampaign operation middleware
func (siw *ServerInterfaceWrapper) RemoveTestCampaign(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "testCampaignId" -------------
	var testCampaignId string

	err = runtime.BindStyledParameter("simple", false, "testCampaignId", chi.URLParam(r, "testCampaignId"), &testCampaignId)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "testCampaignId", Err: err})
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.RemoveTestCampaign(w, r, testCampaignId)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// GetTestCampaign operation middleware
func (siw *ServerInterfaceWrapper) GetTestCampaign(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "testCampaignId" -------------
	var testCampaignId string

	err = runtime.BindStyledParameter("simple", false, "testCampaignId", chi.URLParam(r, "testCampaignId"), &testCampaignId)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "testCampaignId", Err: err})
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetTestCampaign(w, r, testCampaignId)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// StartPipeline operation middleware
func (siw *ServerInterfaceWrapper) StartPipeline(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "testCampaignId" -------------
	var testCampaignId string

	err = runtime.BindStyledParameter("simple", false, "testCampaignId", chi.URLParam(r, "testCampaignId"), &testCampaignId)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "testCampaignId", Err: err})
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.StartPipeline(w, r, testCampaignId)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// GetPipelineHistory operation middleware
func (siw *ServerInterfaceWrapper) GetPipelineHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "testCampaignId" -------------
	var testCampaignId string

	err = runtime.BindStyledParameter("simple", false, "testCampaignId", chi.URLParam(r, "testCampaignId"), &testCampaignId)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "testCampaignId", Err: err})
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetPipelineHistory(w, r, testCampaignId)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// LoadSpecification operation middleware
func (siw *ServerInterfaceWrapper) LoadSpecification(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "testCampaignId" -------------
	var testCampaignId string

	err = runtime.BindStyledParameter("simple", false, "testCampaignId", chi.URLParam(r, "testCampaignId"), &testCampaignId)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "testCampaignId", Err: err})
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.LoadSpecification(w, r, testCampaignId)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{})
}

type ChiServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r chi.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = chi.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/pipelines/{pipelineId}", wrapper.GetPipeline)
	})
	r.Group(func(r chi.Router) {
		r.Put(options.BaseURL+"/pipelines/{pipelineId}", wrapper.RestartPipeline)
	})
	r.Group(func(r chi.Router) {
		r.Put(options.BaseURL+"/pipelines/{pipelineId}/canceled", wrapper.CancelPipeline)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/specifications/{specificationId}", wrapper.GetSpecification)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/test-campaigns", wrapper.GetTestCampaigns)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/test-campaigns", wrapper.CreateTestCampaign)
	})
	r.Group(func(r chi.Router) {
		r.Delete(options.BaseURL+"/test-campaigns/{testCampaignId}", wrapper.RemoveTestCampaign)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/test-campaigns/{testCampaignId}", wrapper.GetTestCampaign)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/test-campaigns/{testCampaignId}/pipeline", wrapper.StartPipeline)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/test-campaigns/{testCampaignId}/pipelines", wrapper.GetPipelineHistory)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/test-campaigns/{testCampaignId}/specification", wrapper.LoadSpecification)
	})

	return r
}