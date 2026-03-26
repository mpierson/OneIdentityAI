# Common Language Runtime (CLR) Type Identifiers

| Type                                                             | Exposed Interface                                             |
|:---------------------------------------------------|:-------------------------------------------------------------|
| VI.DB.Filters.FullTextFilterFormatter                                    | VI.DB.Filters.IFilterFormatter                               |
| VI.DB.Filters.SqlFilterFormatter                                         | VI.DB.Filters.IFilterFormatter                               |
| VI.DB.Filters.WildCardFilterFormatter                                    | VI.DB.Filters.IFilterFormatter                               |
| VI.DB.Filters.JobHDBFilterFormatter                                      | VI.DB.Filters.IFilterFormatter                               |
| VI.Projector.Connection.SystemConnection                                 | VI.Projector.Connection.ISystemConnection                    |
| VI.Projector.ADS.ADSConnectionParameterDescriptor                        | VI.Projector.Connection.ISystemConnectionParameterDescriptor |
| VI.Projector.Database.DatabaseConnectionParameterDescriptor              | VI.Projector.Connection.ISystemConnectionParameterDescriptor |
| VI.Projector.ADS.ProjectorADSConnectorDescriptor                         | VI.Projector.Connection.ISystemConnectorDescriptor           |
| VI.Projector.Database.DatabaseConnectorDescriptor                        | VI.Projector.Connection.ISystemConnectorDescriptor           |
| VI.Projector.Powershell.PoshConnectorDescriptor                          | VI.Projector.Connection.ISystemConnectorDescriptor           |
| VI.Projector.Connection.Scope.SystemScopeDefinition                      | VI.Projector.Connection.Scope.ISystemScopeDefinition         |
| VI.Projector.Filter.SystemObjectFilterJoin                               | VI.Projector.Filter.ISystemObjectFilterJoin                  |
| VI.Projector.ProjectorShell                                              | VI.Projector.IProjectorShell                                 |
| VI.Projector.Mapping.SystemMap                                           | VI.Projector.Mapping.ISystemMap                              |
| VI.Projector.Mapping.ConditionBasedSystemMappingCondition                | VI.Projector.Mapping.ISystemMappingCondition                 |
| VI.Projector.Mapping.Rules.SinglePropertyComparisonRule                  | VI.Projector.Mapping.ISystemMappingRule                      |
| VI.Projector.Mapping.Rules.MembersRule                                   | VI.Projector.Mapping.ISystemMappingRule                      |
| VI.Projector.Projection.ConditionBasedProjectionCondition                | VI.Projector.Projection.IProjectionCondition                 |
| VI.Projector.Projection.ProjectionConfiguration                          | VI.Projector.Projection.IProjectionConfiguration             |
| VI.Projector.Projection.ProjectionStartInfo                              | VI.Projector.Projection.IProjectionStartInfo                 |
| VI.Projector.Projection.ProjectionStep                                   | VI.Projector.Projection.IProjectionStep                      |
| VI.Projector.Projection.SchemaMethodAssignment                           | VI.Projector.Projection.ISchemaMethodAssignment              |
| VI.Projector.ADS.ProjectorADConnectorSchema                              | VI.Projector.Schema.ISchema                                  |
| VI.Projector.Powershell.PoshSchema                                       | VI.Projector.Schema.ISchema                                  |
| VI.Projector.Database.DatabaseSchema                                     | VI.Projector.Schema.ISchema                                  |
| VI.Projector.Schema.GenericSchemaClass                                   | VI.Projector.Schema.ISchemaClass                             |
| VI.Projector.Database.DatabaseSchemaClass                                | VI.Projector.Schema.ISchemaClass                             |
| VI.Projector.Schema.SchemaMethod                                         | VI.Projector.Schema.ISchemaMethod                            |
| VI.Projector.Powershell.Schema.PoshSchemaMethod                          | VI.Projector.Schema.ISchemaMethod                            |
| VI.Projector.Database.DatabaseSchemaMethod                               | VI.Projector.Schema.ISchemaMethod                            |
| VI.Projector.Schema.Properties.ReferenceResolutionSchemaProperty         | VI.Projector.Schema.ISchemaProperty                          |
| VI.Projector.Database.DatabaseSchemaCombinedPkProperty                   | VI.Projector.Schema.ISchemaProperty                          |
| VI.Projector.Schema.Properties.BitMaskSchemaProperty                     | VI.Projector.Schema.ISchemaProperty                          |
| VI.Projector.Database.DatabaseSchemaScopeReferenceProperty               | VI.Projector.Schema.ISchemaProperty                          |
| VI.Projector.Schema.Properties.MNSchemaTypeToMembersProperty             | VI.Projector.Schema.ISchemaProperty                          |
| VI.Projector.Schema.Properties.SchemaWalkerProperty                      | VI.Projector.Schema.ISchemaProperty                          |
| VI.Projector.Database.DatabaseSchemaProperty                             | VI.Projector.Schema.ISchemaProperty                          |
| VI.Projector.Powershell.PoshSchemaProperty                               | VI.Projector.Schema.ISchemaProperty                          |
| VI.Projector.Database.DatabaseSchemaRevisionProperty                     | VI.Projector.Schema.ISchemaProperty                          |
| VI.Projector.Schema.Properties.VirtualDependencyControlSchemaProperty    | VI.Projector.Schema.ISchemaProperty                          |
| VI.Projector.ADS.ProjectorADConnectorSchemaProperty                      | VI.Projector.Schema.ISchemaProperty                          |
| VI.Projector.Schema.Properties.SchemaPropertyJoin                        | VI.Projector.Schema.ISchemaProperty                          |
| VI.Projector.Schema.Properties.ConstantSchemaProperty                    | VI.Projector.Schema.ISchemaProperty                          |
| VI.Projector.Schema.Properties.ScriptedSchemaProperty                    | VI.Projector.Schema.ISchemaProperty                          |
| VI.Projector.Schema.Properties.DataValueConverterSchemaProperty          | VI.Projector.Schema.ISchemaProperty                          |
| VI.Projector.Schema.Properties.SchemaPropertyConverter                   | VI.Projector.Schema.ISchemaProperty                          |
| VI.Projector.Schema.Properties.MultiKeyResolutionSchemaProperty          | VI.Projector.Schema.ISchemaProperty                          |
| VI.Projector.Schema.References.SchemaPropertyReference                   | VI.Projector.Schema.ISchemaPropertyReference                 |
| VI.Projector.Schema.References.VirtualSchemaPropertyReference            | VI.Projector.Schema.ISchemaPropertyReference                 |
| VI.Projector.ADS.ProjectorADConnectorSchemaType                          | VI.Projector.Schema.ISchemaType                              |
| VI.Projector.Powershell.PoshSchemaType                                   | VI.Projector.Schema.ISchemaType                              |
| VI.Projector.Database.DatabaseSchemaType                                 | VI.Projector.Schema.ISchemaType                              |
| VI.Projector.Variables.SystemVariable                                    | VI.Projector.Variables.ISystemVariable                       |
| VI.Projector.Variables.SystemVariableSet                                 | VI.Projector.Variables.ISystemVariableSet                    |
| VI.Projector.Projection.ProjectionStepQuota                              |                                                              |
| VI.Projector.Projection.SystemObjectMatchingSets                         |                                                              |
| VI.Projector.ADS.ProjectorADConnectorReferenceTargetDetectorDN           |                                                              |
| VI.Projector.Schema.Converter.SchemaValueOrderConnectionstringConverter  |                                                              |
| VI.Projector.Connection.Scope.SystemScopeFilter                          |                                                              |
| VI.Projector.AdsiTools.Common.Converter.LargeIntToDateConverter          |                                                              |
| VI.Projector.AdsiTools.Common.Converter.IPv4AddressConverter             |                                                              |
| VI.Projector.Projection.SystemObjectMatchingSet                          |                                                              |
| VI.Projector.Database.DatabaseSchemaXObjectKeyReferenceTargetDetector    |                                                              |
| VI.Projector.AdsiTools.Common.Converter.LogonHoursConverter              |                                                              |
| VI.Projector.AdsiTools.Common.Converter.TimeSpanConverter                |                                                              |
