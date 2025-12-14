# Hash Index page 72
This is a simple implementation of key-value database with an in-memory hash index. The data is saved in a file in a binary format. 

Data is saved with fixed-length prefix. For each couple (key, value) the data is saved in the following format.

```CSS
[key_length(int16): 2 bytes][value_length(int16): 2 bytes][key raw string][value raw string]
```

The data is stored in the ``database.data`` file, at each start the hash index is restored from the file.

To run the code you can run the following command:

```bash
go run hashindex_db.go
```
