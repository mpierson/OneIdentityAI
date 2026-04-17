# Creating synchronization projects with SPEd

Detailed instructions for using SPEd to create Identity Manager synchronization projects are available at the [SPEd Github page](https://github.com/mpierson/OneIdentityAI/raw/refs/heads/main/OneIM/sync-project-editor/README.md).  The list of steps in this document should be followed carefully.

Notes:

- In general, it is preferred to run one SPEd command at a time, and not in batches. 
- When creating synchronization project objects, most SPEd commands return the GUID associated with the created object; capturing the GUID for future use is recommended.
- some SPEd commands take time to complete (file upload, start synchronization); it is preferred to wait for output from the command before assuming failure
- the Identity Manager database can be queried directly to read tables related to synchronization, but updates should only be performed via SPEd commands
