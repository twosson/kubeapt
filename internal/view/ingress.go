package view

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/twosson/kubeapt/internal/cluster"
	"github.com/twosson/kubeapt/internal/content"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"strings"
)

type IngressDetails struct {
}

var _ View = (*IngressDetails)(nil)

func (ing *IngressDetails) Content(ctx context.Context, object runtime.Object, clusterClient cluster.ClientInterface) ([]content.Content, error) {
	ingress, ok := object.(*v1beta1.Ingress)
	if !ok {
		return nil, errors.Errorf("expected object to be Ingress, it was %T", object)
	}

	return []content.Content{
		ingressTLSTable(ingress),
		ingressRulesTable(ingress),
	}, nil
}

func ingressTLSTable(ingress *v1beta1.Ingress) *content.Table {
	table := content.NewTable("TLS")

	table.Columns = []content.TableColumn{
		tableCol("Secret"),
		tableCol("Hosts"),
	}

	for _, tls := range ingress.Spec.TLS {
		table.AddRow(content.TableRow{
			"Secret": content.NewStringText(tls.SecretName),
			"Hosts":  content.NewStringText(strings.Join(tls.Hosts, ", ")),
		})
	}

	return &table
}

func ingressRulesTable(ingress *v1beta1.Ingress) *content.Table {
	table := content.NewTable("Rules")

	table.Columns = []content.TableColumn{
		tableCol("Host"),
		tableCol("Paths"),
	}

	for _, rule := range ingress.Spec.Rules {
		var paths []string

		if rule.HTTP != nil {
			for _, path := range rule.HTTP.Paths {
				pathRoute := path.Path
				backend := fmt.Sprintf("%s:%s", path.Backend.ServiceName, path.Backend.ServicePort.String())
				paths = append(paths, fmt.Sprintf("%s (%s)", pathRoute, backend))
			}

			table.AddRow(content.TableRow{
				"Host":  content.NewStringText(rule.Host),
				"Paths": content.NewStringText(strings.Join(paths, ", ")),
			})
		}

	}

	return &table
}
