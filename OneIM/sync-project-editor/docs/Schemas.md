# Schemas

A schema defines the objects, properties, and methods associated with the data in a system.  Each synchronization project must have two schemas: one schema representing the data in Identity Manager and another schema representing data in the target system.

Typical steps required for each synchronization project:

1. create two new schemas: one for Identity Manager and another for the target system
2. create schema types for each new schema, representing the objects that will be synchronized, e.g. Person, Department, Account, etc.
3. add schema properties to each type, representing attributes that will be involved in synchronization, e.g. Person.FirstName (Identity Manager) and Employee.FirstName (target system)
4. add schema methods to each type, representing the actions that can be performed in each system
5. add schema classes for each type, defining sub-sets of each type of object, e.g. Person-All, Person-Active


## Create a new schema

To create a new schema for a custom target system:

```bash
sped -C my_db.yaml schema --shell '4A82024A-2211-4D36-96CB-9C078B1E5E93' \
            insert-target-schema  --name 'Target System A'  --system-id 'TargetSystemA'
```

Parameters

- shell: UID\_DPRShell of the synchronization project
- name (n): name of the new schema
- system-id: unique identifier for the target system, e.g. tenant ID or FQDN, used by Identity Manager to uniquely distinguish multiple systems of the same type

If successful, SPEd will return the id (UID\_DPRSchema) of the new schema.

To create the Identity Manager schema:

```bash
sped -C my_db.yaml schema --shell '4A82024A-2211-4D36-96CB-9C078B1E5E93' \
            insert-oneim-schema --name 'Identity Manager'
```

Parameters

- shell: UID\_DPRShell of the synchronization project
- name (n): name of the new schema

If successful, SPEd will return the id (UID\_DPRSchema) of the new schema.


## Create a new schema type

Each type of object in a system to be synchronized is represented in the schema as a schema type (DPR\_SchemaType).

```bash
sped -C my_db.yaml schema-type --schema-id '9E51EFE9-761C-4D53-8733-9476051262BC' \
        insert --name "Person"
```

Parameters

- schema-id: UID\_DPRSchema of the parent schema
- name (n): name of the new schema type

The schema type name should correspond to the name used in the corresponding system. For example, in the Identity Manager schema, use the view or table name as the schema type name.

Each schema requires at least one schema type.


## Add properties to a schema type

Each attribute involved in synchronization is represented by a schema property. 

Add properties to each schema type with the _schema-property insert_ command:

```bash
sped -C my_db.yaml schema-property --schema-type-id 'BBE236A6-67B9-4D9D-A49D-89EE5DF2F0E3' \
        insert --name 'employeeId' --data-type String --is-key true --is-secret false
```

Parameters:

- schema-type-id: UID\_DPRSchemaType of the parent type
- name (n): name of the new schema property
- data-type: data type of the new property (_String_, _Integer_, _Boolean_, etc)
- is-key: true if the property is a key field that uniquely identifies a record for the schema type
- is-secret: true if values of this property should be treated as secret


The schema property name should correspond to the name used in the corresponding system. For example, in the Identity Manager schema, use the column name as the schema property name.
For custom target systems, the schema property name should correspond to the attribute name that appears in the connector definition XML file.

The following schema property data types are supported, shown with corresponding SQL Server column types:

String
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



**Note**: SPEd provides an easy way to add properties to the Identity Manager schema type.  The _add-oneim-properties_ sub-command adds all columns of the schema type's corresponding table or view as properties:

```bash
sped -C my_db.yaml schema-type add-oneim-properties --id 'BCBFAD07-4D94-4214-B766-4C8DF84092A6'
```

Parameters:

- id: UID\_DPRSchemaType of the parent type


## Add one or more methods to a schema type

The actions that can be performed on a schema type are represented by schema methods.


```bash
sped -C my_db.yaml schema-method --schema-type-id 'BCBFAD07-4D94-4214-B766-4C8DF84092A6' \
        insert --name 'Update'
```

Parameters:

- schema-type-id: UID\_DPRSchemaType of the parent type
- name (n): method (Insert, Update, ...)

List of available methods:

- Insert
- Update
- Delete
- MarkAsOutstanding (Identity Manager schemas only)
- UnmarkAsOutstanding (Identity Manager schemas only)

In most cases, Identity Manager schemas will support all methods listed above.  

SPEd provides a sub-command to add multiple methods to a schema:

```bash
sped -C my_db.yaml schema-type --schema-id '9E51EFE9-761C-4D53-8733-9476051262BC' \
        add-methods --id 'BCBFAD07-4D94-4214-B766-4C8DF84092A6' \
                    --methods "Insert Delete" 
```

For Identity Manager schemas, use the _all_ flag to add all available methods:

```bash
sped -C my_db.yaml schema-type --schema-id '9E51EFE9-761C-4D53-8733-9476051262BC' \
        add-methods --id 'EB452DE0-E324-4D0B-BB45-1E636016D426' --all
```



## Add class to schema type

Schema classes are used to define subsets of schema data.  In most cases, SPEd should be used only to create the default 'all' class for a schema type, i.e. class with no filter.  Each schema type requires at least one class.

Use the _schema-type add-default-class_ command to add the default class to a schema type:

```bash
sped -C my_db.yaml schema-type add-default-class --id 'EB452DE0-E324-4D0B-BB45-1E636016D426'
```

Parameters:

- id: UID\_DPRSchemaType of the parent schema type

The _add-default-class_ sub-command returns the UID\_DPRSchemaClass of the new class.


Use the _schema-class show_ command to view details of the class created by _schema-type add-default-class_:

```bash
sped -C my_db.yaml schema-class --schema-type-id 'EB452DE0-E324-4D0B-BB45-1E636016D426' \
        show --id 'CCC-7DBEA18F6D654028A2CC31A0400065AD'
```

Parameters:

- schema-type-id: UID\_DPRSchemaType of the parent schema type
- id: UID_\DPRSchemaClass of the new class

