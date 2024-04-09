# s3 set content type

a small program to set the content type of an object in a s3 compatible bucket using as few api calls as possible (2 - including check that everything went well) -- tested on backblaze.

## getting started

``` bash
$ go build

$ ./s3_setct
s3_setct
Set content type of an object in s3.

Ex:
OBJECTBUCKET='bucketname'
OBJECTTYPE='application/epub+zip'
OBJECTURI='https://s3.us-west-002.backblazeb2.com'
OBJECTREGION='us-west-002'
OBJECTKEYID='<s3-keyId>'
OBJECTKEY='<s3-key>'

echo 'some_ebook.epub' | s3_setct
One or more environment variables not set
```

<!-- LocalWords: contenttype Content-Type s3 object storage backblaze
     LocalWords: PUT POST REPLACE in place
     LocalWords: Github readme
 -->
