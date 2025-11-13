# AI Agent instructions for Identity Manager

## Tips for building an Identity Manager environment

TODO: docker instructions

## Tips for interacting with an existing instance

Establish a connection to the database:

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


