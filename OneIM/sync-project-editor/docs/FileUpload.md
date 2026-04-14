# File Upload

SPEd facilitates the upload of custom connector DLL files to the appropriate Identity Manager job server, using the _file insert_ command.  This SPEd command leverages Identity Manager's file management capability: file content is inserted into the database (QBMFileRevision) and distributed to server via the Job Server service.


```bash
sped -C my_db.yaml file insert --file ~/projects/connector/bin/debug/my-connector.dll --file-version '0.1'
```

Parameters

- file: full path of file to be inserted into Identity Manager database
- file-version: version string associated with file; should correspond to the version present in the connector DLL file metadata

**Note**:  Identity Manager's internal processing of the new file may take 5 to 10 minutes to complete. 


