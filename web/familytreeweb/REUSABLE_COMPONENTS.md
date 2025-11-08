# Reusable Components Guide

This document describes the reusable components and hooks available in the familytree web application.

## 1. useApiData Hook

A custom hook for fetching and managing API data with built-in error handling, loading states, and data validation.

### Location

`src/hooks/useApiData.ts`

### Usage

```typescript
import { useApiData, isArray } from '../hooks/useApiData';
import * as api from '../services/api';

function MyComponent() {
  const { data: persons, isLoading, error, refetch } = useApiData({
    fetchFn: api.listPersons,
    initialData: [],
    validateData: isArray<api.Person>,
  });

  return (
    <LoadingState
      isLoading={isLoading}
      error={error}
      isEmpty={persons.length === 0}
      onRetry={refetch}
    >
      {/* Render your data here */}
      {persons.map(person => (
        <div key={person.id}>{person.firstName}</div>
      ))}
    </LoadingState>
  );
}
```

### Features

- Automatic data fetching on mount
- Built-in loading state management
- Error handling with optional callback
- Data validation with type guards
- Manual refetch capability
- Safe fallback to initial data on error

### Parameters

- `fetchFn`: Function that returns a Promise with the data
- `initialData`: Initial value for the data (e.g., `[]` for arrays)
- `validateData`: Optional type guard function to validate response
- `onError`: Optional callback for error handling

### Returns

- `data`: The fetched data (or initial data on error)
- `isLoading`: Boolean indicating loading state
- `error`: Error message string or null
- `refetch`: Function to manually trigger data fetch

---

## 2. LoadingState Component

A reusable component that handles loading, error, and empty states with consistent UI.

### Location

`src/components/LoadingState.tsx`

### Usage

```typescript
import { LoadingState } from '../components';

function MyComponent() {
  const [data, setData] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);

  return (
    <LoadingState
      isLoading={isLoading}
      error={error}
      isEmpty={data.length === 0}
      emptyMessage="No items found"
      onRetry={refetchData}
    >
      {/* Your content when data is loaded */}
      <Table>
        {data.map(item => <TableRow key={item.id} />)}
      </Table>
    </LoadingState>
  );
}
```

### Features

- Displays spinner during loading
- Shows error message with optional retry button
- Shows empty state message when no data
- Only renders children when data is available

### Props

- `isLoading`: Boolean for loading state
- `error`: Error message string or null
- `isEmpty`: Boolean indicating if data is empty
- `emptyMessage`: Custom message for empty state (default: "No data available")
- `onRetry`: Optional callback for retry button
- `children`: Content to render when data is available

---

## 3. Complete Example: Refactored Page Component

Here's how to combine both reusable components in a page:

```typescript
import { useApiData, isArray } from '../hooks/useApiData';
import { LoadingState } from '../components';
import * as api from '../services/api';

export default function PersonsPage() {
  const [filter, setFilter] = useState('');

  // Use the reusable hook for data fetching
  const {
    data: persons,
    isLoading,
    error,
    refetch
  } = useApiData({
    fetchFn: api.listPersons,
    initialData: [],
    validateData: isArray<api.Person>,
  });

  // Filter logic
  const filteredPersons = persons.filter(person =>
    !filter ||
    person.firstName.toLowerCase().includes(filter.toLowerCase()) ||
    person.lastName.toLowerCase().includes(filter.toLowerCase())
  );

  return (
    <Flex direction="column" gap="6">
      <Heading size="8">Persons</Heading>

      {/* Reusable loading/error/empty state handler */}
      <LoadingState
        isLoading={isLoading}
        error={error}
        isEmpty={filteredPersons.length === 0}
        emptyMessage="No persons found. Add one to get started!"
        onRetry={refetch}
      >
        {/* Table or content only renders when data is available */}
        <Table.Root>
          <Table.Body>
            {filteredPersons.map(person => (
              <Table.Row key={person.id}>
                <Table.Cell>{person.firstName} {person.lastName}</Table.Cell>
              </Table.Row>
            ))}
          </Table.Body>
        </Table.Root>
      </LoadingState>
    </Flex>
  );
}
```

---

## 4. Benefits of These Reusable Components

### Code Reduction

- Eliminates repetitive try-catch blocks
- Standardizes error handling across pages
- Reduces boilerplate for loading states

### Consistency

- Same loading spinner appearance everywhere
- Consistent error message display
- Uniform empty state handling

### Type Safety

- Built-in TypeScript support
- Type guards for data validation
- Prevents runtime errors from invalid data

### Maintainability

- Centralized logic for common patterns
- Easier to update loading/error UX globally
- Reduced code duplication

---

## 5. Migration Guide

### Before (Manual approach)

```typescript
const [data, setData] = useState([]);
const [isLoading, setIsLoading] = useState(true);
const [error, setError] = useState(null);

useEffect(() => {
  const loadData = async () => {
    try {
      setIsLoading(true);
      setError(null);
      const result = await api.getData();
      setData(Array.isArray(result) ? result : []);
    } catch (err) {
      setError(err.message);
      setData([]);
    } finally {
      setIsLoading(false);
    }
  };
  loadData();
}, []);

return (
  <>
    {isLoading && <Spinner />}
    {error && <Text color="red">{error}</Text>}
    {!isLoading && !error && data.length === 0 && <Text>No data</Text>}
    {!isLoading && !error && data.length > 0 && (
      // Render data
    )}
  </>
);
```

### After (Using reusable components)

```typescript
const { data, isLoading, error, refetch } = useApiData({
  fetchFn: api.getData,
  initialData: [],
  validateData: isArray,
});

return (
  <LoadingState isLoading={isLoading} error={error} isEmpty={data.length === 0} onRetry={refetch}>
    {/* Render data */}
  </LoadingState>
);
```

---

## 6. Current Status

### Pages Already Fixed

- ✅ DashboardPage - Manual array validation
- ✅ PersonsPage - Manual array validation
- ✅ RelationshipsPage - Manual array validation
- ✅ ExportPage - Manual array validation

### Ready for Refactoring

All the above pages can now be refactored to use the new reusable components for cleaner, more maintainable code.

### How to Refactor

1. Replace manual `useState` and `useEffect` with `useApiData` hook
2. Replace manual loading/error checks with `<LoadingState>` component
3. Remove redundant `Array.isArray()` checks (handled by hook)
4. Keep business logic (filtering, pagination) separate

This will reduce code by ~30-40 lines per page while improving consistency and maintainability.
