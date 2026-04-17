# Workflows (DPRProjectionConfig, DPRProjectionConfigStep)

Workflow objects define the mechanics of a synchronization event. Each workflow contains one or more workflow steps, where a step represents the synchronization of one pair of schema types.  

The following steps are typically required to create a functioning workflow:

1. create the workflow container
2. create workflow step(s), one for each object type
3. assign required actions to each step


## Create a workflow

```bash
sped -C my_db.yaml workflow --shell '4A82024A-2211-4D36-96CB-9C078B1E5E93' \
        insert --name 'Full Synchronization' --direction 'ToTheLeft'
```

Parameters

- shell: UID\_DPRSHell of the synchronization project
- name: name of the new workflow
- direction: supported synchronization direction, one of Inherite, ToTheLeft, or ToTheRight

Synchronization direction notes:

_Inherite_
: Workflow supports both directions; direction will be determined by the Start Info(s) that reference the workflow

_ToTheLeft_
: Workflow supports synchronization from right to left, i.e. in to Identity Manager

_ToTheRight_
: Workflow supports synchronization from left to right, i.e. in to the target system



## Add workflow steps

Each step in a workflow synchronizes a schema type between the left and right systems.  Every workflow requires at least one step, and typically includes a step for each mapping.

Use the _workflow-step insert_ command to add a step to a workflow for a given schema type map:

```bash
sped -C my_db.yaml workflow-step --workflow-id '74ACD0C3-57AB-4F8B-8586-14F759757C49' \
        insert --name 'Person' \
               --map-id 'B72BC648-937B-495F-9240-F1E04FDAD276' \
               --source-is-authoritative true
```

Parameters

- workflow-id: UID\_DPRProjectionConfig of the parent workflow
- name: name of the new step
- map-id: UID\_DPRSystemMap of the schema class map
- source-is-authoritative: true if the incoming data is authoritative, for example, a human resource system


Synchronization actions should be configured for one or more of these data comparison scenarios (see Match Sets below):

- objects in both sides of the map are the same, 
- object exists in the left system but not the right, 
- object exists in the right system but not the left, and 
- object exists in both systems but one or more attributes are different.

Use the _add-method_ sub-command to add an action to a workflow step:

```bash
sped -C my_db.yaml workflow-step 
        add-schema-method --id 'CCC-86B89729D8974C4CB015B230043BE172' \
                          --side Left \
                          --method Insert \
                          --match-set 'DifferenceLeftToRight'
```

Parameters

- id: UID\_DPRProjectionConfigStep of the target step
- side: specify the map side on which the action will apply (Left: Identity Manager, Right: target system)
- method: name of method (one of: Insert, Update, Delete, MarkAsOutstanding, UnMarkAsOutstanding)
- match-set: name of comparison scenario (see below)

Valid match set scenarios:

- **DifferenceLeftToRight**: object exists in target system but not in Identity Manager
- **DifferenceRightToLeft**: object exists in Identity Manager but not in target system
- **IntersectionWithoutDifferences**: object is the same on both sides
- **IntersectionWithDifferences**: object exists on both sides, but one or more attributes are different


Typical actions and match set scenario combinations:

| Scenario                       | Side                    | Action                     |
|:-------------------------------|:------------------------|:---------------------------|
| DifferenceLeftToRight          | Left (Identity Manager) | Insert                     |
| DifferenceRightToLeft          | Left (Identity Manager) | MarkAsOutstanding          |
| IntersectionWithoutDifferences | n/a                     | n/a                        |
| IntersectionWithDifferences    | Left (Identity Manager) | Update, UnMarkAsOutstanding |


