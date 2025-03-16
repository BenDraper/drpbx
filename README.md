## Overview


There are two apps here in the folders:

- file-monitor
- storage-service

For convenience, I've put both of these in the same project but, were this a real-world application, they would be two different projects with their own repositories etc.

File Monitor checks the input directory for changes and relays them to the storage service via HTTP

Storage Service listens for requests to create, update or delete files in the configured output directory and does so.

The other folders are to act as input and output. There's also a test files folder which I used for convenience when moving files into/out of the input folder.

## Running Go checks

Because each "project" has its own go file it's necessary to cd into the top level folder for each project. You can then run go vet ./... and go test ./... as you see fit.

## Running the services

Storage Service

```./run-storage-service  -folder=[folder, defaults to ./output-folder] -port=[port, defaults to 8080]```

File Monitor:

```./run-file-monitor -folder=[folder, defaults to ./input-folder] -address=[address of writer, defaults to http://localhost:8080]```

NOTE: Run these in the order Storage Service THEN File Monitor
NOTE: Just running the run files using the defaults should work 
NOTE: In each case the directory must exist and be accessible prior to running the services.

## Using the services

Simply add to, update or remove a file from the specified directory and see the change reflected in the output directory.


## Future Improvements

- Monitor folders recursively and send the whole structure
- Currently if changes are made to the input directory whilst the monitor is down those changes are missed. Upon startup of file-monitor get the current state of output directory and set that as "oldFiles" so that the folders can start in a synchronised state.
- Some more graceful way of handling what to do if the directory specified in the input doesn't exist. At the moment it just crashes. Maybe we create the directory?
- Error handling in general is pretty basic. It just fatals most of the time. A more comprehensive strategy should be established alongside stakeholders. Maybe retries, for example.
- Currently calls to the storage service are happening on one thread. Could make calls on file-manager.go lines 62, 66 and 70 goroutines and add a wait group to ensure synchronisation. For this prototype version this is overkill, however, as the number of changes in one cycle is likely to be small.
- Delete just uses plain text body. Might want to use JSON or something.
- Could check http method for extra security.
- Might be possible to batch requests, especially delete requests, to make things faster.
- Unit tests just check for core functionality. It might be a good idea to benchmark some of them for memory usage to ensure the multi-part file sender works to minimise memory usage.
- For the sake of time I didn't write a unit test for Local.Update() because it doesn't have any functionality outside Create and Delete. In reality I'd still test this because it's possible that how this is implemented could change in future and it would be good to capture its current functionality in tests.
- Because these apps are so IO heavy, it's hard to unit test them very well. It would be nice to have a set of integration tests where you spin up both services on one machine and run a test which creates, updates and deletes files in the target folder. You can also more easily simulate IO errors there. 


## Notes

- The kind of mocking seen in server_test.go is probably a bit less simple than merely using a custom writeToFileFunc but gomock mocks include all the testing you need for free so it makes the actual testing function easier to understand, especially as most developers should be familiar with mocking of this sort.
- The local_test.go tests aren't pure unit tests as they have side effects on a temporary directory in the test object. As this app is very IO heavy, I don't have much of a problem with this.