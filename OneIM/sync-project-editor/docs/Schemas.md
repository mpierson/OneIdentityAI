# Schemas

A schema defines the objects, properties, and methods associated with the data in a system.  Each synchronization project must have two schemas: one schema representing the data in Identity Manager and another schema representing data in the target system.

Typical steps required for each synchronization project:

1. create two new schemas: one for Identity Manager and another for the target system
2. create schema types for each new schema, representing the objects that will be synchronized, e.g. Person, Department, Account, etc.
3. add schema properties to each type, representing attributes that will be involved in synchronization, e.g. Person.FirstName (Identity Manager) and Employee.FirstName (target system)
4. add schema methods to each type, representing the actions that can be performed in each system
5. add schema classes for each type, defining sub-sets of each type of object, e.g. Person-All, Person-Active


## Create a new schema


```bash
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

Use the values in [Common Language Runtime Type Identifiers](CLRIdentifiers.md) to determine the appropriate .Net CLR identifier of the system's schema, or use the SPEd command _clr_ to lookup the identifier.  The CLR type must expose the _VI.Projector.Schema.ISchema_ interface.

When creating a schema object for a custom Powershell connector, use the CLR type _VI.Projector.Powershell.PoshSchema_. 



SPEd provides a shortcut to create the Identity Manager schema:

```bash
sped -C my_db.yaml schema -s '4A82024A-2211-4D36-96CB-9C078B1E5E93' insert-oneim-schema --name 'Identity Manager'
```

Parameters

- shell: UID\_DPRShell of the synchronization project
- name (n): name of the new schema

The required .Net CLR identifier and system id will be assigned to the new schema.


## Create a new schema type

Each type of object in a system to be synchronized is represented in the schema as a schema type.

```bash
sped -C my_db.yaml schema-type --schema-id '9E51EFE9-761C-4D53-8733-9476051262BC' \
        insert --name "Person" \
               --clr-name 'VI.Projector.Database.DatabaseSchemaType'
```

Parameters

- schema-id: UID\_DPRSchema of the parent schema
- name (n): name of the new schema type
- clr-name: .Net CLR identifier of the schema type

The schema type name should correspond to the name used in the corresponding system. For example, in the Identity Manager schema, use the view or table name as the schema type name.

Use the values in [Common Language Runtime Type Identifiers](CLRIdentifiers.md) to determine the appropriate .Net CLR identifier of the schema object, or use the SPEd command _clr_ to lookup the identifier.  The CLR type must expose the _VI.Projector.Schema.ISchemaType_ interface.
Identity Manager schema types should use the _VI.Projector.Database.DatabaseSchemaType_ CLR identifier.

Each schema requires at least one schema type.

**Note**: every Identity Manager schema should include the QBMVTableRevision type, to support revision tracking.

## Add properties to a schema type

Each attribute involved in synchronization is represented by a schema property. 

```bash
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

Use the values in [Common Language Runtime Type Identifiers](CLRIdentifiers.md) to determine the appropriate .Net CLR identifier of the schema property, or use the SPEd command _clr_ to lookup the identifier.  The CLR type must expose the _VI.Projector.Schema.ISchemaProperty_ interface.
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

```bash
sped -C my_db.yaml schema-type add-oneim-properties --id 'BCBFAD07-4D94-4214-B766-4C8DF84092A6'
```

Parameters:

- id: UID\_DPRSchemaType of the parent type


## Add one or more methods to a schema type

The actions that can be performed on a schema type are represented by schema methods.  Schemas that are read-only do not require a method.

```bash
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

Use the values in [Common Language Runtime Type Identifiers](CLRIdentifiers.md) to determine the appropriate .Net CLR identifier of the schema method, or use the SPEd command _clr_ to lookup the identifier.  The CLR type must expose the _VI.Projector.Schema.ISchemaMethod_ interface.
Identity Manager schema properties should use the _VI.Projector.Database.DatabaseSchemaMethod_ CLR identifier.

SPEd provides a sub-command to add multiple methods to a schema:

```bash
sped -C my_db.yaml schema-type --schema-id '9E51EFE9-761C-4D53-8733-9476051262BC' \
        add-methods --id 'BCBFAD07-4D94-4214-B766-4C8DF84092A6' \
                    --methods "Insert Delete" \
                    --clr-name "VI.Projector.Powershell.Schema.PoshSchemaMethod"
```

For Identity Manager schemas, use the _all_ flag to add all available methods:

```bash
sped -C my_db.yaml schema-type --schema-id '9E51EFE9-761C-4D53-8733-9476051262BC' \
        add-methods --id 'EB452DE0-E324-4D0B-BB45-1E636016D426' \
                    --all \
                    --clr-name "VI.Projector.Database.DatabaseSchemaMethod"
```



## Add one or more classes to a schema type

Schema classes are used to define subsets of schema data.  In most cases, SPEd should be used only to create the default 'all' class for a schema type, i.e. class with no filter.  Each schema type requires at least one class.

In most cases, use the _schema-type add-default-class_ command to add the default 'all' class to a schema type:

```bash
sped -C my_db.yaml schema-type add-default-class --id 'EB452DE0-E324-4D0B-BB45-1E636016D426'
```

Parameters:

- id: UID\_DPRSchemaType of the parent schema type


If additional classes are required, the name and CLR type can be specified with the _schema-class insert_ command:

```bash
sped -C my_db.yaml schema-class --schema-type-id 'BBE236A6-67B9-4D9D-A49D-89EE5DF2F0E3' \
        insert --name 'Employee (ALL)' \
               --clr-name 'VI.Projector.Schema.GenericSchemaClass'
```

Parameters:

- schema-type-id: UID\_DPRSchemaType of the parent schema type
- name (n): name of new class
- clr-name: .Net CLR identifier of the schema class


Use the values in [Common Language Runtime Type Identifiers](CLRIdentifiers.md) to determine the appropriate .Net CLR identifier of the schema class, or use the SPEd command _clr_ to lookup the identifier.  The CLR type must expose the _VI.Projector.Schema.ISchemaClass_ interface.
Identity Manager schema properties should use the _VI.Projector.Database.DatabaseSchemaClass_ CLR identifier.


