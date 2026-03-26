# Connections (DPRSystemConnection)

System connection details, e.g. host address and credentials, are stored in a Connection object.  Both schemas in a synchronization project require a connection object.

## Identity Manager connection

To create a connection object for the Identity Manager system:

```bash
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

```bash
sped -C my_db.yaml connection --shell '4A82024A-2211-4D36-96CB-9C078B1E5E93' \
        insert-target-system-connection --connector-type 'VI.Projector.Powershell.PoshConnectorDescriptor' \
                                        --connection-string $SYSTEM_CONSTRING
```

Parameters

- shell: UID\_DPRSHell of the synchronization project
- connector-type: CLR id of the target system connector
- connection-string: system-specific connection string for the target system

The connection string will be passed to the target system connector at runtime, and typically contains host connection details.  For custom connectors, the connection string also contains a compressed version of the XML connector definition.

Sample connection string format for a Powershell connector:

`ClassName=MyCustomConnector;CommaSeparatedDLLNames=MyConnector.dll;ConnectionPoolSize=1;DefinitionXml=<compressed xml>;FolderContainingDLLs[V]=CP_Posh_FolderContainingDLLs;Hostname[V]=CP_Posh_Hostname;Username[V]=CP_Posh_Username;Password[V]=CP_Posh_Password;Namespace=com.acme.myconnector;SystemId=MyConnector`

The `[V]` designation implies that the connection parameter will be defined as a system variable.

Compressing the connector definition XML is a three step process:

1. connector definition XML is Base64 encoded
2. Base64 encoded XML is compressed via the .Net _System.IO.Compression.DeflateStream_ class
3. Compressed bits are once again Base64 encoded

Sample code for encoding connector definition XML:

```Powershell
# load connector definition XML
$cxml_str = [System.IO.File]::ReadAllText("/home/mpierson/connector/connector-definition.xml");


# base 64 encode XML
$cxml_bytes = [System.Text.Encoding]::UTF8.GetBytes($cxml_str)
$cxml_b64 = [System.Convert]::ToBase64String($cxml_bytes)

# compress 
$ms_out = New-Object System.IO.MemoryStream
$zs = New-Object System.IO.Compression.DeflateStream($ms_out, [System.IO.Compression.CompressionMode]::Compress)
$ms_in = New-Object System.IO.StreamWriter($zs, [System.Text.Encoding]::UTF8)
$ms_in.Write($cxml_b64)
$ms_in.Flush()

# Base64 encode compressed data
$z_bytes = $ms_out.ToArray()
$z_b64 = [System.Convert]::ToBase64String($ms_out.ToArray())

$z_b64
```

