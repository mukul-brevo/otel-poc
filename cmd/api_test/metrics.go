package main

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type statusCodeMetric struct {
	sc    metric.Int64Counter
	sc2XX metric.Int64Counter
	sc3XX metric.Int64Counter
	sc4XX metric.Int64Counter
	sc5XX metric.Int64Counter
}

func newStatusCodeMetric(mp metric.MeterProvider) (m *statusCodeMetric, err error) {
	m = new(statusCodeMetric)
	m.sc, err = mp.Meter(appname).Int64Counter(
		fmt.Sprintf("%s.%s", appname, "http_response_total"),
		metric.WithDescription("Count of HTTP responses by status code"),
		metric.WithUnit("count"),
	)
	if err != nil {
		return nil, errors.Wrap(err, "2xx int64 counter")
	}
	m.sc2XX, err = mp.Meter(appname).Int64Counter(
		fmt.Sprintf("%s.%s", appname, "2XX"),
		metric.WithDescription("HTTP Traffic in 2XX requests count"),
		metric.WithUnit("count"),
	)
	if err != nil {
		return nil, errors.Wrap(err, "2xx int64 counter")
	}
	m.sc3XX, err = mp.Meter(appname).Int64Counter(
		fmt.Sprintf("%s.%s", appname, "3XX"),
		metric.WithDescription("HTTP Traffic in 3XX requests count"),
		metric.WithUnit("count"),
	)
	if err != nil {
		return nil, errors.Wrap(err, "3xx int64 counter")
	}
	m.sc4XX, err = mp.Meter(appname).Int64Counter(
		fmt.Sprintf("%s.%s", appname, "4XX"),
		metric.WithDescription("HTTP Traffic in 4XX requests count"),
		metric.WithUnit("count"),
	)
	if err != nil {
		return nil, errors.Wrap(err, "4xx int64 counter")
	}
	m.sc5XX, err = mp.Meter(appname).Int64Counter(
		fmt.Sprintf("%s.%s", appname, "5XX"),
		metric.WithDescription("HTTP Traffic in 5XX requests count"),
		metric.WithUnit("count"),
	)
	if err != nil {
		return nil, errors.Wrap(err, "5xx int64 counter")
	}
	return m, nil
}

func (h *statusCodeMetric) recordMetric(ctx context.Context, clientID int64, statusCode int, method, path string) {
	switch {
	case statusCode >= 200 && statusCode < 300:
		h.sc2XX.Add(ctx, 1, metric.WithAttributes(attribute.Int64("client_id", clientID)))
	case statusCode >= 300 && statusCode < 400:
		h.sc3XX.Add(ctx, 1, metric.WithAttributes(attribute.Int64("client_id", clientID)))
	case statusCode >= 400 && statusCode < 500:
		h.sc4XX.Add(ctx, 1, metric.WithAttributes(attribute.Int64("client_id", clientID)))
	case statusCode >= 500 && statusCode < 600:
		h.sc5XX.Add(ctx, 1, metric.WithAttributes(attribute.Int64("client_id", clientID)))
	default:
		return
	}
	h.sc.Add(ctx, 1, metric.WithAttributes(
		attribute.Int64("client_id", clientID),
		attribute.Int("status_code", statusCode),
		attribute.String("method", method),
		attribute.String("path", path),
	))
}
