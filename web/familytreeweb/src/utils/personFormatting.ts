import type { Person, Relationship } from '../services/api';

/**
 * Format a person's full name for display
 */
export function formatPersonName(person: Person): string {
  if (!person.firstName && !person.lastName) {
    return 'Unknown';
  }
  return `${person.firstName || ''} ${person.lastName || ''}`.trim();
}

export function formatDate(date: string | undefined): string {
  if (!date) return 'N/A';
  try {
    return new Date(date).toLocaleDateString();
  } catch {
    return date;
  }
}

/**
 * Format a person's birth year for display
 */
export function formatBirthYear(person: Person): string {
  if (!person.birthDate) return 'Unknown';
  return new Date(person.birthDate).getFullYear().toString();
}

/**
 * Format a person's lifespan for display
 */
export function formatLifespan(person: Person): string {
  const birthYear = person.birthDate ? new Date(person.birthDate).getFullYear() : '?';
  const deathYear = person.deathDate ? new Date(person.deathDate).getFullYear() : '';

  if (deathYear) {
    return `${birthYear} - ${deathYear}`;
  }
  return `${birthYear} - Present`;
}

/**
 * Get a detailed description of a person for tooltips or detailed views
 */
export function getPersonDetails(person: Person): string[] {
  const details: string[] = [];

  if (person.birthDate) {
    details.push(`Born: ${new Date(person.birthDate).toLocaleDateString()}`);
  }
  if (person.deathDate) {
    details.push(`Died: ${new Date(person.deathDate).toLocaleDateString()}`);
  }
  if (person.gender) {
    details.push(`Gender: ${person.gender}`);
  }
  if (person.email) {
    details.push(`Email: ${person.email}`);
  }
  if (person.phone) {
    details.push(`Phone: ${person.phone}`);
  }

  return details;
}

/**
 * Format a relationship description
 */
export function formatRelationship(relationship: Relationship): string {
  return `${relationship.type} relationship`;
}

/**
 * Get relationship details for display
 */
export function getRelationshipDetails(relationship: Relationship): string[] {
  const details: string[] = [];

  details.push(`Type: ${relationship.type}`);

  if (relationship.startDate) {
    details.push(`Started: ${new Date(relationship.startDate).toLocaleDateString()}`);
  }
  if (relationship.endDate) {
    details.push(`Ended: ${new Date(relationship.endDate).toLocaleDateString()}`);
  }
  if (relationship.notes) {
    details.push(`Notes: ${relationship.notes}`);
  }

  return details;
}
