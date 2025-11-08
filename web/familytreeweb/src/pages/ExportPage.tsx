import { useState, useEffect } from 'react';
import {
  Flex,
  Card,
  Heading,
  Text,
  Button,
  Select,
  TextArea,
  Spinner,
  Callout,
  Table,
  Box,
} from '@radix-ui/themes';
import { DownloadIcon, CodeIcon, InfoCircledIcon } from '@radix-ui/react-icons';
import * as api from '../services/api';

export default function ExportPage() {
  const [persons, setPersons] = useState<api.Person[]>([]);
  const [selectedExport, setSelectedExport] = useState<string>('all');
  const [format, setFormat] = useState<string>('json');
  const [exportedData, setExportedData] = useState<string>('');
  const [isLoading, setIsLoading] = useState(false);
  const [isExporting, setIsExporting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadPersons();
  }, []);

  const loadPersons = async () => {
    try {
      setIsLoading(true);
      setError(null);
      const data = await api.listPersons();
      setPersons(Array.isArray(data) ? data : []);
    } catch (err) {
      console.error('Failed to load persons:', err);
      setError(err instanceof Error ? err.message : 'Failed to load persons');
      setPersons([]);
    } finally {
      setIsLoading(false);
    }
  };

  const handleExport = async () => {
    try {
      setIsExporting(true);
      setError(null);
      setExportedData('');

      let data: string;
      if (selectedExport === 'all') {
        data = await api.exportAllData(format);
      } else {
        data = await api.exportPerson(selectedExport, format);
      }

      setExportedData(data);
    } catch (err) {
      console.error('Failed to export:', err);
      setError(err instanceof Error ? err.message : 'Failed to export data');
    } finally {
      setIsExporting(false);
    }
  };

  const handleDownload = () => {
    if (!exportedData) return;

    const exportName = selectedExport === 'all' ? 'familytree-all' : `person-${selectedExport}`;
    const filename = `${exportName}-${format}.${format === 'json' ? 'json' : 'txt'}`;
    const blob = new Blob([exportedData], {
      type: format === 'json' ? 'application/json' : 'text/plain',
    });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = filename;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  };

  const handleCopy = async () => {
    if (!exportedData) return;

    try {
      await navigator.clipboard.writeText(exportedData);
    } catch (err) {
      console.error('Failed to copy to clipboard:', err);
    }
  };

  const formatDescriptions = {
    json: 'JSON format - Standard structured data format for easy parsing and import',
    gedcom: 'GEDCOM format - Standard genealogy data format (future support)',
    csv: 'CSV format - Comma-separated values for spreadsheet applications (future support)',
  };

  return (
    <Flex direction="column" gap="6">
      <Flex justify="between" align="center">
        <Heading size="8">Export Family Tree Data</Heading>
      </Flex>

      <Callout.Root color="blue">
        <Callout.Icon>
          <InfoCircledIcon />
        </Callout.Icon>
        <Callout.Text>
          Export your family tree data in various formats. You can export all data or select
          specific persons to export.
        </Callout.Text>
      </Callout.Root>

      <Card>
        <Flex direction="column" gap="4">
          <Heading size="5">Export Configuration</Heading>

          {isLoading ? (
            <Flex justify="center" py="4">
              <Spinner size="3" />
            </Flex>
          ) : error && persons.length === 0 ? (
            <Text color="red" size="3">
              {error}
            </Text>
          ) : (
            <Flex direction="column" gap="4">
              <Flex direction="column" gap="2">
                <Text size="2" weight="bold">
                  Select Data
                </Text>
                <Select.Root value={selectedExport} onValueChange={setSelectedExport}>
                  <Select.Trigger placeholder="Choose what to export..." />
                  <Select.Content>
                    <Select.Item value="all">All Family Tree Data</Select.Item>
                    <Select.Separator />
                    {persons.map(person => (
                      <Select.Item key={person.id} value={person.id || ''}>
                        {person.firstName} {person.lastName}
                      </Select.Item>
                    ))}
                  </Select.Content>
                </Select.Root>
              </Flex>

              <Flex direction="column" gap="2">
                <Text size="2" weight="bold">
                  Export Format
                </Text>
                <Select.Root value={format} onValueChange={setFormat}>
                  <Select.Trigger />
                  <Select.Content>
                    <Select.Item value="json">JSON</Select.Item>
                    <Select.Item value="gedcom" disabled>
                      GEDCOM (Coming Soon)
                    </Select.Item>
                    <Select.Item value="csv" disabled>
                      CSV (Coming Soon)
                    </Select.Item>
                  </Select.Content>
                </Select.Root>
                <Text size="1" color="gray">
                  {formatDescriptions[format as keyof typeof formatDescriptions]}
                </Text>
              </Flex>

              <Button size="3" onClick={handleExport} disabled={isExporting || !selectedExport}>
                <CodeIcon />
                {isExporting ? 'Exporting...' : 'Export'}
              </Button>
            </Flex>
          )}
        </Flex>
      </Card>

      {error && exportedData === '' && (
        <Callout.Root color="red">
          <Callout.Text>{error}</Callout.Text>
        </Callout.Root>
      )}

      {exportedData && (
        <Card>
          <Flex direction="column" gap="4">
            <Flex justify="between" align="center">
              <Heading size="5">Exported Configuration</Heading>
              <Flex gap="2">
                <Button size="2" variant="soft" onClick={handleCopy}>
                  Copy to Clipboard
                </Button>
                <Button size="2" onClick={handleDownload}>
                  <DownloadIcon /> Download
                </Button>
              </Flex>
            </Flex>

            <Box
              style={{
                border: '1px solid var(--gray-a6)',
                borderRadius: 'var(--radius-3)',
                overflow: 'hidden',
              }}
            >
              <TextArea
                value={exportedData}
                readOnly
                style={{
                  minHeight: '400px',
                  fontFamily: 'monospace',
                  fontSize: '13px',
                  resize: 'vertical',
                }}
              />
            </Box>

            <Text size="2" color="gray">
              {exportedData.split('\n').length} lines •{' '}
              {new Blob([exportedData]).size.toLocaleString()} bytes
            </Text>
          </Flex>
        </Card>
      )}

      <Card>
        <Flex direction="column" gap="4">
          <Heading size="5">Format Reference</Heading>
          <Table.Root variant="surface">
            <Table.Header>
              <Table.Row>
                <Table.ColumnHeaderCell>Format</Table.ColumnHeaderCell>
                <Table.ColumnHeaderCell>Use Case</Table.ColumnHeaderCell>
                <Table.ColumnHeaderCell>Description</Table.ColumnHeaderCell>
              </Table.Row>
            </Table.Header>
            <Table.Body>
              <Table.Row>
                <Table.Cell>
                  <Text weight="bold">JSON</Text>
                </Table.Cell>
                <Table.Cell>General Purpose</Table.Cell>
                <Table.Cell>Standard JSON format for easy import/export</Table.Cell>
              </Table.Row>
              <Table.Row>
                <Table.Cell>
                  <Text weight="bold">GEDCOM</Text>
                </Table.Cell>
                <Table.Cell>Genealogy Software</Table.Cell>
                <Table.Cell>Standard genealogy format (coming soon)</Table.Cell>
              </Table.Row>
              <Table.Row>
                <Table.Cell>
                  <Text weight="bold">CSV</Text>
                </Table.Cell>
                <Table.Cell>Spreadsheets</Table.Cell>
                <Table.Cell>Comma-separated values (coming soon)</Table.Cell>
              </Table.Row>
            </Table.Body>
          </Table.Root>
        </Flex>
      </Card>
    </Flex>
  );
}
