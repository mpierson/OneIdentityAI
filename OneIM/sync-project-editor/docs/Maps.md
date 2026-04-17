# System Maps (DPRSystemMap, DPRSystemMappingRule)

Use a System Map to define mapping of attributes between a target system schema class and Identity Manager schema class. In most cases, each schema class will have one map to a schema class in the other system.  Convention is to refer to the Identity Manager as the left side of the map, and the target system schema as the right side of the map.

A System Mapping Rule is used to define the relationship between an attribute on the left and an attribute on the right.  Each System Map should have at least one, and usually many System Mapping Rules.

## Creating a system map


```bash
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



## Adding mapping rules


```bash
sped -C my_db.yaml mapping-rule --map-id 'B72BC648-937B-495F-9240-F1E04FDAD276' \
        insert --name "FirstName_firstName" \
               --left-property "FirstName" --right-property "firstName" \
               --direction "Inherite"
```

Parameters

- map-id: UID\_DPRSystemMap of the parent map
- name: name of the mapping rule
- left-property: name of the attribute in the left schema
- right-property: name of the attribute in the right schema
- direction: mapping direction valid for this rule, one of Inherite, ToTheLeft, ToTheRight, DoNotMap (default is Inherite)

**Note**: data types of the left and right properties must be the same.

Mapping directions:

_Inherite_
: Mapping rule will inherit the direction associated with the parent DPRSystemMap.

_ToTheLeft_
: Mapping rule is only considered when synchronizing data from right to left, i.e. in to Identity Manager; this direction is suitable for attributes that are never updated in the target system.

_ToTheRight_
: Mapping rule is only considered when synchronizing data from left to right, i.e. in to the target system; this direction is suitable for attributes that are always generated in Identity Manager, e.g. UID\_Person.

_DoNotMap_
: Mapping rule is not used.


### Key-based mapping rules

In scenarios where the property mapping is not direct, e.g. a value conversion or value lookup is needed, a virtual attribute may be necessary.

When a value lookup is needed to map a target system property to an Identity Manager property, use the following sped command:

```bash
sped -C my_db.yaml mapping-rule --map-id 'B72BC648-937B-495F-9240-F1E04FDAD276' \
        add-key-based-rule --name "UID_PersonHead to manager" \
               --left-property "UID_PersonHead" \
               --right-property "manager" \
               --lookup-table Person \
               --right-key-attribute PersonnelNumber \
               --left-key-attribute UID_Person
```

Parameters

- map-id: UID\_DPRSystemMap of the parent map
- name: name of the mapping rule
- left-property: name of the attribute in the left schema
- right-property: name of the attribute in the right schema
- lookup-table: Identity Manager table containing the value used in the map
- right-key-attribute: name of attribute in lookup table that contains the value in right side of the mapping (target system)
- left-key-attribute: name of attribute in lookup table that will be used in the left side of the mapping (Identity Manager)

The sample code above will create a mapping from _manager_ in the target system to _UID\_PersonHead_ in Identity Manager, with a lookup of _Person.UID\_Person_ using the manager's _PersonnelNumber_.

**Note**: all tables referenced in key-based mapping rules must be added to the Identity Manager schema using the SPEd command _schema-type insert_.  Attributes referenced in the mapping should be added as schema properties using the SPEd command _schema-property insert_ or _schema-type add-oneim-properties_ commands.


## Adding matching rules

Similar to schema mapping rules, object matching rules define the attribute(s) on each side of the map that should be used to correlate objects in Identity Manager with objects in the target system. Properties are often included in a map as part of a mapping rule _and_ a matching rule.


```bash
sped -C my_db.yaml mapping-rule --map-id 'B72BC648-937B-495F-9240-F1E04FDAD276' \
        insert-matching-rule --name "EmployeeIdKey" \
               --left-property "PersonnelNumber" --right-property "EmployeeId" \
               --add-mapping-rule true
```

Parameters

- map-id: UID\_DPRSystemMap of the parent map
- name: name of the matching rule
- left-property: name of the attribute in the left schema
- right-property: name of the attribute in the right schema
- add-mapping-rule: if true, a standard mapping rule will be created for the same attributes identified in the matching rule




