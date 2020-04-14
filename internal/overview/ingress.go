package overview

import (
	"context"
	"github.com/pkg/errors"
	"github.com/twosson/kubeapt/internal/content"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"strings"
)

type IngressSummary struct{}

var _ View = (*IngressSummary)(nil)

func NewIngressSummary() *IngressSummary {
	return &IngressSummary{}
}

func (js *IngressSummary) Content(ctx context.Context, object runtime.Object, c Cache) ([]content.Content, error) {
	ingress, err := retrieveIngress(object)
	if err != nil {
		return nil, err
	}

	detail, err := printIngressSummary(ingress)
	if err != nil {
		return nil, err
	}

	summary := content.NewSummary("Details", []content.Section{detail})
	return []content.Content{
		&summary,
	}, nil
}

type IngressDetails struct{}

var _ View = (*IngressDetails)(nil)

func NewIngressDetails() *IngressDetails {
	return &IngressDetails{}
}

func (ing *IngressDetails) Content(ctx context.Context, object runtime.Object, c Cache) ([]content.Content, error) {
	ingress, err := retrieveIngress(object)
	if err != nil {
		return nil, err
	}

	return []content.Content{
		ingressTLSTable(ingress),
		ingressRulesTable(ingress),
	}, nil
}

func ingressTLSTable(ingress *v1beta1.Ingress) *content.Table {
	table := content.NewTable("TLS", "TLS is not configured for this Ingress")

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
	table := content.NewTable("Rules", "Rules are not configured for this Ingress")

	table.Columns = tableCols("Host", "Path", "Backend")

	for _, rule := range ingress.Spec.Rules {
		if rule.HTTP != nil {
			for _, path := range rule.HTTP.Paths {
				backendText := backendStringer(&path.Backend)
				table.AddRow(content.TableRow{
					"Host":    content.NewStringText(rule.Host),
					"Path":    content.NewStringText(path.Path),
					"Backend": content.NewLinkText(backendText, gvkPath("v1", "Service", path.Backend.ServiceName)),
				})
			}
		}
	}

	return &table
}

func retrieveIngress(object runtime.Object) (*v1beta1.Ingress, error) {
	ingress, ok := object.(*v1beta1.Ingress)
	if !ok {
		return nil, errors.Errorf("expected object to be an Ingress, it was %T", object)
	}

	return ingress, nil
}
