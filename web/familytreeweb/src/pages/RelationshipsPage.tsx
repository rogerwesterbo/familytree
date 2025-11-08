import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import {
  Flex,
  Card,
  Heading,
  Text,
  Button,
  Table,
  Badge,
  Spinner,
  TextField,
  Box,
} from '@radix-ui/themes';
import { PlusIcon, MagnifyingGlassIcon, ReloadIcon } from '@radix-ui/react-icons';
import * as api from '../services/api';

export default function RelationshipsPage() {
  const [relationships, setRelationships] = useState<api.Relationship[]>([]);
  const [filteredRelationships, setFilteredRelationships] = useState<api.Relationship[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [filter, setFilter] = useState('');
  const [currentPage, setCurrentPage] = useState(1);
  const itemsPerPage = 10;

  useEffect(() => {
    loadRelationships();
  }, []);

  useEffect(() => {
    if (filter.trim()) {
      const filtered = relationships.filter(
        rel =>
          rel.type.toLowerCase().includes(filter.toLowerCase()) ||
          rel.fromPersonId.toLowerCase().includes(filter.toLowerCase()) ||
          rel.toPersonId.toLowerCase().includes(filter.toLowerCase())
      );
      setFilteredRelationships(filtered);
    } else {
      setFilteredRelationships(relationships);
    }
    setCurrentPage(1);
  }, [filter, relationships]);

  const loadRelationships = async () => {
    try {
      setIsLoading(true);
      setError(null);
      const data = await api.listRelationships();
      const validData = Array.isArray(data) ? data : [];
      setRelationships(validData);
      setFilteredRelationships(validData);
    } catch (err) {
      console.error('Failed to load relationships:', err);
      setError(err instanceof Error ? err.message : 'Failed to load relationships');
      setRelationships([]);
      setFilteredRelationships([]);
    } finally {
      setIsLoading(false);
    }
  };

  // Pagination
  const totalPages = Math.ceil(filteredRelationships.length / itemsPerPage);
  const startIndex = (currentPage - 1) * itemsPerPage;
  const endIndex = startIndex + itemsPerPage;
  const currentRelationships = filteredRelationships.slice(startIndex, endIndex);

  const getRelationshipTypeBadgeColor = (type: string): 'blue' | 'green' | 'orange' | 'purple' => {
    const colors: Record<string, 'blue' | 'green' | 'orange' | 'purple'> = {
      parent: 'blue',
      child: 'blue',
      spouse: 'green',
      sibling: 'orange',
      partner: 'purple',
    };
    return colors[type.toLowerCase()] || 'blue';
  };

  return (
    <Flex direction="column" gap="6">
      <Flex justify="between" align="center">
        <Heading size="8">Relationships</Heading>
        <Flex gap="2">
          <Button size="3" variant="soft" onClick={loadRelationships}>
            <ReloadIcon /> Refresh
          </Button>
          <Button size="3">
            <PlusIcon /> Add Relationship
          </Button>
        </Flex>
      </Flex>

      <Card>
        <Flex direction="column" gap="4">
          <Flex justify="between" align="center">
            <Text size="2" color="gray">
              Manage relationships between persons in the family tree
            </Text>
            <Box style={{ width: '300px' }}>
              <TextField.Root
                size="2"
                placeholder="Filter relationships..."
                value={filter}
                onChange={e => setFilter(e.target.value)}
              >
                <TextField.Slot>
                  <MagnifyingGlassIcon height="14" width="14" />
                </TextField.Slot>
              </TextField.Root>
            </Box>
          </Flex>

          {isLoading && (
            <Flex justify="center" py="8">
              <Spinner size="3" />
            </Flex>
          )}

          {error && (
            <Text color="red" size="3">
              {error}
            </Text>
          )}

          {!isLoading && !error && currentRelationships.length === 0 && (
            <Box py="8">
              <Text color="gray" align="center">
                {filter
                  ? 'No relationships match your filter.'
                  : 'No relationships configured yet.'}
              </Text>
            </Box>
          )}

          {!isLoading && !error && currentRelationships.length > 0 && (
            <>
              <Text size="2" color="gray">
                Showing {startIndex + 1}-{Math.min(endIndex, filteredRelationships.length)} of{' '}
                {filteredRelationships.length} relationship
                {filteredRelationships.length !== 1 ? 's' : ''}
                {filter && ` matching "${filter}"`}
              </Text>

              <Table.Root variant="surface">
                <Table.Header>
                  <Table.Row>
                    <Table.ColumnHeaderCell>Type</Table.ColumnHeaderCell>
                    <Table.ColumnHeaderCell>From Person</Table.ColumnHeaderCell>
                    <Table.ColumnHeaderCell>To Person</Table.ColumnHeaderCell>
                    <Table.ColumnHeaderCell>Start Date</Table.ColumnHeaderCell>
                    <Table.ColumnHeaderCell>Actions</Table.ColumnHeaderCell>
                  </Table.Row>
                </Table.Header>

                <Table.Body>
                  {currentRelationships.map(relationship => (
                    <Table.Row key={relationship.id}>
                      <Table.Cell>
                        <Badge color={getRelationshipTypeBadgeColor(relationship.type)}>
                          {relationship.type}
                        </Badge>
                      </Table.Cell>
                      <Table.Cell>
                        <Text size="2">{relationship.fromPersonId}</Text>
                      </Table.Cell>
                      <Table.Cell>
                        <Text size="2">{relationship.toPersonId}</Text>
                      </Table.Cell>
                      <Table.Cell>
                        <Text size="2" color="gray">
                          {relationship.startDate
                            ? new Date(relationship.startDate).toLocaleDateString()
                            : 'N/A'}
                        </Text>
                      </Table.Cell>
                      <Table.Cell>
                        <Link
                          to={`/relationships/${encodeURIComponent(relationship.id || '')}`}
                          style={{ textDecoration: 'none' }}
                        >
                          <Button size="1" variant="soft">
                            View
                          </Button>
                        </Link>
                      </Table.Cell>
                    </Table.Row>
                  ))}
                </Table.Body>
              </Table.Root>

              {totalPages > 1 && (
                <Flex justify="center" gap="2" mt="2">
                  <Button
                    size="2"
                    variant="soft"
                    disabled={currentPage === 1}
                    onClick={() => setCurrentPage(p => p - 1)}
                  >
                    Previous
                  </Button>
                  <Flex align="center" gap="1">
                    {Array.from({ length: Math.min(totalPages, 7) }, (_, i) => {
                      let page;
                      if (totalPages <= 7) {
                        page = i + 1;
                      } else if (currentPage <= 4) {
                        page = i + 1;
                      } else if (currentPage >= totalPages - 3) {
                        page = totalPages - 6 + i;
                      } else {
                        page = currentPage - 3 + i;
                      }
                      return (
                        <Button
                          key={page}
                          size="2"
                          variant={page === currentPage ? 'solid' : 'soft'}
                          onClick={() => setCurrentPage(page)}
                        >
                          {page}
                        </Button>
                      );
                    })}
                  </Flex>
                  <Button
                    size="2"
                    variant="soft"
                    disabled={currentPage === totalPages}
                    onClick={() => setCurrentPage(p => p + 1)}
                  >
                    Next
                  </Button>
                </Flex>
              )}
            </>
          )}
        </Flex>
      </Card>
    </Flex>
  );
}
