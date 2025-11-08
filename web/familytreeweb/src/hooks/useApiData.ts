import { useState, useEffect, useCallback } from 'react';

interface UseApiDataOptions<T> {
  fetchFn: () => Promise<T>;
  initialData?: T;
  validateData?: (data: unknown) => data is T;
  onError?: (error: Error) => void;
}

interface UseApiDataResult<T> {
  data: T;
  isLoading: boolean;
  error: string | null;
  refetch: () => Promise<void>;
}

/**
 * Reusable hook for fetching and managing API data with proper error handling
 * and array validation.
 *
 * @example
 * const { data: persons, isLoading, error, refetch } = useApiData({
 *   fetchFn: api.listPersons,
 *   initialData: [],
 *   validateData: (data): data is Person[] => Array.isArray(data),
 * });
 */
export function useApiData<T>({
  fetchFn,
  initialData,
  validateData,
  onError,
}: UseApiDataOptions<T>): UseApiDataResult<T> {
  const [data, setData] = useState<T>(initialData as T);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchData = useCallback(async () => {
    try {
      setIsLoading(true);
      setError(null);
      const result = await fetchFn();

      // Validate data if validator provided
      if (validateData && !validateData(result)) {
        throw new Error('Invalid data format received from API');
      }

      setData(result);
    } catch (err) {
      console.error('Failed to fetch data:', err);
      const errorMessage = err instanceof Error ? err.message : 'Failed to load data';
      setError(errorMessage);

      // Reset to initial data on error
      if (initialData !== undefined) {
        setData(initialData);
      }

      if (onError && err instanceof Error) {
        onError(err);
      }
    } finally {
      setIsLoading(false);
    }
  }, [fetchFn, validateData, initialData, onError]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  return {
    data,
    isLoading,
    error,
    refetch: fetchData,
  };
}

/**
 * Type guard to validate array data
 */
export function isArray<T>(data: unknown): data is T[] {
  return Array.isArray(data);
}
