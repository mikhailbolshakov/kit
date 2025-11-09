package monitoring

import "gitlab.com/algmib/kit"

const (
	ErrCodePrometheusRegisterGoMetrics      = "MON-001"
	ErrCodePrometheusRegisterProcessMetrics = "MON-002"
	ErrCodePrometheusHttpServer             = "MON-003"
	ErrCodePrometheusInvalidPort            = "MON-004"
	ErrCodePrometheusRegisterAppMetrics     = "MON-005"
)

var (
	ErrPrometheusRegisterGoMetrics = func(cause error) error {
		return kit.NewAppErrBuilder(ErrCodePrometheusRegisterGoMetrics, "").Wrap(cause).Err()
	}
	ErrPrometheusRegisterProcessMetrics = func(cause error) error {
		return kit.NewAppErrBuilder(ErrCodePrometheusRegisterProcessMetrics, "").Wrap(cause).Err()
	}
	ErrPrometheusHttpServer = func(cause error) error {
		return kit.NewAppErrBuilder(ErrCodePrometheusHttpServer, "").Wrap(cause).Err()
	}
	ErrPrometheusInvalidPort = func(port string) error {
		return kit.NewAppErrBuilder(ErrCodePrometheusInvalidPort, "invalid port").F(kit.KV{"port": port}).Err()
	}
	ErrPrometheusRegisterAppMetrics = func(cause error) error {
		return kit.NewAppErrBuilder(ErrCodePrometheusRegisterAppMetrics, "").Wrap(cause).Err()
	}
)
