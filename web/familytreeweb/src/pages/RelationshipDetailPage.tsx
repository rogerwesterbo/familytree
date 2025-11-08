import { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
  Flex,
  Card,
  Heading,
  Text,
  Button,
  Spinner,
  Badge,
  Box,
  AlertDialog,
} from '@radix-ui/themes';
import { ArrowLeftIcon, TrashIcon, Pencil1Icon } from '@radix-ui/react-icons';
import * as api from '../services/api';

export default function RelationshipDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [relationship, setRelationship] = useState<api.Relationship | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showDeleteDialog, setShowDeleteDialog] = useState(false);

  useEffect(() => {
    if (id) {
      loadRelationship(id);
    }
  }, [id]);

  const loadRelationship = async (relationshipId: string) => {
    try {
      setIsLoading(true);
      setError(null);
      const data = await api.getRelationship(relationshipId);
      setRelationship(data);
    } catch (err) {
      console.error('Failed to load relationship:', err);
      setError(err instanceof Error ? err.message : 'Failed to load relationship');
    } finally {
      setIsLoading(false);
    }
  };

  const handleDelete = async () => {
    if (!relationship?.id) return;

    try {
      await api.deleteRelationship(relationship.id);
      navigate('/relationships');
    } catch (err) {
      console.error('Failed to delete relationship:', err);
      setError(err instanceof Error ? err.message : 'Failed to delete relationship');
    }
  };

  if (isLoading) {
    return (
      <Flex justify="center" align="center" style={{ minHeight: '400px' }}>
        <Spinner size="3" />
      </Flex>
    );
  }

  if (error || !relationship) {
    return (
      <Flex direction="column" gap="4">
        <Button variant="soft" onClick={() => navigate('/relationships')}>
          <ArrowLeftIcon /> Back to Relationships
        </Button>
        <Card>
          <Text color="red" size="3">
            {error || 'Relationship not found'}
          </Text>
        </Card>
      </Flex>
    );
  }

  return (
    <Flex direction="column" gap="6">
      <Flex justify="between" align="center">
        <Flex align="center" gap="4">
          <Button variant="soft" onClick={() => navigate('/relationships')}>
            <ArrowLeftIcon />
          </Button>
          <Heading size="8">Relationship Details</Heading>
        </Flex>
        <Flex gap="2">
          <Button variant="soft" color="blue">
            <Pencil1Icon /> Edit
          </Button>
          <Button variant="soft" color="red" onClick={() => setShowDeleteDialog(true)}>
            <TrashIcon /> Delete
          </Button>
        </Flex>
      </Flex>

      <Card size="3">
        <Flex direction="column" gap="6">
          <Flex direction="column" gap="2">
            <Text size="2" weight="bold" color="gray">
              Relationship Type
            </Text>
            <Badge size="2" color="blue">
              {relationship.type}
            </Badge>
          </Flex>

          <Box style={{ height: '1px', backgroundColor: 'var(--gray-a5)' }} />

          <Flex direction="column" gap="4">
            <Flex direction="column" gap="2">
              <Text size="2" weight="bold" color="gray">
                From Person ID
              </Text>
              <Text size="3">{relationship.fromPersonId}</Text>
            </Flex>

            <Flex direction="column" gap="2">
              <Text size="2" weight="bold" color="gray">
                To Person ID
              </Text>
              <Text size="3">{relationship.toPersonId}</Text>
            </Flex>
          </Flex>

          {relationship.startDate && (
            <>
              <Box style={{ height: '1px', backgroundColor: 'var(--gray-a5)' }} />
              <Flex direction="column" gap="2">
                <Text size="2" weight="bold" color="gray">
                  Start Date
                </Text>
                <Text size="3">{new Date(relationship.startDate).toLocaleDateString()}</Text>
              </Flex>
            </>
          )}

          {relationship.endDate && (
            <Flex direction="column" gap="2">
              <Text size="2" weight="bold" color="gray">
                End Date
              </Text>
              <Text size="3">{new Date(relationship.endDate).toLocaleDateString()}</Text>
            </Flex>
          )}

          {relationship.notes && (
            <>
              <Box style={{ height: '1px', backgroundColor: 'var(--gray-a5)' }} />
              <Flex direction="column" gap="2">
                <Text size="2" weight="bold" color="gray">
                  Notes
                </Text>
                <Text size="3">{relationship.notes}</Text>
              </Flex>
            </>
          )}
        </Flex>
      </Card>

      <AlertDialog.Root open={showDeleteDialog} onOpenChange={setShowDeleteDialog}>
        <AlertDialog.Content>
          <AlertDialog.Title>Delete Relationship</AlertDialog.Title>
          <AlertDialog.Description>
            Are you sure you want to delete this {relationship.type} relationship? This action
            cannot be undone.
          </AlertDialog.Description>
          <Flex gap="3" mt="4" justify="end">
            <AlertDialog.Cancel>
              <Button variant="soft" color="gray">
                Cancel
              </Button>
            </AlertDialog.Cancel>
            <AlertDialog.Action>
              <Button variant="solid" color="red" onClick={handleDelete}>
                Delete Relationship
              </Button>
            </AlertDialog.Action>
          </Flex>
        </AlertDialog.Content>
      </AlertDialog.Root>
    </Flex>
  );
}
