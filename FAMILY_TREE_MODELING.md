# Family Tree Modeling Guide

## Overview

This document explains how to model family relationships in the familytree application using the React Flow graph visualization.

## Data Model

### Person
```typescript
{
  id: string;
  firstName: string;
  lastName: string;
  birthDate?: string;
  deathDate?: string;
  gender?: string;
  email?: string;
  phone?: string;
}
```

### Relationship
```typescript
{
  id: string;
  type: string;           // "spouse", "parent", "child", "sibling"
  fromPersonId: string;   // First person in relationship
  toPersonId: string;     // Second person in relationship
  startDate?: string;     // When relationship began (marriage, birth, etc.)
  endDate?: string;       // When relationship ended (divorce, death, etc.)
  notes?: string;         // Additional details
}
```

## Graph Structure

The family tree uses a **hybrid node-edge approach**:

1. **Person Nodes** - Circular/rectangular nodes showing individual details
2. **Relationship Nodes** - Intermediate nodes showing relationship metadata
3. **Edges** - Connect persons to relationships and relationships to persons/children

```
Person A ─────→ [Relationship] ─────→ Person B
                      ↓
                   Child C
```

## Modeling Patterns

### 1. Marriage/Partnership (Spouse Relationship)

Two people getting married or partnered:

**Persons:**
- John Doe (id: "p1")
- Jane Smith (id: "p2")

**Relationship:**
```json
{
  "id": "r1",
  "type": "spouse",
  "fromPersonId": "p1",
  "toPersonId": "p2",
  "startDate": "2010-06-15",
  "notes": "Married at City Hall"
}
```

**Graph Visualization:**
```
John Doe ─────→ [Spouse: 2010-06-15] ─────→ Jane Smith
                 "Married at City Hall"
```

### 2. Marriage with Divorce

**Relationship:**
```json
{
  "id": "r1",
  "type": "spouse",
  "fromPersonId": "p1",
  "toPersonId": "p2",
  "startDate": "2010-06-15",
  "endDate": "2020-03-10",
  "notes": "Divorced amicably"
}
```

### 3. Parent-Child Relationship

**Persons:**
- John Doe (id: "p1")
- Sarah Doe (id: "p3", child)

**Relationship (Option A - Parent perspective):**
```json
{
  "id": "r2",
  "type": "parent",
  "fromPersonId": "p1",
  "toPersonId": "p3",
  "startDate": "2012-09-20",
  "notes": "Birth date"
}
```

**Relationship (Option B - Child perspective):**
```json
{
  "id": "r2",
  "type": "child",
  "fromPersonId": "p3",
  "toPersonId": "p1",
  "startDate": "2012-09-20"
}
```

> **Recommendation:** Use "parent" type consistently with fromPersonId = parent, toPersonId = child for clarity.

### 4. Complete Family (Parents + Children)

**Persons:**
- John Doe (id: "p1", father)
- Jane Smith (id: "p2", mother)
- Sarah Doe (id: "p3", daughter)
- Mike Doe (id: "p4", son)

**Relationships:**

1. Marriage:
```json
{
  "id": "r1",
  "type": "spouse",
  "fromPersonId": "p1",
  "toPersonId": "p2",
  "startDate": "2010-06-15"
}
```

2. Father → Daughter:
```json
{
  "id": "r2",
  "type": "parent",
  "fromPersonId": "p1",
  "toPersonId": "p3",
  "startDate": "2012-09-20"
}
```

3. Mother → Daughter:
```json
{
  "id": "r3",
  "type": "parent",
  "fromPersonId": "p2",
  "toPersonId": "p3",
  "startDate": "2012-09-20"
}
```

4. Father → Son:
```json
{
  "id": "r4",
  "type": "parent",
  "fromPersonId": "p1",
  "toPersonId": "p4",
  "startDate": "2015-03-15"
}
```

5. Mother → Son:
```json
{
  "id": "r5",
  "type": "parent",
  "fromPersonId": "p2",
  "toPersonId": "p4",
  "startDate": "2015-03-15"
}
```

**Graph Visualization:**
```
          [Spouse Rel]
         /            \
    John Doe        Jane Smith
      |    \        /    |
      |     \      /     |
      |   [Parent Rels]  |
      |       /    \     |
      |      /      \    |
   Sarah Doe        Mike Doe
```

### 5. Alternative: Children from Marriage Relationship

Instead of creating separate parent-child relationships, you could link children directly from the spouse relationship:

**Graph Visualization:**
```
John Doe ───→ [Spouse Rel] ───→ Jane Smith
                   ↓
                   ├─→ Sarah Doe
                   └─→ Mike Doe
```

To implement this in code, when rendering a spouse relationship, query for children where both parents have parent relationships to the same child.

### 6. Siblings

**Option A - Explicit Sibling Relationship:**
```json
{
  "id": "r6",
  "type": "sibling",
  "fromPersonId": "p3",
  "toPersonId": "p4"
}
```

**Option B - Derive from Common Parents:**
Query for people who share the same parent relationships (more reliable for half-siblings).

### 7. Adoption

Use the `notes` field to document adoption details:

```json
{
  "id": "r7",
  "type": "parent",
  "fromPersonId": "p5",
  "toPersonId": "p6",
  "startDate": "2018-11-10",
  "notes": "Adopted from ABC Adoption Agency. Legal adoption finalized on 2018-11-10."
}
```

### 8. Multiple Marriages (Blended Family)

**Persons:**
- John Doe (id: "p1")
- Jane Smith (id: "p2", first wife)
- Mary Johnson (id: "p7", second wife)
- Sarah Doe (id: "p3", child from first marriage)
- Tom Doe (id: "p8", child from second marriage)

**Relationships:**

1. First Marriage (ended):
```json
{
  "type": "spouse",
  "fromPersonId": "p1",
  "toPersonId": "p2",
  "startDate": "2010-06-15",
  "endDate": "2016-12-20"
}
```

2. Second Marriage:
```json
{
  "type": "spouse",
  "fromPersonId": "p1",
  "toPersonId": "p7",
  "startDate": "2018-08-05"
}
```

3. Children from each marriage with appropriate parent relationships.

## Best Practices

### 1. Consistent Relationship Direction

- **Parent-Child:** Always use `fromPersonId` = parent, `toPersonId` = child
- **Spouse:** Order doesn't matter, but be consistent (alphabetical, age, etc.)
- **Sibling:** Order doesn't matter, but avoid duplicates (don't create both A→B and B→A)

### 2. Use startDate and endDate

- **Marriage:** startDate = wedding date, endDate = divorce/death
- **Parent-Child:** startDate = birth/adoption date, endDate typically not used
- **Sibling:** startDate = birth of younger sibling

### 3. Rich Notes Field

Use the notes field for:
- Adoption details
- Custody arrangements
- Location information (married in Paris, born in NYC)
- Special circumstances

### 4. Relationship Types

Standard types:
- `"spouse"` or `"married"` - Marriage/partnership
- `"parent"` - Parent to child
- `"child"` - Child to parent (alternative)
- `"sibling"` - Brother/sister
- `"adoptive-parent"` - Adopted parent
- `"step-parent"` - Step parent
- `"guardian"` - Legal guardian

### 5. Handling Complex Scenarios

**Twins:** Create two separate persons with same birthDate, link both to parents

**Half-siblings:** They share one parent but not both

**Step-siblings:** Share no biological parents, but parents are married

**In-laws:** Create person nodes for in-laws and spouse relationships

## Graph Layout Tips

1. **Generations:** Position people vertically by generation (grandparents top, grandchildren bottom)
2. **Families:** Group family units horizontally
3. **Relationship Nodes:** Position between the connected persons
4. **Children:** Position below the parent relationship, centered
5. **Use React Flow's layout algorithms** or implement custom hierarchical layout

## Example API Usage

### Creating a Family

```typescript
// 1. Create persons
const john = await api.createPerson({
  firstName: "John",
  lastName: "Doe",
  birthDate: "1985-05-15",
  gender: "male"
});

const jane = await api.createPerson({
  firstName: "Jane",
  lastName: "Smith",
  birthDate: "1987-08-20",
  gender: "female"
});

const sarah = await api.createPerson({
  firstName: "Sarah",
  lastName: "Doe",
  birthDate: "2012-09-20",
  gender: "female"
});

// 2. Create marriage relationship
await api.createRelationship({
  type: "spouse",
  fromPersonId: john.id,
  toPersonId: jane.id,
  startDate: "2010-06-15",
  notes: "Married at City Hall"
});

// 3. Create parent-child relationships
await api.createRelationship({
  type: "parent",
  fromPersonId: john.id,
  toPersonId: sarah.id,
  startDate: "2012-09-20"
});

await api.createRelationship({
  type: "parent",
  fromPersonId: jane.id,
  toPersonId: sarah.id,
  startDate: "2012-09-20"
});
```

## Visual Styling

The GraphViewPage automatically color-codes relationships:

- **Pink:** Spouse/Marriage relationships
- **Green:** Parent-Child relationships
- **Orange:** Sibling relationships
- **Blue:** Other relationships

## Future Enhancements

Consider implementing:

1. **Automatic Layout:** Hierarchical layout algorithm for generations
2. **Child Linking:** Automatically connect children from spouse relationships
3. **Timeline View:** Show relationships chronologically
4. **Family Groups:** Collapsible family unit nodes
5. **Pedigree Chart:** Traditional pedigree/ancestor chart view
6. **Descendant Chart:** Tree showing all descendants from an ancestor
7. **DNA/Genetic Links:** Mark biological vs. non-biological relationships

## Summary

The key to modeling family trees effectively is:

1. **Use Relationship Nodes** as first-class entities, not just edges
2. **Store rich metadata** in relationships (dates, notes, type)
3. **Be consistent** with relationship direction and types
4. **Link children** either through separate parent relationships or from spouse nodes
5. **Use the graph** to visualize complex multi-generational families

This approach allows you to:
- Track detailed relationship information (when married, divorced, etc.)
- Handle complex scenarios (multiple marriages, adoptions, etc.)
- Visualize the entire family network interactively
- Query relationships programmatically (find all children, siblings, etc.)
