# XML Support

The report renderer exposes an ```xpath``` filter function for use in templates.
It is designed to run the specified xpath query on the provided XML string. In
cases where more than one result is returned, one can use a for loop to iterate
over results.

Example may be found under ```examples/xpath.tpl```.

