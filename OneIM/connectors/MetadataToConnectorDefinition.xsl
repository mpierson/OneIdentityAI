<?xml version='1.0' encoding="UTF-8"?>
<!--

  Transform One Identity Manager connector description to a PowerShell connector definition

  Author: M Pierson
  Date: Nov 2025
  Version: 0.1

  Use ConectorConfiguration.exe to generate source XML

 -->
<xsl:stylesheet version="3.0" xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
                              xmlns:xs="http://www.w3.org/2001/XMLSchema" 
                              xmlns:ois="http://www.oneidentity.com/IdentityManager/Connector/XSL" 
                              exclude-result-prefixes="ois xs">
  <xsl:output omit-xml-declaration="yes" indent="yes" cdata-section-elements="CustomCommand" />


 <!-- IdentityTransform -->
 <xsl:template match="/ | @* | node()">
   <xsl:copy> <xsl:apply-templates select="@* | node()" /> </xsl:copy>
 </xsl:template>


 <xsl:template match="Connector">
     <PowershellConnectorDefinition Version="1.0">
         <xsl:attribute name="Id" select="concat(@namespace, '.', @className)" />
         <xsl:attribute name="Description" select="Description" />

         <PluginAssemblies />
           <ConnectionParameters>
               <ConnectionParameter Name="FolderContainingDLLs" Description="Path to folder containing dlls" />
               <ConnectionParameter Name="CommaSeparatedDLLNames" Description="Names of connector dlls" />
               <ConnectionParameter Name="Namespace" Description="Namespace of the connector dll" />
               <ConnectionParameter Name="ClassName" Description="Main class name of the connector dll" />

               <xsl:apply-templates select="Configuration/Parameter" mode="ConnectionParameter" />
           </ConnectionParameters>

          <Initialization>
            <CustomCommands>
              <CustomCommand Name="Connect">
                <![CDATA[
                param( 
                   [parameter(Mandatory =$true, ValueFromPipelineByPropertyName =$true)] 
                   [ValidateNotNullOrEmpty()] 
                   [String]$FolderContainingDLLs, 

                   [parameter(Mandatory =$true, ValueFromPipelineByPropertyName =$true)]
                   [ValidateNotNullOrEmpty()]
                   [String]$CommaSeparatedDLLNames,

                   [parameter(Mandatory =$true, ValueFromPipelineByPropertyName =$true)]
                   [ValidateNotNullOrEmpty()]
                   [string]$Namespace,

                   [parameter(Mandatory =$true, ValueFromPipelineByPropertyName =$true)]
                   [ValidateNotNullOrEmpty()]
                   [string]$ClassName,

                  ]]>
                  <xsl:apply-templates select="Configuration/Parameter" mode="CommandParameter" />
                 <xsl:text><![CDATA[
                    )

                    $path = $FolderContainingDLLs
                    $arDLLs = $CommaSeparatedDLLNames.Split(",")
                    foreach ($dll in $arDLLs)
                    {
                        $fullName = Join-Path -Path $path -ChildPath $dll.Trim()
                        [Reflection.Assembly]::LoadFile($fullName)
                    }

                    $FQTN = $Namespace + "." + $ClassName
                    $global:connector = New-Object -TypeName $FQTN
                    $global:namespace = $Namespace

                 $params = @{
                 ]]></xsl:text>
                     <xsl:apply-templates select="Configuration/Parameter" mode="hash-setter" />
                 <xsl:text><![CDATA[
                     }

                     $global:connector.Configure($params)
                     $global:connector.Connect()
             ]]></xsl:text>

              </CustomCommand>
              <CustomCommand Name="Disconnect"><![CDATA[ $global:connector.Disconnect() ]]></CustomCommand> 

              <xsl:apply-templates select="Classes/Class" mode="CustomCommands" />

          </CustomCommands>

           <PredefinedCommands />

          <EnvironmentInitialization>
              <Connect>
                  <CommandSequence>
                      <Item Command="Connect" Order="0">
                          <SetParameter Param="FolderContainingDLLs" Source="ConnectionParameter" Value="FolderContainingDLLs" />
                          <SetParameter Param="CommaSeparatedDLLNames" Source="ConnectionParameter" Value="CommaSeparatedDLLNames" />
                          <SetParameter Param="Namespace" Source="ConnectionParameter" Value="Namespace" />
                          <SetParameter Param="ClassName" Source="ConnectionParameter" Value="ClassName" />
                          <xsl:apply-templates select="Configuration/Parameter" mode="ItemSetParameter" />
                      </Item>
                  </CommandSequence>
              </Connect>
              <Disconnect>
                  <CommandSequence> <Item Command="Disconnect" Order="0" /> </CommandSequence>
              </Disconnect>
          </EnvironmentInitialization>
      </Initialization>

      <Schema>
          <xsl:apply-templates select="Classes/Class" mode="Schema" />
      </Schema>

     </PowershellConnectorDefinition>
 </xsl:template>

 <xsl:template match="Parameter" mode="ConnectionParameter">
     <ConnectionParameter IsSensibleData="false">
         <xsl:attribute name="Name" select="@name" />
         <xsl:attribute name="Description" select="Description" />
     </ConnectionParameter>
 </xsl:template>
 <xsl:template match="Parameter" mode="ItemSetParameter">
     <SetParameter Source="ConnectionParameter">
         <xsl:attribute name="Param" select="@name" />
         <xsl:attribute name="Value" select="@name" />
     </SetParameter>
 </xsl:template>
 <xsl:template match="Parameter" mode="CommandParameter">
     <xsl:text><![CDATA[
                   [parameter(Mandatory =$true, ValueFromPipelineByPropertyName =$true)]
                   [ValidateNotNullOrEmpty()]
     ]]></xsl:text><xsl:value-of select="concat(
                    '[', ois:get-powershell-type(@type), ']$', @name,
                    if ( position() &lt; last() ) then ',' else ''
            )" />
 </xsl:template>
 <xsl:template match="Parameter" mode="configure-argument">
       <xsl:value-of select="concat( '$', @name, if (not(position() = last())) then ', ' else '')" />
 </xsl:template>
 <xsl:template match="Parameter" mode="hash-setter">
     <xsl:value-of select="concat( '&#xa;  &quot;', @name, '&quot; = $', @name )"/>
 </xsl:template>

 <xsl:template match="Class" mode="CustomCommands">
     <xsl:variable name="class-name" select="ois:get-class-name(@name)" />

      <!-- GET ALL -->
      <CustomCommand>
          <xsl:attribute name="Name" select="@listMethod" />
          <xsl:value-of select="concat('$global:connector.', @listMethod, '()&#xa;')" />
      </CustomCommand>              

      <!-- GET -->
      <CustomCommand>
          <xsl:attribute name="Name" select="@getMethod" />
          <![CDATA[
          param(
              [parameter(Mandatory =$true, ValueFromPipelineByPropertyName =$true)]
              [ValidateNotNullOrEmpty()]
              [String]$Id
          )
          ]]> 
          <xsl:value-of select="concat('$global:connector.', @getMethod, '($Id)&#xa;')" />
      </CustomCommand>              

      <!-- CREATE -->
      <CustomCommand>
          <xsl:attribute name="Name" select="@createMethod" />
          <![CDATA[
          param(
          ]]> 
          <xsl:apply-templates select="Attribute" mode="command-argument" />
         <xsl:text><![CDATA[
          )  

             $attrs = @{
         ]]></xsl:text>
         <xsl:apply-templates select="Attribute" mode="hash-setter" />
         <xsl:text><![CDATA[
             }

         ]]></xsl:text>
         <xsl:value-of select="concat('$global:connector.', @createMethod, '($Id, $attrs)&#xa;')" />
      </CustomCommand>              

      <!-- UPDATE -->
      <CustomCommand>
          <xsl:attribute name="Name" select="@updateMethod" />
          <![CDATA[
          param(
          ]]> 
          <xsl:apply-templates select="Attribute" mode="command-argument" />
         <xsl:text><![CDATA[
             )  

             $attrs = @{
         ]]></xsl:text>
         <xsl:apply-templates select="Attribute" mode="hash-setter" />
         <xsl:text><![CDATA[
             }

         ]]></xsl:text>
         <xsl:value-of select="concat('$global:connector.', @updateMethod, '($Id, $attrs)&#xa;')" />
      </CustomCommand>              

      <!-- DELETE -->
      <CustomCommand>
          <xsl:attribute name="Name" select="@deleteMethod" />
         <xsl:text><![CDATA[

          param(
              [parameter(Mandatory =$true, ValueFromPipelineByPropertyName =$true)]
              [ValidateNotNullOrEmpty()]
              [String]$Id
          )

         ]]></xsl:text>
         <xsl:value-of select="concat('$global:connector.Delete', $class-name, '($Id)&#xa;')" />
      </CustomCommand>              
 </xsl:template>

 <xsl:template match="Attribute" mode="command-argument">
     <xsl:value-of select="concat(
                            '&#xa;[parameter(Mandatory =',
                            if ( @required='True' ) then '$true' else '$false',
                            ', ValueFromPipelineByPropertyName =$true)] &#xa;'
        )"/>
     <xsl:text><![CDATA[
                   [ValidateNotNullOrEmpty()]
       ]]></xsl:text><xsl:value-of select="concat(
                    '[', ois:get-powershell-type(@type), ']$', @name,
                    if ( position() &lt; last() ) then ',' else ''
            )" />
 </xsl:template>
 <xsl:template match="Attribute" mode="hash-setter">
     <!-- double/float values are not supported by PoSh connector, so needs conversion -->
     <xsl:choose>
         <xsl:when test="@type='double' or @type='float' or @type='decimal'">
             <xsl:value-of select="concat( '&#xa;  &quot;', @name, '&quot; = Format-Number -Number $', @name )"/>
         </xsl:when>
         <xsl:otherwise>
             <xsl:value-of select="concat( '&#xa;  &quot;', @name, '&quot; = $', @name )"/>
         </xsl:otherwise>
     </xsl:choose>
 </xsl:template>


 <xsl:template match="Class" mode="Schema">
     <xsl:variable name="class-name" select="ois:get-class-name(@name)" />

     <Class>
         <xsl:attribute name="Name" select="$class-name" />
         <Properties>
             <xsl:apply-templates select="Attribute" mode="schema-property" />
        </Properties>
        <ReadConfiguration>
            <ListingCommand>
                <xsl:attribute name="Command" select="@listMethod" />
            </ListingCommand>
            <CommandSequence>
                <Item Order="1">
                    <xsl:attribute name="Command" select="@getMethod" />
                </Item>
            </CommandSequence>
        </ReadConfiguration>

      <MethodConfiguration>
        <Method Name="Insert">
          <CommandSequence>
              <Item Order="1">
                  <xsl:attribute name="Command" select="@createMethod" />
              </Item>
          </CommandSequence>
        </Method>
        <Method Name="Update">
          <CommandSequence>
              <Item Order="1">
                  <xsl:attribute name="Command" select="@updateMethod" />
              </Item>
          </CommandSequence>
        </Method>
        <Method Name="Delete">
          <CommandSequence>
              <Item Order="1">
                  <xsl:attribute name="Command" select="@deleteMethod" />
              </Item>
          </CommandSequence>
        </Method>
      </MethodConfiguration>
     </Class>

 </xsl:template>
 <xsl:template match="Attribute" mode="schema-property">
     <xsl:variable name="class-name" select="ois:get-class-name(../@name)" />
     <Property>
         <xsl:attribute name="Name" select="@name" />
         <xsl:attribute name="DataType" select="ois:get-powershell-type(@type)" />
         <xsl:attribute name="IsUniqueKey" select="if ( @name = 'Id' ) then 'true' else 'false'" />
         <xsl:attribute name="IsMandatory" select="ois:get-boolean(@required)" />
         <CommandMappings>
             <xsl:if test="@name = 'Id'">
                 <Map>
                     <xsl:attribute name="Parameter" select="@name" />
                     <xsl:attribute name="ToCommand" select="../@getMethod" />
                 </Map>
                 <Map>
                     <xsl:attribute name="Parameter" select="@name" />
                     <xsl:attribute name="ToCommand" select="../@updateMethod" />
                 </Map>
                 <Map>
                     <xsl:attribute name="Parameter" select="@name" />
                     <xsl:attribute name="ToCommand" select="../@deleteMethod" />
                 </Map>
             </xsl:if>
             <Map>
                 <xsl:attribute name="Parameter" select="@name" />
                 <xsl:attribute name="ToCommand" select="../@updateMethod" />
             </Map>
         </CommandMappings>
         <ReturnBindings>
             <xsl:if test="@name = 'Id'">
                 <Bind>
                     <xsl:attribute name="Path" select="@name" />
                     <xsl:attribute name="CommandResultOf" select="../@listMethod" />
                 </Bind>
             </xsl:if>
             <Bind>
                 <xsl:attribute name="Path" select="@name" />
                 <xsl:attribute name="CommandResultOf" select="../@getMethod" />
             </Bind>
             <Bind>
                 <xsl:attribute name="Path" select="@name" />
                 <xsl:attribute name="CommandResultOf" select="../@createMethod" />
             </Bind>
             <Bind>
                 <xsl:attribute name="Path" select="@name" />
                 <xsl:attribute name="CommandResultOf" select="../@updateMethod" />
             </Bind>
         </ReturnBindings>
     </Property>
 </xsl:template>



 <xsl:function name="ois:get-class-name" as="xs:string">
    <xsl:param name="fqcn" as="xs:string"/>
    <xsl:value-of select="tokenize($fqcn,'\.')[last()]" />
 </xsl:function>

 <xsl:function name="ois:get-powershell-type" as="xs:string">
    <xsl:param name="csharp-type" as="xs:string"/>
    <xsl:value-of select="if ( $csharp-type = 'string' ) 
                          then 'String'
                          else if ($csharp-type = 'int' )
                          then 'Int'
                          else if ($csharp-type = 'double' )
                          then 'String'
                          else $csharp-type 
        " />
 </xsl:function>

 <xsl:function name="ois:get-boolean" as="xs:string">
    <xsl:param name="v" as="xs:string"/>
    <xsl:value-of select="lower-case($v)" />
 </xsl:function>

</xsl:stylesheet>


