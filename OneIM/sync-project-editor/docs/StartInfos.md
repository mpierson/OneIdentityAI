# Start Infos (DPRProjectionStartInfo)

Start Info objects define a scheduled synchronization event, including the following attributes:

 - workflow
 - variable set
 - schedule
 - root object


```bash
sped -C my_db.yaml start-info --shell 'CCC-19F48527609980498D5E843FF49BB8AD' \
        insert --name 'Full Synchronization' \
               --workflow-id 'CCC-7202478647387649AFE0B1E7F5351C22' \
               --direction 'ToTheLeft'
```

Parameters

- shell: UID\_DPRSHell of the synchronization project
- name: name of the new start info
- workflow-id: UID\_DPRProjectionConfig of the workflow to be used for synchronization
- direction: synchronization direction, ToTheLeft or ToTheRight


Mapping direction notes:

_ToTheLeft_
: Synchronize from right to left, i.e. data flows into Identity Manager

_ToTheRight_
: Synchronize from left to right, i.e. data flows into the target system


## Schedules

Assign a schedule for synchronization using the _add-schedule_ sub-command:

```bash
sped -C my_db.yaml start-info --shell 'CCC-19F48527609980498D5E843FF49BB8AD' \
        add-schedule --id 'BC7DBA15-9B97-453C-ADC0-513027CA9E63' 
                     --type 'Month' \
                     --frequency 6 \
                     --time-zone UTC \
                     --start-time '01:00'
```

Parameters

- id: UID\_DPRProjectionStartInfo of the start info object
- type: type of schedule (Hour, Day, Week, Month, Year); default is _Day_
- frequency: how often to run the schedule, in terms of type, e.g. once every **6** months; default is _1_
- time-zone: short name of time zone for schedule; default is _UTC_
- start-time: time of day to run scheduled synchronization, in **hh:mm** 24hr format; default is midnight _00:00_


Add a schedule using defaults (every day at midnight UTC):

```bash
sped -C my_db.yaml start-info --shell 'CCC-19F48527609980498D5E843FF49BB8AD' \
        add-schedule --id 'BC7DBA15-9B97-453C-ADC0-513027CA9E63'
```

## Root Object

Synchronizations may target a specific system or domain in Identity Manager, e.g. synchronization of an Active Directory domain. Other synchronizations are not specific to a target system, e.g. synchronization of a Human Resource system into Identity Manager.  All scheduled synchronization events require a target object.

For synchronization of an Active Directory domain, LDAP domain, or a generic target system represented in Identity Manager's UNS tables, the root object will correspond to the systems UNSRoot record (ADSDomain, LDAPDomain, UNSRootB, ...).


```bash
sped -C my_db.yaml start-info --shell 'CCC-19F48527609980498D5E843FF49BB8AD' \
        add-root-object --id '73C6842D-1CA5-468A-880E-5EF0C32DF4EA' \
                        --server-name 'IAMS03' \
                        --root-object-key '<Key><T>UNSRootB</T><P>aa669e4f-3d82-4882-9bb1-d88f3e412a3c</P></Key>'
```

Parameters

- shell: UID\_DPRShell of synchronization project
- id: UID\_DPRProjectionStartInfo of the start info object
- root-object-key: XObjectKey of the root object (ADSDomain, LDAPDomain, UNSRootB, etc.)
- server-name: name of the Identity Manager job server that will perform the synchronization

See **Identifying a job server for synchronization** below for help selecting an appropriate job server.


For projects that synchronize a Human Resource system, a target table is used instead of a single target record. For this type of synchronization use the _table-name_ flag instead of the _root-object-key_ flag:

```bash
sped -C my_db.yaml start-info --shell 'CCC-19F48527609980498D5E843FF49BB8AD' \
        add-root-object --id '73C6842D-1CA5-468A-880E-5EF0C32DF4EA' \
                        --server-name 'IAMS03' \
                        --table-name 'Person'
```

Parameters

- shell: UID\_DPRShell of synchronization project
- id: UID\_DPRProjectionStartInfo of the start info object
- table-name: name of table in Identity Manager schema that is primary target of synchronization, for Human Resource synchronization use _Person_
- server-name: name of the Identity Manager job server that will perform the synchronization


## Identifying a job server for synchronization

Use the following SQL to identify job servers suitable for synchronization using a custom connector:

```sql
select Ident_Server from QBMServer s
where s.IsQBMServiceInstalled=1
and exists (
   select top 1 1 from QBMServerHasServerTag sht
    join QBMServerTag t on t.UID_QBMServerTag = sht.UID_QBMServerTag
    where sht.UID_QBMServer = s.UID_QBMServer and t.Ident_QBMServerTag='Powershell Connector'
)

order by s.LastJobFetchTime desc
```

## Running synchronization

The _start-info run_ command initiates a synchronization event in Identity Manager, and if successful returns a synchronization job identifier.  The job identifier can be used to fetch status of the synchronization using the _start-info sync-status_ command.

```bash
sped -C my_db.yaml start-info --shell 'CCC-90426A03CA40354E930643DB36C87870' \
        run --id 'CCC-1AB568FD2BF4C04782887CE2F5015DAA'
```

Parameters

- shell: UID\_DPRShell of synchronization project
- id: UID\_DPRProjectionStartInfo of the start info object

**Note**: this command typically takes 1-2 minutes to complete, and the synchronization may take many minutes or hours to complete.


To check the status of a synchronization:

```bash
sped -C my_db.yaml start-info --shell 'CCC-90426A03CA40354E930643DB36C87870' \
        sync-status --id 'CCC-1AB568FD2BF4C04782887CE2F5015DAA' --job-id 'd0df1e18-aa6d-468f-94dc-ef5da4141dd7'
```

Parameters

- shell: UID\_DPRShell of synchronization project
- id: UID\_DPRProjectionStartInfo of the start info object
- job-id: UID\_Job of a running synchronization, as returned by the _start-info run_ command

**Note**: the _sync-status_ command returns "Success" for a successfully complated synchronization; it may take up to a minute for synchronization status to be available.
