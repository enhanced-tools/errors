# Server with 2 endpoints

In this example you can see how to split error messages being sent to client and seen on server log.

After running the server, you have available 2 endpoints:

- `GET /add/a/b` - returns `a + b` as a sum of 2 integers
- `GET /divide/a/b` - returns `a / b` as a division of 2 integers

You can perform curl requests to see how it works:

```bash
curl http://localhost:3000/add/1/2 # correct request
curl http://localhost:3000/divide/4/2 # correct request
curl http://localhost:3000/add/x/2 # error
curl http://localhost:3000/divide/4/0 # error
```


