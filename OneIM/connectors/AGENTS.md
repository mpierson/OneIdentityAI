# Instructions For Custom Identity Manager Connectors

Custom Identity Manager connectors are used to integrate Identity Manager with identity sources (e.g. HR system) or systems that require provisioning of identity information (target systems).

There are typically three components of a custom connector:

 - a **.Net shared library** that is capable of performing CRUD operations on the connected system
 - an **XML connector definition** that contains a recipe for using the shared library
 - an **Identity Manager synchronization project** that defines the mapping and data flow to and from the connected system

Additional details of Identity Manager custom connectors [here](https://github.com/OneIdentity/IdentityManager.PoSh-Connector-Guide/blob/main/README.md).


## Connector Shared Libary

The custom connector library should include a public class that implements _ois.oneim.ConnectorBase.ConnectorBase_, as defined in [ConnectorBase.cs](ConnectorBase.cs).

### Connector Configuration

The _Configure_ method is used to setup target system connection details.  Configuration parameters are passed in key/value pairs.   The connector's static _ConfigParameters()_ should return metadata for the parameters to be passed to the _Configure()_ method.

For example, a connector that requires host, port, and service account credentials for the 'ACME' target system would implement the _ConfigParameters()_ method as follows:

``` csharp

public static ParameterDef[] ConfigParameters()
{
    return new ParameterDef[]
    {
        new ParameterDef { Name = "acme_host", Type = "string", Description = "hostname or IP address of ACME server", Required = true },
        new ParameterDef { Name = "acme_port", Type = "int", Description = "server port, default is 443" },
        new ParameterDef { Name = "acme_username", Type = "string", Description = "username for ACME authentication", Required = true },
        new ParameterDef { Name = "acme_password", Type = "string", Description = "password for ACME authentication", Required = true }
    };
}

```

### Data Types

For each data type in the connected system to be exposed to Identity Manager, the connector class should implement methods to support CRUD operations in the target system:

 - create a single object, given a _System.Collections.Hashtable_ of name/value pairs
 - return list of all objects, returning a _System.Collections.Generic.List_
 - return a single object
 - update a single object, given a _System.Collections.Hashtable_ of name/value pairs
 - delete a single object

Each data type should be represented by a class, with a public field for each supported attribute in the connected system.

**Note**: Identity Manager supports the following attribute types:

- String, 
- Bool, 
- Int, 
- DateTime

Other data types exposed by the connected system should be converted to String within the connector.

Each data type class should implement the _ois.oneim.ConnectorBase.ConnectorData_ interface.  The connector's static method  _DataClasses()_ returns metadata for each of the supported data types.


For example, for an _Employee_ data type, the connector should implement methods like:

``` csharp

   // fetch all employee records from target system
    List<Employee> GetAllEmployees(){}

    // fetch one employee record
    Employee GetEmployee(string Id){} 

    // create one employee record
    void CreateEmployee(string Id, Hashtable attributes){}

    // update one employee record
    void UpdateEmployee(string Id, Hashtable attributes){}

    // delete one employee record
    void DeleteEmployee(string Id);

```

The metadata for this example _Employee_ class would be:

``` csharp

       var emplMetaData = new DataDef {
            ClassName       = "ois.oneim.sample.Employee",
            CreateMethod    = "CreateEmployee",
            ListMethod      = "GetAllEmployees",
            GetMethod       = "GetEmployee",
            UpdateMethod    = "UpdateEmployee",
            DeleteMethod    = "DeleteEmployee",
            Attributes = new AttributeDef[] {
                new AttributeDef { Name="Id", Type="string", Description="Unique identifier" },
                new AttributeDef { Name="FirstName", Type="string", Required=true },
                new AttributeDef { Name="LastName", Type="string", Required=true },
                new AttributeDef { Name="MiddleName", Type="string" },
                new AttributeDef { Name="JobTitle", Type="string", Description="Job title or position" },
                new AttributeDef { Name="Manager", Type="string", Description="Id of direct manager" },
                new AttributeDef { Name="Status", Type="string", Description="Status, A = active, I = inactive" },
            }
       };
```

## Connector Definition File

An XML connector definition file is used by Identity Manager to determine how to utilize functionality in the connector library.

The [Connector Guide](https://github.com/OneIdentity/IdentityManager.PoSh-Connector-Guide/blob/main/README.md#the-powershell-connector-xml-definition-format) describes the format of the connector definition file in detail.

A two-step process is available to automatically generate the connector definition file:

1. use the _cme_ tool (Connector Metadata Extractor) to extract metadata from the connector library, in XML format
2. apply the _MetadataToConnectorDefinition.xsl_ stylesheet to convert the connector metadata to a connector definition file using a XSLT v3 processor

Usage of the _cme_ tool:

``` bash
> dotnet run cme.dll
Usage: cme <fully qualified class name> <absolute path to connector library DLL> <connector description>
```


Use the following command to generate the definition file in a Linux shell:

``` bash
dotnet run cme.dll ois.oneim.sample.SampleConnector \
                   /home/one_identity/sample_connector/SampleConnector.dll \
                   "Sample connector for demonstration only" |  \
    saxonhe-xslt -xsl:MetadataToConnectorDefinition.xsl -s:- > connector-definition.xml
```
