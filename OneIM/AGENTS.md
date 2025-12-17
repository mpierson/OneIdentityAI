# AI Agent instructions for Identity Manager

## Environment and Dependencies

Requires .Net 8 SDK

Microsoft SQL Server driver, e.g. Microsoft.Data.SqlClient.dll

The following Identity Manager DLLs are typically required:

- VI.Base.dll
- VI.DB.dll
- \[\*\].Customizer.dll (QER.Customizer.dll, CPL.Customizer.dll, etc.)


## Tips for interacting with an existing instance

Establish a connection to OneIM database:

``` csharp
using VI.DB;
using VI.DB.Entities;

// connection string is MS SQL standard
ConnectData connectData = DbApp.Instance.Connect(
    new ViSqlFactory(), 
    "Data Source=host_and_port;Initial Catalog=db_name;User ID=svc_account;Password=password");

// connect using system account, e.g. viAdmin
connectData.Connection.Authenticate("Module=DialogUser;User=user_name;Password=password");

ISession session = connectData.Connection.Session;
```

## Managing Identities

Identities are stored in the Person table.  See [Person](Person.md) for details.

Important columns in the Person table:

| Column  | Description | Notes   | Type       |
|-----------------------------|------------------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-----------|
| CentralAccount              | Central user account               | Identity's unique central user account that can be mapped to any target system.                                                                                                                                 | nvarchar  |
| DefaultEmailAddress         | Default email address              | Default email address, that can possibly be mapped to the actual originating mailboxes using templates.                                                                                                         | nvarchar  |
| Description                 | Description                        | Description.                                                                                                                                                                                                    | nvarchar  |
| EmployeeType                | Employee type                      | More precise classification of the identity in terms of their contractual relationship to the company.                                                                                                          | varchar   |
| EntryDate                   | Entry date                         | Identity's date of joining the company.                                                                                                                                                                         | datetime  |
| ExitDate                    | Leaving date                       | Identity's company leaving date.                                                                                                                                                                                | datetime  |
| FirstName                   | First name                         | Identity's first name.                                                                                                                                                                                          | nvarchar  |
| IsExternal                  | External                           | Defines whether the identity is external.                                                                                                                                                                       | bit       |
| IsInActive                  | Disabled permanently               | Specifies whether the identity is actively used. If an identity is permanently deactivated, all its permissions as a %Globals.QIM_ProductNameShort% user are revoked.                                           | bit       |
| IsVIP                       | VIP                                | Defines whether the identity is a VIP.                                                                                                                                                                          | bit       |
| LastName                    | Last name                          | Last name of identity.                                                                                                                                                                                          | nvarchar  |
| MiddleName                  | Middle name                        | Middle name of identity.                                                                                                                                                                                        | nvarchar  |
| PersonalTitle               | Job description                    | Identity's service name.                                                                                                                                                                                        | nvarchar  |
| PersonnelNumber             | Personnel number                   | Identity's personnel number                                                                                                                                                                                     | nvarchar  |
| UID_Person                  | Identity                           | Unique identity identifier.                                                                                                                                                                                     | varchar   |
| UID_PersonHead              | Manager                            | Manager identifier.                                                                                                                                                                                             | varchar   |


### Create a new identity

``` csharp
IEntity person = await s.Source().CreateNewAsync("Person");
await Task.WhenAll(
    person.PutValueAsync("FirstName", "Claude"),
    person.PutValueAsync("LastName", "Van Damme"),
    person.PutValueAsync("Description", "identity created via object model")
);

using (IUnitOfWork uow = s.StartUnitOfWork()) {
    await uow.PutAsync(person);
    await uow.CommitAsync();
}
```

### Fetch identities

``` csharp
// fetch using a UID
string UID_Person = "e78a48aa-e101-4ec6-b9c9-a5d8f221df71";
IEntity person2 = await session.Source().GetAsync("Person", UID_Person);
Console.WriteLine(person2.Display);


// fetch using criteria

var q = Query.From("Person").Where(p => p.Column("LastName") == "Smith");

int count = await session.Source().GetCountAsync(q.SelectCount());
Console.WriteLine(count);

IEntityCollection persons = await session.Source().GetCollectionAsync(q.SelectAll());
foreach (IEntity p in persons) {
    Console.WriteLine(p.Display);
}

// fetch all
IEntityCollection allPersons = await session.Source().GetCollectionAsync(Query.From("Person").SelectAll());
foreach (IEntity p in allPersons) {
    Console.WriteLine(p.Display);
}
```

### Update an identity

``` csharp

await person2.PutValueAsync("MiddleName", "Jane");
using (IUnitOfWork uow = session.StartUnitOfWork()) {
    await uow.PutAsync(person2);
    await uow.CommitAsync();
}

```

### Delete identities

``` csharp
// mark users for delete, triggering deferred delete after 30 days
q = Query.From("Person").Where(p => p.Column("LastName") == "Van Damme");
IEntityCollection personsToDelete_safe = await session.Source().GetCollectionAsync(
                                    q.SelectAll(), EntityCollectionLoadType.Bulk);
using (IUnitOfWork uow = session.StartUnitOfWork()) {
    foreach (IEntity p in personsToDelete_safe) {
        Console.WriteLine(p.Display);
        p.MarkForDeletion();
        await uow.PutAsync(p);
    }
    await uow.CommitAsync();
}

// delete
q = Query.From("Person").Where(p => p.Column("LastName") == "Van Damme");
IEntityCollection personsToDelete = await session.Source().GetCollectionAsync(
                                    q.SelectAll(), EntityCollectionLoadType.Bulk);
using (IUnitOfWork uow = session.StartUnitOfWork()) {
    foreach (IEntity p in personsToDelete) {
        Console.WriteLine(p.Display);
        p.MarkForDeletionWithoutDelay();
        await uow.PutAsync(p);
    }
    await uow.CommitAsync();
}
```

## Password Policies

Password policies are stored in the QBMPwdPolicy table.  See [QBMPwdPolicy](QBPwdPolicy.md) for details.

Important columns in the QBMPwdPolicy table:

| Column  | Description | Notes   | Type       |
|-----------------------------|------------------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-----------|
| AdditionalPwdRequirements                                                                                                                                    | Additional requirements                 | Verbal description of the additional password requirements that will be tested in the test script.                                                  | nvarchar |
| DefaultInitialPassword                                                                                                                                       | Initial password                        | The initial password for newly created user accounts. If no password is given when the account is first created then the default password is taken. | varchar  |
| Description                                                                                                                                                  | Description                             | Description.                                                                                                                                        | nvarchar |
| DisplayName                                                                                                                                                  | Display name                            | Display name.                                                                                                                                       | nvarchar |
| HistoryLen                                                                                                                                                   | Password history                        | Number of passwords to save. If the value 10 is given, the last 10 user passwords are saved.                                                        | int      |
| IsLowerLetterNotAllowed                                                                                                                                      | Do not generate lowercase letters       | The password must not include lowercase alphabetic characters. The generation of a random password takes place without this.                        | bit      |
| IsNumberNotAllowed                                                                                                                                           | Do not generate digits                  | The password must not include digits. The generation of a random password takes place without this.                                                 | bit      |
| IsSpecialNotAllowed                                                                                                                                          | Do not generate special characters      | The password must not include special characters. The generation of a random password takes place without this.                                     | bit      |
| IsUpperLetterNotAllowed                                                                                                                                      | Do not generate uppercase letters       | The password must not include uppercase alphabetic characters. The generation of a random password takes place without this.                        | bit      |
| MandatoryCharacterClasses                                                                                                                                    | Required number of character classes    | Number of character classes that must be used to satisfy the password policy. If the value is 0, all configured character classes must be used.     | int      |
| MaxAge                                                                                                                                                       | Validity period (max. # days)           | After the time period has run out the user has to change the password.  The value is only taken into account when connecting to %Globals.QIM_ProductNameShort% frontends.                                                            | int                                     |
| MaxBadAttempts                                                                                                                                               | Max. failed logins                      | If the user reaches the number the user account is locked.  The value is only taken into account when connecting to %Globals.QIM_ProductNameShort% frontends with a system user or identity-based authentication module. | int                                     |
| MaxLen                                                                                                                                                       | Max. length                             | Maximum length of the password.                                                                                                                     | int      |
| MaxRepeatCount                                                                                                                                               | Max. identical characters at all        | Defines the maximum number of identical characters which may occur in the password.                                                                 | int      |
| MaxRepeatLen                                                                                                                                                 | Max. identical characters in succession | Defines the maximum number of identical characters which can be repeated in succession.                                                             | int      |
| MinLen                                                                                                                                                       | Min. length                             | Minimum length of the password.                                                                                                                     | int      |
| MinLetters                                                                                                                                                   | Min. number letters                     | Defines the minimum number of alphabetic characters that must be in a password.                                                                     | int      |
| MinLettersLowerCase                                                                                                                                          | Min. number lowercase                   | Defines the minimum number of lowercase alphabetic characters that must be in a password.                                                           | int      |
| MinLettersUpperCase                                                                                                                                          | Min. number uppercase                   | Defines the minimum number of uppercase alphabetic characters that must be in a password.                                                           | int      |
| MinNumbers                                                                                                                                                   | Min. number digits                      | Defines the minimum number of digits that must be in a password.                                                                                    | int      |
| MinPasswordQuality                                                                                                                                           | Min. password strength                  | Defines the minimal strength of a password.                                                                                                         | int      |
| MinSpecialChar                                                                                                                                               | Min. number special characters          | Defines the minimum number of special characters that must be in a password. Special characters are: `$!"#%&()*+,-./:;<=>?@\_{}~`                     | int      |
| SpecialCharsAllowed                                                                                                                                          | Allowed special characters              | List of allowed special characters.                                                                                                                 | nvarchar |
| SpecialCharsDenied                                                                                                                                           | Prohibited special characters           | List of denied special characters.                                                                                                                  | nvarchar |
| UID_QBMPwdPolicy                                                                                                                                             | Password policy                         | Password policy identifier.                                                                                                                         | varchar  |


### Create a password policy

``` csharp

IEntity policy = await session.Source().CreateNewAsync("QBMPwdPolicy");
await Task.WhenAll(
    policy.PutValueAsync("DisplayName", "Sample Password Policy"),
    policy.PutValueAsync("Description", "Policy with sample complexity requirements"),
    policy.PutValueAsync("MinLen", 8),
    policy.PutValueAsync("MinLetters", 1),
    policy.PutValueAsync("MinLettersLowerCase", 1),
    policy.PutValueAsync("MinLettersUpperCase", 1),
    policy.PutValueAsync("MinNumbers", 1),
    policy.PutValueAsync("MinSpecialChar", 1),
    policy.PutValueAsync("MandatoryCharacterClasses", 2)
);

using (IUnitOfWork uow = session.StartUnitOfWork()) {
    await uow.PutAsync(policy);
    await uow.CommitAsync();
}

```

### Generate a password 

To generate a password that conforms to a policy:

``` csharp
using VI.DB.Passwords;

int requiredLength = 16;
string UID_QBMPwdPolicy = policy.GetValue("UID_QBMPwdPolicy");

IPasswordManager passwordManagerImpl = session.Resolve<IPasswordManager>();
var pwdGenerator = await passwordManagerImpl.GetPolicyAsync(UID_QBMPwdPolicy);
var securePwd = pwdGenerator.CreatePassword(requiredLength);
string plaintextPwd = new System.Net.NetworkCredential(string.Empty, securePwd).Password;

Console.WriteLine( plaintextPwd );

```

### Reset passwords

``` csharp

// load identities from database
var q = Query.From("Person").Where(p => p.Column("PersonalTitle") == "Relief Pitcher");
int count = await session.Source().GetCountAsync(q.SelectCount());
Console.WriteLine(count);

// reset each user's password to random value, according to policy
IEntityCollection persons = await session.Source().GetCollectionAsync(
                                q.SelectAll(), EntityCollectionLoadType.Bulk);
using (IUnitOfWork uow = session.StartUnitOfWork()) {
    foreach (IEntity p in persons) {
        Console.WriteLine(p.Display);

        var pwd = pwdGenerator.CreatePassword(requiredLength);
        await p.PutValueAsync("CentralPassword",
                new System.Net.NetworkCredential(string.Empty, pwd).Password);
        await uow.PutAsync(p);
    }
    await uow.CommitAsync();
}
```

## Managing Organization Structures

Organization structures, including departments, locations, and cost centers are managed similarily in Identity Manager.

Organization structure tables:

- _Department_: organization's departments
- _Locality_: office locations or region
- _ProfitCenter_: cost centers or accounting units

All three organization types support the following:

- hierarchichal inheritance, via foreign key to parent, e.g. _Department.UID\_ParentDepartment_
- assignment of structure to employee via foreign key column, e.g. _Person.UID\_Department_
- assignment of employees via membership table, e.g. _PersonInDepartment_
- assignment of entitlements via assignment table, e.g. _DepartmentHasADSGroup_
- assignment of manager via foreign key column, e.g. _Department.UID\_PersonHead_


### Create a new organization object

``` csharp

// create Overdue Accounts department within the Accounting parent
string UID_ParentDepartment = "xxx"; // ID of Accounting department
IEntity department = await session.Source().CreateNewAsync("Department");
await Task.WhenAll(
    department.PutValueAsync("DepartmentName", "Overdue Accounts"),
    department.PutValueAsync("Description", "Manage delinquent accounts"),
    department.PutValueAsync("UID_ParentDepartment", UID_ParentDepartment)  // parent is Accounting
);

using (IUnitOfWork uow = session.StartUnitOfWork()) {
    await uow.PutAsync(department);
    await uow.CommitAsync();
}

```

### Fetch organization object

``` csharp

// fetch all departments
IEntityCollection allDepartments = await session.Source().GetCollectionAsync(Query.From("Department").SelectAll()
);
foreach (IEntity d in allDepartments) {
    Console.WriteLine(d.Display);
}

// fetch using a UID
string UID_Department = "a60e34fe-e415-43f5-b169-e07ee18340ea";
IEntity department = await session.Source().GetAsync("Department", UID_Department);
Console.WriteLine(department.Display);


// fetch using SQL criteria
var q = Query.From("Department").Where(p => p.Column("DepartmentName") == "Accounting");
int count = await session.Source().GetCountAsync(q.SelectCount());
Console.WriteLine(count);

IEntityCollection departments = await session.Source().GetCollectionAsync(q.SelectAll());
foreach (IEntity d in departments) {
    Console.WriteLine(d.Display);
}

```

### Manage employee assignment

Employees can be assigned to a primary structure, e.g. a primary location where the employee typically works, plus one or more secondary assignments, e.g. additional locations where the employee has privileges.

``` csharp

// update primary department of employees
string UID_Department_Accounting = "xxx";
var q = Query.From("Person").Where(p => p.Column("PersonalTitle") == "Accountant");
IEntityCollection accountants = await session.Source().GetCollectionAsync(
        q.SelectAll(), EntityCollectionLoadType.Bulk);
using (IUnitOfWork uow = session.StartUnitOfWork()) {
    foreach (IEntity a in accountants) {
        Console.WriteLine(a.Display);
        await p.PutValueAsync("UID_Department", UID_Department_Accounting);
        await uow.PutAsync(p);
    }
    await uow.CommitAsync();
}


// add employees to a secondary department
string UID_Department_Compliance = "xxx";
q = Query.From("Person").Where(p => p.Column("LastName") == "Van Damme");
IEntityCollection complianceOfficer = await session.Source().GetCollectionAsync(
        q.SelectAll(), EntityCollectionLoadType.Bulk);
using (IUnitOfWork uow = session.StartUnitOfWork()) {
    foreach (IEntity p in complianceOfficer) {
        Console.WriteLine(p.Display);

        IEntity pid = await session.Source().CreateNewAsync("PersonInDepartment");
        await pid.PutValueAsync("UID_Department", UID_Department_Compliance);
        await pid.PutValueAsync("UID_Person", p.GetValue("UID_Person"));
        await uow.PutAsync(pid);
    }
    await uow.CommitAsync();
}

```


### Schema

Key columns for organization structure tables:

#### Department

Table name: Department

| Column  | Description | Notes   | Type       |
|-----------------------------|------------------------------------|-----------------------------------------------------|--------|
| DepartmentName              | Department name         | Department's unique name.   | nvarchar  |
| Description                 | Description                        |                                                                                                                                                                                            | nvarchar  |
| FullPath                 | Full name of department, including parent hierarchy                        |                                                                                                                                                                                            | nvarchar  |
| UID\_Department      | Unique identifier |  | varchar |
| UID\_ParentDepartment      | Unique identifier of parent department | Can be null, if object is at top of heriarcy  | varchar |


#### Location

Location records should include the following fields:

- Ident\_Locality : name of location
- City
- UID\_Country : country
- UID_\_State : state


Table name: Locality

| Column  | Description | Notes   | Type       |
|-----------------------------|------------------------------------|-----------------------------------------------------|--------|
| Ident\_Locality     | Location name         | Location's unique name.   | nvarchar  |
| Description         | Description       |                               | nvarchar  |
| FullPath            | Full name of Location, including parent hierarchy |     | nvarchar  |
| City                | City of location |   | varchar |
| UID\_Locality      | Unique identifier |  | varchar |
| UID\_ParentLocality | Unique identifier of parent Location | Can be null, if object is at top of heriarcy  | varchar |
| UID\_DialogCountry  | Unique identifier of country | DialogCountry table  | varchar |
| UID\_DialogState  | Unique identifier of state or province | DialogState table, which includes UID\_Country as key  | varchar |

*Note*: the DialogCountry table includes the following fields:

- _UID\_DialogCountry_ unique identifier
- _CountryName_ common name of country
- _ISO3166\_2_ two letter ISO 3166 designated country code \([list](https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2)\)
- _ISO3166\_3_ three letter ISO 3166 designated country code \([list](https://en.wikipedia.org/wiki/ISO_3166-1_alpha-3)\)


Sample code to insert a new location record, including state and country foreign keys:

```csharp

// fetch country record for our location (Canada)
var q = Query.From("DialogCountry").Where(c => c.Column("CountryName") == "Canada");
int count = await session.Source().GetCountAsync(q.SelectCount());
Console.WriteLine(count);

IEntityCollection countries = await session.Source().GetCollectionAsync(q.SelectAll());
IEntity canada = countries.First();



// fetch province record for our location (Ontario)
q = Query.From("DialogState").Where(
        p => p.Column("UID_DialogCountry") == canada.GetValue("UID_DialogCountry")
             &&  p.Column("Ident_DialogState") == "Ontario");
count = await session.Source().GetCountAsync(q.SelectCount());

IEntityCollection provinces = await session.Source().GetCollectionAsync(q.SelectAll());
IEntity ontario = provinces.First();



// insert Locality record
IEntity toronto = await session.Source().CreateNewAsync("Locality");
await Task.WhenAll(
        toronto.PutValueAsync("Ident_Locality", "Toronto - HQ"),
        toronto.PutValueAsync("Description", "Toronto headquarters"),
        toronto.PutValueAsync("City", "Toronto"),
        toronto.PutValueAsync("UID_DialogCountry", canada.GetValue("UID_DialogCountry")),
        toronto.PutValueAsync("UID_DialogState", ontario.GetValue("UID_DialogState"))
);

using (IUnitOfWork uow = session.StartUnitOfWork()) {
    await uow.PutAsync(toronto);
    await uow.CommitAsync();
}

```




#### Cost Center

Table name: ProfitCenter

| Column  | Description | Notes   | Type       |
|-----------------------------|------------------------------------|-----------------------------------------------------|--------|
| AccountNumber              | Accounting identifier of cost center    | Cost center's unique name.   | nvarchar  |
| Description                 | Description                        |                                                                                                                                                                                            | nvarchar  |
| FullPath                 | Full name of Profit Center, including parent hierarchy                        |                                                                                                                                                                                            | nvarchar  |
| UID\_ProfitCenter      | Unique identifier |  | varchar |
| UID\_ParentProfitCenter      | Unique identifier of parent ProfitCenter | Can be null, if object is at top of heriarcy  | varchar |




## Business Roles


Business roles provide [RBAC](https://en.wikipedia.org/wiki/Role-based_access_control) for identities.  The _Org_ table is used to store business roles in Identity Manager.

Relevent schema tables:

- _Org_ : primary table 
- _OrgRoot_ : role class
- _PersonInOrg_ : user's membership in a role
- _OrgHas\*_ : entitlement assignment to role, e.g. _OrgHasADSGroup_ 
- _DynamicGroup_ : attribute-based assignment of users to role

The role class (_OrgRoot_) defines role behaviour, e.g. are entitlements inherited top-down or bottom-up, and is required for all business roles.  Roles are assigned to users via the _PersonInOrg_ table, and entitlements are associated with roles via the OrgHas* tables. Each entitlement type will have an OrgHas* table for role assignment, e.g. Active Directory groups are assigned to roles via the _OrgHasADSGroup_ table.  The _DynamicGroup_ table is used to define attribute-based rules for role assignment, e.g. assignment of a role based on a user's job title.

### Dependencies

The following _additional_ Identity Manager DLLs are typically required:

- RMB.Customizer.dll



### Schema

Table name: Key attributes of Org table

| Column  | Description | Notes   | Type       |
|-----------------------------|------------------------------------|-----------------------------------------------------|--------|
| Ident\_Org     | Role name         | Role's unique name.   | nvarchar  |
| Description    | Description                        |              | nvarchar  |
| FullPath       | Full name of role, including role class and parent hierarchy  | Read-only   | nvarchar  |
| ShortName      | Abbreviated name                       |              | nvarchar  |
| UID\_Org       | Unique identifier |  | varchar |
| UID\_ParentOrg | Unique identifier of parent role | Can be null, if object is at top of hierarchy  | varchar |
| UID\_PersonHead | Unique identifier of role manager |  | varchar |






### Create a new role class

A role class defines two aspects of roles in the class:

- are entitlements inherited top-down or bottom-up
- what object types can be assigned to the role

The assignment types for role classes are stored in the _OrgRootAssign_ table.  All necessary _OrgRootAssign_ records are automatically created for _OrgRoot_ records.  The _IsAssignmentAllowed_ and _IsDirectAssignmentAllowed_ flags in the _OrgRootAssign_ table indicate if the object type can be assigned in general and assigned directly to roles of the given class.  Assignment target values come from the _BaseTreeAssign_ table, e.g. _QER-AsgnBT-Person_, _ADS-AsgnBT-ADSGroup_ .


```csharp

IEntity roleClass = await session.Source().CreateNewAsync("OrgRoot");
await Task.WhenAll(
    roleClass.PutValueAsync("Ident_OrgRoot", "My Role Class"),
    roleClass.PutValueAsync("Description", "Sample role class"),
    roleClass.PutValueAsync("IsTopDown", true),
);

using (IUnitOfWork uow = session.StartUnitOfWork()) {
    await uow.PutAsync(roleClass);
    await uow.CommitAsync();
}


// define what objects can be assigned to the role, using OrgRootAssign table
var q = Query.From("OrgRootAssign").Where(a => a.Column("UID_OrgRoot") == roleClass.GetValue("UID_OrgRoot"));
IEntityCollection assignments = await session.Source().GetCollectionAsync(
        q.SelectAll(), EntityCollectionLoadType.Bulk);
using (IUnitOfWork uow = session.StartUnitOfWork()) {
    foreach (IEntity a in assignments) {

        string assignTarget = a.GetValue("UID_BaseTreeAssign");
        Console.WriteLine(assignTarget);

        if ( assignTarget == "QER-AsgnBT-Person" ) {
            await a.PutValueAsync("IsAssignmentAllowed", true);
        }
        if ( assignTarget == "ADS-AsgnBT-ADSGroup" ) {
            await a.PutValueAsync("IsAssignmentAllowed", true);
            await a.PutValueAsync("IsDirectAssignmentAllowed", true);
        }
        await uow.PutAsync(a);
    }
    await uow.CommitAsync();
}

```


### Create a new business role

Role name and class are required properties.

```csharp

IEntity role = await session.Source().CreateNewAsync("Org");
await Task.WhenAll(
    role.PutValueAsync("Ident_Org", "My Role"),
    role.PutValueAsync("Description", "Sample role"),
    role.PutValueAsync("UID_OrgRoot", roleClass.GetValue("UID_OrgRoot"))
);

using (IUnitOfWork uow = session.StartUnitOfWork()) {
    await uow.PutAsync(role);
    await uow.CommitAsync();
}
```
