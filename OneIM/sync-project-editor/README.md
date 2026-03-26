# Synchronization Project Editor (SPEd)

For AI agents (or humans), SPEd is a command line tool capable of viewing, inserting, and updating Identity Manager synchronization projects.

**Warning**: do not use in production environments. This tool does not use Identity Manager system accounts and by-passes the Identity Manager permissions model.  Use only on non-production environments.

Items that **should always be confirmed** before creating a synchronization project:

- list of objects in the target system that will be synchronized to and from Identity Manager
- corresponding list of tables in the Identity Manager schema that will be synchronized to and from the target system
- restricted operations and read-only object attributes
- location of custom connector DLL and the connector definition XML
- connection string and SQL credentials for Identity Manager


Creating a synchronization project includes the following steps:

1. create project (DPRShell), see [Project.md](docs/Project.md)
2. create default variable set and variables (DPRSystemVariableSet, DPRSystemVariable), see [Variables.md](docs/Variables.md)
3. update project's default variable set (DPRShell.UID\_DPRSystemVariableSetDef)
4. create Identity Manager schema, types, methods, properties, and classes, see [Schemas.md](docs/Schemas.md)
5. create Identity Manager connection (DPRSystemConnection), see [Connections.md](docs/Connections.md)
6. create target system schema, types, methods, properties, and classes, see [Schemas.md](docs/Schemas.md)
7. create target system connection, see [Connections.md](docs/Connections.md)
8. create schema map(s), including mapping rules and matching rules, see [Maps.md](docs/Maps.md)
9. create workflow, and workflow steps, see [Workflows.md](docs/Workflows.md)
10. create start info and schedule, see [StartInfos.md](docs/StartInfos.md)
11. mark project construction as complete, see [Project.md](docs/Project.md)


## Updating objects with SPEd

Updating object with SPEd is done with a JSON format payload.

For example, to update synchronization project (shell) with a default set of variables, use the _shell_ command with the _update_ sub-command and a JSON payload:

```bash
sped -C my_db.yaml shell update --id '4A82024A-2211-4D36-96CB-9C078B1E5E93' \
                                --content '{"UID_DPRSystemVariableSetDef": "AAAA-BBB-..."}'
```

## SPEd configuration file

For convenience, Identity Manager connection parameters can be supplied to SPEd as a YAML configuration file, using the _-C_ command line flag.

The following command line flags can be provided as attributes in a YAML file:

- host 
- port
- database
- user
- password

**Note**: the password can be provided as an environment variable _SPED\_PASSWORD_ to avoid storing in YAML.


Sample YAML file:
```yaml
host: services-uscentral.skytap.com
port: 9050
database: ACME
user: svc_1im_sql
```



