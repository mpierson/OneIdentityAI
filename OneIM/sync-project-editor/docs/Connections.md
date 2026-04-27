# Connections (DPRSystemConnection)

System connection details, e.g. host address and credentials, are stored in a Connection object.  Both schemas in a synchronization project require a connection object.

## Identity Manager connection

To create a connection object for the Identity Manager system:

```bash
sped -C my_db.yaml connection --shell '4A82024A-2211-4D36-96CB-9C078B1E5E93' \
            insert-oneim-connection --parameters $ONEIM_CONSTRING
```

Parameters

- shell: UID\_DPRSHell of the synchronization project
- parameters: connection parameters required by Identity Manager

Identity Manager connection strings for synchronization typically take the form:

`Authentication=ProjectorAuthenticator;data source=<host name>;DBFactory="VI.DB.ViSqlFactory, VI.DB";initial catalog=<database name>;integrated security=False;user id=<username>;password[S]=<password>;pooling=False`

Attributes in a typical Identity Manager connection string:

- _Authentication_: type of authentication used in synchronization; custom connectors should use `ProjectorAuthenticator`
- _data source_: host name of SQL Server server hosting the Identity Manager database
- _DBFactory_: .Net SQL connection factory; custom connectors should use `VI.DB.ViSqlFactory, VI.DB`
- _initial catalog_: name of Identity Manager database
- _integrated security_: true if authentication is performed via Kerberos using service account; typically `False`
- _user id_: service account name
- _password_: service account password
- _pooling_: true if the SQL Server driver should pool connections; should be `False`

**Note**: all attributes listed above are mandantory.

## Target system connection

To create a connection object for the target system:

```bash
sped -C my_db.yaml connection --shell '4A82024A-2211-4D36-96CB-9C078B1E5E93' \
        insert-target-system-connection --parameters $MY_SYSTEM_PARAMS
```

Parameters

- shell: UID\_DPRSHell of the synchronization project
- parameters: system-specific connection parameters for the target system

The connection parameters will be passed to the target system connector at runtime, and typically contains host connection details.  For custom connectors, the parameters also contains a compressed version of the XML connector definition.

Format for a Powershell connector:

`SystemId=MyConnector;Namespace=com.acme.myconnector;ClassName=MyCustomConnector;CommaSeparatedDLLNames=MyConnector.dll;ConnectionPoolSize=1;DefinitionXml=<compressed xml>;FolderContainingDLLs[V]=CP_Posh_FolderContainingDLLs;[Other parameters required by connector, e.g. host, port, user name, password ...]`

The `[V]` designation implies that the connection parameter will be defined as a system variable and is included in the default variable set via SPEd.

Parameters required for a target system connection:

- _SystemId_: unique identifier of the target system, e.g. FQDN of service
- _ClassName_: name of the class that implements `ois.oneim.ConnectorBase.ConnectorBase.ConnectorInterface`
- _Namespace_: .Net class namespace of class identified in _ClassName_
- _CommaSeparatedDLLNames_: file names of connector DLL plus any dependencies
- _ConnectionPoolSize_: maximum size of connection pool, if implemented by custom connector; typically `1`
- _DefinitionXml_: connector definition XML, Base64 encoded and compressed, described below
- _FolderContainingDLLs_: local path of folder on job server that will be used to store connector DLLs; for most Identity Manager job servers folder will be `c:\Program Files\One Identity\Identity Manager`

Additional connection parameters required by connector should be included in the connection string, e.g. host name of server, service account credentials.

### Connector definition XML

SPEd can be used to encode the connector definition XML for use in the target system connection parameters:

```bash
sped connection compress-connector-definition -xml '<PowershellConnectorDefinition Version="1.0" Description="my connector"...'
```

Parameters

- xml: connector definition, in XML format

