# AI Agent instructions for Identity Manager

## Tips for building an Identity Manager environment

TODO: docker instructions

Install the tool binaries:

``` powershell
.\OneIdentityManager.9.3.1\Setup\InstallManager.Cli.exe `
    -r E:\OneIdentityManager.9.3.1\  `
    -m install -fo `
    -mod QBM QER ADS CPL ATT POL RMB RPS TSB `
    -i 'C:\Program Files\One Identity\TestInstall' `
    -d Documentation Client Client\Administration Client\Configuration Client\DevelopmentAndTesting Client\Monitoring

```

TODO: ConfigWizard.exe details


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

