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


