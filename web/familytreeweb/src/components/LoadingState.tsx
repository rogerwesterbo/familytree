import { Flex, Spinner, Text, Card, Button } from '@radix-ui/themes';
import { ReloadIcon } from '@radix-ui/react-icons';

interface LoadingStateProps {
  isLoading: boolean;
  error: string | null;
  onRetry?: () => void;
  emptyMessage?: string;
  isEmpty?: boolean;
  children?: React.ReactNode;
}

/**
 * Reusable component for handling loading, error, and empty states
 *
 * @example
 * <LoadingState isLoading={isLoading} error={error} isEmpty={persons.length === 0} onRetry={refetch}>
 *   <PersonsList persons={persons} />
 * </LoadingState>
 */
export function LoadingState({
  isLoading,
  error,
  onRetry,
  emptyMessage = 'No data available',
  isEmpty = false,
  children,
}: LoadingStateProps) {
  if (isLoading) {
    return (
      <Flex justify="center" py="8">
        <Spinner size="3" />
      </Flex>
    );
  }

  if (error) {
    return (
      <Card>
        <Flex direction="column" gap="3" align="center" py="4">
          <Text color="red" size="3">
            {error}
          </Text>
          {onRetry && (
            <Button size="2" variant="soft" onClick={onRetry}>
              <ReloadIcon /> Retry
            </Button>
          )}
        </Flex>
      </Card>
    );
  }

  if (isEmpty) {
    return (
      <Card>
        <Flex direction="column" gap="3" align="center" py="4">
          <Text color="gray" size="3">
            {emptyMessage}
          </Text>
          {onRetry && (
            <Button size="2" variant="soft" onClick={onRetry}>
              <ReloadIcon /> Refresh
            </Button>
          )}
        </Flex>
      </Card>
    );
  }

  return <>{children}</>;
}
