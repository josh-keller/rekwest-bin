# Request Bin Clone

## Todos

0. Testing
   - Tests for the backend
1. Integrate Mongo with what's working already
   - API: Abtract out the data storage part first into its own package
   - Then swap mongo in for that
2. MongoDB
   - One collection for all bins
   - Each bin is a document
   - Each document has an array of requests
   - Figure out how to add and slice at the same time
3. Frontend
   - Look at templating in Go
   - Create template for the various pages
4. Deployment
5. Automatic deletion (cron job)
