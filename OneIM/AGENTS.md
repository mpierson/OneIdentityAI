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

Create a new identity:

``` csharp
IEntity dbPerson = await s.Source().CreateNewAsync("Person");
await Task.WhenAll(
    dbPerson.PutValueAsync("FirstName", "Claude"),
    dbPerson.PutValueAsync("LastName", "Van Damme"),
    dbPerson.PutValueAsync("Description", "identity created via object model")
);

using (IUnitOfWork uow = s.StartUnitOfWork()) {
    await uow.PutAsync(dbPerson);
    await uow.CommitAsync();
}
```


