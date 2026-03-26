# Workflows (DPRProjectionConfig, DPRProjectionConfigStep)

Workflow objects define the mechanics of a synchronization event. Each workflow contains one or more workflow steps, where a step represents the synchronization of one pair of schema types.  

The following steps are typically required to create a functioning workflow:

1. create the workflow container
2. add the Identity Manager and target system connections to workflow
3. create workflow step(s), one for each object type
4. assign required actions to each step


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



## Add connections to the workflow

Each step in the workflow requires a connection to both the left and right systems, but first any connections in scope must be associated with the workflow.

```bash
sped -C my_db.yaml workflow --shell '4A82024A-2211-4D36-96CB-9C078B1E5E93' \
        add-connection --id '74ACD0C3-57AB-4F8B-8586-14F759757C49' \
                       --connection-id 'FC39CF49-3D68-4251-9004-7458A1E61334'
```

Parameters

- shell: UID\_DPRSHell of the synchronization project
- id: UID\_DPRProjectionConfig of the workflow
- connection-id: UID\_DPRSystemConnection of the connection object


SPEd provides an easy way to add all connections in the project to an existing workflow:

```bash
sped -C my_db.yaml workflow --shell '4A82024A-2211-4D36-96CB-9C078B1E5E93' \
        add-all-connections --id '74ACD0C3-57AB-4F8B-8586-14F759757C49'
```

## Add workflow steps

One step in a workflow synchronizes a schema type between the left and right systems.  

The following components are required when creating a new step:

- parent workflow
- connection to the left (Identity Manager) system
- connection to the right (target) system
- system map between corresponding schema classes in each system
- match set objects (see below)

In most cases, the default system connections and default match sets are appropriate, so creation of a new step requires only the parent workflow and system map:

```bash
sped -C my_db.yaml workflow-step --workflow-id '74ACD0C3-57AB-4F8B-8586-14F759757C49' \
        insert --name 'Person' \
               --map-id 'B72BC648-937B-495F-9240-F1E04FDAD276' \
               --include-default-match-sets \
               --use-default-connections 
```

Parameters

- workflow-id: UID\_DPRProjectionConfig of the parent workflow
- name: name of the new step
- map-id: UID\_DPRSystemMap of the schema class map
- include-default-match-sets: use a new match set collection with all four matching scenarios (see below)
- use-default-connections: use the Identity Manager and target system connections

Note: if the parent workflow is associated with more than one Identity Manager connection or more than one target system connection, do not use the use-default-connections flag -- provide the connections explicitly (see `sped workflow-step insert -h`).

Synchronization actions can be configured for each of these four data comparison scenarios (see Match Sets below):

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


### Match Sets

In most cases, it is not necessary to manually create Match Sets.  When possible, the `--include-default-match-sets` flag should be used when creating a workflow step (see above).  Notes are provided below for synchronizations with unique requirements.

The synchronization may encounter four scenarios when comparing objects the two systems: 

- objects are the same, object exists in the left system but not the right, 
- object exists in the right system but not the left, and 
- object exists in both systems but one or more attributes are different.  

Each of these scenarios is represented by a Match Set.  Each workflow step should be assigned a collection of Match Sets, representing the scenarios that should be considered in a synchronization event.


Insert a new collection of Match Set objects:

```bash
sped -C my_db.yaml match-sets insert --name "SystemA_FullSync_Employee"
```

Parameters

- name: name of the new match set collection


To add all the default match sets (objects are the same, object exists in the left system but not the right, object exists in the right system but not the left, and object exists in both systems but one or more attributes are different) to a collection:

```bash
sped -C my_db.yaml match-sets add-default-sets --id 'BDB472A2-6788-4277-B456-23BCDE9A89BC'
```

Parameters

- id: UID\_DPRSystemObjectMatchSets of parent collection


To create a single match set:

```bash
sped -C my_db.yaml match-set insert --name 'DifferenceLeftToRight'
```

Parameters

- name: name of the new match set, must be one of DifferenceLeftToRight, DifferenceRightToLeft, IntersectionWithoutDifferences, IntersectionWithDifferences


