import { getValidAccessToken } from './auth';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:15000';

export interface Person {
  id?: string;
  firstName: string;
  lastName: string;
  birthDate?: string;
  deathDate?: string;
  gender?: string;
  email?: string;
  phone?: string;
  createdAt?: string;
  updatedAt?: string;
}

export interface Relationship {
  id?: string;
  type: string; // parent, child, spouse, sibling, etc.
  fromPersonId: string;
  toPersonId: string;
  startDate?: string;
  endDate?: string;
  notes?: string;
  createdAt?: string;
  updatedAt?: string;
}

export interface SearchResult {
  type: 'person' | 'relationship';
  person?: Person;
  relationship?: Relationship;
  highlight?: string;
}

export interface SearchResponse {
  results: SearchResult[] | null;
  total: number;
  query: string;
}

class ApiError extends Error {
  status: number;
  response?: unknown;

  constructor(status: number, message: string, response?: unknown) {
    super(message);
    this.name = 'ApiError';
    this.status = status;
    this.response = response;
  }
}

async function apiRequest<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
  const token = await getValidAccessToken();

  if (!token) {
    throw new ApiError(401, 'Not authenticated');
  }

  const url = `${API_BASE_URL}${endpoint}`;

  const response = await fetch(url, {
    ...options,
    headers: {
      Authorization: `Bearer ${token}`,
      'Content-Type': 'application/json',
      ...options.headers,
    },
  });

  if (!response.ok) {
    let errorMessage = `API request failed: ${response.statusText}`;
    let errorData: unknown;

    try {
      errorData = await response.json();
      if (errorData && typeof errorData === 'object' && 'error' in errorData) {
        errorMessage = String(errorData.error);
      }
    } catch {
      // Ignore JSON parse errors
    }

    throw new ApiError(response.status, errorMessage, errorData);
  }

  // Handle 204 No Content
  if (response.status === 204) {
    return undefined as T;
  }

  return response.json();
}

// Person endpoints
export async function listPersons(): Promise<Person[]> {
  return apiRequest<Person[]>('/v1/persons');
}

export async function getPerson(id: string): Promise<Person> {
  return apiRequest<Person>(`/v1/persons/${encodeURIComponent(id)}`);
}

export async function createPerson(person: Person): Promise<Person> {
  return apiRequest<Person>('/v1/persons', {
    method: 'POST',
    body: JSON.stringify(person),
  });
}

export async function updatePerson(id: string, person: Person): Promise<Person> {
  return apiRequest<Person>(`/v1/persons/${encodeURIComponent(id)}`, {
    method: 'PUT',
    body: JSON.stringify(person),
  });
}

export async function deletePerson(id: string): Promise<void> {
  return apiRequest<void>(`/v1/persons/${encodeURIComponent(id)}`, {
    method: 'DELETE',
  });
}

// Relationship endpoints
export async function listRelationships(): Promise<Relationship[]> {
  return apiRequest<Relationship[]>('/v1/relationships');
}

export async function getRelationship(id: string): Promise<Relationship> {
  return apiRequest<Relationship>(`/v1/relationships/${encodeURIComponent(id)}`);
}

export async function createRelationship(relationship: Relationship): Promise<Relationship> {
  return apiRequest<Relationship>('/v1/relationships', {
    method: 'POST',
    body: JSON.stringify(relationship),
  });
}

export async function updateRelationship(
  id: string,
  relationship: Relationship
): Promise<Relationship> {
  return apiRequest<Relationship>(`/v1/relationships/${encodeURIComponent(id)}`, {
    method: 'PUT',
    body: JSON.stringify(relationship),
  });
}

export async function deleteRelationship(id: string): Promise<void> {
  return apiRequest<void>(`/v1/relationships/${encodeURIComponent(id)}`, {
    method: 'DELETE',
  });
}

// Get relationships for a specific person
export async function getPersonRelationships(personId: string): Promise<Relationship[]> {
  return apiRequest<Relationship[]>(`/v1/persons/${encodeURIComponent(personId)}/relationships`);
}

// Search endpoint
export async function search(
  query: string,
  types?: ('person' | 'relationship')[]
): Promise<SearchResponse> {
  const params = new URLSearchParams({ q: query });

  if (types && types.length > 0) {
    types.forEach(type => params.append('type', type));
  }

  return apiRequest<SearchResponse>(`/v1/search?${params.toString()}`);
}

// Export endpoints
export async function exportAllData(format: string = 'json'): Promise<string> {
  const token = await getValidAccessToken();

  if (!token) {
    throw new ApiError(401, 'Not authenticated');
  }

  const response = await fetch(`${API_BASE_URL}/v1/export?format=${format}`, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });

  if (!response.ok) {
    throw new ApiError(response.status, `Failed to export data`);
  }

  return response.text();
}

export async function exportPerson(id: string, format: string = 'json'): Promise<string> {
  const token = await getValidAccessToken();

  if (!token) {
    throw new ApiError(401, 'Not authenticated');
  }

  const response = await fetch(
    `${API_BASE_URL}/v1/export/person/${encodeURIComponent(id)}?format=${format}`,
    {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    }
  );

  if (!response.ok) {
    throw new ApiError(response.status, `Failed to export person ${id}`);
  }

  return response.text();
}

export { ApiError };
