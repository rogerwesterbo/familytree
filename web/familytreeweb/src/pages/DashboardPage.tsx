import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { Flex, Card, Heading, Text, Grid, Box, Spinner, Badge, Button } from '@radix-ui/themes';
import {
  PersonIcon,
  GlobeIcon,
  RocketIcon,
  LightningBoltIcon,
  ReloadIcon,
} from '@radix-ui/react-icons';
import * as api from '../services/api';
import * as adminApi from '../services/admin-api';

export default function DashboardPage() {
  const [persons, setPersons] = useState<api.Person[]>([]);
  const [relationships, setRelationships] = useState<api.Relationship[]>([]);
  const [systemStats, setSystemStats] = useState<adminApi.SystemStats | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      setIsLoading(true);
      setError(null);
      const [personsData, relationshipsData, system] = await Promise.all([
        api.listPersons().catch(() => []),
        api.listRelationships().catch(() => []),
        adminApi.getSystemStats().catch(() => null),
      ]);
      console.log('DashboardPage - System stats:', system);
      setPersons(Array.isArray(personsData) ? personsData : []);
      setRelationships(Array.isArray(relationshipsData) ? relationshipsData : []);
      setSystemStats(system);
    } catch (err) {
      console.error('Failed to load data:', err);
      setError(err instanceof Error ? err.message : 'Failed to load dashboard data');
    } finally {
      setIsLoading(false);
    }
  };

  // Calculate metrics
  const totalPersons = persons.length;
  const totalRelationships = relationships.length;

  // Relationship type distribution
  const relationshipTypes = Array.isArray(relationships)
    ? relationships.reduce(
        (acc, rel) => {
          acc[rel.type] = (acc[rel.type] || 0) + 1;
          return acc;
        },
        {} as Record<string, number>
      )
    : {};

  const topRelationshipTypes = Object.entries(relationshipTypes)
    .sort(([, a], [, b]) => (b as number) - (a as number))
    .slice(0, 5);

  return (
    <Flex direction="column" gap="6">
      <Flex justify="between" align="center">
        <Heading size="8">Dashboard</Heading>
        <Flex gap="3" align="center">
          <Text size="2" color="gray">
            {new Date().toLocaleDateString('en-US', {
              weekday: 'long',
              year: 'numeric',
              month: 'long',
              day: 'numeric',
            })}
          </Text>
          <Button size="3" variant="soft" onClick={loadData} disabled={isLoading}>
            <ReloadIcon />
            Refresh
          </Button>
        </Flex>
      </Flex>

      {isLoading ? (
        <Flex justify="center" py="8">
          <Spinner size="3" />
        </Flex>
      ) : error ? (
        <Card>
          <Text color="red" size="3">
            {error}
          </Text>
        </Card>
      ) : (
        <>
          <Grid columns={{ initial: '1', sm: '2', md: '4' }} gap="4">
            <Card>
              <Flex direction="column" gap="3">
                <Flex align="center" gap="2">
                  <PersonIcon width="20" height="20" />
                  <Heading size="4">Persons</Heading>
                </Flex>
                <Flex direction="column" gap="1">
                  <Text size="7" weight="bold">
                    {totalPersons}
                  </Text>
                  <Text size="2" color="gray">
                    Total persons
                  </Text>
                </Flex>
              </Flex>
            </Card>

            <Card>
              <Flex direction="column" gap="3">
                <Flex align="center" gap="2">
                  <GlobeIcon width="20" height="20" />
                  <Heading size="4">Relationships</Heading>
                </Flex>
                <Flex direction="column" gap="1">
                  <Text size="7" weight="bold">
                    {totalRelationships}
                  </Text>
                  <Text size="2" color="gray">
                    Total relationships
                  </Text>
                </Flex>
              </Flex>
            </Card>

            <Card>
              <Flex direction="column" gap="3">
                <Flex align="center" gap="2">
                  <RocketIcon width="20" height="20" />
                  <Heading size="4">API Calls</Heading>
                </Flex>
                <Flex direction="column" gap="1">
                  <Text size="7" weight="bold">
                    {systemStats?.query_log?.total_queries?.toLocaleString() || '0'}
                  </Text>
                  <Text size="2" color="gray">
                    Total API requests
                  </Text>
                </Flex>
              </Flex>
            </Card>

            <Card>
              <Flex direction="column" gap="3">
                <Flex align="center" gap="2">
                  <LightningBoltIcon width="20" height="20" />
                  <Heading size="4">Cache Hit Rate</Heading>
                </Flex>
                <Flex direction="column" gap="1">
                  <Text
                    size="7"
                    weight="bold"
                    color={systemStats?.query_log?.enabled ? 'green' : 'gray'}
                  >
                    {systemStats?.query_log?.enabled
                      ? `${(systemStats.query_log.cache_hit_rate * 100).toFixed(1)}%`
                      : 'Disabled'}
                  </Text>
                  <Text size="2" color="gray">
                    {systemStats?.query_log?.enabled ? 'Cache efficiency' : 'Enable in settings'}
                  </Text>
                </Flex>
              </Flex>
            </Card>
          </Grid>

          <Grid columns={{ initial: '1', md: '2' }} gap="4">
            <Card>
              <Flex direction="column" gap="4">
                <Heading size="5">Relationship Type Distribution</Heading>
                {topRelationshipTypes.length === 0 ? (
                  <Text size="2" color="gray">
                    No relationships found
                  </Text>
                ) : (
                  <Flex direction="column" gap="3">
                    {topRelationshipTypes.map(([type, count]) => (
                      <Flex key={type} justify="between" align="center">
                        <Flex align="center" gap="2">
                          <Badge
                            color={
                              type === 'parent' || type === 'child'
                                ? 'blue'
                                : type === 'spouse' || type === 'partner'
                                  ? 'green'
                                  : type === 'sibling'
                                    ? 'orange'
                                    : 'purple'
                            }
                          >
                            {type}
                          </Badge>
                          <Text size="2">
                            {count} relationship{count !== 1 ? 's' : ''}
                          </Text>
                        </Flex>
                        <Box
                          style={{
                            width: '100px',
                            height: '8px',
                            backgroundColor: 'var(--gray-a3)',
                            borderRadius: '4px',
                            overflow: 'hidden',
                          }}
                        >
                          <Box
                            style={{
                              width: `${(count / totalRelationships) * 100}%`,
                              height: '100%',
                              backgroundColor: 'var(--accent-9)',
                            }}
                          />
                        </Box>
                      </Flex>
                    ))}
                  </Flex>
                )}
              </Flex>
            </Card>

            <Card>
              <Flex direction="column" gap="4">
                <Heading size="5">Recent Persons</Heading>
                {persons.length === 0 ? (
                  <Text size="2" color="gray">
                    No persons found
                  </Text>
                ) : (
                  <Flex direction="column" gap="3">
                    {persons.slice(0, 5).map(person => (
                      <Box
                        key={person.id}
                        style={{
                          padding: '12px',
                          borderRadius: '6px',
                          backgroundColor: 'var(--gray-a2)',
                        }}
                      >
                        <Flex justify="between" align="center">
                          <Flex direction="column" gap="1">
                            <Link to={`/persons`} style={{ textDecoration: 'none' }}>
                              <Text size="2" weight="medium">
                                {person.firstName} {person.lastName}
                              </Text>
                            </Link>
                            <Text size="1" color="gray">
                              {person.email || 'No email'}
                            </Text>
                          </Flex>
                          <Badge color="blue">
                            {person.birthDate ? new Date(person.birthDate).getFullYear() : 'N/A'}
                          </Badge>
                        </Flex>
                      </Box>
                    ))}
                  </Flex>
                )}
              </Flex>
            </Card>
          </Grid>
        </>
      )}
    </Flex>
  );
}
