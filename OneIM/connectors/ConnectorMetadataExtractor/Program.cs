using System;
using System.Reflection;
using System.Collections.Generic;


if (args.Length < 3) {
    Console.WriteLine("Usage: cme <fully qualified class name> <absolute path to connector library DLL> <connector description>");
    return;
}


string className = args[0];
string dll_path = args[1];
string description = args[2];

// extract namespace and class name
string[] tokens = className.Split('.');
string name = tokens[tokens.Length - 1];
string ns = className.Substring(0, className.Length - name.Length - 1);



// See https://aka.ms/new-console-template for more information
Console.WriteLine($"<Connector className=\"{name}\" namespace=\"{ns}\">");
Console.WriteLine($" <Description>{description}</Description>");

Assembly connectorAssembly = Assembly.LoadFile(dll_path);

var conn = connectorAssembly.GetType($"{ns}.{name}");
if ( conn is not null ) 
{

    Console.WriteLine(" <Configuration>");
    var m = conn.GetMethod("ConfigParameters");
    if ( m is not null) {
        var result = m.Invoke(null, null);

        if (result is not null) {

            var e = (object[])result;
            foreach (var d in e) {
                Console.WriteLine($"  <Parameter name=\"{attr(d,"Name")}\" type=\"{attr(d, "Type")}\" required=\"{attr(d,"Required")}\">");
                if ( !string.IsNullOrEmpty(attr(d, "Description").ToString()) ){
                    Console.WriteLine($"    <Description>{attr(d, "Description")}</Description>");
                }
                Console.WriteLine("  </Parameter>");
            }
        }
    }
    Console.WriteLine(" </Configuration>");



    Console.WriteLine(" <Classes>");
    m = conn.GetMethod("DataClasses");
    if ( m is not null) {
        var result = m.Invoke(null, null);

        if (result is not null) {
            foreach (var d in (object[])result) {

                Console.WriteLine("  <Class ");
                Console.WriteLine($"    name = \"{attr(d, "ClassName")}\"");

                Console.WriteLine($"    createMethod = \"{attr(d, "CreateMethod")}\"");
                Console.WriteLine($"    listMethod   = \"{attr(d, "ListMethod")}\"");
                Console.WriteLine($"    getMethod    = \"{attr(d, "GetMethod")}\"");
                Console.WriteLine($"    updateMethod = \"{attr(d, "UpdateMethod")}\"");
                Console.WriteLine($"    deleteMethod = \"{attr(d, "DeleteMethod")}\"");

                Console.WriteLine("  >");

                if ( !string.IsNullOrEmpty(attr(d, "Description").ToString()) ){
                    Console.WriteLine($"    <Description>{attr(d, "Description")}</Description>");
                }

                // get array of attr def objects
                var attrDefs = attr(d, "Attributes");
                if (attrDefs is not null) {
                    foreach (var def in (object[])attrDefs) {

                        Console.WriteLine("    <Attribute ");
                        Console.WriteLine($"      name     = \"{attr(def, "Name")}\"");
                        Console.WriteLine($"      type     = \"{attr(def, "Type")}\"");
                        Console.WriteLine($"      required = \"{attr(def, "Required")}\"");
                        Console.WriteLine("    >");

                        if ( !string.IsNullOrEmpty(attr(def, "Description").ToString()) ){
                            Console.WriteLine($"     <Description>{attr(def, "Description")}</Description>");
                        }

                        Console.WriteLine("    </Attribute>");
                    }
                }


                var pm = conn.GetMethod("DataTypeAttributes");
                pm = null;
                if ( pm is not null) {
                    var attrs = pm.Invoke(null, new object[] {d});
                    if (attrs is not null) {
                        foreach (var a in (object[])attrs) {
                            Console.WriteLine($"   <Attribute name=\"{attr(a,"Name")}\" type=\"{attr(a, "Type")}\" required=\"{attr(a, "Required")}\">");
                            if ( !string.IsNullOrEmpty(attr(a, "Description").ToString()) ){
                                Console.WriteLine($"     <Description>{attr(a, "Description")}</Description>");
                            }
                            Console.WriteLine("   </Attribute>");
                        }
                    }
                }

                Console.WriteLine("  </Class>");
            }
        }
    }
    Console.WriteLine(" </Classes>");
}


 
Console.WriteLine("</Connector>");

object attr(object o, string name)
{
    var t = o.GetType();

    var a = t.GetField(name);
    if ( a is null ) {
        Console.WriteLine($"field {name} not found in {t}");
        return "";
    }

    var v = a.GetValue(o);
    if ( v is not null ) return v;

    return "";
}

