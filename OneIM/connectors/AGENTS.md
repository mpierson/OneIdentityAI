# Creating Custom Identity Manager Connectors

Custom Identity Manager connectors are used to integrate Identity Manager with identity sources (e.g. HR system) or systems that require provisioning of identity information (target systems).

There are typically three components of a custom connector:

 - a **.Net shared library** that is capable of performing CRUD operations on the connected system
 - an **XML connector definition** that contains a recipe for using the shared library
 - an **Identity Manager synchronization project** that defines the mapping and data flow to and from the connected system

Additional details of Identity Manager custom connectors [here](https://github.com/OneIdentity/IdentityManager.PoSh-Connector-Guide/blob/main/README.md).

Once the shared library, connector definition, and synchronization project are complete, the following manual steps will be required:

- copy the shared library, connector definition, and any dependencies to the target Identity Manager job server
- update the synchronization project variables to reflect the location of the shared library and connector definition files


# Connector Shared Libary

The custom connector library should include a public class that implements _ois.oneim.ConnectorBase.ConnectorBase.ConnectorInterface_, as defined in [ConnectorBase.cs](ConnectorBase.cs).

For best results, a single, .Net 8 compliant DLL should be created.

## Connector Configuration

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

**Note**: individual parameters are preferred over multi-value combined parameters, e.g. use host, port, database name, instead of a single connection string.



## Data Types

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



# Connector Definition File

An XML connector definition file is used by Identity Manager to determine how to utilize functionality in the connector library.

The [Connector Guide](https://github.com/OneIdentity/IdentityManager.PoSh-Connector-Guide/blob/main/README.md#the-powershell-connector-xml-definition-format) describes the format of the connector definition file in detail.

An XSL template is available that will convert library metadata to a connector definition file.

Sample library metadata, in XML format:

``` xml

<?xml version="1.0" encoding="utf-8"?>
<Connector className="MyConnector" namespace="ois.oneim">
  <Description>Sample connector library for xyz target system</Description>

  <!-- Configuration section includes one Parameter element for each parameter to be passed to the Configure method of the connector.  e.g. host name, service account credentials.
  <Configuration>
    <Parameter name="connected_host" type="string">
      <Description>Hostname or IP address of connected system</Description>
    </Parameter>
    ...
  </Configuration>

  <!-- 
     A Class element is required for each data type in the connected system that is supported by the connector. The _listMethod_, _createMethod_, _getMethod_, _updateMethod_, and _deleteMethod_ attributes indicate the name of the corresponding method in the connector. 
     Each attribute of the data type is represented by an Attribute element, including XML attributes for the name, type, and if a value is required by the connected system.  If the attribute must not be empty when creating or updating the connected system, the XML attribute _required_ should be _True_.
  -->
  
  <Classes>
    <Class name="ois.oneim.frappe.Employee"
           createMethod="CreateEmployee"
           listMethod="GetAllEmployees"
           getMethod="GetEmployee"
           updateMethod="UpdateEmployee"
           deleteMethod="DeleteEmployee">
      <Attribute name="Id" type="string" required="False" />
      ...
    </Class>
  </Classes>
</Connector>
```

Use an XSLT v3 processor to convert the connector's XML metadata to a connector definition file, using the provided _MetadataToConnectorDefinition.xsl_ template.

``` bash
    cat connector-metadata | saxonhe-xslt -xsl:MetadataToConnectorDefinition.xsl -s:- > connector-definition.xml
```

# Synchronization Project

The SPEd utility is used to create and manage synchronization projects.


