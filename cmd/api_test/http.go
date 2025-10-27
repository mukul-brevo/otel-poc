package main

import (
	"context"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/DTSL/golang-libraries/errors"
	"github.com/DTSL/golang-libraries/httpjson"
	"go.opentelemetry.io/otel/metric"
)

type httpTestHandler struct {
	recordStatusCodeMetric func(ctx context.Context, statusCode int, method, path string)
}

func newhttpTestHandler(mp metric.MeterProvider) (*httpTestHandler, error) {
	statusCodeMetricRecorder, err := newStatusCodeMetric(mp)
	if err != nil {
		return nil, errors.Wrap(err, "new status code metric")
	}
	return &httpTestHandler{
		recordStatusCodeMetric: statusCodeMetricRecorder.recordMetric,
	}, nil
}

func (h *httpTestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := h.handle(ctx, w, r)
	if err != nil {
		httpjson.WriteResponse(ctx, w, http.StatusInternalServerError, nil)
	}
}

func (h *httpTestHandler) handle(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {
	// generate random duration between 0 and 2 seconds
	delay := time.Duration(rand.Intn(2001)) * time.Millisecond
	time.Sleep(delay)

	sc, resp := h.getResp(r)
	defer h.recordStatusCodeMetric(ctx, sc, r.Method, r.URL.Path)
	err = httpjson.WriteResponse(ctx, w, sc, resp)
	if err != nil {
		return errors.Wrap(err, "write response")
	}
	return nil
}

func (h *httpTestHandler) getResp(r *http.Request) (int, any) {
	defaultStatusCode := http.StatusOK
	defaultResp := map[string]string{"message": http.StatusText(defaultStatusCode)}
	statusCodeStr := r.URL.Query().Get("sc")
	if statusCodeStr == "" {
		return defaultStatusCode, defaultResp
	}
	statusCode, err := strconv.Atoi(statusCodeStr)
	if err != nil {
		return defaultStatusCode, defaultResp
	}
	if http.StatusText(statusCode) == "" {
		return defaultStatusCode, defaultResp
	}
	return statusCode, map[string]string{"message": http.StatusText(statusCode)}
}
