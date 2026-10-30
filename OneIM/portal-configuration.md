# Portal Configuration

Many aspects of the API server web applications can be configured without code changes, via JSON file stored in the database as a QBMFileRevision record with name _bin\CCC.CompositionApi.global.json_ and can be managed via the Software Loader tool. Use the deploy target _Server\Web\BusinessApiServer_ for this file.


The JSON includes an array of _Scopes_, with each entry corresponding to an application in the ApiServer.  Configuration applying globally to all applications uses the scope _imx_ and the end user portal scope has the name _portal_.  A _Data_ element contains name-value pairs for each configuration.

See also [Configuring API projects and web applications](https://docs.oneidentity.com/bundle/one-identity-manager_web-application-configuration_10.0/page/sources/webapp/wprconfigureallwebapplications.htm)

Common configurations:

| Scope         | Key                                                                       | Description                                                                 |
|---------------|---------------------------------------------------------------------------|-----------------------------------------------------------------------------|
| imx           | ServerLevelConfig/CompanyLogoUrl                                          | Uri of logo to be displayed in all applications                             |
| imx           | ServerLevelConfig/DefaultHtmlTheme                                          | Name of theme to applied to users who have not chosen a theme           |
| imx           | ServerLevelConfig/ProductIdentifier                                       | Name of application, e.g. "Identity Manager"                                |
| portal        | ServerConfig/PersonConfig/VI_Employee_MasterData_OrganizationalAttributes | List of editable identity attributes (database column names, Person table)) |
| passwordreset | ProjectLevelConfig/CaptchaCaseInsensitive                                 | Boolean value, indicating if captchas should be case insensitive            |



Sample JSON content:

```json
{
  "Scopes": [
    {
      "AppliesTo": "portal",
      "Ignore": false,
      "Data": {
        "ServerConfig/PersonConfig/VI_Employee_MasterData_OrganizationalAttributes": [
          "UID_Org",
          "UID_Department",
          "UID_ProfitCenter",
          "IsExternal",
          "IdentityType",
          "UID_PersonMasterIdentity",
          "EmployeeType",
          "EntryDate",
          "ExitDate",
          "UID_PersonHead",
          "IsInActive",
          "IsTemporaryDeactivated",
          "LeaveOfAbsenceReason",
          "DeactivationEnd"
        ]
      }
    },
    {
      "AppliesTo": "imx",
      "Ignore": false,
      "Data": {
        "ServerLevelConfig/CompanyLogoUrl": "https://www.alouette.com/data/63-aalv2/images/manual/logo_en.png",
        "ServerLevelConfig/DefaultHtmlTheme": "bootcamp-theme"
      }
    },
    {
      "AppliesTo": "passwordreset",
      "Ignore": false,
      "Data": {
        "ProjectLevelConfig/CaptchaCaseInsensitive": true
      }
    }
  ],
  "Ignore": false
}
```

## Custom Logo

Although the portal supports using a logo hosted elsewhere, e.g. on a customer's public facing web site, it is often more convenient to host the logo within the Identity Manager portal.  An easy way to do this is to package the logo image in a custom theme.  For example, the ACME logo would be packaged in the same zip file as acme-theme.css and served at html/acme-theme/acme.png.
