import { render, screen } from '@testing-library/react';
import React from 'react';
import { CloudWatchDatasource } from '../datasource';
import { VariableQuery, VariableQueryType } from '../types';
import VariableQueryEditor, { Props } from './VariableQueryEditor';

const props: Props = {
  onChange: (query) => {},
  query: {} as unknown as VariableQuery,
  datasource: {
    getRegions: async () => Promise.resolve([]),
    getNamespaces: async () => Promise.resolve([]),
    getVariables: () => [],
    getMetrics: async (namespace: string) => Promise.resolve([]),
    getDimensionKeys: async (namespace: string, region: string) => Promise.resolve([]),
  } as unknown as CloudWatchDatasource,
  onRunQuery: () => {},
};

describe('VariableEditor', () => {
  describe('and a new variable is created', () => {
    it('should trigger a query using the first query type in the array', (done) => {
      props.onChange = jest.fn();
      render(<VariableQueryEditor {...props} />);
      const querySelect = screen.queryByLabelText('Query Type');
      expect(querySelect).toBeInTheDocument();
      const regionSelect = screen.queryByLabelText('Region');
      expect(regionSelect).not.toBeInTheDocument();
    });
  });
  describe('and an existing variable is edited', () => {
    it('should trigger new query using the saved query type', (done) => {
      props.query = {
        queryType: VariableQueryType.Metrics,
        namespace: 'foo',
        region: 'bar',
        metricName: '',
        dimensionKey: '',
        dimensionFilters: '',
        ec2Filters: '',
        instanceID: '',
        attributeName: '',
        resourceType: '',
        tags: '',
        refId: '',
      };
      props.onChange = jest.fn();

      render(<VariableQueryEditor {...props} />);
      screen.debug();
      const regionSelect = screen.queryByLabelText('Region');
      expect(regionSelect).toBeInTheDocument();
      const namespaceSelect = screen.queryByLabelText('Namespace');
      expect(namespaceSelect).toBeInTheDocument();
    });
  });
  describe('and a different query type is selected', () => {});
});
