package services

import (
	"context"
	"fmt"

	"github.com/grafana/grafana/pkg/services/ngalert/models"
)

const defaultAlertmanagerConfigJSON = `
{
	"template_files": null,
	"alertmanager_config": {
		"route": {
			"receiver": "grafana-default-email",
			"group_by": [
				"..."
			],
			"routes": [{
				"receiver": "grafana-default-email",
				"object_matchers": [["a", "=", "b"]]
			}]
		},
		"templates": null,
		"receivers": [{
			"name": "grafana-default-email",
			"grafana_managed_receiver_configs": [{
				"uid": "",
				"name": "email receiver",
				"type": "email",
				"disableResolveMessage": false,
				"settings": {
					"addresses": "\u003cexample@email.com\u003e"
				},
				"secureFields": {}
			}]
		}, {
			"name": "a new receiver",
			"grafana_managed_receiver_configs": [{
				"uid": "",
				"name": "email receiver",
				"type": "email",
				"disableResolveMessage": false,
				"settings": {
					"addresses": "\u003canother@email.com\u003e"
				},
				"secureFields": {}
			}]
		}]
	}
}
`

type fakeAMConfigStore struct {
	config models.AlertConfiguration
}

func newFakeAMConfigStore() *fakeAMConfigStore {
	return &fakeAMConfigStore{
		config: models.AlertConfiguration{
			AlertmanagerConfiguration: defaultAlertmanagerConfigJSON,
			ConfigurationVersion:      "v1",
			Default:                   true,
			OrgID:                     1,
		},
	}
}

func (f *fakeAMConfigStore) GetLatestAlertmanagerConfiguration(ctx context.Context, query *models.GetLatestAlertmanagerConfigurationQuery) error {
	query.Result = &f.config
	query.Result.OrgID = query.OrgID
	return nil
}

func (f *fakeAMConfigStore) SaveAlertmanagerConfiguration(ctx context.Context, cmd *models.SaveAlertmanagerConfigurationCmd) error {
	f.config = models.AlertConfiguration{
		AlertmanagerConfiguration: cmd.AlertmanagerConfiguration,
		ConfigurationVersion:      cmd.ConfigurationVersion,
		Default:                   cmd.Default,
		OrgID:                     cmd.OrgID,
	}
	return nil
}

type fakeProvisioningStore struct {
	records map[int64]map[string]models.Provenance
}

func newFakeProvisioningStore() *fakeProvisioningStore {
	return &fakeProvisioningStore{
		records: map[int64]map[string]models.Provenance{},
	}
}

func (f *fakeProvisioningStore) GetProvenance(ctx context.Context, o models.ProvisionableInOrg) (models.Provenance, error) {
	if val, ok := f.records[o.ResourceOrgID()]; ok {
		if prov, ok := val[o.ResourceID()]; ok {
			return prov, nil
		}
	}
	return models.ProvenanceNone, nil
}

func (f *fakeProvisioningStore) SetProvenance(ctx context.Context, o models.ProvisionableInOrg, p models.Provenance) error {
	orgID := o.ResourceOrgID()
	if _, ok := f.records[orgID]; !ok {
		f.records[orgID] = map[string]models.Provenance{}
	}
	f.records[orgID][o.ResourceID()] = p
	return nil
}

type failingProvisioningStore struct{}

func (f *failingProvisioningStore) GetProvenance(ctx context.Context, o models.ProvisionableInOrg) (models.Provenance, error) {
	return models.ProvenanceNone, fmt.Errorf("failed to store provenance")
}

func (f *failingProvisioningStore) SetProvenance(ctx context.Context, o models.ProvisionableInOrg, p models.Provenance) error {
	return fmt.Errorf("failed to set provenance")
}

type nopTransactionManager struct{}

func newNopTransactionManager() *nopTransactionManager {
	return &nopTransactionManager{}
}

func (n *nopTransactionManager) InTransaction(ctx context.Context, work func(ctx context.Context) error) error {
	return work(ctx)
}
