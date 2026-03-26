# Creating Custom Identity Manager Connectors

Custom Identity Manager connectors are used to integrate Identity Manager with identity sources (e.g. HR system) or systems that require provisioning of identity information (target systems).

There are typically three components of a custom connector:

- **.Net shared library** that is capable of performing CRUD operations on the connected system, see [Connector Shared Library](ConnectorSharedLibrary.md)
- **XML connector definition** that contains a recipe for using the shared library, see [Connector Definition XML](ConnectorDefinition.md)
- **Identity Manager synchronization project** that defines the mapping and data flow to and from the connected system, see [Synchronization Project Editor](../sync-project-editor/README.md)

Additional details of Identity Manager custom connectors [here](https://github.com/OneIdentity/IdentityManager.PoSh-Connector-Guide/blob/main/README.md).





