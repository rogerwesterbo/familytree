import { useCallback, useEffect, useState } from 'react';
import {
  ReactFlow,
  Background,
  Controls,
  MiniMap,
  addEdge,
  useNodesState,
  useEdgesState,
  Panel,
  type Connection,
  type Edge,
  type Node,
} from '@xyflow/react';
import '@xyflow/react/dist/style.css';
import { Box, Button, Flex, Heading } from '@radix-ui/themes';
import { PlusIcon } from '@radix-ui/react-icons';
import * as api from '../services/api';
import { useTheme } from '../contexts';

/**
 * FAMILY TREE MODELING GUIDE
 *
 * This graph uses a hybrid node-edge approach for family trees:
 *
 * 1. PERSON NODES: Represent individuals with all their details (name, dates, etc.)
 *
 * 2. RELATIONSHIP NODES: Intermediate nodes representing connections between people
 *    - Contains rich data: type, startDate, endDate, notes
 *    - Different types: spouse/married, parent-child, sibling, etc.
 *    - Can be color-coded by type
 *
 * 3. EDGES: Connect persons to relationships
 *    - Person A → Relationship → Person B
 *    - Relationship → Child Person (for offspring)
 *
 * MODELING PATTERNS:
 *
 * A. MARRIAGE/PARTNERSHIP:
 *    Parent1 → [Spouse Relationship] → Parent2
 *    - Relationship type: "spouse" or "married"
 *    - startDate: when they met/married
 *    - endDate: divorce/separation (if applicable)
 *    - notes: additional details
 *
 * B. CHILDREN FROM RELATIONSHIP:
 *    Parent1 → [Spouse Rel] → Parent2
 *                  ↓
 *                Child
 *    - Create parent-child relationships from each parent to child
 *    - Or link child directly from the spouse relationship node
 *
 * C. PARENT-CHILD:
 *    Parent → [Parent Relationship] → Child
 *    - Relationship type: "parent" (or "child" from child's perspective)
 *    - Can track adoption dates, custody info in notes
 *
 * D. SIBLINGS:
 *    Sibling1 → [Sibling Relationship] → Sibling2
 *    - Relationship type: "sibling"
 *    - Or derive from common parents
 *
 * ADVANCED: For children from marriages, you could:
 * 1. Query child relationships where fromPersonId or toPersonId is a parent
 * 2. Position child nodes below the relationship node
 * 3. Connect relationship → child with edges
 */

interface PersonNodeData extends api.Person {
  [key: string]: unknown;
}

interface RelationshipNodeData extends api.Relationship {
  [key: string]: unknown;
}

// Custom node types
const nodeTypes = {
  person: PersonNode,
  relationship: RelationshipNode,
};

function PersonNode({ data }: { data: PersonNodeData }) {
  return (
    <Box
      style={{
        padding: '12px 16px',
        borderRadius: '8px',
        background: 'var(--accent-3)',
        border: '2px solid var(--accent-8)',
        minWidth: '180px',
      }}
    >
      <div style={{ fontWeight: 'bold', marginBottom: '4px' }}>
        {data.firstName} {data.lastName}
      </div>
      {data.birthDate && (
        <div style={{ fontSize: '12px', color: 'var(--gray-11)' }}>
          Born: {new Date(data.birthDate).toLocaleDateString()}
        </div>
      )}
      {data.gender && (
        <div style={{ fontSize: '12px', color: 'var(--gray-11)' }}>{data.gender}</div>
      )}
    </Box>
  );
}

// Relationship Group Node - acts as a container for the two parents and children
function RelationshipNode({ data }: { data: RelationshipNodeData }) {
  const getRelationshipColor = (type: string) => {
    switch (type.toLowerCase()) {
      case 'spouse':
      case 'married':
        return { bg: 'var(--pink-2)', border: 'var(--pink-7)', header: 'var(--pink-9)' };
      case 'parent':
      case 'child':
        return { bg: 'var(--green-2)', border: 'var(--green-7)', header: 'var(--green-9)' };
      case 'sibling':
        return { bg: 'var(--orange-2)', border: 'var(--orange-7)', header: 'var(--orange-9)' };
      default:
        return { bg: 'var(--blue-2)', border: 'var(--blue-7)', header: 'var(--blue-9)' };
    }
  };

  const colors = getRelationshipColor(data.type);

  return (
    <Box
      style={{
        padding: '8px',
        borderRadius: '8px',
        background: colors.bg,
        border: `2px dashed ${colors.border}`,
        minWidth: '450px',
        minHeight: '280px',
      }}
    >
      <div
        style={{
          fontWeight: 'bold',
          fontSize: '12px',
          textAlign: 'center',
          marginBottom: '8px',
          padding: '4px 8px',
          background: colors.header,
          color: 'white',
          borderRadius: '4px',
        }}
      >
        {data.type.toUpperCase()}
      </div>
      <div
        style={{
          display: 'flex',
          flexDirection: 'column',
          gap: '4px',
          fontSize: '10px',
          color: 'var(--gray-11)',
        }}
      >
        {data.startDate && <div>Start: {new Date(data.startDate).toLocaleDateString()}</div>}
        {data.endDate && <div>End: {new Date(data.endDate).toLocaleDateString()}</div>}
        {data.notes && (
          <div style={{ fontStyle: 'italic', marginTop: '4px', fontSize: '9px' }}>{data.notes}</div>
        )}
      </div>
    </Box>
  );
}

export default function GraphViewPage() {
  const { theme } = useTheme();
  const [nodes, setNodes, onNodesChange] = useNodesState<Node>([]);
  const [edges, setEdges, onEdgesChange] = useEdgesState<Edge>([]);
  const [loading, setLoading] = useState(true);

  const buildGraph = useCallback(
    (personsData: api.Person[], relationshipsData: api.Relationship[]) => {
      const newNodes: Node[] = [];
      const newEdges: Edge[] = [];

      // Track which persons are part of relationships (will be inside groups)
      const personsInRelationships = new Set<string>();
      const relationshipGroups = new Map<string, { rel: api.Relationship; index: number }>();

      // First pass: identify persons in relationships and create relationship groups
      let relationshipIndex = 0;
      relationshipsData.forEach(rel => {
        const relNodeId = `rel-${rel.id}`;
        relationshipGroups.set(relNodeId, { rel, index: relationshipIndex });

        // Mark these persons as being in a relationship group
        personsInRelationships.add(rel.fromPersonId);
        personsInRelationships.add(rel.toPersonId);

        relationshipIndex++;
      });

      // Create relationship group nodes (these will contain the parent persons)
      relationshipGroups.forEach(({ rel, index }, relNodeId) => {
        const groupX = 50 + (index % 3) * 550;
        const groupY = 50 + Math.floor(index / 3) * 400;

        newNodes.push({
          id: relNodeId,
          type: 'relationship',
          position: { x: groupX, y: groupY },
          data: { ...rel } as RelationshipNodeData,
          style: {
            width: 450,
            height: 280,
          },
        });
      });

      // Create person nodes
      // Persons in relationships are positioned inside their group, others positioned independently
      let standalonePersonIndex = 0;

      personsData.forEach(person => {
        const personId = person.id!;
        const personNodeId = `person-${personId}`;

        // Find which relationship(s) this person is in
        const parentRelationships = Array.from(relationshipGroups.entries()).filter(
          ([, { rel }]) => rel.fromPersonId === personId || rel.toPersonId === personId
        );

        if (parentRelationships.length > 0) {
          // Position inside the first relationship group
          const [relNodeId, { rel }] = parentRelationships[0];

          // Position parents side by side inside the group
          const isFirstPerson = rel.fromPersonId === personId;
          const relativeX = isFirstPerson ? 50 : 250;
          const relativeY = 80;

          newNodes.push({
            id: personNodeId,
            type: 'person',
            position: { x: relativeX, y: relativeY },
            data: { ...person } as PersonNodeData,
            parentId: relNodeId, // This makes it a child of the group
            extent: 'parent' as const, // Constrains movement within parent
          });
        } else {
          // Standalone person (not in a relationship yet)
          const x = 50 + (standalonePersonIndex % 4) * 200;
          const y = 600 + Math.floor(standalonePersonIndex / 4) * 200;

          newNodes.push({
            id: personNodeId,
            type: 'person',
            position: { x, y },
            data: { ...person } as PersonNodeData,
          });

          standalonePersonIndex++;
        }
      });

      // Find children from relationships
      relationshipGroups.forEach(({ rel }, relNodeId) => {
        // Look for persons who are children of both parents in this relationship
        const children = relationshipsData.filter(
          childRel =>
            childRel.type.toLowerCase() === 'parent' &&
            (childRel.fromPersonId === rel.fromPersonId || childRel.fromPersonId === rel.toPersonId)
        );

        // Find common children (children of both parents)
        const childrenIds = new Map<string, number>();
        children.forEach(childRel => {
          const childId = childRel.toPersonId;
          childrenIds.set(childId, (childrenIds.get(childId) || 0) + 1);
        });

        // Children with count=2 are children of both parents
        const commonChildren = Array.from(childrenIds.entries())
          .filter(([, count]) => count === 2)
          .map(([childId]) => childId);

        // Create edges from relationship group to children
        commonChildren.forEach((childId, index) => {
          newEdges.push({
            id: `edge-${relNodeId}-child-${childId}`,
            source: relNodeId,
            target: `person-${childId}`,
            animated: false,
            label: 'child',
            style: { stroke: 'var(--green-9)', strokeWidth: 2 },
            type: 'smoothstep',
            sourceHandle: 'bottom',
          });

          // Update child position to be below the relationship group
          const childNodeIndex = newNodes.findIndex(n => n.id === `person-${childId}`);
          if (childNodeIndex !== -1) {
            const groupNode = newNodes.find(n => n.id === relNodeId);
            if (groupNode) {
              newNodes[childNodeIndex].position = {
                x: groupNode.position.x + 100 + index * 150,
                y: groupNode.position.y + 350,
              };
            }
          }
        });
      });

      setNodes(newNodes);
      setEdges(newEdges);
    },
    [setNodes, setEdges]
  );

  // Load data
  const loadData = useCallback(async () => {
    try {
      setLoading(true);
      const [personsData, relationshipsData] = await Promise.all([
        api.listPersons(),
        api.listRelationships(),
      ]);

      // Ensure we have arrays
      const persons = Array.isArray(personsData) ? personsData : [];
      const relationships = Array.isArray(relationshipsData) ? relationshipsData : [];

      console.log('Loaded persons:', persons.length, 'relationships:', relationships.length);
      buildGraph(persons, relationships);
    } catch (error) {
      console.error('Failed to load data:', error);
      // Set empty graph on error
      setNodes([]);
      setEdges([]);
    } finally {
      setLoading(false);
    }
  }, [buildGraph, setEdges, setNodes]);

  useEffect(() => {
    loadData();
  }, [loadData]);

  const onConnect = useCallback(
    (params: Connection) => {
      // Allow connections from relationship groups to persons (for children)
      // or between persons and relationships
      const newEdge = {
        ...params,
        type: 'smoothstep',
        animated: false,
        style: { stroke: 'var(--green-9)', strokeWidth: 2 },
      };
      setEdges(eds => addEdge(newEdge, eds));
    },
    [setEdges]
  );

  const addPersonNode = () => {
    const newId = `person-new-${Date.now()}`;
    const newNode: Node = {
      id: newId,
      type: 'person',
      position: { x: Math.random() * 500, y: Math.random() * 500 },
      data: {
        id: newId,
        firstName: 'New',
        lastName: 'Person',
      },
    };
    setNodes(nds => [...nds, newNode]);
  };

  const addRelationshipNode = () => {
    const newId = `rel-new-${Date.now()}`;
    const newNode: Node = {
      id: newId,
      type: 'relationship',
      position: { x: Math.random() * 500 + 200, y: Math.random() * 500 + 100 },
      data: {
        id: newId,
        type: 'parent',
      },
    };
    setNodes(nds => [...nds, newNode]);
  };

  if (loading) {
    return (
      <Flex direction="column" gap="4" p="4">
        <Heading size="6">Family Tree Graph</Heading>
        <div>Loading...</div>
      </Flex>
    );
  }

  return (
    <Flex direction="column" style={{ height: '100%', width: '100%' }}>
      <Box p="4" style={{ borderBottom: '1px solid var(--gray-6)' }}>
        <Flex justify="between" align="center">
          <Heading size="6">Family Tree Graph</Heading>
          <Flex gap="2">
            <Button onClick={loadData} variant="soft">
              Refresh
            </Button>
          </Flex>
        </Flex>
      </Box>

      <Box style={{ flexGrow: 1, position: 'relative' }}>
        <ReactFlow
          nodes={nodes}
          edges={edges}
          onNodesChange={onNodesChange}
          onEdgesChange={onEdgesChange}
          onConnect={onConnect}
          nodeTypes={nodeTypes}
          colorMode={theme}
          fitView
          attributionPosition="bottom-left"
          connectionLineStyle={{ stroke: 'var(--green-9)', strokeWidth: 2 }}
          defaultEdgeOptions={{
            type: 'smoothstep',
            animated: false,
            style: { strokeWidth: 2 },
          }}
        >
          <Background />
          <Controls />
          <MiniMap />
          <Panel position="top-right">
            <Flex direction="column" gap="2" p="2">
              <Button onClick={addPersonNode} size="2">
                <PlusIcon /> Add Person
              </Button>
              <Button onClick={addRelationshipNode} size="2" variant="soft">
                <PlusIcon /> Add Relationship
              </Button>
            </Flex>
          </Panel>
        </ReactFlow>
      </Box>
    </Flex>
  );
}
