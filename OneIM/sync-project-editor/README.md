# Synchronization Project Editor (SPEd)

For AI agents (or humans), SPEd is a command line tool capable of viewing, inserting, and updating Identity Manager synchronization projects.

**Warning**: do not use in production environments. This tool does not use Identity Manager system accounts and by-passes the Identity Manager permissions model.  Use only on non-production environments.

Creating a sync project typically includes the following steps:

1. create project (DPRShell)
2. create default variable set and variables (DPRSystemVariableSet, DPRSystemVariable)
2. update project's default variable set (DPRShell.UID\_DPRSystemVariableSetDef)
2. create Identity Manager schema, types, methods, properties, and classes
3. create Identity Manager connection (DPRSystemConnection)
3. create target system schema, types, methods, properties, and classes
4. create target system connection
5. create schema map(s), including mapping rules and matching rules
5. create workflow, and workflow steps
5. create start info and schedule
6. mark project construction as complete


## Updating objects with SPEd

Updating object with SPEd is done with a JSON format payload.

For example, to update synchronization project (shell) with a default set of variables, use the _shell_ command with the _update_ sub-command and a JSON payload:

```
sped -C my_db.yaml shell update --id '4A82024A-2211-4D36-96CB-9C078B1E5E93' \
                                --content '{"UID_DPRSystemVariableSetDef": "AAAA-BBB-..."}'
```



# Project 

Top level container for sync project objects.

Class:
  DPRShell

## Create a new Shell

```
sped -C my_db.yaml shell insert --name TestProject
```

Parameters

- name (n): name of the sync project


# Variables

Variables are used to store environment specific data, e.g. host names, target system credentials.

Classes:

- DPRSystemVariableSet: container for individual variables
- DPRSystemVariable: variable's name-value pair, plus metadata


## Create a new variable set

```
sped -C my_db.yaml variable-set --shell '4A82024A-2211-4D36-96CB-9C078B1E5E93' insert --name "SampleVariableSet"
```

Parameters

- shell (s): UID\_DPRShell of the parent sync project
- name (n): name of the new variable set


## Add variables to a set


```
sped -C my_db.yaml variable --variable-set B99711CF-27EE-484E-AA0C-392D5F76D78A \
        insert -n "Host" --value "10.0.0.100"
```

Parameters

- variable-set: UID\_DPRSystemVariableSet of the parent collection
- name (n): name of the new variable
- value: value of variable (optional)
- secret: value is sensitive


Use the _secret_ option to create a variable that contains a password or other sensitive information:

```
sped -C my_db.yaml variable --variable-set B99711CF-27EE-484E-AA0C-392D5F76D78A \
        insert -n "Password" --secret --value $HOST_PASSWORD
```

# Schemas

A schema defines the objects, properties, and methods associated with the data in a system.  Each synchronization project must have two schemas: one schema representing the data in Identity Manager and another schema representing data in the target system.

Typical steps required for each synchronization project:

1. create two new schemas: one for Identity Manager and another for the target system
2. create schema types for each new schema, representing the objects that will be synchronized, e.g. Person, Department, Account, etc.
3. add schema properties to each type, representing attributes that will be involved in synchronization, e.g. Person.FirstName (Identity Manager) and Employee.FirstName (target system)
4. add schema methods to each type, representing the actions that can be performed in each system
5. add schema classes for each type, defining sub-sets of each type of object, e.g. Person-All, Person-Active


## Create a new schema


```
sped -C my_db.yaml schema -shell '4A82024A-2211-4D36-96CB-9C078B1E5E93' \
            insert  --name 'Target System A' \
                    --clr-name 'VI.Projector.Powershell.PoshSchema' \
                    --system-id 'TargetSystemA'
```

Parameters

- shell: UID\_DPRShell of the synchronization project
- name (n): name of the new schema
- clr-name: .Net CLR identifier of the schema 
- system-id: unique identifier for the target system, e.g. tenant ID or FQDN

Use the values in the Appendix below to determine the appropriate .Net CLR identifier of the system's schema, or use the SPEd command _clr_ to lookup the identifier.  The CLR type must expose the _VI.Projector.Schema.ISchema_ interface.

When creating a schema object for a custom Powershell connector, use the CLR type _VI.Projector.Powershell.PoshSchema_. 

SPEd provides a shortcut to create the Identity Manager schema:

```
sped -C my_db.yaml schema -s '4A82024A-2211-4D36-96CB-9C078B1E5E93' insert-oneim-schema --name 'Identity Manager'
```

Parameters

- shell: UID\_DPRShell of the synchronization project
- name (n): name of the new schema

The required .Net CLR identifier and system id will be assigned to the new schema.


## Create a new schema type

Each type of object in a system to be synchronized is represented in the schema as a schema type.

```
sped -C my_db.yaml schema-type --schema-id '9E51EFE9-761C-4D53-8733-9476051262BC' \
        insert --name "Person" \
               --clr-name 'VI.Projector.Database.DatabaseSchemaType'
```

Parameters

- schema-id: UID\_DPRSchema of the parent schema
- name (n): name of the new schema type
- clr-name: .Net CLR identifier of the schema type

The schema type name should correspond to the name used in the corresponding system. For example, in the Identity Manager schema, use the view or table name as the schema type name.

Use the values in the Appendix below to determine the appropriate .Net CLR identifier of the schema object, or use the SPEd command _clr_ to lookup the identifier.  The CLR type must expose the _VI.Projector.Schema.ISchemaType_ interface.
Identity Manager schema types should use the _VI.Projector.Database.DatabaseSchemaType_ CLR identifier.

Each schema requires at least one schema type.


## Add properties to a schema type

Each attribute involved in synchronization is represented by a schema property. 

```
sped -C my_db.yaml schema-property --schema-type-id 'BBE236A6-67B9-4D9D-A49D-89EE5DF2F0E3' \
        insert --name 'firstName' \
               --clr-name 'VI.Projector.Powershell.PoshSchemaProperty' \
               --data-type string
```

Parameters:

- schema-type-id: UID\_DPRSchemaType of the parent type
- name (n): name of the new schema property
- clr-name: .Net CLR identifier of the schema property 
- data-type: data type of the new property (_string_, _Integer_, _Boolean_, etc)


The schema property name should correspond to the name used in the corresponding system. For example, in the Identity Manager schema, use the column name as the schema property name.

Use the values in the Appendix below to determine the appropriate .Net CLR identifier of the schema property, or use the SPEd command _clr_ to lookup the identifier.  The CLR type must expose the _VI.Projector.Schema.ISchemaProperty_ interface.
Identity Manager schema properties should use the _VI.Projector.Database.DatabaseSchemaProperty_ CLR identifier.

The following schema property data types are supported, shown with corresponding SQL Server column types:

string
: VARCHAR
: NVARCHAR
: NCHAR
: CHAR

Binary
: VARBINARY

Integer
: INT

Integer
: BIT

Float
: FLOAT

DateTime
: DATETIME

Boolean
: BOOL



Each schema type requires at least one property.


SPEd provides an easy way to add properties to the Identity Manager schema type.  The _add-oneim-properties_ sub-command adds all columns of the schema type's corresponding table or view as properties:

```
sped -C my_db.yaml schema-type --schema-id '9E51EFE9-761C-4D53-8733-9476051262BC' \
        add-oneim-properties --id 'BCBFAD07-4D94-4214-B766-4C8DF84092A6'
```

## Add one or more methods to a schema type

The actions that can be performed on a schema type are represented by schema methods.  Schemas that are read-only do not require a method.

```
sped -C my_db.yaml schema-method --schema-type-id 'BCBFAD07-4D94-4214-B766-4C8DF84092A6' \
        insert --name 'Update' --clr-name "VI.Projector.Powershell.Schema.PoshSchemaMethod"
```

Parameters:

- schema-type-id: UID\_DPRSchemaType of the parent type
- name (n): method (Insert, Update, ...)
- clr-name: .Net CLR identifier of the schema method

List of available methods:

- Insert
- Update
- Delete
- MarkAsOutstanding (Identity Manager schemas only)
- UnmarkAsOutstanding (Identity Manager schemas only)

Use the values in the Appendix below to determine the appropriate .Net CLR identifier of the schema method, or use the SPEd command _clr_ to lookup the identifier.  The CLR type must expose the _VI.Projector.Schema.ISchemaMethod_ interface.
Identity Manager schema properties should use the _VI.Projector.Database.DatabaseSchemaMethod_ CLR identifier.

SPEd provides a sub-command to add multiple methods to a schema:

```
sped -C my_db.yaml schema-type --schema-id '9E51EFE9-761C-4D53-8733-9476051262BC' \
        add-methods --id 'BCBFAD07-4D94-4214-B766-4C8DF84092A6' \
                    --methods "Insert Delete" \
                    --clr-name "VI.Projector.Powershell.Schema.PoshSchemaMethod"
```

For Identity Manager schemas, use the _all_ flag to add all available methods:

```
sped -C my_db.yaml schema-type --schema-id '9E51EFE9-761C-4D53-8733-9476051262BC' \
        add-methods --id 'EB452DE0-E324-4D0B-BB45-1E636016D426' \
                    --all \
                    --clr-name "VI.Projector.Database.DatabaseSchemaMethod"
```



## Add one or more classes to a schema type

Schema classes are used to define subsets of schema data.  In most cases, SPEd should be used only to create the default 'all' class for a schema type, i.e. class with no filter.  Each schema type requires at least one class.

```
sped -C my_db.yaml schema-class --schema-type-id 'BBE236A6-67B9-4D9D-A49D-89EE5DF2F0E3' \
        insert --name 'Employee (ALL)' \
               --clr-name 'VI.Projector.Schema.GenericSchemaClass'
```

Parameters:

- schema-type-id: UID\_DPRSchemaType of the parent schema type
- name (n): name of new class
- clr-name: .Net CLR identifier of the schema class


Use the values in the Appendix below to determine the appropriate .Net CLR identifier of the schema class, or use the SPEd command _clr_ to lookup the identifier.  The CLR type must expose the _VI.Projector.Schema.ISchemaClass_ interface.
Identity Manager schema properties should use the _VI.Projector.Database.DatabaseSchemaClass_ CLR identifier.


SPEd also provides a command to add the default 'all' class to a schema type:

```
sped -C my_db.yaml schema-type --schema-id '9E51EFE9-761C-4D53-8733-9476051262BC' \
        add-default-class --id 'EB452DE0-E324-4D0B-BB45-1E636016D426'
```

# Connections

System connection details, e.g. host address and credentials, are stored in a Connection object.  Both schemas in a synchronization project require a connection object.

## Identity Manager connection

To create a connection object for the Identity Manager system:

```
sped -C my_db.yaml connection --shell '4A82024A-2211-4D36-96CB-9C078B1E5E93' \
            insert-oneim-connection --connection-string $ONEIM_CONSTRING
```

Parameters

- shell: UID\_DPRSHell of the synchronization project
- connection-string: connection string for the Identity Manager database

Identity Manager connection strings for synchronization typically take the form:

`Authentication=ProjectorAuthenticator;data source=<host name>;DBFactory="VI.DB.ViSqlFactory, VI.DB";initial catalog=<database name>;integrated security=False;user id=<username>;password[S]=<password>;pooling=False`

## Target system connection

To create a connection object for the target system:

```
sped -C my_db.yaml connection --shell '4A82024A-2211-4D36-96CB-9C078B1E5E93' \
        insert-target-system-connection --connector-type 'VI.Projector.Powershell.PoshConnectorDescriptor' \
                                        --connection-string $SYSTEM_CONSTRING
```

Parameters

- shell: UID\_DPRSHell of the synchronization project
- connector-type: CLR id of the target system connector
- connection-string: system-specific connection string for the target system

The connection string will be passed to the target system connector at runtime, and typically contains host connection details.  For Powershell connectors, the connection string also contains the base64-encoded XML connector definition.

Sample connection string format for a Powershell connector:

`ClassName=MyCustomConnector;CommaSeparatedDLLNames=MyConnector.dll;ConnectionPoolSize=1;DefinitionXml=<base64-encoded xml>;FolderContainingDLLs[V]=CP_Posh_FolderContainingDLLs;Hostname[V]=CP_Posh_Hostname;Username[V]=CP_Posh_Username;Password[V]=CP_Posh_Password;Namespace=com.acme.myconnector;SystemId=MyConnector`

The `[V]` designation implies that the connection parameter will be defined as a system variable.

# System Maps

Use a System Map to define mapping of attributes between a target system schema class and Identity Manager schema class. In most cases, each schema class will have one map to a schema class in the other system.  Convention is to refer to the Identity Manager as the left side of the map, and the target system schema as the right side of the map.

A System Mapping Rule is used to define the relationship between an attribute on the left and an attribute on the right.  Each System Map should have at least one, and usually many System Mapping Rules.

## Creating a system map


```
sped -C my_db.yaml map --shell '4A82024A-2211-4D36-96CB-9C078B1E5E93' \
        insert --name 'Person_Employee' \
               --left-schema-class-id 'BE98E9B0-37F1-41D8-BFF2-EE0F34F03E9C' \
               --right-schema-class-id 'BBF3338E-2185-4E0F-B01D-550930FED369' \
               --direction 'BothDirections'
```

Parameters

- shell: UID\_DPRSHell of the synchronization project
- left-schema-class-id: UID\_DPRSchemaClass of class on the left side (Identity Manager) of the map
- right-schema-class-id: UID\_DPRSchemaClass of class on the right side (target system) of the map
- direction: direction of data flow, one of _ToTheLeft_, _ToTheRight_, or the default _BothDirections_


SPEd also supports creating a map using schema class names:

```
sped -C my_db.yaml map --shell '4A82024A-2211-4D36-96CB-9C078B1E5E93' \
        insert-by-name --name 'Person_Employee2' \
                       --left "Person_Master" --right "Employee_Master"
```

Parameters

- shell: UID\_DPRSHell of the synchronization project
- left: name of class on the left side (Identity Manager) of the map
- right: name of class on the right side (target system) of the map
- direction: direction of data flow, one of _ToTheLeft_, _ToTheRight_, or the default _BothDirections_


## Adding mapping rules


```
sped -C my_db.yaml mapping-rule --map-id 'B72BC648-937B-495F-9240-F1E04FDAD276' \
        insert --name "FirstName_firstName" \
               --left-property "FirstName" --right-property "firstName" \
               --clr-name 'VI.Projector.Mapping.Rules.SinglePropertyComparisonRule'
```

Parameters

- map-id: UID\_DPRSystemMap of the parent map
- name: name of the mapping rule
- left-property: name of the attribute in the left schema
- right-property: name of the attribute in the right schema
- clr-name: .Net CLR identifier of the mapping rule (default is _VI.Projector.Mapping.Rules.SinglePropertyComparisonRule_)

Use the values in the Appendix below to determine the appropriate .Net CLR identifier of the mapping rule, or use the SPEd command _clr_ to lookup the identifier.  The CLR type must expose the _VI.Projector.Mapping.ISystemMappingRule_ interface.  Simple mapping rules, e.g. map of single value attributes, should use _VI.Projector.Mapping.Rules.SinglePropertyComparisonRule_.


## Adding matching rules

Similar to schema mapping rules, object matching rules define the attribute(s) on each side of the map that should be used to correlate objects in Identity Manager with objects in the target system. Properties are often included in a map as part of a mapping rule _and_ a matching rule.


```
sped -C my_db.yaml mapping-rule --map-id 'B72BC648-937B-495F-9240-F1E04FDAD276' \
        insert-matching-rule --name "EmployeeIdKey" \
               --left-property "PersonnelNumber" --right-property "EmployeeId" \
               --clr-name 'VI.Projector.Mapping.Rules.SinglePropertyComparisonRule'
```

Parameters

- map-id: UID\_DPRSystemMap of the parent map
- name: name of the matching rule
- left-property: name of the attribute in the left schema
- right-property: name of the attribute in the right schema
- clr-name: .Net CLR identifier of the mapping rule (default: _VI.Projector.Mapping.Rules.SinglePropertyComparisonRule_)




# Workflows

Workflow objects define the mechanics of a synchronization event. Each workflow contains one or more workflow steps, where a step represents the synchronization of one pair of schema types.  

The following steps are typically required to create a functioning workflow:

1. create the workflow container
2. add the Identity Manager and target system connections to workflow
3. create workflow step(s), one for each object type
4.  assign required actions to each step


## Create a workflow

```
sped -C my_db.yaml workflow --shell '4A82024A-2211-4D36-96CB-9C078B1E5E93' insert --name 'Full Synchronization'
```

Parameters

- shell: UID\_DPRSHell of the synchronization project
- name: name of the new workflow


## Add connections to the workflow

Each step in the workflow requires a connection to both the left and right systems, but first any connections in scope must be associated with the workflow.

```
sped -C my_db.yaml workflow --shell '4A82024A-2211-4D36-96CB-9C078B1E5E93' \
        add-connection --id '74ACD0C3-57AB-4F8B-8586-14F759757C49' \
                       --connection-id 'FC39CF49-3D68-4251-9004-7458A1E61334'
```

Parameters

- shell: UID\_DPRSHell of the synchronization project
- id: UID\_DPRProjectionConfig of the workflow
- connection-id: UID\_DPRSystemConnection of the connection object


SPEd provides an easy way to add all connections in the project to an existing workflow:

```
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

```
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
- object exists in the right system but not the left, 
- and object exists in both systems but one or more attributes are different.

Use the _add-method_ sub-command to add an action to a workflow step:

```
sped -C my_db.yaml workflow-step 
        add-schema-method --id 'CCC-86B89729D8974C4CB015B230043BE172' \
                          --side Left \
                          --method Insert \
                          --match-set 'DifferenceLeftToRight'
```

Parameters

- id: UID\_DPRProjectionConfigStep of the target step
- side: specify the map side on which the action will apply (Left: Identity Manager, Right: target system)
- method: name of method (Insert, Update, Delete, MarkAsOutstanding, UnMarkAsOutstanding)
- match-set: name of comparison scenario (see below)

Valid match set scenarios:

- **DifferenceLeftToRight**: object exists in target system but not in Identity Manager
- **DifferenceRightToLeft**: object exists in Identity Manager but not in target system
- **IntersectionWithoutDifferences**: object is the same on both sides
- **IntersectionWithDifferences**: object exists on both sides, but one or more attributes are different



### Match Sets

In most cases, it is not necessary to manually create Match Sets.  When possible, the `--include-default-match-sets` flag should be used when creating a workflow step (see above).  Notes are provided below for synchronizations with unique requirements.

The synchronization may encounter four scenarios when comparing objects the two systems: objects are the same, object exists in the left system but not the right, object exists in the right system but not the left, and object exists in both systems but one or more attributes are different.  Each of these scenarios is represented by a Match Set.  Each workflow step should be assigned a collection of Match Sets, representing the scenarios that should be considered in a synchronization event.


Insert a new collection of Match Set objects:

```
sped -C my_db.yaml match-sets insert --name "SystemA_FullSync_Employee"
```

Parameters

- name: name of the new match set collection


To add all the default match sets (objects are the same, object exists in the left system but not the right, object exists in the right system but not the left, and object exists in both systems but one or more attributes are different) to a collection:

```
sped -C my_db.yaml match-sets add-default-sets --id 'BDB472A2-6788-4277-B456-23BCDE9A89BC'
```

Parameters

- id: UID\_DPRSystemObjectMatchSets of parent collection


To create a single match set:

```
sped -C my_db.yaml match-set insert --name 'DifferenceLeftToRight'
```

Parameters

- name: name of the new match set, must be one of DifferenceLeftToRight, DifferenceRightToLeft, IntersectionWithoutDifferences, IntersectionWithDifferences


# Start Infos

Start Info objects define a scheduled synchronization event, including the following attributes:

 - workflow
 - variable set
 - schedule
 - root object, if needed


```
sped -C my_db.yaml start-info --shell 'CCC-19F48527609980498D5E843FF49BB8AD' \
        insert --name 'Full Synchronization' \
               --variable-set-id 'B99711CF-27EE-484E-AA0C-392D5F76D78A' \
               --workflow-id 'CCC-7202478647387649AFE0B1E7F5351C22'
```

Parameters

- shell: UID\_DPRSHell of the synchronization project
- name: name of the new start info
- variable-set-id: UID\_DPRSystemVariableSet of the variables to be used by synchronization
- workflow-id: UID\_DPRProjectionConfig of the workflow to be used for synchronization


To use the default variables assigned to the project, use the _use-default-variables_ flag.  

The synchronization workflow can be identified by name using the _workflow-name_ flag.


```
sped -C my_db.yaml start-info --shell 'CCC-19F48527609980498D5E843FF49BB8AD' \
        insert --name 'Full Synchronization' \
               --use-default-variables \
               --workflow-name 'FullSynchronization'
```


## Schedules

Assign a schedule for synchronization using the _add-schedule_ sub-command:

```
sped -C my_db.yaml start-info --shell 'CCC-19F48527609980498D5E843FF49BB8AD' \
        add-schedule --id 'BC7DBA15-9B97-453C-ADC0-513027CA9E63' 
                     --type 'Month' \
                     --frequency 6 \
                     --time-zone UTC \
                     --start-time '01:00'
```

Parameters

- id: UID\_DPRProjectionStartInfo of the start info object
- type: type of schedule (Hour, Day, Week, Month, Year); default is _Day_
- frequency: how often to run the schedule, in terms of type, e.g. once every **6** months; default is _1_
- time-zone: short name of time zone for schedule; default is _UTC_
- start-time: time of day to run scheduled synchronization, in **hh:mm** 24hr format; default is midnight _00:00_


Add a schedule using defaults (every day at midnight UTC):

```
sped -C my_db.yaml start-info --shell 'CCC-19F48527609980498D5E843FF49BB8AD' \
        add-schedule -i 'BC7DBA15-9B97-453C-ADC0-513027CA9E63'
```

## Root Object

Some synchronizations target a specific system or domain in Identity Manager, e.g. synchronization of an Active Directory domain. Other synchronizations are not specific to a target system, e.g. synchronization of a Human Resource system into Identity Manager.  All scheduled synchronization events require a target object.

For synchronization of an Active Directory domain, LDAP domain, or a generic target system represented in Identity Manager's UNS tables, the root object will correspond to the systems UNSRoot record (ADSDomain, LDAPDomain, UNSRootB, ...).

Add this type of root object to the start info with the _add-root-object_ sub-command:

```
sped -C my_db.yaml start-info --shell 'CCC-19F48527609980498D5E843FF49BB8AD' \
        add-root-object --id '73C6842D-1CA5-468A-880E-5EF0C32DF4EA' \
                        --root-object-key '<Key><T>UNSRootB</T><P>aa669e4f-3d82-4882-9bb1-d88f3e412a3c</P></Key>' \
                        --connection-id 'CCC-1D6726110E33C941BBF9EE0C0480DB29' \
                        --variable-set-id 'CCC-C6DEBD8334E97C4BB709639DF649FBD1' \
                        --server-id '850649CD-003E-40CB-A1FD-F5D9C5C89529' 
```

Parameters

- id: UID\_DPRProjectionStartInfo of the start info object
- root-object-key: XObjectKey of the root object (ADSDomain, LDAPDomain, UNSRootB, etc.)
- connection-id: UID\_DPRSystemConnection of the connection to the target system associated with root object
- variable-set-id: UID\_DPRSystemVariableSet of the variables to be used with root object
- server-id: UID\_QBMServer of the Identity Manager job server that will perform the synchronization


Use the _use-default-connection_ flag to use the default target system connection to build the root object.  Use the _use-default-variables_ flag to use the default variable set.  Use the _server-name_ flag to reference a job server by name.

```
sped -C my_db.yaml start-info --shell 'CCC-19F48527609980498D5E843FF49BB8AD' \
        add-root-object --id '73C6842D-1CA5-468A-880E-5EF0C32DF4EA' \
                        --use-default-connection \
                        --use-default-variables \
                        --server-name 'IAMS03' \
                        --root-object-key '<Key><T>UNSRootB</T><P>aa669e4f-3d82-4882-9bb1-d88f3e412a3c</P></Key>'
```

For synchronization projects that use a target table instead of a target system, use the _table-name_ flag instead of the _root-object-key_ flag:

```
sped -C my_db.yaml start-info --shell 'CCC-19F48527609980498D5E843FF49BB8AD' \
        add-root-object --id '73C6842D-1CA5-468A-880E-5EF0C32DF4EA' \
                        --use-default-connection \
                        --use-default-variables \
                        --server-name 'IAMS03' \
                        --table-name 'Person'
```

# Appendix: Common Language Runtime Type Identifiers

| CLR ID                                            | Description                                                             | Notes                                              |
|---------------------------------------------------|-------------------------------------------------------------------------|----------------------------------------------------|
| VI.Projector.ProjectorShell                       | Synchronization project definition                                       |                                                    |
| VI.Projector.Variables.SystemVariableSet          | Collection of project variables                                         |                                                    |
| VI.Projector.Variables.SystemVariable             | Sync project variable                                                   |                                                    |
| VI.Projector.Database.DatabaseSchema              | Schema definition for database system                                   | interface: VI.Projector.Schema.ISchema             |
| VI.Projector.Powershell.PoshSchema                | Schema definition for custom Powershell connected system                | interface: VI.Projector.Schema.ISchema             |
| VI.Projector.ADS.ProjectorADConnectorSchema       | Schema definition for AD system                                         | interface: VI.Projector.Schema.ISchema             |
| VI.Projector.Database.DatabaseSchemaType          | Schema definition for an object in a database system                    | interface: VI.Projector.Schema.ISchemaType         |
| VI.Projector.Powershell.PoshSchema                | Schema definition for an object in a custom Powershell connected system | interface: VI.Projector.Schema.ISchemaType         |
| VI.Projector.ADS.ProjectorADConnectorSchema       | Schema definition for an object in a AD system                          | interface: VI.Projector.Schema.ISchemaType         |
| VI.Projector.Database.DatabaseSchemaClass         | Schema definition for a object class in a database system               | interface: VI.Projector.Schema.ISchemaClass        |
| VI.Projector.Schema.GenericSchemaClass            | Schema definition for a object class in a connected system              | interface: VI.Projector.Schema.ISchemaClass        |
| VI.Projector.Connection.ISystemConnection         | Connection to a system                                                  |                                                    |
| VI.Projector.Database.DatabaseConnectorDescriptor | Connector metadata for a database system                                | VI.Projector.Connection.ISystemConnectorDescriptor |
| VI.Projector.ADS.ProjectorADSConnectorDescriptor  | Connector metadata for an AD domain                                     | VI.Projector.Connection.ISystemConnectorDescriptor |
| VI.Projector.Powershell.PoshConnectorDescriptor   | Connector metadata for a custom Powershell connected system             | VI.Projector.Connection.ISystemConnectorDescriptor |
| VI.Projector.Mapping.SystemMap                    | Mapping of an object class between connected systems                     |                                                    |
| VI.Projector.Projection.ProjectionConfiguration   | Synchronization workflow                                                  |                                                    |
| VI.Projector.Projection.ProjectionStep            | One step in a synchronization workflow                                     |                                                    |
| VI.Projector.Projection.ProjectionStartInfo       | Metadata for synchronization execution                                        |                                                    |
