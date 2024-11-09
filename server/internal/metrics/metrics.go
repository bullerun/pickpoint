package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var issuedOrder = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "orders_issued_total",
		Help: "Total number of orders issued.",
	},
)

func init() {
	prometheus.MustRegister(issuedOrder)
}

func IssueOrder(numberOrdersIssued float64) {
	issuedOrder.Add(numberOrdersIssued)
}
