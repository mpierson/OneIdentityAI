using System.Collections;
using System.Text.Json.Serialization;

namespace ois.oneim.ConnectorBase {

public interface ConnectorData {
    // primary unique identifier
    string Id { get; set; }

    // map of all attribute name/value pairs, including Id
    Hashtable GetAttributes();
}


public class ParameterDef {
    public required string Name { get; set; }
    public required string Type { get; set; }
    public string? Description { get; set; }
    public bool Required {get; set;}
}

public class AttributeDef {
    public required string Name { get; set; }
    public required string Type { get; set; }
    public string? Description { get; set; }
    public bool Required {get; set;}
}

// Metadata for a data type supported by connector
// e.g.
/*
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
*/
public class DataDef {
    // fully qualified name of Class, e.g. ois.oneim.ConnectorBase.Employee
    public required string ClassName { get; set; }

    public required string CreateMethod { get; set; }
    public required string ListMethod { get; set; }
    public required string GetMethod { get; set; }
    public required string UpdateMethod { get; set; }
    public required string DeleteMethod { get; set; }

    // returns list of attributes available in the data type
    public required AttributeDef[] Attributes {get; set;}
}


public interface ConnectorInterface
{
    // text describing the connector and indicating the target system
    static abstract string Description();

    // returns list of parameters accepted by the Configure method
    static abstract ParameterDef[] ConfigParameters();

    // returns list of data classes supported by connector
    static abstract DataDef[] DataClasses();


    // configure connector with parameters required to connect to target
    void Configure(Hashtable parameters);

    // perform any initialization or authentication with target system
    // <returns>true if successful</returns>
    bool Connect();

    // perform cleanup 
    void Disconnect();


}

} // namespace
