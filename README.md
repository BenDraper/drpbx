## Overview


There are two apps here in the folders:

file-monitor
storage-service

For convenience, I've put both of these in the same project but, were this a real-world application, they would be two different projects with their own repositories etc.

File Monitor checks the input directory for changes and relays them to the storage service via HTTP

Storage Service listens for requests to create, update or delete files in the configured output directory and does so.

The other folders are to act as input and output.

## Running the services

File monitor:

```./run-file-monitor -folder=[folder, defaults to ./input-folder] -address=[address of writer, defaults to http://localhost:8080]```

Storage Service

```./run-storage-service  -folder=[folder, defaults to ./output-folder] -port=[port, defaults to 8080]```

NOTE: in each case the directory must exist and be accessible prior to running the services.

## Using the services

Simply add to, update or remove a file from the specified directory and see the change reflected in the output directory.


## Future Improvements

- Monitor folders recursively and send the whole structure
- Currently if changes are made to the input directory whilst the monitor is down those changes are missed. Upon startup of file-monitor get the current state of output directory and set that as "oldFiles" so that the folders can start in a synchronised state.
- Some more graceful way of handling what to do if the directory specified in the input doesn't exist. At the moment it just crashes. Maybe we create the directory?
- Error handling in general is pretty basic. It just fatals most of the time. A more comprehensive strategy should be established alongside stakeholders. Maybe retries, for example.
- Currently calls to the storage service are happening on one thread. Could make calls on file-manager.go lines 62, 66 and 70 goroutines and add a wait group to ensure synchronisation. For this prototype version this is overkill, however, as the number of changes in one cycle is likely to be small.