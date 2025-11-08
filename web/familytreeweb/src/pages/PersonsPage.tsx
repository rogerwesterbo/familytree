import { useEffect, useState } from 'react';
import {
  Flex,
  Card,
  Heading,
  Text,
  Button,
  Table,
  Spinner,
  TextField,
  Box,
  AlertDialog,
  IconButton,
} from '@radix-ui/themes';
import {
  PlusIcon,
  MagnifyingGlassIcon,
  Pencil1Icon,
  TrashIcon,
  ReloadIcon,
} from '@radix-ui/react-icons';
import * as api from '../services/api';

export default function PersonsPage() {
  const [persons, setPersons] = useState<api.Person[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [filter, setFilter] = useState('');
  const [currentPage, setCurrentPage] = useState(1);
  const itemsPerPage = 15;

  // CRUD states
  const [deletingPerson, setDeletingPerson] = useState<api.Person | null>(null);

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

  const handleDeletePerson = async () => {
    if (!deletingPerson || !deletingPerson.id) return;

    try {
      await api.deletePerson(deletingPerson.id);
      await loadPersons();
      setDeletingPerson(null);
    } catch (err) {
      console.error('Failed to delete person:', err);
      setError(err instanceof Error ? err.message : 'Failed to delete person');
    }
  };

  // Apply filters
  const filteredPersons = Array.isArray(persons)
    ? persons.filter(person => {
        const matchesFilter =
          !filter.trim() ||
          person.firstName.toLowerCase().includes(filter.toLowerCase()) ||
          person.lastName.toLowerCase().includes(filter.toLowerCase()) ||
          (person.email && person.email.toLowerCase().includes(filter.toLowerCase()));
        return matchesFilter;
      })
    : [];

  // Pagination
  const totalPages = Math.ceil(filteredPersons.length / itemsPerPage);
  const startIndex = (currentPage - 1) * itemsPerPage;
  const endIndex = startIndex + itemsPerPage;
  const currentPersons = filteredPersons.slice(startIndex, endIndex);

  // Reset to page 1 when filter changes
  useEffect(() => {
    setCurrentPage(1);
  }, [filter]);

  return (
    <Flex direction="column" gap="6">
      <Flex justify="between" align="center">
        <Heading size="8">Persons</Heading>
        <Flex gap="2">
          <Button size="3" variant="soft" onClick={loadPersons}>
            <ReloadIcon /> Refresh
          </Button>
          <Button size="3">
            <PlusIcon /> Add Person
          </Button>
        </Flex>
      </Flex>

      <Card>
        <Flex direction="column" gap="4">
          <Flex justify="between" align="center" gap="3">
            <Text size="2" color="gray">
              Manage persons in the family tree
            </Text>
            <Box style={{ width: '250px' }}>
              <TextField.Root
                size="2"
                placeholder="Filter persons..."
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

          {!isLoading && !error && currentPersons.length === 0 && (
            <Box py="8">
              <Text color="gray" align="center">
                {persons.length === 0 ? 'No persons found.' : 'No persons match your filter.'}
              </Text>
            </Box>
          )}

          {!isLoading && !error && currentPersons.length > 0 && (
            <>
              <Text size="2" color="gray">
                Showing {startIndex + 1}-{Math.min(endIndex, filteredPersons.length)} of{' '}
                {filteredPersons.length} person{filteredPersons.length !== 1 ? 's' : ''}
                {filter && ' (filtered)'}
              </Text>

              <Table.Root variant="surface">
                <Table.Header>
                  <Table.Row>
                    <Table.ColumnHeaderCell>First Name</Table.ColumnHeaderCell>
                    <Table.ColumnHeaderCell>Last Name</Table.ColumnHeaderCell>
                    <Table.ColumnHeaderCell>Birth Date</Table.ColumnHeaderCell>
                    <Table.ColumnHeaderCell>Email</Table.ColumnHeaderCell>
                    <Table.ColumnHeaderCell>Actions</Table.ColumnHeaderCell>
                  </Table.Row>
                </Table.Header>

                <Table.Body>
                  {currentPersons.map((person, index) => (
                    <Table.Row key={person.id || index}>
                      <Table.Cell>
                        <Text weight="medium">{person.firstName}</Text>
                      </Table.Cell>
                      <Table.Cell>
                        <Text weight="medium">{person.lastName}</Text>
                      </Table.Cell>
                      <Table.Cell>
                        <Text size="2" color="gray">
                          {person.birthDate
                            ? new Date(person.birthDate).toLocaleDateString()
                            : 'N/A'}
                        </Text>
                      </Table.Cell>
                      <Table.Cell>
                        <Text size="2" color="gray">
                          {person.email || 'N/A'}
                        </Text>
                      </Table.Cell>
                      <Table.Cell>
                        <Flex gap="2">
                          <IconButton size="1" variant="ghost" color="blue">
                            <Pencil1Icon />
                          </IconButton>
                          <IconButton
                            size="1"
                            variant="ghost"
                            color="red"
                            onClick={() => setDeletingPerson(person)}
                          >
                            <TrashIcon />
                          </IconButton>
                        </Flex>
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

      {/* Delete Person Confirmation */}
      <AlertDialog.Root
        open={!!deletingPerson}
        onOpenChange={open => !open && setDeletingPerson(null)}
      >
        <AlertDialog.Content>
          <AlertDialog.Title>Delete Person</AlertDialog.Title>
          <AlertDialog.Description>
            Are you sure you want to delete{' '}
            <strong>
              {deletingPerson?.firstName} {deletingPerson?.lastName}
            </strong>
            ? This action cannot be undone.
          </AlertDialog.Description>
          <Flex gap="3" mt="4" justify="end">
            <AlertDialog.Cancel>
              <Button variant="soft" color="gray">
                Cancel
              </Button>
            </AlertDialog.Cancel>
            <AlertDialog.Action>
              <Button variant="solid" color="red" onClick={handleDeletePerson}>
                Delete Person
              </Button>
            </AlertDialog.Action>
          </Flex>
        </AlertDialog.Content>
      </AlertDialog.Root>
    </Flex>
  );
}
