# Variables (DPRSystemVariableSet, DPRSystemVariable)

Variables are used to store environment specific data, e.g. host names, target system credentials.

Classes:

- DPRSystemVariableSet: container for individual variables
- DPRSystemVariable: variable's name-value pair, plus metadata


## Create a new variable set

```bash
sped -C my_db.yaml variable-set --shell '4A82024A-2211-4D36-96CB-9C078B1E5E93' insert --name "SampleVariableSet"
```

Parameters

- shell (s): UID\_DPRShell of the parent sync project
- name (n): name of the new variable set


## Add variables to a set


```bash
sped -C my_db.yaml variable --variable-set B99711CF-27EE-484E-AA0C-392D5F76D78A \
        insert -n "Host" --value "10.0.0.100"
```

Parameters

- variable-set: UID\_DPRSystemVariableSet of the parent collection
- name (n): name of the new variable
- value: value of variable (optional)
- secret: value is sensitive


Use the _secret_ option to create a variable that contains a password or other sensitive information:

```bash
sped -C my_db.yaml variable --variable-set B99711CF-27EE-484E-AA0C-392D5F76D78A \
        insert -n "Password" --secret --value $HOST_PASSWORD
```

