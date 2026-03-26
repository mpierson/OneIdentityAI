# Connector Definition File

An XML connector definition file is used by Identity Manager to access methods in the connector library.

The [Connector Guide](https://github.com/OneIdentity/IdentityManager.PoSh-Connector-Guide/blob/main/README.md#the-powershell-connector-xml-definition-format) describes the format of the connector definition file in detail.

## Connector metadata 

Metadata from a connector's shared library can be used to generate an Identity Manager connector definition.

Sample metadata, in XML format:

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

An XSL template is available that will convert library metadata to a connector definition file: [MetadataToConnectorDefinition.xsl](MetadataToConnectorDefinition.xsl) .

Use an XSLT v3 processor to convert the connector's XML metadata to a connector definition file:

``` bash
    cat connector-metadata.xml | saxonhe-xslt -xsl:MetadataToConnectorDefinition.xsl -s:- > connector-definition.xml
```


