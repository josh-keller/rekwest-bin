# Request Bin Clone

## Features:
- Create a bin (gives you an end point)
- Log requests made to the bin
- Display the last 20 requests made to the bin

### Other Requirements:
- Bins are automatically deleted after 48 hours
- Each bin has a limit of 20 requests (FIFO)

## ERD (Actually just Mongo)

- Database in MongoDB
  - 1 collection
  - each bin is a document
  - each bin has a list of requests

### Schema:
```
bin: {
  binId: string,
  created_at: timestamp,
  requests: [
    {
      method: string,
      host: string,
      path: string,
      created: timestamp,
      parameters: {},
      headers: {},
      body: string,
    },
    {
      ...
    }
  ]
}
```



## Todos

0. Testing
   - Tests for the backend
1. MongoDB
   - One collection for all bins
   - Each bin is a document
   - Each document has an array of requests
   - Figure out how to add and slice at the same time
2. Integrate Mongo with what's working already
   - API: Abtract out the data storage part first into its own package
   - Then swap mongo in for that
3. Frontend
   - Look at templating in Go
   - Create template for the various pages
4. Deployment
5. Automatic deletion (cron job)

```


```
