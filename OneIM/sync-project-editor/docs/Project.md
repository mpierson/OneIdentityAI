# Synchronization Project (DPRShell)

Top level container for sync project objects.

Class:
  DPRShell

## Create a new Shell

```bash
sped -C my_db.yaml shell insert --name TestProject
```

Parameters

- name (n): name of the sync project

## Mark synchronization project as complete

The _IsFinalized_ attribute of a synchronization project's DPRShell record indicates the state of the project.  A value of _3_ indicates the project is ready to be used.

Update the project with `IsFinalized = 3` only when all other steps are complete, and project has been validated. 

```bash
sped -C my_db.yaml shell update --id '4A82024A-2211-4D36-96CB-9C078B1E5E93' \
                                --content '{"IsFinalized": 3}'
```

