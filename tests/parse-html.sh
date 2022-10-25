#!/bin/bash
echo $(curl -i -X POST -H "Content-Type: multipart/form-data" -F "data=@index.html" http://localhost:8000/process/html) | cat -A
